package apis

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// UserProfile lưu thông tin user
type UserProfile struct {
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	Avatar    string `json:"avatar,omitempty"`
	Bio       string `json:"bio,omitempty"`
	CreatedAt string `json:"createdAt"`
	IsPrivate bool
}

// ProfileHandler quản lý profile
type ProfileHandler struct {
	Users map[int]UserProfile // key = user_id
}

// RegisterRoutes đăng ký các endpoint profile
func (h *ProfileHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/users/{user_id}", h.GetProfile).Methods("GET")
	router.HandleFunc("/me", h.UpdateProfile).Methods("PATCH")
	router.HandleFunc("/users", h.SearchUsers).Methods("GET")
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get profile of a user by user_id
// @Tags profile
// @Produce json
// @Param user_id path int true "User ID"
// @Param Authorization header string false "Bearer token"
// @Success 200 {object} UserProfile
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{user_id} [get]
func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["user_id"]
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"Invalid user ID"}`, http.StatusBadRequest)
		return
	}

	user, exists := h.Users[userID]
	if !exists {
		http.Error(w, `{"error":"User not found"}`, http.StatusNotFound)
		return
	}

	// Demo: nếu profile private và không phải chính chủ
	if user.IsPrivate {
		http.Error(w, `{"error":"Private profile"}`, http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// UpdateProfile godoc
// @Summary Update own profile
// @Description Update your own profile
// @Tags profile
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param body body UserProfile true "Profile data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /me [patch]
func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var req UserProfile
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid data"}`, http.StatusBadRequest)
		return
	}

	// Demo: giả sử user hiện tại là user_id=1
	currentUser, exists := h.Users[1]
	if !exists {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusForbidden)
		return
	}

	if req.Username != "" {
		currentUser.Username = req.Username
	}
	if req.Avatar != "" {
		currentUser.Avatar = req.Avatar
	}
	if req.Bio != "" {
		currentUser.Bio = req.Bio
	}

	h.Users[1] = currentUser
	json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated"})
}

// SearchUsers godoc
// @Summary Search users
// @Description Search users by query
// @Tags profile
// @Produce json
// @Param search query string false "Search query"
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort field"
// @Param Authorization header string false "Bearer token"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /users [get]
func (h *ProfileHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("search")
	offsetStr := r.URL.Query().Get("offset")
	limitStr := r.URL.Query().Get("limit")

	offset, _ := strconv.Atoi(offsetStr)
	limit, _ := strconv.Atoi(limitStr)

	usersList := []UserProfile{}
	for _, u := range h.Users {
		if q == "" || containsIgnoreCase(u.Username, q) {
			usersList = append(usersList, u)
		}
	}

	// áp limit, offset
	end := offset + limit
	if end > len(usersList) {
		end = len(usersList)
	}
	if offset > len(usersList) {
		offset = len(usersList)
	}

	resp := map[string]interface{}{
		"users": usersList[offset:end],
		"total": len(usersList),
	}
	json.NewEncoder(w).Encode(resp)
}

// containsIgnoreCase kiểm tra substring không phân biệt hoa thường
func containsIgnoreCase(s, substr string) bool {
	return len(substr) == 0 || (len(s) >= len(substr) &&
		strings.Contains(strings.ToLower(s), strings.ToLower(substr)))
}
