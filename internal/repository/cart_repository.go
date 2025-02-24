package repository

import (
	"context"
	"errors"

	"github.com/Auxesia23/toko-online/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartRepository interface {
	Create(ctx context.Context, input models.CartInput, userID uint) (models.CartResponse, error)
	GetList(ctx context.Context, userID uint) ([]models.CartResponse, error)
	Update(ctx context.Context, input models.CartUpdate,cartID uuid.UUID, userID uint) (models.CartResponse, error)
	Delete(ctx context.Context, cartID uuid.UUID, userID uint)(error)
}

type CartRepo struct {
	DB *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &CartRepo{
		DB: db,
	}
}

func (repo *CartRepo) Create(ctx context.Context, input models.CartInput, userID uint) (models.CartResponse, error) {
	var existingCart models.Cart

	err := repo.DB.WithContext(ctx).Where("product_id = ? AND user_id = ?", input.ProductID, userID).First(&existingCart).Error
	if err == nil {
		existingCart.Quantity += input.Quantity
		err = repo.DB.WithContext(ctx).Save(&existingCart).Error
		if err != nil {
			return models.CartResponse{}, err
		}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		newCart := models.Cart{
			UserID:    userID,
			ProductID: input.ProductID,
			Quantity:  input.Quantity,
		}
		err = repo.DB.WithContext(ctx).Create(&newCart).Error
		if err != nil {
			return models.CartResponse{}, err
		}
		existingCart = newCart
	} else {
		return models.CartResponse{}, err
	}

	err = repo.DB.WithContext(ctx).Preload("Product").First(&existingCart, existingCart.ID).Error
	if err != nil {
		return models.CartResponse{}, err
	}

	product := models.CartProduct{
		ID:       &existingCart.Product.ID,
		Name:     &existingCart.Product.Name,
		Price:    &existingCart.Product.Price,
		Stock:    &existingCart.Product.Stock,
		ImageUrl: &existingCart.Product.ImageUrl,
	}

	response := models.CartResponse{
		ID:       &existingCart.ID,
		Quantity: &existingCart.Quantity,
		Product:  &product,
	}

	return response, nil
}

func (repo *CartRepo) GetList(ctx context.Context, userID uint) ([]models.CartResponse, error) {
	var carts []models.Cart
	err := repo.DB.WithContext(ctx).Preload("Product").Where("user_id = ?", userID).Find(&carts).Error
	if err != nil {
		return []models.CartResponse{}, err
	}

	var response []models.CartResponse
	for _, cart := range carts {
		product := models.CartProduct{
			ID:       &cart.Product.ID,
			Name:     &cart.Product.Name,
			Price:    &cart.Product.Price,
			Stock:    &cart.Product.Stock,
			ImageUrl: &cart.Product.ImageUrl,
		}

		response = append(response, models.CartResponse{
			ID:       &cart.ID,
			Quantity: &cart.Quantity,
			Product:  &product,
		})
	}
	return response, nil

}

func (repo *CartRepo) Update(ctx context.Context, input models.CartUpdate,cartID uuid.UUID, userID uint) (models.CartResponse, error) {
	var existingCart models.Cart

	err := repo.DB.WithContext(ctx).Preload("Product").First(&existingCart, cartID).Error
	if err != nil {
		return models.CartResponse{}, err
	}

	if existingCart.User.ID != userID {
		return models.CartResponse{}, errors.New("forbiden")
	}

	existingCart.Quantity = input.Quantity
	err = repo.DB.WithContext(ctx).Save(&existingCart).Error
	if err != nil {
		return models.CartResponse{}, err
	}

	product := models.CartProduct{
		ID:       &existingCart.Product.ID,
		Name:     &existingCart.Product.Name,
		Price:    &existingCart.Product.Price,
		Stock:    &existingCart.Product.Stock,
		ImageUrl: &existingCart.Product.ImageUrl,
	}

	response := models.CartResponse{
		ID:       &existingCart.ID,
		Quantity: &existingCart.Quantity,
		Product:  &product,
	}

	return response, nil
}


func (repo *CartRepo) Delete(ctx context.Context, cartID uuid.UUID, userID uint) error {
	var existingCart models.Cart 
	err := repo.DB.WithContext(ctx).Preload("User").First(&existingCart,&cartID).Error
	if err != nil {
		return err
	}

	if existingCart.User.ID != userID {
		return errors.New("forbiden")
	}

	err = repo.DB.WithContext(ctx).Delete(&existingCart).Error
	if err != nil {
		return err
	}

	return nil

}