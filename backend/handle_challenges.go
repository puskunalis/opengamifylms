package opengamifylms

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/puskunalis/opengamifylms/store/db"
	"go.uber.org/zap"
)

func checkChallengeCompletion(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(r.PathValue("userID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing user ID: %v", err)
			return
		}

		// Get the number of completed courses for the user
		completedCourses, err := courseStore.GetUserCourses(r.Context(), userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error getting user completed courses", zap.Error(err))
			return
		}

		// Check if the user has completed 5 courses
		if len(completedCourses) >= 5 {
			// Award the challenge badge to the user
			if err := courseStore.AddUserBadge(r.Context(), db.AddUserBadgeParams{
				UserID:  userID,
				BadgeID: 1, // Assuming the challenge badge ID is 1
			}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Error("error adding user badge", zap.Error(err))
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	})
}
