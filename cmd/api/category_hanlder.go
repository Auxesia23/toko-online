package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Auxesia23/toko-online/internal/models"
)

func(app *application) CreateCategoryHanlder(w http.ResponseWriter, r *http.Request){
	var category models.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}

	newCategory, err := app.Category.Create(context.Background(),category)
	if err != nil {
		http.Error(w,"Category Already Exist",http.StatusConflict)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCategory)
}