package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Auxesia23/toko-online/internal/models"
)

func (app *application) GetUser(w http.ResponseWriter, r *http.Request) {
	email, ok := r.Context().Value(UserEmailContextKey).(string)
	if !ok {
		http.Error(w, "Invalid user claims", http.StatusUnauthorized)
		return
	}

	user, err := app.User.GetByEmail(context.Background(), email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	respone := models.UserResponse{
		Email: user.Email,
		Name:  user.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respone)
}
