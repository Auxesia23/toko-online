package repository

import (
	"context"
	"gorm.io/gorm"

	"github.com/Auxesia23/toko-online/internal/models"
)

type ProductRepository interface{
	Create(ctx context.Context, product models.Product)(models.Product, error)
}

type ProductRepo struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository{
	return &ProductRepo{
		DB : db,
	}
}

func (repo *ProductRepo) Create(ctx context.Context, product models.Product) (models.Product, error) {
	err := repo.DB.WithContext(ctx).Create(&product).Error
	if err != nil {
		return models.Product{}, err
	}
	return product, nil
}