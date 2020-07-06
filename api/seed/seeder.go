package seed

import (
	"github.com/jinzhu/gorm"
	"github.com/zackartz/code-share/api/models"
	"log"
)

var users = []models.User{
	models.User{
		Username: "Steven victor",
		Email:    "steven@gmail.com",
		Password: "password",
	},
	models.User{
		Username: "Martin Luther",
		Email:    "luther@gmail.com",
		Password: "password",
	},
}

func Load(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&models.User{}, &models.CodeSnippet{}).Error
	if err != nil {
		log.Fatalf("cannot drop table %v", err)
	}
	err = db.Debug().AutoMigrate(&models.CodeSnippet{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("attaching foreign key error %v", err)
	}

	for i, _ := range users {
		err := db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table %v", err)
		}
	}
}
