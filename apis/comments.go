package apis

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

// Comment represents a comment
type Comment struct {
	CommentID int    `json:"comment_id"`
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	Avatar    string `json:"avatar,omitempty"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	IsDeleted bool   `json:"isDeleted"`
}

// CommentRequest represents request body for creating/updating comment
type CommentRequest struct {
	Content string `json:"content"`
}

// CommentResponse represents generic response
type CommentResponse struct {
	CommentID int    `json:"comment_id,omitempty"`
	Message   string `json:"message,omitempty"`
	Error     string `json:"error,omitempty"`
}

// GetCommentsResponse represents response for GET comments
type GetCommentsResponse struct {
	Comments []Comment `json:"comments"`
	Total    int       `json:"total"`
}

// CommentsHandler handles comment endpoints
type CommentsHandler struct {
	mu       sync.Mutex
	comments map[int][]Comment // post_id -> list of comments
	nextID   int
}

// NewCommentsHandler constructor
func NewCommentsHandler() *CommentsHandler {
	return &CommentsHandler{
		comments: make(map[int][]Comment),
		nextID:   1,
	}
}

// RegisterRoutes register routes
func (h *CommentsHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/posts/{post_id}/comments", h.GetComments).Methods("GET")
	router.HandleFunc("/posts/{post_id}/comments", h.CreateComment).Methods("POST")
	router.HandleFunc("/comments/{comment_id}", h.UpdateComment).Methods("PUT")
	router.HandleFunc("/comments/{comment_id}", h.DeleteComment).Methods("DELETE")
}

// @Summary Get Comments
// @Description Get comments of a post
// @Tags comments
// @Accept json
// @Produce json
// @Param post_id path int true "Post ID"
// @Param Authorization header string false "Bearer token"
// @Success 200 {object} GetCommentsResponse
// @Failure 404 {object} CommentResponse
// @Router /posts/{post_id}/comments [get]
func (h *CommentsHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, _ := strconv.Atoi(vars["post_id"])

	h.mu.Lock()
	defer h.mu.Unlock()

	comments, ok := h.comments[postID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(CommentResponse{Error: "Post not found"})
		return
	}

	resp := GetCommentsResponse{
		Comments: comments,
		Total:    len(comments),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// @Summary Create Comment
// @Description Create a new comment for a post
// @Tags comments
// @Accept json
// @Produce json
// @Param post_id path int true "Post ID"
// @Param Authorization header string true "Bearer token"
// @Param body body CommentRequest true "Comment body"
// @Success 201 {object} CommentResponse
// @Failure 400 {object} CommentResponse
// @Router /posts/{post_id}/comments [post]
func (h *CommentsHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, _ := strconv.Atoi(vars["post_id"])

	var req CommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Content == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommentResponse{Error: "Invalid content"})
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	comment := Comment{
		CommentID: h.nextID,
		UserID:    1, // giả lập user
		Username:  "user1",
		Content:   req.Content,
		CreatedAt: "2025-08-15T00:00:00Z",
		UpdatedAt: "2025-08-15T00:00:00Z",
		IsDeleted: false,
	}
	h.nextID++

	h.comments[postID] = append(h.comments[postID], comment)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CommentResponse{
		CommentID: comment.CommentID,
		Message:   "Comment created",
	})
}

// @Summary Update Comment
// @Description Update a comment
// @Tags comments
// @Accept json
// @Produce json
// @Param comment_id path int true "Comment ID"
// @Param Authorization header string true "Bearer token"
// @Param body body CommentRequest true "Comment body"
// @Success 200 {object} CommentResponse
// @Failure 403 {object} CommentResponse
// @Failure 404 {object} CommentResponse
// @Router /comments/{comment_id} [put]
func (h *CommentsHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID, _ := strconv.Atoi(vars["comment_id"])

	var req CommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Content == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommentResponse{Error: "Invalid content"})
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	found := false
	for postID, commentList := range h.comments {
		for i, c := range commentList {
			if c.CommentID == commentID {
				// giả lập check quyền
				c.Content = req.Content
				c.UpdatedAt = "2025-08-15T01:00:00Z"
				h.comments[postID][i] = c
				found = true
				break
			}
		}
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(CommentResponse{Error: "Comment not found"})
		return
	}

	json.NewEncoder(w).Encode(CommentResponse{Message: "Comment updated"})
}

// @Summary Delete Comment
// @Description Soft delete a comment
// @Tags comments
// @Accept json
// @Produce json
// @Param comment_id path int true "Comment ID"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} CommentResponse
// @Failure 403 {object} CommentResponse
// @Failure 404 {object} CommentResponse
// @Router /comments/{comment_id} [delete]
func (h *CommentsHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID, _ := strconv.Atoi(vars["comment_id"])

	h.mu.Lock()
	defer h.mu.Unlock()

	found := false
	for postID, commentList := range h.comments {
		for i, c := range commentList {
			if c.CommentID == commentID {
				c.IsDeleted = true
				h.comments[postID][i] = c
				found = true
				break
			}
		}
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(CommentResponse{Error: "Comment not found"})
		return
	}

	json.NewEncoder(w).Encode(CommentResponse{Message: "Comment soft deleted"})
}
