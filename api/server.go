package api

import (
	"github.com/zackartz/code-share/api/controllers"
	"os"
)

var server = controllers.Server{}

func Run() {
	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	// seed.Load(server.DB)

	server.Run(":1337")
}
