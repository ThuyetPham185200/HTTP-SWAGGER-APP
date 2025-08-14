package main

import (
	"fmt"
	"net/http"

	"http-swagger-app/apis"

	_ "http-swagger-app/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Swagger with net/http
// @version 1.0
// @description This is a sample Swagger API with net/http
// @host localhost:8080
// @BasePath /
func main() {
	mux := http.NewServeMux()

	// Auth Handler
	authHandler := &apis.AuthHandler{
		Users: make(map[string]apis.User),
	}
	authHandler.RegisterRoutes(mux)

	// Profile Handler
	profileHandler := &apis.ProfileHandler{}
	profileHandler.RegisterRoutes(mux)

	// Profile Handler
	postHandler := &apis.PostsHandler{}
	postHandler.RegisterRoutes(mux)

	// Swagger
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	fmt.Println("Server started at :8080")
	fmt.Println("Swagger: http://localhost:8080/swagger/index.html")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println("Server stopped:", err)
	}
}
