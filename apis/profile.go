package apis

import (
	"encoding/json"
	"net/http"
)

// Profile struct for returning user profile data
type Profile struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Avatar    string `json:"avatar"`
	Bio       string `json:"bio"`
	CreatedAt string `json:"createdAt"`
}

type ProfileUpdate struct {
	Avatar   *string `json:"avatar,omitempty"`
	Username *string `json:"username,omitempty"`
	Bio      *string `json:"bio,omitempty"`
}

type ProfileHandler struct{}

// RegisterRoutes registers all profile endpoints
func (h *ProfileHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/users/", h.GetUserProfile) // /users/{user_id}
	mux.HandleFunc("/me", h.UpdateOwnProfile)   // PATCH /me
	mux.HandleFunc("/users", h.SearchUsers)     // /users?search=...
}

// @Summary Get user profile
// @Description Get profile by user_id (may require auth for private profiles)
// @Tags profile
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} Profile
// @Failure 403 {object} map[string]string "Private profile"
// @Router /users/{user_id} [get]
func (h *ProfileHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	// Mock example
	profile := Profile{
		UserID:    "123",
		Username:  "john_doe",
		Avatar:    "https://example.com/avatar.jpg",
		Bio:       "Hello there!",
		CreatedAt: "2025-08-14T12:00:00Z",
	}
	json.NewEncoder(w).Encode(profile)
}

// @Summary Update own profile
// @Description Update avatar, username, or bio
// @Tags profile
// @Accept json
// @Produce json
// @Param data body ProfileUpdate true "Profile update data"
// @Success 200 {object} map[string]string "message: Profile updated"
// @Failure 400 {object} map[string]string "error: Invalid data"
// @Router /me [patch]
func (h *ProfileHandler) UpdateOwnProfile(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"message": "Profile updated"}`))
}

// @Summary Search users
// @Description Search users by query with pagination and sorting
// @Tags profile
// @Produce json
// @Param search query string false "Search query"
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort by field"
// @Success 200 {object} map[string]interface{} "users list and total count"
// @Failure 400 {object} map[string]string "error: Invalid parameters"
// @Router /users [get]
func (h *ProfileHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{
		"users": []map[string]string{
			{"user_id": "123", "username": "john_doe", "avatar": "https://example.com/a.jpg"},
		},
		"total": 1,
	}
	json.NewEncoder(w).Encode(resp)
}
