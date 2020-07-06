package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/zackartz/code-share/api/auth"
	"github.com/zackartz/code-share/api/models"
	"github.com/zackartz/code-share/api/responses"
	"github.com/zackartz/code-share/api/utils/formaterror"
	"io/ioutil"
	"net/http"
)

func (s *Server) CreateSnippet(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	snippet := models.CodeSnippet{}
	err = json.Unmarshal(body, &snippet)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	authorId, err := auth.ExtractTokenID(r)
	if err != nil {
		authorId = 0
	}
	snippet.AuthorID = authorId
	snippet.Prepare()
	err = snippet.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	Validate(snippet, s.DB)
	snippetCreated, err := snippet.SaveCode(s.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())

		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%s", r.Host, r.RequestURI, snippetCreated.Slug))
	responses.JSON(w, http.StatusCreated, snippetCreated)
}

func (s *Server) GetSnippetBySlug(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]
	rid, err := auth.ExtractTokenID(r)
	if err != nil {
		rid = 0
	}
	user := models.User{}
	userGotten, err := user.FindByID(s.DB, uint32(rid))
	if err != nil {
		userGotten = &models.User{Admin: false}
	}
	snippet := models.CodeSnippet{}
	snippetGotten, err := snippet.FindSnippetBySlug(s.DB, slug)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	if snippetGotten.Private == true {
		if snippetGotten.AuthorID != rid {
			if userGotten.Admin != true {
				responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
			}
		}
	}
	responses.JSON(w, http.StatusOK, snippetGotten)
}

func Validate(s models.CodeSnippet, db *gorm.DB) bool {
	s.Prepare()
	_, err := s.FindSnippetBySlug(db, s.Slug)
	if err == nil {
		Validate(s, db)
	}
	return true
}
