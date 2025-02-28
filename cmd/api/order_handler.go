package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Auxesia23/toko-online/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (app *application) PreviewOrderHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "Invalid user claims", http.StatusUnauthorized)
		return
	}

	var input models.OrderInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	preview, err := app.Order.Preview(context.Background(), userID, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&preview)
}

func (app *application) CreateOrderHanlder(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "Invalid user claims", http.StatusUnauthorized)
		return
	}

	var input models.OrderInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderID, err := app.Order.Create(context.Background(), userID, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := models.OrderID{
		ID: &orderID,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&response)
}

func (app *application) GetListOrderhanlder(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "Invalid user claims", http.StatusUnauthorized)
		return
	}

	orders, err := app.Order.GetList(context.Background(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&orders)
}

func (app *application) GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "Invalid user claims", http.StatusUnauthorized)
		return
	}

	orderIDStr := chi.URLParam(r, "id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := app.Order.GetByID(context.Background(), userID, orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&order)
}

func (app *application) CreatePaymentHandler(w http.ResponseWriter, r *http.Request) {
	orderIDStr := chi.URLParam(r, "id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	payment, err := app.Order.CreatePayment(context.Background(), orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&payment)
}

func (app *application) MidtransWebhookHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	orderIDStr, ok := data["order_id"].(string)
	if !ok {
		http.Error(w, "Invalid order_id format", http.StatusBadRequest)
		return
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		http.Error(w, "Invalid order_id UUID", http.StatusBadRequest)
		return
	}

	status, ok := data["transaction_status"].(string)
	if !ok {
		http.Error(w, "Invalid transaction_status format", http.StatusBadRequest)
		return
	}

	err = app.Order.UpdatePaymentStatus(context.Background(), orderID, status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
