package opengamifylms

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/puskunalis/opengamifylms/store/db"
	
	"go.uber.org/zap"
)

type CreateEnrollmentRequest struct {
	UserID   int64 `json:"userId"`
	CourseID int64 `json:"courseId"`
}

var _ Validator = (*CreateEnrollmentRequest)(nil)

func (r CreateEnrollmentRequest) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	if r.UserID == 0 {
		problems["user_id"] = "user_id is required"
	}

	if r.CourseID == 0 {
		problems["course_id"] = "course_id is required"
	}

	return problems
}

func createEnrollment(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, problems, err := decodeValid[CreateEnrollmentRequest](r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding request, problems: %v, error: %v", problems, err)
			return
		}

		err = courseStore.EnrollUser(r.Context(), db.EnrollUserParams{
			UserID:   req.UserID,
			CourseID: req.CourseID,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("failed to enroll user", zap.Error(err))
			return
		}

		w.WriteHeader(http.StatusCreated)
	})
}

func getEnrollment(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enrollmentID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding enrollment id")
			return
		}

		enrollment, err := courseStore.GetEnrollment(r.Context(), int64(enrollmentID))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "error: %v", err)
			return
		}

		err = encode(w, r, http.StatusOK, enrollment)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}
