package opengamifylms

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/puskunalis/opengamifylms/store/db"
	
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

func createElement(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params db.CreateElementParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding request body: %v", err)
			return
		}

		element, err := courseStore.CreateElement(r.Context(), params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error creating element", zap.Error(err))
			return
		}

		if err := encode(w, r, http.StatusCreated, element); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}

func createVideoElement(logger *zap.Logger, courseStore *db.Queries, minioEndpoint string, minioAccessKeyID string, minioSecretAccessKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		submoduleID, err := strconv.ParseInt(r.PathValue("submoduleID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing submodule ID: %v", err)
			return
		}

		// Parse the multipart form
		err = r.ParseMultipartForm(200 << 20) // 200 MB max file size
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing multipart form: %v", err)
			return
		}

		// Get the video file from the form
		file, header, err := r.FormFile("video")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error getting video file: %v", err)
			return
		}
		defer file.Close()

		// Initialize Minio client
		minioClient, err := minio.New(minioEndpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(minioAccessKeyID, minioSecretAccessKey, ""),
			Secure: false,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error initializing Minio client", zap.Error(err))
			return
		}

		// Set up Minio bucket and object details
		bucketName := "course-elements-videos"

		// Create the bucket if it doesn't exist
		err = minioClient.MakeBucket(r.Context(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			// Check if the bucket already exists
			exists, errBucketExists := minioClient.BucketExists(r.Context(), bucketName)
			if errBucketExists != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Error("error creating bucket", zap.Error(err))
				return
			}

			if !exists {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Error("bucket not found after creating")
				return
			}
		}

		objectName := uuid.New().String()
		contentType := header.Header.Get("Content-Type")

		// Upload the video file to Minio
		info, err := minioClient.PutObject(r.Context(), bucketName, objectName, file, header.Size, minio.PutObjectOptions{ContentType: contentType})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error uploading video to Minio", zap.Error(err))
			return
		}

		maxOrder, err := courseStore.GetMaxElementOrderForSubmodule(r.Context(), submoduleID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error retrieving max order for submodule", zap.Error(err))
			return
		}

		// Create a new element with the dummy URL
		element, err := courseStore.CreateElement(r.Context(), db.CreateElementParams{
			SubmoduleID: submoduleID,
			Type:        "video",
			Content:     info.Key,
			Order:       maxOrder + 1,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error creating element", zap.Error(err))
			return
		}

		if err := encode(w, r, http.StatusCreated, element); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}

func getVideoElement(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		elementID, err := strconv.ParseInt(r.PathValue("elementID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing element ID: %v", err)
			return
		}

		// Get the element from the database
		element, err := courseStore.GetElementByID(r.Context(), elementID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error getting element", zap.Error(err))
			return
		}

		// Initialize Minio client
		endpoint := "minio:9000"
		accessKeyID := "user"
		secretAccessKey := "password"
		useSSL := false

		minioClient, err := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: useSSL,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error initializing Minio client", zap.Error(err))
			return
		}

		// Set up Minio bucket and object details
		bucketName := "course-elements-videos"
		objectName := element.Content

		// Get the video object from Minio
		object, err := minioClient.GetObject(r.Context(), bucketName, objectName, minio.GetObjectOptions{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error getting video object from Minio", zap.Error(err))
			return
		}
		defer object.Close()

		// Set the appropriate headers for video streaming
		w.Header().Set("Content-Type", "video/mp4")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, objectName))

		// Stream the video to the response writer
		if _, err := io.Copy(w, object); err != nil {
			logger.Error("error streaming video", zap.Error(err))
			return
		}
	})
}

func getElement(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		elementID, err := strconv.ParseInt(r.PathValue("elementID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing element ID: %v", err)
			return
		}

		element, err := courseStore.GetElementByID(r.Context(), elementID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error getting element", zap.Error(err))
			return
		}

		if err := encode(w, r, http.StatusOK, element); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}

func getElementsBySubmoduleID(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		submoduleID, err := strconv.ParseInt(r.PathValue("submoduleID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing submodule ID: %v", err)
			return
		}

		elements, err := courseStore.GetElementsBySubmoduleID(r.Context(), submoduleID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error getting elements", zap.Error(err))
			return
		}

		if err := encode(w, r, http.StatusOK, elements); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("encoding error", zap.Error(err))
			return
		}
	})
}

func updateElementOrder(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req UpdateOrderRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error decoding request: %v", err)
			return
		}

		submoduleID, err := strconv.ParseInt(r.PathValue("submoduleID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing submodule ID: %v", err)
			return
		}

		ids, orders := req.updateOrderRequestToSlices()

		params := db.UpdateElementOrderBatchParams{
			SubmoduleID: submoduleID,
			Ids:         ids,
			Orders:      orders,
		}

		err = courseStore.UpdateElementOrderBatch(r.Context(), params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error updating element order", zap.Error(err))
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func deleteElement(logger *zap.Logger, courseStore *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		elementID, err := strconv.ParseInt(r.PathValue("elementID"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error parsing element ID: %v", err)
			return
		}

		if err := courseStore.DeleteElement(r.Context(), elementID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("error deleting element", zap.Error(err))
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
