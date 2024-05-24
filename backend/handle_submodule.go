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

func createSubmodule(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var submodule db.Submodule
		err := json.NewDecoder(r.Body).Decode(&submodule)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding request body: %v", err)
			return
		}

		moduleID, err := strconv.ParseInt(r.PathValue("moduleID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing module ID: %v", err)
			return
		}

		submodule.ModuleID = moduleID
		newSubmodule, err := courseStore.CreateSubmodule(r.Context(), db.CreateSubmoduleParams{
			ModuleID: moduleID,
			Title:    submodule.Title,
			XpReward: submodule.XpReward,
			Order:    submodule.Order,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error creating submodule", zap.Error(err))
			return
		}

		err = encode(w, r, http.StatusCreated, newSubmodule)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}

func updateSubmoduleOrder(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req UpdateOrderRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding request: %v", err)
			return
		}

		moduleID, err := strconv.ParseInt(r.PathValue("moduleID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing module ID: %v", err)
			return
		}

		ids, orders := req.updateOrderRequestToSlices()

		params := db.UpdateSubmoduleOrderBatchParams{
			ModuleID: moduleID,
			Ids:      ids,
			Orders:   orders,
		}

		err = courseStore.UpdateSubmoduleOrderBatch(r.Context(), params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error updating submodule order", zap.Error(err))
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func deleteSubmodule(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		submoduleID, err := strconv.ParseInt(r.PathValue("submoduleID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing submodule ID: %v", err)
			return
		}

		err = courseStore.DeleteSubmodule(r.Context(), submoduleID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error deleting submodule", zap.Error(err))
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func getSubmodule(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		submoduleID, err := strconv.ParseInt(r.PathValue("submoduleID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing submodule ID: %v", err)
			return
		}

		submodule, err := courseStore.GetSubmodule(r.Context(), submoduleID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error retrieving submodule", zap.Error(err))
			return
		}

		err = encode(w, r, http.StatusOK, submodule)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}

func getSubmodulesByModuleID(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		moduleID, err := strconv.ParseInt(r.PathValue("moduleID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing module ID: %v", err)
			return
		}

		submodules, err := courseStore.GetSubmodulesByModuleID(r.Context(), moduleID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error retrieving submodules", zap.Error(err))
			return
		}

		err = encode(w, r, http.StatusOK, submodules)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}

func getUserCompletedSubmodules(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(r.PathValue("userID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing user ID: %v", err)
			return
		}

		completedSubmodules, err := courseStore.GetUserCompletedSubmodulesByUserID(r.Context(), userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error retrieving user completed submodules", zap.Error(err))
			return
		}

		err = encode(w, r, http.StatusOK, completedSubmodules)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}

func markSubmoduleCompleted(logger *zap.Logger, pool *pgxpool.Pool, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(r.PathValue("userID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing user ID: %v", err)
			return
		}

		submoduleID, err := strconv.ParseInt(r.PathValue("submoduleID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing submodule ID: %v", err)
			return
		}

		err = courseStore.CreateUserCompletedSubmodule(r.Context(), db.CreateUserCompletedSubmoduleParams{
			UserID:      userID,
			SubmoduleID: submoduleID,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error marking submodule as completed", zap.Error(err))
			return
		}

		submodule, err := courseStore.GetSubmodule(r.Context(), submoduleID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error retrieving submodule by id")
			return
		}

		err = updateUserXPAndCheckThreshold(r.Context(), pool, courseStore, int64(userID), submodule.XpReward)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "failed to update user xp")
			return
		}

		w.WriteHeader(http.StatusCreated)
	})
}

func checkSubmoduleCompleted(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(r.PathValue("userID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing user ID: %v", err)
			return
		}

		submoduleID, err := strconv.ParseInt(r.PathValue("submoduleID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing submodule ID: %v", err)
			return
		}

		completed, err := courseStore.CheckIfUserCompletedSubmodule(r.Context(), db.CheckIfUserCompletedSubmoduleParams{
			UserID:      userID,
			SubmoduleID: submoduleID,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error checking submodule completion", zap.Error(err))
			return
		}

		err = encode(w, r, http.StatusOK, map[string]bool{"completed": completed})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}

func getUserCourseProgress(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(r.PathValue("userID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing user ID: %v", err)
			return
		}

		courseID, err := strconv.ParseInt(r.PathValue("courseID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing course ID: %v", err)
			return
		}

		modules, err := courseStore.GetModulesByCourseID(r.Context(), courseID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error retrieving modules", zap.Error(err))
			return
		}

		completedSubmodules, err := courseStore.GetUserCompletedSubmodulesByUserID(r.Context(), userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error retrieving user completed submodules", zap.Error(err))
			return
		}

		type ModuleProgress struct {
			ModuleID int64   `json:"module_id"`
			Progress float64 `json:"progress"`
		}

		var moduleProgresses []ModuleProgress
		for _, module := range modules {
			submodules, err := courseStore.GetSubmodulesByModuleID(r.Context(), module.ID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Error("error retrieving submodules", zap.Error(err))
				return
			}

			completedCount := 0
			for _, submodule := range submodules {
				for _, completedSubmodule := range completedSubmodules {
					if completedSubmodule.SubmoduleID == submodule.ID {
						completedCount++
						break
					}
				}
			}

			if len(submodules) == 0 {
				moduleProgresses = append(moduleProgresses, ModuleProgress{
					ModuleID: module.ID,
					Progress: 1,
				})
				continue
			}

			progress := float64(completedCount) / float64(len(submodules))
			moduleProgresses = append(moduleProgresses, ModuleProgress{
				ModuleID: module.ID,
				Progress: progress,
			})
		}

		allModulesCompleted := true
		for _, mp := range moduleProgresses {
			if mp.Progress != 1 {
				allModulesCompleted = false
				break
			}
		}

		if allModulesCompleted {
			err := courseStore.SetUserCourseComplete(r.Context(), courseID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Error("error setting course as complete", zap.Error(err))
				return
			}
		}

		err = encode(w, r, http.StatusOK, moduleProgresses)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}
