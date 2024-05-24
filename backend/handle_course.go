package opengamifylms

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/puskunalis/opengamifylms/store/db"
	"go.uber.org/zap"
)

type CreateCourseRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (r CreateCourseRequest) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	if len(r.Title) < 3 {
		problems["title"] = "title must be at least 3 characters long"
	}

	if len(r.Description) < 10 {
		problems["description"] = "description must be at least 10 characters long"
	}

	return problems
}

func createCourse(logger *zap.Logger, courseStore *db.Queries, jwtSecretKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, problems, err := decodeValid[CreateCourseRequest](r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding request, problems: %v, error: %v", problems, err)
			return
		}

		// Extract the instructor ID from the JWT token
		claims, err := getAllClaimsFromToken(r.Header, jwtSecretKey)
		if err != nil {
			logger.Error("decode token error", zap.Error(err))
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "invalid or missing token")
			return
		}

		c, err := courseStore.CreateCourse(r.Context(), db.CreateCourseParams{
			Title:        req.Title,
			Description:  req.Description,
			InstructorID: claims.ID,
			Icon:         fmt.Sprintf("https://picsum.photos/seed/%d/768/432", rand.Int()), // TODO
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error creating course", zap.Error(err))
			return
		}

		err = encode(w, r, http.StatusCreated, c)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}

func getCourse(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding id")
			return
		}

		c, err := courseStore.GetCourseByID(r.Context(), int64(id))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "error: %v", err)
			return
		}

		err = encode(w, r, http.StatusOK, c)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}

func getAllCourses(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		courses, err := courseStore.GetAllAvailableAndPublishedCourses(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error getting all courses", zap.Error(err))
			return
		}

		err = encode(w, r, http.StatusOK, courses)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}

type UpdateCourseRequest struct {
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	InstructorID int64  `json:"instructor_id"`
}

func (r UpdateCourseRequest) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	if len(r.Title) < 3 {
		problems["title"] = "title must be at least 3 characters long"
	}

	if len(r.Description) < 10 {
		problems["description"] = "description must be at least 10 characters long"
	}

	return problems
}

func updateCourse(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding id")
			return
		}

		req, problems, err := decodeValid[UpdateCourseRequest](r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding request, problems: %v, error: %v", problems, err)
			return
		}

		req.ID = int64(id)

		updatedCourse, err := courseStore.UpdateCourse(r.Context(), db.UpdateCourseParams{
			ID:           req.ID,
			Title:        req.Title,
			Description:  req.Description,
			InstructorID: req.InstructorID,
		})
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "error: %v", err)
			return
		}

		err = encode(w, r, http.StatusOK, updatedCourse)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}

func publishCourse(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding id")
			return
		}

		course, err := courseStore.GetCourseByID(r.Context(), int64(id))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "error: %v", err)
			return
		}

		if course.Available && course.Published {
			w.WriteHeader(http.StatusNotModified)
			fmt.Fprintf(w, "course already published")
			return
		}

		course.Available = true
		course.Published = true

		updatedCourse, err := courseStore.UpdateCourse(r.Context(), db.UpdateCourseParams{
			ID:           course.ID,
			Title:        course.Title,
			Description:  course.Description,
			InstructorID: course.InstructorID,
			Published:    true,
			Available:    true,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error updating course: %v", err)
			return
		}

		err = encode(w, r, http.StatusOK, updatedCourse)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}

type ChangeCourseAvailabilityRequest struct {
	Available bool `json:"available"`
}

func changeCourseAvailability(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding id")
			return
		}

		course, err := courseStore.GetCourseByID(r.Context(), int64(id))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "error: %v", err)
			return
		}

		req, err := decode[ChangeCourseAvailabilityRequest](r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding request, error: %v", err)
			return
		}

		course.Available = req.Available

		updatedCourse, err := courseStore.UpdateCourse(r.Context(), db.UpdateCourseParams{
			ID:           course.ID,
			Title:        course.Title,
			Description:  course.Description,
			InstructorID: course.InstructorID,
			Published:    course.Published,
			Available:    req.Available,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error updating course: %v", err)
			return
		}

		err = encode(w, r, http.StatusOK, updatedCourse)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}

func deleteCourse(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding id")
			return
		}

		err = courseStore.DeleteCourse(r.Context(), int64(id))
		if err != nil {
			logger.Error("error deleting", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func getModule(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		moduleID, err := strconv.Atoi(r.PathValue("moduleID"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding module ID")
			return
		}

		module, err := courseStore.GetModuleByID(r.Context(), int64(moduleID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error getting module by ID", zap.Error(err))
			return
		}

		err = encode(w, r, http.StatusOK, module)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}

func getCoursesForInstructor(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		instructorID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding instructor ID")
			return
		}

		courses, err := courseStore.GetCoursesForInstructor(r.Context(), int64(instructorID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error getting courses for instructor", zap.Error(err))
			return
		}

		err = encode(w, r, http.StatusOK, courses)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}
