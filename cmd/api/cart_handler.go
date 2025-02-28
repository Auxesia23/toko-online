package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Auxesia23/toko-online/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (app *application) CreatecartHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "Invalid user claims", http.StatusUnauthorized)
		return
	}
	log.Println(userID)

	var input models.CartInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = app.Cart.Create(context.Background(), input, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *application) GetCartsHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "Invalid user claims", http.StatusUnauthorized)
		return
	}

	carts, err := app.Cart.GetList(context.Background(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(carts)
}

func (app *application) DeleteCartHandler(w http.ResponseWriter, r *http.Request){
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "Invalid user claims", http.StatusUnauthorized)
		return
	}

	cartIdStr := chi.URLParam(r, "id")
	cartID, err := uuid.Parse(cartIdStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = app.Cart.Delete(context.Background(),cartID,userID)
	if err != nil {
		http.Error(w, "Cart not found",http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) IncrementCartHandler(w http.ResponseWriter, r *http.Request){
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "Invalid user claims", http.StatusUnauthorized)
		return
	}

	cartIdStr := chi.URLParam(r, "id")
	cartID, err := uuid.Parse(cartIdStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = app.Cart.Increment(context.Background(),cartID,userID)
	if err != nil {
		http.Error(w,err.Error(),http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (app *application) DecrementsCartHandler(w http.ResponseWriter, r *http.Request){
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "Invalid user claims", http.StatusUnauthorized)
		return
	}

	cartIdStr := chi.URLParam(r, "id")
	cartID, err := uuid.Parse(cartIdStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = app.Cart.Decrement(context.Background(),cartID,userID)
	if err != nil {
		http.Error(w,err.Error(),http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}