package apis

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

// Follow represents a user follow
type Follow struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar,omitempty"`
}

// FollowResponse represents a generic follow response
type FollowResponse struct {
	Followers []Follow `json:"followers,omitempty"`
	Following []Follow `json:"following,omitempty"`
	Total     int      `json:"total,omitempty"`
	Message   string   `json:"message,omitempty"`
	Error     string   `json:"error,omitempty"`
}

// FollowsHandler handles follow endpoints
type FollowsHandler struct {
	mu        sync.Mutex
	followers map[int][]Follow // key = user_id
	following map[int][]Follow // key = user_id
}

// NewFollowsHandler constructor
func NewFollowsHandler() *FollowsHandler {
	return &FollowsHandler{
		followers: make(map[int][]Follow),
		following: make(map[int][]Follow),
	}
}

// RegisterRoutes register routes
func (h *FollowsHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/me/followers", h.GetMyFollowers).Methods("GET")
	router.HandleFunc("/me/following", h.GetMyFollowing).Methods("GET")
	router.HandleFunc("/users/{user_id}/followers", h.GetFollowers).Methods("GET")
	router.HandleFunc("/users/{user_id}/following", h.GetFollowing).Methods("GET")
	router.HandleFunc("/users/{target_user_id}/follow", h.FollowUser).Methods("POST")
	router.HandleFunc("/users/{target_user_id}/follow", h.UnfollowUser).Methods("DELETE")
}

// @Summary Get My Followers
// @Description Get list of my followers
// @Tags follows
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Success 200 {object} FollowResponse
// @Failure 401 {object} FollowResponse
// @Router /me/followers [get]
func (h *FollowsHandler) GetMyFollowers(w http.ResponseWriter, r *http.Request) {
	// TODO: implement: giả lập userID = 1
	h.GetFollowersByUserID(w, 1)
}

// @Summary Get My Following
// @Description Get list of users I am following
// @Tags follows
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Success 200 {object} FollowResponse
// @Failure 401 {object} FollowResponse
// @Router /me/following [get]
func (h *FollowsHandler) GetMyFollowing(w http.ResponseWriter, r *http.Request) {
	// TODO: implement: giả lập userID = 1
	h.GetFollowingByUserID(w, 1)
}

// @Summary Get Followers
// @Description Get followers of a user
// @Tags follows
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param Authorization header string false "Bearer token"
// @Success 200 {object} FollowResponse
// @Failure 404 {object} FollowResponse
// @Router /users/{user_id}/followers [get]
func (h *FollowsHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["user_id"])
	h.GetFollowersByUserID(w, userID)
}

// @Summary Get Following
// @Description Get users a user is following
// @Tags follows
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param Authorization header string false "Bearer token"
// @Success 200 {object} FollowResponse
// @Failure 404 {object} FollowResponse
// @Router /users/{user_id}/following [get]
func (h *FollowsHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["user_id"])
	h.GetFollowingByUserID(w, userID)
}

func (h *FollowsHandler) GetFollowersByUserID(w http.ResponseWriter, userID int) {
	h.mu.Lock()
	defer h.mu.Unlock()

	followers, ok := h.followers[userID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(FollowResponse{Error: "User not found"})
		return
	}
	json.NewEncoder(w).Encode(FollowResponse{
		Followers: followers,
		Total:     len(followers),
	})
}

func (h *FollowsHandler) GetFollowingByUserID(w http.ResponseWriter, userID int) {
	h.mu.Lock()
	defer h.mu.Unlock()

	following, ok := h.following[userID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(FollowResponse{Error: "User not found"})
		return
	}
	json.NewEncoder(w).Encode(FollowResponse{
		Following: following,
		Total:     len(following),
	})
}

// @Summary Follow User
// @Description Follow a user
// @Tags follows
// @Accept json
// @Produce json
// @Param target_user_id path int true "Target User ID"
// @Param Authorization header string true "Bearer token"
// @Success 201 {object} FollowResponse
// @Failure 400 {object} FollowResponse
// @Router /users/{target_user_id}/follow [post]
func (h *FollowsHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetID, _ := strconv.Atoi(vars["target_user_id"])

	// TODO: giả lập userID = 1
	currentID := 1

	h.mu.Lock()
	defer h.mu.Unlock()

	// kiểm tra đã follow chưa
	for _, u := range h.following[currentID] {
		if u.UserID == targetID {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(FollowResponse{Error: "Already following"})
			return
		}
	}

	user := Follow{UserID: targetID, Username: "user" + strconv.Itoa(targetID)}
	h.following[currentID] = append(h.following[currentID], user)
	h.followers[targetID] = append(h.followers[targetID], Follow{UserID: currentID, Username: "user1"})

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(FollowResponse{Message: "Followed"})
}

// @Summary Unfollow User
// @Description Unfollow a user
// @Tags follows
// @Accept json
// @Produce json
// @Param target_user_id path int true "Target User ID"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} FollowResponse
// @Failure 403 {object} FollowResponse
// @Router /users/{target_user_id}/follow [delete]
func (h *FollowsHandler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetID, _ := strconv.Atoi(vars["target_user_id"])

	currentID := 1

	h.mu.Lock()
	defer h.mu.Unlock()

	followingList := h.following[currentID]
	found := false
	for i, u := range followingList {
		if u.UserID == targetID {
			// remove from slice
			h.following[currentID] = append(followingList[:i], followingList[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(FollowResponse{Error: "Unauthorized"})
		return
	}

	// remove from followers of target
	followerList := h.followers[targetID]
	for i, u := range followerList {
		if u.UserID == currentID {
			h.followers[targetID] = append(followerList[:i], followerList[i+1:]...)
			break
		}
	}

	json.NewEncoder(w).Encode(FollowResponse{Message: "Unfollowed"})
}
