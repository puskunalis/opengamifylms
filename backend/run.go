package opengamifylms

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/puskunalis/opengamifylms/store/db"

	"github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func Run(
	ctx context.Context,
	args []string,
	getenv func(string) string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() //nolint:errcheck

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.InitialInterval = time.Second
	expBackoff.MaxInterval = 3 * time.Second
	expBackoff.MaxElapsedTime = 60 * time.Second

	pool, err := pgxpool.New(ctx, getenv("DB_CONN_STR"))
	if err != nil {
		logger.Error("unable to connect to db", zap.Error(err))
		return err
	}
	defer pool.Close()

	pingDB := func() error {
		if err := pool.Ping(ctx); err != nil {
			logger.Error("db ping failed", zap.Error(err))
			return err
		}
		return nil
	}

	err = backoff.Retry(pingDB, backoff.WithContext(expBackoff, ctx))
	if err != nil {
		logger.Fatal("db ping backoff failed")
	}

	logger.Info("database connection and ping successful")

	courseStore := db.New(pool)

	srv := NewServer(
		logger,
		pool,
		courseStore,
		getenv("MINIO_ENDPOINT"),
		getenv("MINIO_ACCESS_KEY_ID"),
		getenv("MINIO_SECRET_ACCESS_KEY"),
		getenv("JWT_SECRET_KEY"),
		CustomSystemSettings{
			Title:          getenv("CUSTOM_SYSTEM_TITLE"),
			PrimaryColor:   getenv("CUSTOM_SYSTEM_PRIMARY_COLOR"),
			SecondaryColor: getenv("CUSTOM_SYSTEM_SECONDARY_COLOR"),
		},
	)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(getenv("HOST"), getenv("PORT")),
		Handler: srv,
	}
	go func() {
		log.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	// Starts a server on /ready
	go func() {
		mux := http.NewServeMux()

		mux.Handle("GET /ready", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "OK")
		}))

		readinessServer := &http.Server{
			Addr:    net.JoinHostPort(getenv("HOST"), getenv("READINESS_PORT")),
			Handler: mux,
		}

		log.Printf("readiness server listening on %s\n", readinessServer.Addr)
		if err := readinessServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving on readiness server: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()

	wg.Wait()

	return nil
}
