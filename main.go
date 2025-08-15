package main

import (
	"fmt"
	"net/http"

	"http-swagger-app/apis"

	_ "http-swagger-app/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Swagger with net/http
// @version 1.0
// @description This is a sample Swagger API with net/http
// @host localhost:8080
// @BasePath /
func main() {
	// Dùng gorilla/mux router
	router := mux.NewRouter()

	// Auth Handler
	authHandler := &apis.AuthHandler{
		Users: make(map[string]apis.User),
	}
	authHandler.RegisterRoutes(router)

	// Profile Handler
	profileHandler := &apis.ProfileHandler{}
	profileHandler.RegisterRoutes(router)

	// Posts Handler
	postHandler := &apis.PostsHandler{}
	postHandler.RegisterRoutes(router)

	// Posts Handler
	reactHandler := &apis.ReactionsHandler{}
	reactHandler.RegisterRoutes(router)

	// Swagger
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	fmt.Println("Server started at :8080")
	fmt.Println("Swagger: http://localhost:8080/swagger/index.html")
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println("Server stopped:", err)
	}
}
