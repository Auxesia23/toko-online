package database

import (
	"fmt"
	"log"

	"github.com/Auxesia23/toko-online/internal/env"
	"github.com/Auxesia23/toko-online/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	// Ganti dengan konfigurasi database PostgreSQL kamu
	host := env.GetString("DB_HOST", "localhost")
	port := env.GetString("DB_PORT", "5432")
	user := env.GetString("DB_USER", "root")
	password := env.GetString("DB_PASS", "root")
	dbname := env.GetString("DB_NAME", "db")
	sslmode := "disable"

	// Buat DSN (Data Source Name) untuk PostgreSQL
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	// Koneksi ke PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate tabel berdasarkan model
	err = db.AutoMigrate(&models.User{}, &models.Product{}, &models.Category{},
		&models.Cart{}, &models.Order{}, &models.OrderItem{}, &models.Payment{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	return db
}
