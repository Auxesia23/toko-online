package database

import (
	"github.com/Auxesia23/toko-online/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.User{},&models.Product{},models.Category{},models.Cart{},models.Order{},models.OrderItem{},models.Payment{})

	return db
}
