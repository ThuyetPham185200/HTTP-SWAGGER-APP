package apis

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type AuthHandler struct {
	Users map[string]User // key = username hoặc email
}

type User struct {
	ID        int
	Username  string
	Email     string
	Password  string // lưu plain-text chỉ demo, thực tế phải hash
	IsDeleted bool
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Login    string `json:"login"` // username hoặc email
	Password string `json:"password"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// RegisterRoutes registers all auth endpoints
func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/auth/register", h.Register)
	mux.HandleFunc("/auth/login", h.Login)
	mux.HandleFunc("/auth/me/password", h.ChangePassword)
	mux.HandleFunc("/auth/me", h.DeleteAccount)
}

// Register godoc
// @Summary Register a new user
// @Description Creates a new account
// @Tags auth
// @Accept  json
// @Produce  json
// @Param body body RegisterRequest true "Register data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, `{"error":"Invalid data"}`, http.StatusBadRequest)
		return
	}

	// fake ID
	newID := len(h.Users) + 1
	user := User{
		ID:       newID,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
	if h.Users == nil {
		h.Users = make(map[string]User)
	}
	h.Users[strings.ToLower(req.Username)] = user
	h.Users[strings.ToLower(req.Email)] = user

	resp := map[string]interface{}{
		"user_id": newID,
		"token":   "fake-jwt-token-" + time.Now().Format("150405"),
	}
	json.NewEncoder(w).Encode(resp)
}

// Login godoc
// @Summary Login user
// @Description Login using username or email
// @Tags auth
// @Accept  json
// @Produce  json
// @Param body body LoginRequest true "Login data"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid data"}`, http.StatusBadRequest)
		return
	}

	user, exists := h.Users[strings.ToLower(req.Login)]
	if !exists || user.Password != req.Password || user.IsDeleted {
		http.Error(w, `{"error":"Invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	resp := map[string]string{
		"token": "fake-jwt-token-" + time.Now().Format("150405"),
	}
	json.NewEncoder(w).Encode(resp)
}

// ChangePassword godoc
// @Summary Change password
// @Description Change password for the current user
// @Tags auth
// @Accept  json
// @Produce  json
// @Param body body ChangePasswordRequest true "Password data"
// @Success 200 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /me/password [put]
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	// Demo: giả sử user hiện tại là "alice"
	currentUser, exists := h.Users["alice"]
	if !exists {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusForbidden)
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid data"}`, http.StatusBadRequest)
		return
	}

	if currentUser.Password != req.OldPassword {
		http.Error(w, `{"error":"Invalid old password"}`, http.StatusForbidden)
		return
	}

	currentUser.Password = req.NewPassword
	h.Users["alice"] = currentUser
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password updated"})
}

// DeleteAccount godoc
// @Summary Soft delete current account
// @Description Mark account as deleted
// @Tags auth
// @Produce  json
// @Success 200 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /me [delete]
func (h *AuthHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	// Demo: giả sử user hiện tại là "alice"
	currentUser, exists := h.Users["alice"]
	if !exists {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusForbidden)
		return
	}
	currentUser.IsDeleted = true
	h.Users["alice"] = currentUser

	json.NewEncoder(w).Encode(map[string]string{"message": "Account soft deleted"})
}
