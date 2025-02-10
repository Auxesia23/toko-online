package main

import (
	"log"

	"github.com/Auxesia23/toko-online/internal/database"
	"github.com/Auxesia23/toko-online/internal/repository"
	"github.com/Auxesia23/toko-online/internal/env"
	"github.com/joho/godotenv"

)
func main(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config{
		addr: env.GetString("PORT",":8000"),
	}
	db := database.InitDB()
	userRepo := repository.NewUserRepository(db)


	app := &application{
		Config: cfg,
		User: userRepo,
	}
	
	mux := app.mount()
	log.Fatal(app.run(mux))
}