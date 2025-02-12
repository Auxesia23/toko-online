package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Auxesia23/toko-online/internal/models"
)

func (app *application) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	email, ok := r.Context().Value(UserEmailContextKey).(string)
	if !ok {
		http.Error(w, "Invalid user claims", http.StatusUnauthorized)
		return
	}

	user, err := app.User.GetByEmail(context.Background(), email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	response := models.UserResponse{
		Email: &user.Email,
		Name:  &user.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (app *application) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	email, ok := r.Context().Value(UserEmailContextKey).(string)
	if !ok {
		http.Error(w, "Invalid user claims", http.StatusUnauthorized)
		return
	}

	var input models.UserUpdate
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	

	user, err := app.User.GetByEmail(context.Background(), email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	updatedUser, err := app.User.Update(context.Background(), models.User{
		Name: input.Name,
	},user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.UserResponse{
		Name: &updatedUser.Name,
		Email: &updatedUser.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
