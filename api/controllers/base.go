package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/rs/cors"
	"github.com/zackartz/code-share/api/models"
	"log"
	"net/http"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (s *Server) Initialize(DbDriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error
	DBUri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
	s.DB, err = gorm.Open(DbDriver, DBUri)
	if err != nil {
		fmt.Printf("cannot connect to this %s database", DbDriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("Connected to the %s database", DbDriver)
	}

	s.DB.Debug().AutoMigrate(&models.User{}, &models.CodeSnippet{})

	s.Router = mux.NewRouter()

	s.initializeRoutes()
}

func (s *Server) Run(addr string) {
	fmt.Println("Listening to port", addr)
	c := cors.AllowAll()
	log.Fatal(http.ListenAndServe(addr, c.Handler(s.Router)))
}
