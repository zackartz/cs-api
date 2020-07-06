package controllers

import (
	"github.com/zackartz/code-share/api/responses"
	"net/http"
)

func (s *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome to the Code Share!")
}
