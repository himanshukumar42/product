package db

import (
	"log"

	"github.com/himanshuk42/product/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func Init(url string) Handler {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatalf("couldn't connect to db: %v\n", err)
	}

	db.AutoMigrate(&models.Product{})
	db.AutoMigrate(&models.StockDecreaseLog{})

	return Handler{DB: db}
}
