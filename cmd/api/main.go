package main

import (
	"log"

	"github.com/Auxesia23/toko-online/internal/database"
	"github.com/Auxesia23/toko-online/internal/env"
	"github.com/Auxesia23/toko-online/internal/image"
	"github.com/Auxesia23/toko-online/internal/repository"
	"github.com/Auxesia23/toko-online/internal/utils"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	utils.SetupGoogleOAuth()

	cfg := config{
		addr: env.GetString("PORT", ":8000"),
	}
	db := database.InitDB()
	cld := image.InitCloudinary(env.GetString("CLOUDINARY_URL", ""))

	userRepo := repository.NewUserRepository(db)
	imageRepo := repository.NewImageRepository(cld)
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	cartRepo := repository.NewCartRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	app := &application{
		Config:   cfg,
		User:     userRepo,
		Image:    imageRepo,
		Product:  productRepo,
		Category: categoryRepo,
		Cart:     cartRepo,
		Order:    orderRepo,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
