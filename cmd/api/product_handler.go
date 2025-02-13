package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Auxesia23/toko-online/internal/models"
)

func(app *application) CreateProductHandler(w http.ResponseWriter, r *http.Request){

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
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

	imageUrl, err := app.Image.Upload(context.Background(),file,handler.Filename)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to upload image", http.StatusInternalServerError)
		return
	}

	product, err := app.Product.Create(context.Background(),models.Product{
		Name: name,
		Description: description,
		Price: int32(price),
		Stock: int16(stock),
		ImageUrl: imageUrl,
	})
	if err != nil {
		http.Error(w, "Failed to upload image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)

}

func (app *application) GetProductsListHandler(w http.ResponseWriter, r *http.Request){
	products, err := app.Product.GetList(context.Background())
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}