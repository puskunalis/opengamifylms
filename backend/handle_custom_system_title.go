package opengamifylms

import (
	"net/http"

	"go.uber.org/zap"
)

type CustomSystemSettings struct {
	Title          string `json:"title"`
	PrimaryColor   string `json:"primary_color"`
	SecondaryColor string `json:"secondary_color"`
}

func getCustomSystemSettings(logger *zap.Logger, customSystemSettings CustomSystemSettings) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := encode(w, r, http.StatusOK, customSystemSettings)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}
