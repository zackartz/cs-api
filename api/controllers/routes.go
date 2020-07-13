package controllers

import "github.com/zackartz/code-share/api/middlewares"

func (s *Server) initializeRoutes() {
	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/api/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	// User Routes
	s.Router.HandleFunc("/api/signup", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/api/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/api/user/me", middlewares.SetMiddlewareJSON(s.GetSelf)).Methods("GET")
	s.Router.HandleFunc("/api/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/api/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/api/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	// Snippet Routes
	s.Router.HandleFunc("/api/create", middlewares.SetMiddlewareJSON(s.CreateSnippet)).Methods("POST")
	s.Router.HandleFunc("/api/snippets/{slug}", middlewares.SetMiddlewareJSON(s.GetSnippetBySlug)).Methods("GET")
}
