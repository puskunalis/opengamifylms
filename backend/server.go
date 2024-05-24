package opengamifylms

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/puskunalis/opengamifylms/store/db"
	"go.uber.org/zap"
)

func NewServer(
	logger *zap.Logger,
	pool *pgxpool.Pool,
	courseStore *db.Queries,
	minioEndpoint string,
	minioAccessKeyID string,
	minioSecretAccessKey string,
	jwtSecretKey string,
	customSystemSettings CustomSystemSettings,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		logger,
		pool,
		courseStore,
		minioEndpoint,
		minioAccessKeyID,
		minioSecretAccessKey,
		jwtSecretKey,
		customSystemSettings,
	)
	var handler http.Handler = mux
	return handler
}
