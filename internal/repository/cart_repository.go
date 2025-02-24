package repository

import (
	"context"

	"github.com/Auxesia23/toko-online/internal/models"
	"gorm.io/gorm"
)

type CartRepository interface{
	Create(ctx context.Context, input models.CartInput, userID uint) (models.CartResponse, error)
	GetList(ctx context.Context, userID uint)([]models.CartResponse, error)
}

type CartRepo struct{
	DB *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository{
	return &CartRepo{
		DB: db,
	}
}

func (repo *CartRepo) Create(ctx context.Context, input models.CartInput, userID uint)(models.CartResponse, error){
	cartInput := models.Cart{
		UserID: userID,
		ProductID: input.ProductID,
		Quantity: input.Quantity,
	}

	err := repo.DB.WithContext(ctx).Create(&cartInput).Error
	if err != nil {
		return models.CartResponse{}, err
	}

	err = repo.DB.WithContext(ctx).Preload("Product").First(&cartInput,&cartInput.ID).Error
	if err != nil {
		return models.CartResponse{}, err
	}

	response := models.CartResponse{
		ProductName: &cartInput.Product.Name,
		ProductPrice: &cartInput.Product.Price,
		Quantity: &cartInput.Quantity,
	}

	return response, nil
}

func (repo *CartRepo) GetList(ctx context.Context, userID uint)([]models.CartResponse, error){
	var carts []models.Cart
	err := repo.DB.WithContext(ctx).Preload("Product").Where("user_id = ?", userID).Find(&carts).Error
	if err != nil {
		return []models.CartResponse{}, err
	}

	var response []models.CartResponse
	for _, cart := range carts{
		response = append(response, models.CartResponse{
			ProductName: &cart.Product.Name,
			ProductPrice: &cart.Product.Price,
			Quantity: &cart.Quantity,
		})
	}
	return response, nil

}