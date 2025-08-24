package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(handlers *Handlers) *mux.Router {
	router := mux.NewRouter()

	// Debug route
	router.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Router is working"))
	}).Methods("GET")

	// Public blog routes
	router.HandleFunc("/", handlers.ServeHomePage).Methods("GET")
	router.HandleFunc("/post/{id:[0-9]+}", handlers.ServePostDetail).Methods("GET")

	// Admin dashboard routes
	router.HandleFunc("/dashboard", handlers.ServeDashboard).Methods("GET")
	router.HandleFunc("/dashboard/create", handlers.ServeCreateForm).Methods("GET")
	router.HandleFunc("/dashboard/edit", handlers.ServeEditForm).Methods("GET")

	// API routes
	router.HandleFunc("/api/posts", handlers.CreatePost).Methods("POST")
	router.HandleFunc("/api/posts", handlers.GetAllPosts).Methods("GET")
	router.HandleFunc("/api/posts", handlers.GetPost).Methods("GET")
	router.HandleFunc("/api/posts", handlers.UpdatePost).Methods("PUT")
	router.HandleFunc("/api/posts", handlers.DeletePost).Methods("DELETE")

	// Search API routes
	router.HandleFunc("/api/search", handlers.SearchPosts).Methods("POST")
	router.HandleFunc("/api/search", handlers.SearchPostsGet).Methods("GET")

	return router
}
