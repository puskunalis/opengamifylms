package opengamifylms

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/puskunalis/opengamifylms/store/db"
	
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func createQuiz(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params db.CreateQuizParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding request body: %v", err)
			return
		}

		quiz, err := courseStore.CreateQuiz(r.Context(), params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error creating quiz", zap.Error(err))
			return
		}

		if err := encode(w, r, http.StatusCreated, quiz); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}

func getQuizByID(logger *zap.Logger, courseStore *db.Queries, jwtSecretKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		quizID, err := strconv.ParseInt(r.PathValue("quizID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing quiz ID: %v", err)
			return
		}

		quiz, err := courseStore.GetQuizByID(r.Context(), quizID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error getting quiz", zap.Error(err))
			return
		}

		answers, err := courseStore.GetQuizAnswersByQuizID(r.Context(), quizID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error getting quiz answers", zap.Error(err))
			return
		}

		claims, err := getAllClaimsFromToken(r.Header, jwtSecretKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("failed to get claims from token", zap.Error(err))
			return
		}

		answeredByUser, err := courseStore.CheckIfUserAnsweredQuiz(r.Context(), db.CheckIfUserAnsweredQuizParams{
			UserID: claims.ID,
			QuizID: quiz.ID,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error getting quiz answers", zap.Error(err))
			return
		}

		type QuizAnswer struct {
			ID         int64  `json:"id"`
			AnswerText string `json:"answer_text"`
		}

		type Quiz struct {
			db.Quiz
			Answers        []QuizAnswer `json:"quiz_answers"`
			AnsweredByUser bool         `json:"answered_by_user"`
		}

		q := Quiz{
			Quiz:           quiz,
			Answers:        make([]QuizAnswer, 0, len(answers)),
			AnsweredByUser: answeredByUser,
		}

		for _, ans := range answers {
			q.Answers = append(q.Answers, QuizAnswer{
				ID:         ans.ID,
				AnswerText: ans.AnswerText,
			})
		}

		if err := encode(w, r, http.StatusOK, q); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}

func getAllQuizzesBySubmoduleID(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		submoduleID, err := strconv.ParseInt(r.PathValue("submoduleID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing submodule ID: %v", err)
			return
		}

		quizzes, err := courseStore.GetAllQuizzesBySubmoduleID(r.Context(), submoduleID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error getting quizzes", zap.Error(err))
			return
		}

		if err := encode(w, r, http.StatusOK, quizzes); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}

func updateQuiz(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		quizID, err := strconv.ParseInt(r.PathValue("quizID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing quiz ID: %v", err)
			return
		}

		var params db.UpdateQuizParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding request body: %v", err)
			return
		}
		params.ID = quizID

		if err := courseStore.UpdateQuiz(r.Context(), params); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error updating quiz", zap.Error(err))
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func deleteQuiz(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		quizID, err := strconv.ParseInt(r.PathValue("quizID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing quiz ID: %v", err)
			return
		}

		if err := courseStore.DeleteQuiz(r.Context(), quizID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error deleting quiz", zap.Error(err))
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func createQuizAnswer(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		quizID, err := strconv.ParseInt(r.PathValue("quizID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing quiz ID: %v", err)
			return
		}

		var params db.CreateQuizAnswerParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding request body: %v", err)
			return
		}
		params.QuizID = quizID

		answer, err := courseStore.CreateQuizAnswer(r.Context(), params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error creating quiz answer", zap.Error(err))
			return
		}

		if err := encode(w, r, http.StatusCreated, answer); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}

func getQuizAnswersByQuizID(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		quizID, err := strconv.ParseInt(r.PathValue("quizID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing quiz ID: %v", err)
			return
		}

		answers, err := courseStore.GetQuizAnswersByQuizID(r.Context(), quizID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error getting quiz answers", zap.Error(err))
			return
		}

		if err := encode(w, r, http.StatusOK, answers); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}

func updateQuizAnswer(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		answerID, err := strconv.ParseInt(r.PathValue("answerID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing answer ID: %v", err)
			return
		}

		var params db.UpdateQuizAnswerParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding request body: %v", err)
			return
		}
		params.ID = answerID

		if err := courseStore.UpdateQuizAnswer(r.Context(), params); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error updating quiz answer", zap.Error(err))
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func deleteQuizAnswer(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		answerID, err := strconv.ParseInt(r.PathValue("answerID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing answer ID: %v", err)
			return
		}

		if err := courseStore.DeleteQuizAnswer(r.Context(), answerID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error deleting quiz answer", zap.Error(err))
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func checkQuizAnswers(logger *zap.Logger, pool *pgxpool.Pool, courseStore *db.Queries, jwtSecretKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		quizID, err := strconv.ParseInt(r.PathValue("quizID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing quiz ID: %v", err)
			return
		}

		var receivedAnswers []string
		if err := json.NewDecoder(r.Body).Decode(&receivedAnswers); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding request body: %v", err)
			return
		}

		answers, err := courseStore.GetQuizAnswersByQuizID(r.Context(), quizID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error getting quiz answers")
			return
		}

		answerCheckMap := make(map[string]int)

		for _, ans := range answers {
			if !ans.IsCorrect {
				continue
			}
			answerCheckMap[ans.AnswerText]++
		}
		for _, ans := range receivedAnswers {
			answerCheckMap[ans]--
		}

		allCorrect := true
		for _, v := range answerCheckMap {
			if v != 0 {
				allCorrect = false
				break
			}
		}

		claims, err := getAllClaimsFromToken(r.Header, jwtSecretKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("failed to get claims from token", zap.Error(err))
			return
		}

		userAlreadyAnsweredQuiz, err := courseStore.CheckIfUserAnsweredQuiz(r.Context(), db.CheckIfUserAnsweredQuizParams{
			UserID: claims.ID,
			QuizID: quizID,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error checking quiz answered status", zap.Error(err))
			return
		}

		if userAlreadyAnsweredQuiz {
			w.WriteHeader(http.StatusConflict)
			fmt.Fprintf(w, "quiz already answered")
			return
		}

		dryRun := r.URL.Query().Get("dryRun") == "true"

		if allCorrect && !dryRun {
			quiz, err := courseStore.GetQuizByID(r.Context(), quizID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Error("error getting quiz", zap.Error(err))
				return
			}

			err = courseStore.SetQuizAnsweredByUser(r.Context(), db.SetQuizAnsweredByUserParams{
				UserID: claims.ID,
				QuizID: quizID,
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "failed to set quiz as answered")
				return
			}

			err = updateUserXPAndCheckThreshold(r.Context(), pool, courseStore, claims.ID, quiz.XpReward)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "failed to update user xp")
				return
			}
		}

		if err := encode(w, r, http.StatusOK, map[string]bool{"correct": allCorrect}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}
