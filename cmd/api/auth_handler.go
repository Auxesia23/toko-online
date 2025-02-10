package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Auxesia23/toko-online/internal/models"
	"github.com/Auxesia23/toko-online/internal/utils"
)

func (app *application) RegisterHanlder(w http.ResponseWriter, r *http.Request) {
	var userInput models.UserInput
	err := json.NewDecoder(r.Body).Decode(&userInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Periksa apakah email sudah ada di database
	_, err = app.userRepo.GetByEmail(context.Background(), userInput.Email)
	if err == nil {
		http.Error(w, "Email already registered", http.StatusConflict) // 409 Conflict
		return
	}


	// Hash password sebelum disimpan
	hashedPassword, err := utils.HashPassword(userInput.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Buat user baru
	user := models.User{
		Name:     userInput.Name,
		Email:    userInput.Email,
		Password: hashedPassword,
	}

	newUser, err := app.userRepo.Create(context.Background(), user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Response tanpa password
	response := models.UserResponse{
		Name:  newUser.Name,
		Email: newUser.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}


func (app * application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input models.UserLogin
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User
	user, err = app.userRepo.GetByEmail(context.Background(), input.Email)
	if err != nil {
		http.Error(w, "Invalid Credentials", http.StatusConflict) 
		return
	}

	if !utils.CheckPasswordHash(input.Password,user.Password) {
		http.Error(w, "Invalid Credentials", http.StatusConflict) 
		return
	}

	token, err := utils.GenerateToken(&user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
