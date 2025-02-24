package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Auxesia23/toko-online/internal/models"
)

func(app *application) CreatecartHandler(w http.ResponseWriter, r *http.Request){
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "Invalid user claims", http.StatusUnauthorized)
		return
	}

	var input models.CartInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cart, err := app.Cart.Create(context.Background(),input,userID)
	if err != nil {
		http.Error(w,err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cart)
}

func (app *application) GetCartsHandler(w http.ResponseWriter, r *http.Request){
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "Invalid user claims", http.StatusUnauthorized)
		return
	}

	carts, err := app.Cart.GetList(context.Background(),userID)
	if err != nil {
		http.Error(w,err.Error(),http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(carts)
}