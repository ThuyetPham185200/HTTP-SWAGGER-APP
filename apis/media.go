package apis

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

// Media represents an uploaded media
type Media struct {
	ID     int    `json:"media_id"`
	Type   string `json:"type"`
	PostID int    `json:"post_id"`
	URL    string `json:"url"`
}

// MediaResponse represents response for media operations
type MediaResponse struct {
	MediaID int    `json:"media_id,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// MediaHandler handles media endpoints
type MediaHandler struct {
	mu     sync.Mutex
	nextID int
	medias []Media
}

// NewMediaHandler constructor
func NewMediaHandler() *MediaHandler {
	return &MediaHandler{
		nextID: 1,
		medias: make([]Media, 0),
	}
}

// RegisterRoutes registers media routes
func (h *MediaHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/media", h.UploadMedia).Methods("POST")
}

// @Summary Upload Media
// @Description Upload an image or video file associated with a post
// @Tags media
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param type formData string true "Media type: image or video"
// @Param file formData file true "Media file"
// @Param post_id formData int true "ID of the associated post"
// @Success 201 {object} MediaResponse
// @Failure 400 {object} MediaResponse
// @Failure 404 {object} MediaResponse
// @Router /media [post]
func (h *MediaHandler) UploadMedia(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(MediaResponse{Error: "Invalid form data"})
		return
	}

	mediaType := r.FormValue("type")
	if mediaType != "image" && mediaType != "video" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(MediaResponse{Error: "Invalid media type"})
		return
	}

	postIDStr := r.FormValue("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(MediaResponse{Error: "Post not found"})
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(MediaResponse{Error: "File is required"})
		return
	}
	defer file.Close()

	// Save file to disk (in ./uploads/)
	uploadDir := "./uploads"
	os.MkdirAll(uploadDir, os.ModePerm)
	filename := fmt.Sprintf("%d_%s", h.nextID, filepath.Base(handler.Filename))
	dstPath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(MediaResponse{Error: "Cannot save file"})
		return
	}
	defer dst.Close()
	io.Copy(dst, file)

	// Save media info
	media := Media{
		ID:     h.nextID,
		Type:   mediaType,
		PostID: postID,
		URL:    dstPath,
	}
	h.medias = append(h.medias, media)
	h.nextID++

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(MediaResponse{
		MediaID: media.ID,
		Message: "Media uploaded",
	})
}
