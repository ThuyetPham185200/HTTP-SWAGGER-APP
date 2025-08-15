package apis

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// FeedItem represents a feed post
type FeedItem struct {
	PostID       int      `json:"post_id"`
	UserID       int      `json:"user_id"`
	Username     string   `json:"username"`
	Avatar       string   `json:"avatar,omitempty"`
	Content      string   `json:"content"`
	MediaURLs    []string `json:"media_urls,omitempty"`
	CreatedAt    string   `json:"created_at"`
	LikeCount    int      `json:"like_count"`
	CommentCount int      `json:"comment_count"`
	IsLiked      bool     `json:"is_liked"`
}

// FeedResponse represents the response of feeds
type FeedResponse struct {
	Feeds      []FeedItem `json:"feeds"`
	NextCursor string     `json:"next_cursor,omitempty"`
	Error      string     `json:"error,omitempty"`
}

// FeedsHandler handles news feed endpoints
type FeedsHandler struct {
	mu    sync.Mutex
	feeds []FeedItem
}

// NewFeedsHandler constructor
func NewFeedsHandler() *FeedsHandler {
	return &FeedsHandler{
		feeds: make([]FeedItem, 0),
	}
}

// RegisterRoutes register feed routes
func (h *FeedsHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/feeds", h.GetNewsFeed).Methods("GET")
}

// @Summary Get My News Feed
// @Description Get news feed posts
// @Tags feeds
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param before query string false "Timestamp cursor (optional)"
// @Param limit query int false "Number of posts to return"
// @Success 200 {object} FeedResponse
// @Failure 401 {object} FeedResponse
// @Router /feeds [get]
func (h *FeedsHandler) GetNewsFeed(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Giả lập userID = 1
	//currentUserID := 1

	// Lấy query param
	beforeStr := r.URL.Query().Get("before")
	limitStr := r.URL.Query().Get("limit")

	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	var beforeTime time.Time
	if beforeStr != "" {
		t, err := strconv.ParseInt(beforeStr, 10, 64)
		if err == nil {
			beforeTime = time.Unix(t, 0)
		}
	} else {
		beforeTime = time.Now()
	}

	// Lọc feed theo before timestamp
	result := []FeedItem{}
	for _, f := range h.feeds {
		created, _ := time.Parse(time.RFC3339, f.CreatedAt)
		if created.Before(beforeTime) || beforeStr == "" {
			result = append(result, f)
			if len(result) >= limit {
				break
			}
		}
	}

	nextCursor := ""
	if len(result) > 0 {
		last := result[len(result)-1]
		t, _ := time.Parse(time.RFC3339, last.CreatedAt)
		nextCursor = strconv.FormatInt(t.Unix(), 10)
	}

	json.NewEncoder(w).Encode(FeedResponse{
		Feeds:      result,
		NextCursor: nextCursor,
	})
}
