package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Auxesia23/toko-online/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (app *application) CreateProductHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to get image file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	name := r.FormValue("name")
	description := r.FormValue("description")
	price, err := strconv.Atoi(r.FormValue("price"))
	if err != nil {
		http.Error(w, "Invalid stock format", http.StatusBadRequest)
		return
	}

	stock, err := strconv.Atoi(r.FormValue("stock"))
	if err != nil {
		http.Error(w, "Invalid stock format", http.StatusBadRequest)
		return
	}

	imageUrl, err := app.Image.Upload(context.Background(), file, handler.Filename)
	if err != nil {
		http.Error(w, "Failed to upload image", http.StatusInternalServerError)
		return
	}

	category_id, err := strconv.Atoi(r.FormValue("category_id"))
	if err != nil {
		http.Error(w, "Invalid category_id format", http.StatusBadRequest)
		return
	}

	product, err := app.Product.Create(context.Background(), models.Product{
		Name:        name,
		CategoryID:  uint(category_id),
		Description: description,
		Price:       int32(price),
		Stock:       int16(stock),
		ImageUrl:    imageUrl,
	})
	if err != nil {
		http.Error(w, "Failed to create product", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)

}

func (app *application) GetProductsListHandler(w http.ResponseWriter, r *http.Request) {
	products, err := app.Product.GetList(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func (app *application) GetSingleProductHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product, err := app.Product.GetById(context.Background(), id)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)

}

func (app *application) UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := uuid.Parse(idStr)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")

	categoryID, err := strconv.Atoi(r.FormValue("category_id"))
	if err != nil {
		http.Error(w, "Invalid Category id", http.StatusBadRequest)
		return
	}
	priceInt, err := strconv.Atoi(r.FormValue("price"))
	if err != nil {
		http.Error(w, "Invalid price format", http.StatusBadRequest)
		return
	}
	stockInt, err := strconv.Atoi(r.FormValue("stock"))
	if err != nil {
		http.Error(w, "Invalid stock format", http.StatusBadRequest)
		return
	}

	var ImageUrl string
	file, handler, err := r.FormFile("image")
	if err == nil {
		defer file.Close()

		imageUrl, err := app.Image.Upload(context.Background(), file, handler.Filename)
		if err != nil {
			http.Error(w, "Failed to upload image", http.StatusInternalServerError)
			return
		}
		ImageUrl = imageUrl
	}

	product := models.Product{
		ID:          id,
		Name:        name,
		Description: description,
		CategoryID:  uint(categoryID),
		Price:       int32(priceInt),
		Stock:       int16(stockInt),
		ImageUrl:    ImageUrl,
	}

	updatedProduct, err := app.Product.Update(context.Background(), product)
	if err != nil {
		http.Error(w, "Failed to update product", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedProduct)
}

func (app *application) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	idSrt := chi.URLParam(r, "id")
	id, err := uuid.Parse(idSrt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = app.Product.Delete(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
