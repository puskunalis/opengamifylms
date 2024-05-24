package opengamifylms

import (
	"context"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/puskunalis/opengamifylms/store/db"
	"go.uber.org/zap"
)

func addRoutes(
	mux *http.ServeMux,
	l *zap.Logger,
	pool *pgxpool.Pool,
	courseStore *db.Queries,
	minioEndpoint string,
	minioAccessKeyID string,
	minioSecretAccessKey string,
	jwtSecretKey string,
	customSystemSettings CustomSystemSettings,
) {
	mux.Handle("OPTIONS /", all(l, optionsHandler()))

	// User routes
	mux.Handle("POST /api/v1/login", all(l, createUser(l, courseStore, jwtSecretKey)))
	mux.Handle("GET /api/v1/user/{id}", all(l, authUserIDMiddleware(l, getUser(l, courseStore), jwtSecretKey)))
	mux.Handle("GET /api/v1/user/{id}/enrollments", all(l, authUserIDMiddleware(l, getUserCourses(l, courseStore), jwtSecretKey)))
	mux.Handle("PUT /api/v1/user", all(l, updateUser(l, courseStore)))
	mux.Handle("DELETE /api/v1/user/{id}", all(l, authUserIDMiddleware(l, deleteUser(l, courseStore), jwtSecretKey)))
	mux.Handle("GET /api/v1/topUsersByXp", all(l, getTopUsersByXp(l, courseStore)))

	// Course routes
	mux.Handle("POST /api/v1/course", all(l, createCourse(l, courseStore, jwtSecretKey)))
	mux.Handle("GET /api/v1/course/{id}", all(l, getCourse(l, courseStore)))
	mux.Handle("GET /api/v1/course", all(l, getAllCourses(l, courseStore)))
	mux.Handle("PUT /api/v1/course/{id}", all(l, updateCourse(l, courseStore)))
	mux.Handle("PUT /api/v1/course/{id}/publish", all(l, publishCourse(l, courseStore)))
	mux.Handle("PUT /api/v1/course/{id}/availability", all(l, changeCourseAvailability(l, courseStore)))
	mux.Handle("DELETE /api/v1/course/{id}", all(l, deleteCourse(l, courseStore)))

	// Enrollment routes
	mux.Handle("POST /api/v1/enrollment", all(l, createEnrollment(l, courseStore)))
	mux.Handle("GET /api/v1/enrollment/{id}", all(l, getEnrollment(l, courseStore)))

	// Element routes
	mux.Handle("POST /api/v1/elements", all(l, createElement(l, courseStore)))
	mux.Handle("POST /api/v1/submodule/{submoduleID}/elements/video", all(l, createVideoElement(l, courseStore, minioEndpoint, minioAccessKeyID, minioSecretAccessKey)))
	mux.Handle("GET /api/v1/element/video/{elementID}", all(l, getVideoElement(l, courseStore)))
	mux.Handle("GET /api/v1/element/{elementID}", all(l, getElement(l, courseStore)))
	mux.Handle("GET /api/v1/submodule/{submoduleID}/elements", all(l, getElementsBySubmoduleID(l, courseStore)))
	mux.Handle("GET /api/v1/course/{courseID}/module/{moduleID}/submodule/{submoduleID}/elements", all(l, getElementsBySubmoduleID(l, courseStore)))
	mux.Handle("PUT /api/v1/course/{courseID}/module/{moduleID}/submodule/{submoduleID}/elements/order", all(l, updateElementOrder(l, courseStore)))
	mux.Handle("DELETE /api/v1/element/{elementID}", all(l, deleteElement(l, courseStore)))

	// Module routes
	mux.Handle("GET /api/v1/course/{courseID}/modules", all(l, getModulesByCourseID(l, courseStore)))
	mux.Handle("GET /api/v1/module/{moduleID}", all(l, getModule(l, courseStore)))
	mux.Handle("POST /api/v1/course/{courseID}/modules", all(l, createModule(l, courseStore)))
	mux.Handle("PUT /api/v1/course/{courseID}/modules/order", all(l, updateModuleOrder(l, courseStore)))
	mux.Handle("DELETE /api/v1/module/{moduleID}", all(l, deleteModule(l, courseStore)))

	// Submodule routes
	mux.Handle("GET /api/v1/module/{moduleID}/submodules", all(l, getSubmodulesByModuleID(l, courseStore)))
	mux.Handle("GET /api/v1/course/{courseID}/module/{moduleID}/submodules", all(l, getSubmodulesByModuleID(l, courseStore)))
	mux.Handle("GET /api/v1/submodule/{submoduleID}", all(l, getSubmodule(l, courseStore)))
	mux.Handle("POST /api/v1/course/{courseID}/module/{moduleID}/submodules", all(l, createSubmodule(l, courseStore)))
	mux.Handle("PUT /api/v1/course/{courseID}/module/{moduleID}/submodules/order", all(l, updateSubmoduleOrder(l, courseStore)))
	mux.Handle("DELETE /api/v1/submodule/{submoduleID}", all(l, deleteSubmodule(l, courseStore)))

	// Quiz routes
	mux.Handle("POST /api/v1/quizzes", all(l, createQuiz(l, courseStore)))
	mux.Handle("GET /api/v1/quiz/{quizID}", all(l, getQuizByID(l, courseStore, jwtSecretKey)))
	mux.Handle("GET /api/v1/submodule/{submoduleID}/quizzes", all(l, getAllQuizzesBySubmoduleID(l, courseStore)))
	mux.Handle("PUT /api/v1/quiz/{quizID}", all(l, updateQuiz(l, courseStore)))
	mux.Handle("DELETE /api/v1/quiz/{quizID}", all(l, deleteQuiz(l, courseStore)))

	// Quiz Answer routes
	mux.Handle("POST /api/v1/quiz/{quizID}/answers", all(l, createQuizAnswer(l, courseStore)))
	mux.Handle("GET /api/v1/quiz/{quizID}/answers", all(l, getQuizAnswersByQuizID(l, courseStore)))
	mux.Handle("PUT /api/v1/quiz/answers/{answerID}", all(l, updateQuizAnswer(l, courseStore)))
	mux.Handle("DELETE /api/v1/quiz/answers/{answerID}", all(l, deleteQuizAnswer(l, courseStore)))
	mux.Handle("PUT /api/v1/checkQuiz/{quizID}", all(l, checkQuizAnswers(l, pool, courseStore, jwtSecretKey)))

	// Badge routes
	mux.Handle("POST /api/v1/badges", all(l, createBadge(l, courseStore)))
	mux.Handle("GET /api/v1/badge/{id}", all(l, getBadgeByID(l, courseStore)))
	mux.Handle("GET /api/v1/badges", all(l, getAllBadges(l, courseStore)))
	mux.Handle("PUT /api/v1/badge/{id}", all(l, updateBadge(l, courseStore)))
	mux.Handle("DELETE /api/v1/badge/{id}", all(l, deleteBadge(l, courseStore)))

	// User Badge routes
	mux.Handle("POST /api/v1/user/{userID}/badge/{badgeID}", all(l, addUserBadge(l, courseStore)))
	mux.Handle("GET /api/v1/user/{userID}/badges", all(l, getUserBadges(l, courseStore)))

	// User Completed Submodule routes
	mux.Handle("GET /api/v1/user/{userID}/completedSubmodules", all(l, getUserCompletedSubmodules(l, courseStore)))
	mux.Handle("PUT /api/v1/user/{userID}/completedSubmodule/{submoduleID}", all(l, markSubmoduleCompleted(l, pool, courseStore)))
	mux.Handle("GET /api/v1/user/{userID}/completedSubmodule/{submoduleID}", all(l, checkSubmoduleCompleted(l, courseStore)))

	// Challenge routes
	mux.Handle("PUT /api/v1/user/{userID}/challenges", all(l, checkChallengeCompletion(l, courseStore)))

	mux.Handle("GET /api/v1/user/{userID}/courseProgress/{courseID}", all(l, getUserCourseProgress(l, courseStore)))

	mux.Handle("GET /api/v1/user/{id}/courses", all(l, getCoursesForInstructor(l, courseStore)))
	mux.Handle("GET /api/v1/customSystemSettings", all(l, getCustomSystemSettings(l, customSystemSettings)))
}

// Common middleware for all routes
func all(l *zap.Logger, next http.Handler) http.Handler {
	return loggingMiddleware(l, corsMiddleware(next))
}

func authUserIDMiddleware(l *zap.Logger, next http.Handler, jwtSecretKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := getAllClaimsFromToken(r.Header, jwtSecretKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			l.Error("failed to get claims from token", zap.Error(err))
			return
		}

		// Check if the user is accessing their own data or if they have admin privileges
		if !claims.IsAdministrator && r.PathValue("id") != strconv.Itoa(int(claims.ID)) {
			l.Error("user unauthorized")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// Set the user ID and admin status in the request context for further use
		ctx := context.WithValue(r.Context(), "user_id", claims.ID)              //nolint:staticcheck
		ctx = context.WithValue(ctx, "is_administrator", claims.IsAdministrator) //nolint:staticcheck
		r = r.WithContext(ctx)

		// Token is valid and user has appropriate permissions, pass the request to the next handler
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func optionsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func loggingMiddleware(l *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Info("Received request",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
		)

		next.ServeHTTP(w, r)
	})
}
