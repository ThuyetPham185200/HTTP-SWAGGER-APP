package apis

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Handler chứa danh sách usernames
type Handler struct {
	Usernames []string
}

// logRequest in thông tin request
func (h *Handler) logRequest(r *http.Request) {
	fmt.Println("Method:", r.Method)
	fmt.Println("URL Path:", r.URL.Path)
	fmt.Println("Full URL:", r.URL.String())
	fmt.Println("Query Params:")
	for key, values := range r.URL.Query() {
		fmt.Printf("  %s: %v\n", key, values)
	}
	fmt.Println("Headers:")
	for key, values := range r.Header {
		fmt.Printf("  %s: %v\n", key, values)
	}
	fmt.Println("Remote Addr:", r.RemoteAddr)
}

// GetUsername godoc
// @Summary Get all usernames
// @Description return list of usernames
// @Tags username
// @Produce  json
// @Success 200 {array} string
// @Router /username/ [get]
func (h *Handler) GetUsername(w http.ResponseWriter, r *http.Request) {
	h.logRequest(r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h.Usernames)
}

// PostUsername godoc
// @Summary Add new username
// @Description add a username to list
// @Tags username
// @Accept  json
// @Produce  json
// @Param username body map[string]string true "Username JSON"
// @Success 200 {string} string "added successfully"
// @Router /username/ [post]
func (h *Handler) PostUsername(w http.ResponseWriter, r *http.Request) {
	h.logRequest(r)

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var data map[string]string
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	for _, v := range data {
		if v != "" {
			h.Usernames = append(h.Usernames, v)
		}
	}

	fmt.Println("Updated usernames:", h.Usernames)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "added successfully"})
}
