package apis

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Post lưu thông tin bài viết
type Post struct {
	PostID    int    `json:"post_id"`
	UserID    int    `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
	MediaIDs  []int  `json:"media_ids,omitempty"`
	IsDeleted bool   `json:"-"`
}

// PostsHandler quản lý posts
type PostsHandler struct {
	Posts map[int]Post // key = post_id
}

// RegisterRoutes đăng ký các endpoint posts
func (h *PostsHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/posts/{post_id}", h.GetPost).Methods("GET")
	router.HandleFunc("/users/{user_id}/posts", h.GetUserPosts).Methods("GET")
	router.HandleFunc("/me/posts", h.GetOwnPosts).Methods("GET")
	router.HandleFunc("/posts", h.CreatePost).Methods("POST")
	router.HandleFunc("/posts/{post_id}", h.UpdatePost).Methods("PATCH")
	router.HandleFunc("/posts/{post_id}", h.DeletePost).Methods("DELETE")
}

// GetPost godoc
// @Summary Get a post by ID
// @Description Get post detail
// @Tags posts
// @Produce json
// @Param post_id path int true "Post ID"
// @Param Authorization header string false "Bearer token"
// @Success 200 {object} Post
// @Failure 404 {object} map[string]string
// @Router /posts/{post_id} [get]
func (h *PostsHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["post_id"]
	postID, _ := strconv.Atoi(idStr)

	post, exists := h.Posts[postID]
	if !exists || post.IsDeleted {
		http.Error(w, `{"error":"Post not found"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(post)
}

// GetUserPosts godoc
// @Summary Get posts of a user
// @Description Get list of posts by user_id
// @Tags posts
// @Produce json
// @Param user_id path int true "User ID"
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Param Authorization header string false "Bearer token"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Router /users/{user_id}/posts [get]
func (h *PostsHandler) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["user_id"]
	userID, _ := strconv.Atoi(idStr)

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	userPosts := []Post{}
	for _, p := range h.Posts {
		if p.UserID == userID && !p.IsDeleted {
			userPosts = append(userPosts, p)
		}
	}

	if len(userPosts) == 0 {
		http.Error(w, `{"error":"User not found"}`, http.StatusNotFound)
		return
	}

	end := offset + limit
	if end > len(userPosts) {
		end = len(userPosts)
	}
	if offset > len(userPosts) {
		offset = len(userPosts)
	}

	resp := map[string]interface{}{
		"posts": userPosts[offset:end],
		"total": len(userPosts),
	}
	json.NewEncoder(w).Encode(resp)
}

// GetOwnPosts godoc
// @Summary Get own posts
// @Description Get list of posts of current user
// @Tags posts
// @Produce json
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]interface{}
// @Router /me/posts [get]
func (h *PostsHandler) GetOwnPosts(w http.ResponseWriter, r *http.Request) {
	// Demo: current user = user_id 1
	currentUserID := 1
	r = r.WithContext(r.Context()) // for future auth middleware

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	userPosts := []Post{}
	for _, p := range h.Posts {
		if p.UserID == currentUserID && !p.IsDeleted {
			userPosts = append(userPosts, p)
		}
	}

	end := offset + limit
	if end > len(userPosts) {
		end = len(userPosts)
	}
	if offset > len(userPosts) {
		offset = len(userPosts)
	}

	resp := map[string]interface{}{
		"posts": userPosts[offset:end],
		"total": len(userPosts),
	}
	json.NewEncoder(w).Encode(resp)
}

// CreatePost godoc
// @Summary Create a post
// @Description Create a new post
// @Tags posts
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param body body Post true "Post data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /posts [post]
func (h *PostsHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req Post
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Content == "" {
		http.Error(w, `{"error":"Invalid data"}`, http.StatusBadRequest)
		return
	}

	// Demo: fake ID
	newID := len(h.Posts) + 1
	req.PostID = newID
	req.UserID = 1 // current user
	req.CreatedAt = time.Now().Format(time.RFC3339)
	h.Posts[newID] = req

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"post_id": newID,
		"message": "Post created",
	})
}

// UpdatePost godoc
// @Summary Update a post
// @Description Update content or media_ids of a post
// @Tags posts
// @Accept json
// @Produce json
// @Param post_id path int true "Post ID"
// @Param Authorization header string true "Bearer token"
// @Param body body Post true "Post update data"
// @Success 200 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /posts/{post_id} [patch]
func (h *PostsHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["post_id"]
	postID, _ := strconv.Atoi(idStr)

	post, exists := h.Posts[postID]
	if !exists || post.IsDeleted || post.UserID != 1 { // demo currentUserID=1
		http.Error(w, `{"error":"Unauthorized or not the author"}`, http.StatusForbidden)
		return
	}

	var req Post
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid data"}`, http.StatusBadRequest)
		return
	}

	if req.Content != "" {
		post.Content = req.Content
	}
	if req.MediaIDs != nil {
		post.MediaIDs = req.MediaIDs
	}
	h.Posts[postID] = post
	json.NewEncoder(w).Encode(map[string]string{"message": "Post updated"})
}

// DeletePost godoc
// @Summary Soft delete a post
// @Description Mark post as deleted
// @Tags posts
// @Produce json
// @Param post_id path int true "Post ID"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /posts/{post_id} [delete]
func (h *PostsHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["post_id"]
	postID, _ := strconv.Atoi(idStr)

	post, exists := h.Posts[postID]
	if !exists || post.IsDeleted || post.UserID != 1 { // demo currentUserID=1
		http.Error(w, `{"error":"Unauthorized or not the author"}`, http.StatusForbidden)
		return
	}

	post.IsDeleted = true
	h.Posts[postID] = post
	json.NewEncoder(w).Encode(map[string]string{"message": "Post soft deleted"})
}
