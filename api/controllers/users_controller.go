package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/zackartz/code-share/api/auth"
	"github.com/zackartz/code-share/api/models"
	"github.com/zackartz/code-share/api/responses"
	"github.com/zackartz/code-share/api/utils/formaterror"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	_, err = user.FindByEmail(s.DB, user.Email)
	if err == nil {
		formattedError := formaterror.FormatError("email already used")

		responses.ERROR(w, http.StatusBadRequest, formattedError)
		return
	}
	_, err = user.FindByUsername(s.DB, user.Username)
	if err == nil {
		formattedError := formaterror.FormatError("username already used")

		responses.ERROR(w, http.StatusBadRequest, formattedError)
		return
	}
	user.Prepare()
	err = user.Validate("")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	userCreated, err := user.SaveUser(s.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())

		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.ID))
	Return(s, userCreated, w, r)
}

func (s *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	users, err := user.Find(s.DB, 100)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	rid, err := auth.ExtractTokenID(r)
	if err != nil {
		for i := range users {
			users[i].Password = ""
			users[i].Email = ""
		}
		responses.JSON(w, http.StatusOK, users)
		return
	}
	user2, err := user.FindByID(s.DB, rid)
	if err != nil {
		for i := range users {
			users[i].Password = ""
			users[i].Email = ""
		}
		responses.JSON(w, http.StatusOK, users)
		return
	}
	if user2.Admin == true {
		for i := range users {
			users[i].Password = ""
		}
		responses.JSON(w, http.StatusOK, users)
	} else {
		for i := range users {
			users[i].Password = ""
			users[i].Email = ""
		}
		responses.JSON(w, http.StatusOK, users)
	}
}

func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	user := models.User{}
	userGotten, err := user.FindByID(s.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	Return(s, userGotten, w, r)
}

func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != uint32(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	user.Prepare()
	err = user.Validate("update")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	updatedUser, err := user.Update(s.DB, uint32(uid))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	Return(s, updatedUser, w, r)
}

func (s *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	user := models.User{}

	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	if tokenID != 0 && tokenID != uint32(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	rid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	_, err = user.Delete(s.DB, uint32(uid), rid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	responses.JSON(w, http.StatusNoContent, "")
}

func Return(s *Server, userGotten *models.User, w http.ResponseWriter, r *http.Request) {
	var user models.User
	rid, err := auth.ExtractTokenID(r)
	if err != nil {
		userGotten.Email = ""
		userGotten.Password = ""
		responses.JSON(w, http.StatusOK, userGotten)
		return
	}
	user2, err := user.FindByID(s.DB, rid)
	if err != nil {
		userGotten.Email = ""
		userGotten.Password = ""
		responses.JSON(w, http.StatusOK, userGotten)
		return
	}
	if user2.Admin == true {
		userGotten.Password = ""
		responses.JSON(w, http.StatusOK, userGotten)
	} else {
		userGotten.Email = ""
		userGotten.Password = ""
		responses.JSON(w, http.StatusOK, userGotten)
	}
}
