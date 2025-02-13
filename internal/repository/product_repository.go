package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/Auxesia23/toko-online/internal/models"
)

type ProductRepository interface {
	Create(ctx context.Context, product models.Product) (models.ProductResponse, error)
	GetList(ctx context.Context) ([]models.ProductResponse, error)
}

type ProductRepo struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &ProductRepo{
		DB: db,
	}
}

func (repo *ProductRepo) Create(ctx context.Context, product models.Product) (models.ProductResponse, error) {
	err := repo.DB.WithContext(ctx).Create(&product).Error
	if err != nil {
		return models.ProductResponse{}, err
	}

	response := models.ProductResponse{
		ID:          &product.ID,
		Name:        &product.Name,
		Description: &product.Description,
		Price:       &product.Price,
		Stock:       &product.Stock,
		ImageUrl:    &product.ImageUrl,
	}
	return response, nil
}

func (repo *ProductRepo) GetList(ctx context.Context) ([]models.ProductResponse, error) {
	var products []models.Product
	err := repo.DB.WithContext(ctx).Find(&products).Error
	if err != nil {
		return []models.ProductResponse{}, err
	}

	var response []models.ProductResponse

	for _, product := range products {
		response = append(response, models.ProductResponse{
			ID:          &product.ID,
			Name:        &product.Name,
			Description: &product.Description,
			Price:       &product.Price,
			Stock:       &product.Stock,
			ImageUrl:    &product.ImageUrl,
		})
	}

	return response, nil
}
