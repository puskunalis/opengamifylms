package opengamifylms

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/puskunalis/opengamifylms/store/db"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	User  db.User `json:"user"`
	Token string  `json:"token"`
}

func createUser(logger *zap.Logger, courseStore *db.Queries, jwtSecretKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, err := decode[CreateUserRequest](r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding request, error: %v", err)
			return
		}

		// Encrypt the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error hashing password", zap.Error(err))
			return
		}

		u, err := courseStore.GetUserByEmail(r.Context(), req.Email)
		if err != nil {
			u, err = courseStore.CreateUser(r.Context(), db.CreateUserParams{
				Email:    req.Email,
				Password: string(hashedPassword),
				FullName: req.Email,
			})
			if err != nil {
				w.WriteHeader(http.StatusConflict)
				fmt.Fprintf(w, "error: %v", err)
				return
			}
		}

		// Compare the provided password with the encrypted password
		err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "invalid email or password")
			return
		}

		// Generate JWT
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":               u.ID,
			"is_administrator": u.IsAdministrator,
			"email":            u.Email,
			"full_name":        u.FullName,
		})

		// Sign and get the complete encoded token as a string
		tokenString, err := token.SignedString([]byte(jwtSecretKey))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("JWT signing error", zap.Error(err))
			return
		}

		err = encode(w, r, http.StatusOK, CreateUserResponse{
			User:  u,
			Token: tokenString,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}

func getUser(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding id")
			return
		}

		u, err := courseStore.GetUserByID(r.Context(), int64(id))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "error: %v", err)
			return
		}

		err = encode(w, r, http.StatusOK, u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}

func getUserCourses(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding id")
			return
		}

		u, err := courseStore.GetUserCourses(r.Context(), int64(id))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "error: %v", err)
			return
		}

		err = encode(w, r, http.StatusOK, u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}

type UpdateUserRequest struct {
	ID              int64  `json:"id"`
	Email           string `json:"email"`
	FullName        string `json:"full_name"`
	Password        string `json:"password"`
	IsAdministrator bool   `json:"is_administrator"`
	IsInstructor    bool   `json:"is_instructor"`
	Points          int    `json:"points"`
}

var _ Validator = (*UpdateUserRequest)(nil)

func (r UpdateUserRequest) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	if len(r.Email) < 6 {
		problems["email"] = "email must be at least 6 characters long"
	}

	if len(r.Email) > 64 {
		problems["email"] = "email must be at most 64 characters long"
	}

	if len(r.FullName) < 6 {
		problems["full_name"] = "full name must be at least 3 characters long"
	}

	if len(r.FullName) > 128 {
		problems["full_name"] = "full name must be at most 128 characters long"
	}

	return problems
}

func updateUser(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, problems, err := decodeValid[UpdateUserRequest](r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding request, problems: %v, error: %v", problems, err)
			return
		}

		err = courseStore.UpdateUser(r.Context(), db.UpdateUserParams{
			Email:           req.Email,
			FullName:        req.FullName,
			Password:        req.Password,
			IsAdministrator: req.IsAdministrator,
			IsInstructor:    req.IsInstructor,
			Xp:              int32(req.Points),
			ID:              req.ID,
		})
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "error: %v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func deleteUser(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding id")
			return
		}

		err = courseStore.DeleteUser(r.Context(), int64(id))
		if err != nil {
			logger.Error("error deleting", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func getTopUsersByXp(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		users, err := courseStore.GetTopUsersByXp(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error getting top users by XP", zap.Error(err))
			return
		}

		err = encode(w, r, http.StatusOK, users)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}
