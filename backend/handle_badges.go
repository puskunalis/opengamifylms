package opengamifylms

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/puskunalis/opengamifylms/store/db"
	"go.uber.org/zap"
)

func createBadge(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params db.CreateBadgeParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding request body: %v", err)
			return
		}

		badge, err := courseStore.CreateBadge(r.Context(), params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error creating badge", zap.Error(err))
			return
		}

		if err := encode(w, r, http.StatusCreated, badge); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}

func getBadgeByID(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		badgeID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing badge ID: %v", err)
			return
		}

		badge, err := courseStore.GetBadgeByID(r.Context(), badgeID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error getting badge", zap.Error(err))
			return
		}

		if err := encode(w, r, http.StatusOK, badge); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}

func getAllBadges(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		badges, err := courseStore.GetAllBadges(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error getting badges", zap.Error(err))
			return
		}

		if err := encode(w, r, http.StatusOK, badges); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}

func updateBadge(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		badgeID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing badge ID: %v", err)
			return
		}

		var params db.UpdateBadgeParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding request body: %v", err)
			return
		}
		params.ID = badgeID

		if err := courseStore.UpdateBadge(r.Context(), params); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error updating badge", zap.Error(err))
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func deleteBadge(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		badgeID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing badge ID: %v", err)
			return
		}

		if err := courseStore.DeleteBadge(r.Context(), badgeID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error deleting badge", zap.Error(err))
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func addUserBadge(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(r.PathValue("userID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing user ID: %v", err)
			return
		}

		badgeID, err := strconv.ParseInt(r.PathValue("badgeID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing badge ID: %v", err)
			return
		}

		if err := courseStore.AddUserBadge(r.Context(), db.AddUserBadgeParams{
			UserID:  userID,
			BadgeID: badgeID,
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error adding user badge", zap.Error(err))
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func getUserBadges(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(r.PathValue("userID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing user ID: %v", err)
			return
		}

		badges, err := courseStore.GetUserBadges(r.Context(), userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error getting user badges", zap.Error(err))
			return
		}

		if err := encode(w, r, http.StatusOK, badges); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}
