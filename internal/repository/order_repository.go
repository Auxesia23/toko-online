package repository

import (
	"context"
	"errors"

	"github.com/Auxesia23/toko-online/internal/models"
	"github.com/Auxesia23/toko-online/internal/payment"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(ctx context.Context, userID uint, order models.OrderInput) (uuid.UUID, error)
	GetList(ctx context.Context, userID uint) ([]models.OrderResponse, error)
	Preview(ctx context.Context, userID uint, order models.OrderInput) (models.OrderPreview, error)
	GetByID(ctx context.Context, userID uint, orderID uuid.UUID) (models.OrderResponse, error)
	CreatePayment(ctx context.Context, orderID uuid.UUID) (models.Payment, error)
	UpdatePaymentStatus(ctx context.Context, orderID uuid.UUID, midtransStatus string) error
}

type OrderRepo struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &OrderRepo{
		DB: db,
	}
}

func (repo *OrderRepo) Create(ctx context.Context, userID uint, input models.OrderInput) (uuid.UUID, error) {
	tx := repo.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var carts []models.Cart
	if err := tx.Preload("Product").Where("id IN ?", input.Carts).Find(&carts).Error; err != nil {
		tx.Rollback()
		return uuid.UUID{}, err
	}

	var totalPrice int32 = 0
	for _, cart := range carts {
		if cart.Quantity > cart.Product.Stock {
			tx.Rollback()
			return uuid.UUID{}, errors.New("not enough stock for product: " + cart.Product.Name)
		}
		totalPrice += cart.Product.Price * int32(cart.Quantity)
	}

	order := models.Order{
		UserID:     userID,
		TotalPrice: totalPrice,
		Status:     "Menunggu pembayaran",
	}
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return uuid.UUID{}, err
	}

	var orderItems []models.OrderItem
	for _, cart := range carts {
		orderItems = append(orderItems, models.OrderItem{
			OrderID:   order.ID,
			ProductID: cart.ProductID,
			Quantity:  cart.Quantity,
			Price:     cart.Product.Price,
		})

		if err := tx.Model(&models.Product{}).
			Where("id = ?", cart.ProductID).
			UpdateColumn("stock", gorm.Expr("stock - ?", cart.Quantity)).Error; err != nil {
			tx.Rollback()
			return uuid.UUID{}, err
		}
	}

	if err := tx.Create(&orderItems).Error; err != nil {
		tx.Rollback()
		return uuid.UUID{}, err
	}

	if err := tx.Delete(&carts).Error; err != nil {
		tx.Rollback()
		return uuid.UUID{}, err
	}

	tx.Commit()
	return order.ID, nil
}

func (repo *OrderRepo) GetList(ctx context.Context, userID uint) ([]models.OrderResponse, error) {
	var orders []models.Order
	err := repo.DB.WithContext(ctx).Preload("Payment").Preload("OrderItems.Product").Where("user_id = ?", userID).Find(&orders).Error
	if err != nil {
		return []models.OrderResponse{}, err
	}

	var orderResponses []models.OrderResponse
	for _, order := range orders {
		var itemResponse []models.OrderItemResponse
		for _, item := range order.OrderItems {
			itemResponse = append(itemResponse, models.OrderItemResponse{
				Quantity:        &item.Quantity,
				ProductName:     &item.Product.Name,
				ProductPrice:    &item.Product.Price,
				ProductImageUrl: &item.Product.ImageUrl,
			})
		}
		orderResponses = append(orderResponses, models.OrderResponse{
			ID:         &order.ID,
			TotalPrice: &order.TotalPrice,
			Status:     &order.Status,
			CreatedAt:  &order.CreatedAt,
			OrderItems: &itemResponse,
			Payment:    &order.Payment,
		})
	}

	return orderResponses, nil
}

func (repo *OrderRepo) GetByID(ctx context.Context, userID uint, orderID uuid.UUID) (models.OrderResponse, error) {
	var order models.Order
	err := repo.DB.WithContext(ctx).Preload("Payment").Preload("OrderItems.Product").First(&order, &orderID).Error
	if err != nil {
		return models.OrderResponse{}, err
	}

	if order.UserID != userID {
		return models.OrderResponse{}, errors.New("not found")
	}

	var items []models.OrderItemResponse
	for _, item := range order.OrderItems {
		items = append(items, models.OrderItemResponse{
			Quantity:        &item.Quantity,
			ProductName:     &item.Product.Name,
			ProductPrice:    &item.Product.Price,
			ProductImageUrl: &item.Product.ImageUrl,
		})
	}

	response := models.OrderResponse{
		ID:         &order.ID,
		TotalPrice: &order.TotalPrice,
		Status:     &order.Status,
		CreatedAt:  &order.CreatedAt,
		OrderItems: &items,
		Payment:    &order.Payment,
	}

	return response, nil
}

func (repo *OrderRepo) CreatePayment(ctx context.Context, orderID uuid.UUID) (models.Payment, error) {
	var order models.Order
	err := repo.DB.WithContext(ctx).Preload("OrderItems.Product").Preload("User").First(&order, orderID).Error
	if err != nil {
		return models.Payment{}, err
	}

	token, err := payment.CreateMidtransPayment(&order)
	if err != nil {
		return models.Payment{}, err
	}

	payment := models.Payment{
		ID:            uuid.New(),
		OrderID:       order.ID,
		Status:        "Pending",
		MidtransToken: token,
	}
	err = repo.DB.WithContext(ctx).Create(&payment).Error
	if err != nil {
		return models.Payment{}, err
	}
	return payment, nil
}

func (repo *OrderRepo) Preview(ctx context.Context, userID uint, order models.OrderInput) (models.OrderPreview, error) {
	var orderItems []models.Cart
	err := repo.DB.WithContext(ctx).Preload("Product").Where("id IN ?", order.Carts).Find(&orderItems).Error
	if err != nil {
		return models.OrderPreview{}, err
	}

	var totalPrice int32 = 0
	var itemResponses []models.OrderItemResponse

	for _, cart := range orderItems {
		totalPrice += cart.Product.Price * int32(cart.Quantity)
		itemResponses = append(itemResponses, models.OrderItemResponse{
			Quantity:        &cart.Quantity,
			ProductName:     &cart.Product.Name,
			ProductPrice:    &cart.Product.Price,
			ProductImageUrl: &cart.Product.ImageUrl,
		})
	}

	response := models.OrderPreview{
		TotalPrice: &totalPrice,
		OrderItems: &itemResponses,
	}

	return response, nil
}

func (repo *OrderRepo) UpdatePaymentStatus(ctx context.Context, orderID uuid.UUID, midtransStatus string) error {
	var paymentStatus string
	var orderStatus string

	switch midtransStatus {
	case "capture", "settlement":
		paymentStatus = "berhasil"
		orderStatus = "Di Proses" 
	case "pending":
		paymentStatus = "menunggu"
		orderStatus = "Menunggu Pembayaran"
	case "deny", "cancel", "expire", "failure":
		paymentStatus = "gagal"
		orderStatus = "Gagal"
	default:
		paymentStatus = "menunggu"
		orderStatus = "Menunggu Pembayaran"
	}

	tx := repo.DB.WithContext(ctx).Begin()

	if err := tx.Model(&models.Payment{}).
		Where("order_id = ?", orderID).
		Update("status", paymentStatus).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&models.Order{}).
		Where("id = ?", orderID).
		Update("status", orderStatus).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
