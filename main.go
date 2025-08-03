package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	_ "example.com/http-swagger-app/docs" // for swaggo docs
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Swagger with net/http
// @version 1.0
// @description This is a sample Swagger API with net/http
// @host localhost:8080
// @BasePath /

// username list (mock)
var usernames = []string{"alice", "bob", "Cais dits con me may"}

// getUsername godoc
// @Summary Get all usernames
// @Description return list of usernames
// @Tags username
// @Produce  json
// @Success 200 {array} string
// @Router /username/ [get]
func getUsername(w http.ResponseWriter, r *http.Request) {
	// Print method
	fmt.Println("Method:", r.Method)

	// Print URL path and full URL
	fmt.Println("URL Path:", r.URL.Path)
	fmt.Println("Full URL:", r.URL.String())

	// Print query parameters
	fmt.Println("Query Params:")
	for key, values := range r.URL.Query() {
		fmt.Printf("  %s: %v\n", key, values)
	}

	// Print headers
	fmt.Println("Headers:")
	for key, values := range r.Header {
		fmt.Printf("  %s: %v\n", key, values)
	}

	// Print remote address (client IP)
	fmt.Println("Remote Addr:", r.RemoteAddr)

	// Write JSON response
	json.NewEncoder(w).Encode(usernames)
}

// postUsername godoc
// @Summary Add new username
// @Description add a username to list
// @Tags username
// @Accept  json
// @Produce  json
// @Param username body map[string]string true "Username JSON"
// @Success 200 {string} string "added successfully"
// @Router /username/ [post]
func postUsername(w http.ResponseWriter, r *http.Request) {
	// Debug basic request info
	fmt.Println("=== POST /username ===")
	fmt.Println("Method:", r.Method)
	fmt.Println("URL Path:", r.URL.Path)
	fmt.Println("Remote Addr:", r.RemoteAddr)

	// Print all headers
	fmt.Println("Headers:")
	for key, values := range r.Header {
		fmt.Printf("  %s: %v\n", key, values)
	}

	// Read full body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Now decode JSON from raw body
	var data map[string]string
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Extract and append username
	additionalProp1 := data["additionalProp1"]
	additionalProp2 := data["additionalProp2"]
	additionalProp3 := data["additionalProp3"]
	fmt.Println("Property = ", additionalProp1, additionalProp2, additionalProp3)

	usernames = append(usernames, additionalProp1)
	usernames = append(usernames, additionalProp2)
	usernames = append(usernames, additionalProp3)

	// Debug logs to console
	fmt.Println("Updated usernames:", usernames)

	// Send response
	json.NewEncoder(w).Encode("added successfully")
}

func main() {
	http.HandleFunc("/username/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getUsername(w, r)
		} else if r.Method == http.MethodPost {
			postUsername(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	fmt.Println("Open: http://localhost:8080/swagger/index.html")

	http.Handle("/swagger/", httpSwagger.WrapHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("End game! Server stopped with error:", err)
	}
}
