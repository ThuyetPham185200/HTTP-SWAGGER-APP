package apis

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

// ReactionRequest represents the request body for reacting/removing reaction
type ReactionRequest struct {
	ReactionType string `json:"reaction_type"`
}

// ReactionResponse represents generic response
type ReactionResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// GetReactionsResponse represents response for GET /posts/{post_id}/reactions
type GetReactionsResponse struct {
	Count int                 `json:"count"`
	Types []string            `json:"types"`
	Users []map[string]string `json:"users"`
	Total int                 `json:"total"`
}

// ReactionsHandler handles reactions endpoints
type ReactionsHandler struct {
	mu        sync.Mutex
	reactions map[string]map[string]string // post_id -> user_id -> reaction_type
}

// NewReactionsHandler constructor
func NewReactionsHandler() *ReactionsHandler {
	return &ReactionsHandler{
		reactions: make(map[string]map[string]string),
	}
}

// RegisterRoutes register routes with mux
func (h *ReactionsHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/posts/{post_id}/reactions", h.GetReactions).Methods("GET")
	router.HandleFunc("/posts/{post_id}/reactions", h.ReactToPost).Methods("POST")
	router.HandleFunc("/posts/{post_id}/reactions", h.RemoveReaction).Methods("DELETE")
}

// @Summary Get Reactions
// @Description Get reactions of a post
// @Tags reactions
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Param Authorization header string false "Bearer token"
// @Success 200 {object} GetReactionsResponse
// @Failure 404 {object} ReactionResponse
// @Router /posts/{post_id}/reactions [get]
func (h *ReactionsHandler) GetReactions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["post_id"]

	h.mu.Lock()
	defer h.mu.Unlock()

	postReactions, ok := h.reactions[postID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ReactionResponse{Error: "Post not found"})
		return
	}

	count := len(postReactions)
	typeSet := make(map[string]struct{})
	users := []map[string]string{}
	for userID, react := range postReactions {
		typeSet[react] = struct{}{}
		users = append(users, map[string]string{"user_id": userID, "username": userID})
	}

	types := []string{}
	for t := range typeSet {
		types = append(types, t)
	}

	resp := GetReactionsResponse{
		Count: count,
		Types: types,
		Users: users,
		Total: count,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// @Summary React to Post
// @Description Add reaction to a post
// @Tags reactions
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Param Authorization header string true "Bearer token"
// @Param body body ReactionRequest true "Reaction body"
// @Success 201 {object} ReactionResponse
// @Failure 400 {object} ReactionResponse
// @Router /posts/{post_id}/reactions [post]
func (h *ReactionsHandler) ReactToPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["post_id"]

	var req ReactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.ReactionType) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ReactionResponse{Error: "Invalid reaction type"})
		return
	}

	userID := "user1" // giả lập user
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.reactions[postID]; !ok {
		h.reactions[postID] = make(map[string]string)
	}
	h.reactions[postID][userID] = req.ReactionType

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ReactionResponse{Message: "Reaction added"})
}

// @Summary Remove Reaction
// @Description Remove reaction from a post
// @Tags reactions
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Param Authorization header string true "Bearer token"
// @Param body body ReactionRequest false "Reaction body (optional if only 1 type)"
// @Success 200 {object} ReactionResponse
// @Failure 404 {object} ReactionResponse
// @Router /posts/{post_id}/reactions [delete]
func (h *ReactionsHandler) RemoveReaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["post_id"]

	var req ReactionRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	userID := "user1" // giả lập user
	h.mu.Lock()
	defer h.mu.Unlock()

	postReactions, ok := h.reactions[postID]
	if !ok || postReactions[userID] == "" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ReactionResponse{Error: "Reaction not found"})
		return
	}

	delete(postReactions, userID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ReactionResponse{Message: "Reaction removed"})
}
