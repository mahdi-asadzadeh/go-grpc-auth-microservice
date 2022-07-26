package db

import (
	"log"

	"github.com/mahdi-asadzadeh/go-grpc-auth-microservice/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&models.User{})
}

func Init(url string) Handler {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	migrate(db)

	return Handler{db}
}
