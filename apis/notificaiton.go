package apis

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

// Notification represents a user notification
type Notification struct {
	ID           int    `json:"id"`
	Type         string `json:"type"`
	SourceUserID int    `json:"source_user_id,omitempty"`
	PostID       int    `json:"post_id,omitempty"`
	Read         bool   `json:"read"`
	CreatedAt    string `json:"created_at"`
}

// NotificationResponse represents response for list
type NotificationResponse struct {
	Notifications []Notification `json:"notifications,omitempty"`
	Total         int            `json:"total,omitempty"`
	Error         string         `json:"error,omitempty"`
}

// NotificationHandler handles notifications
type NotificationHandler struct {
	mu            sync.Mutex
	notifications []Notification
}

// NewNotificationHandler constructor
func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{
		notifications: make([]Notification, 0),
	}
}

// RegisterRoutes register notification routes
func (h *NotificationHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/notifications", h.GetNotifications).Methods("GET")
	router.HandleFunc("/notifications/{notification_id}", h.MarkAsRead).Methods("PATCH")
}

// @Summary Get Notifications
// @Description Get list of notifications
// @Tags notifications
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Success 200 {object} NotificationResponse
// @Failure 400 {object} NotificationResponse
// @Router /notifications [get]
func (h *NotificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	offsetStr := r.URL.Query().Get("offset")
	limitStr := r.URL.Query().Get("limit")

	offset, _ := strconv.Atoi(offsetStr)
	limit := 10
	if l, err := strconv.Atoi(limitStr); err == nil {
		limit = l
	}

	total := len(h.notifications)
	end := offset + limit
	if end > total {
		end = total
	}

	result := h.notifications[offset:end]

	json.NewEncoder(w).Encode(NotificationResponse{
		Notifications: result,
		Total:         total,
	})
}

// @Summary Mark Notification as Read
// @Description Mark a notification as read
// @Tags notifications
// @Accept json
// @Produce json
// @Param notification_id path int true "Notification ID"
// @Param Authorization header string true "Bearer token"
// @Param body body map[string]bool false "Optional read body"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /notifications/{notification_id} [patch]
func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	vars := mux.Vars(r)
	idStr := vars["notification_id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid notification ID"})
		return
	}

	for i, n := range h.notifications {
		if n.ID == id {
			// Giả lập current user = 1
			if n.SourceUserID != 1 {
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(map[string]string{"error": "Forbidden"})
				return
			}

			h.notifications[i].Read = true
			json.NewEncoder(w).Encode(map[string]string{"message": "Notification marked as read"})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Notification not found"})
}
