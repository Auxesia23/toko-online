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

	hashedPassword, err := utils.HashPassword(userInput.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := models.User{
		Name:     userInput.Name,
		Email:    userInput.Email,
		Password: hashedPassword,
	}

	newUser, err := app.User.Create(context.Background(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func (app *application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input models.UserLogin
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := app.User.Login(context.Background(), input.Email, input.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
