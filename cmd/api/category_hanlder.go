package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Auxesia23/toko-online/internal/models"
	"github.com/go-chi/chi/v5"
)

func (app *application) CreateCategoryHanlder(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newCategory, err := app.Category.Create(context.Background(), category)
	if err != nil {
		http.Error(w, "Category Already Exist", http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCategory)
}

func (app *application) GetCategoryListHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := app.Category.GetList(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(categories)
}

func (app *application) GetCategoryHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category id format", http.StatusBadRequest)
		return
	}
	category, err := app.Category.GetByID(context.Background(), uint(id))
	if err != nil {
		http.Error(w,err.Error(),http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

func (app *application) DeleteCategoryHandler(w http.ResponseWriter, r *http.Request){
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category id format", http.StatusBadRequest)
		return
	}

	err = app.Category.Delete(context.Background(),uint(id))
	if err != nil {
		http.Error(w,err.Error(),http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
