// apis/posts.go
package apis

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// Post model
type Post struct {
	PostID    int    `json:"post_id"`
	UserID    int    `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
	MediaIDs  []int  `json:"media_ids,omitempty"`
}

// CreatePostRequest request body
type CreatePostRequest struct {
	Content  string `json:"content"`
	MediaIDs []int  `json:"media_ids"`
}

// UpdatePostRequest request body
type UpdatePostRequest struct {
	Content  *string `json:"content,omitempty"`
	MediaIDs []int   `json:"media_ids,omitempty"`
}

// PostsHandler handler struct
type PostsHandler struct{}

// RegisterRoutes registers post routes
func (h *PostsHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/posts/{post_id}", h.GetPost).Methods("GET")
	r.HandleFunc("/users/{user_id}/posts", h.GetUserPosts).Methods("GET")
	r.HandleFunc("/me/posts", h.GetOwnPosts).Methods("GET")
	r.HandleFunc("/posts", h.CreatePost).Methods("POST")
	r.HandleFunc("/posts/{post_id}", h.UpdatePost).Methods("PATCH")
	r.HandleFunc("/posts/{post_id}", h.DeletePost).Methods("DELETE")
}

// GetPost godoc
// @Summary Get a post
// @Tags Posts
// @Param post_id path int true "Post ID"
// @Success 200 {object} Post
// @Failure 404 {object} map[string]string
// @Router /posts/{post_id} [get]
func (h *PostsHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID := params["post_id"]
	_ = postID

	// Dummy response
	post := Post{PostID: 1, UserID: 2, Content: "Hello world", CreatedAt: "2025-08-14"}
	json.NewEncoder(w).Encode(post)
}

// GetUserPosts godoc
// @Summary Get posts by user
// @Tags Posts
// @Param user_id path int true "User ID"
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Router /users/{user_id}/posts [get]
func (h *PostsHandler) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["user_id"]
	_ = userID
	posts := []Post{{PostID: 1, UserID: 2, Content: "User post", CreatedAt: "2025-08-14"}}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"posts": posts,
		"total": len(posts),
	})
}

// GetOwnPosts godoc
// @Summary Get own posts
// @Tags Posts
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Success 200 {object} map[string]interface{}
// @Router /me/posts [get]
func (h *PostsHandler) GetOwnPosts(w http.ResponseWriter, r *http.Request) {
	posts := []Post{{PostID: 1, UserID: 1, Content: "My post", CreatedAt: "2025-08-14"}}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"posts": posts,
		"total": len(posts),
	})
}

// CreatePost godoc
// @Summary Create a new post
// @Tags Posts
// @Accept json
// @Produce json
// @Param request body CreatePostRequest true "Post data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /posts [post]
func (h *PostsHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Content) == "" {
		http.Error(w, `{"error": "Invalid data"}`, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"post_id": 1,
		"message": "Post created",
	})
}

// UpdatePost godoc
// @Summary Update a post
// @Tags Posts
// @Accept json
// @Produce json
// @Param post_id path int true "Post ID"
// @Param request body UpdatePostRequest true "Updated post data"
// @Success 200 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /posts/{post_id} [patch]
func (h *PostsHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID := params["post_id"]
	_, _ = strconv.Atoi(postID)

	var req UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid data"}`, http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Post updated",
	})
}

// DeletePost godoc
// @Summary Delete a post
// @Tags Posts
// @Param post_id path int true "Post ID"
// @Success 200 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /posts/{post_id} [delete]
func (h *PostsHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Post soft deleted",
	})
}
