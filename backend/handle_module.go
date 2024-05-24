package opengamifylms

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/puskunalis/opengamifylms/store/db"
	"go.uber.org/zap"
)

func getModulesByCourseID(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		courseID, err := strconv.ParseInt(r.PathValue("courseID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing course ID: %v", err)
			return
		}

		modules, err := courseStore.GetModulesByCourseID(r.Context(), courseID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error getting modules", zap.Error(err))
			return
		}

		err = encode(w, r, http.StatusOK, modules)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}

func createModule(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var module db.CreateModuleParams
		err := json.NewDecoder(r.Body).Decode(&module)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding request body: %v", err)
			return
		}

		courseID, err := strconv.ParseInt(r.PathValue("courseID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing course ID: %v", err)
			return
		}

		module.CourseID = courseID
		newModule, err := courseStore.CreateModule(r.Context(), module)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error creating module", zap.Error(err))
			return
		}

		err = encode(w, r, http.StatusCreated, newModule)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}

func updateModuleOrder(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req UpdateOrderRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding request: %v", err)
			return
		}

		courseID, err := strconv.ParseInt(r.PathValue("courseID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing course ID: %v", err)
			return
		}

		ids, orders := req.updateOrderRequestToSlices()

		params := db.UpdateModuleOrderBatchParams{
			CourseID: courseID,
			Ids:      ids,
			Orders:   orders,
		}

		err = courseStore.UpdateModuleOrderBatch(r.Context(), params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error updating module order", zap.Error(err))
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func deleteModule(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		moduleID, err := strconv.ParseInt(r.PathValue("moduleID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing module ID: %v", err)
			return
		}

		err = courseStore.DeleteModule(r.Context(), moduleID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error deleting module", zap.Error(err))
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
