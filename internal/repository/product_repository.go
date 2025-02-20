package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/Auxesia23/toko-online/internal/models"
	"github.com/google/uuid"
)

type ProductRepository interface {
	Create(ctx context.Context, product models.Product) (models.ProductResponse, error)
	GetList(ctx context.Context) ([]models.ProductResponse, error)
	GetById(ctx context.Context, id uuid.UUID) (models.ProductResponse, error)
	Update(ctx context.Context, product models.ProductResponse) (models.ProductResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
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
		Category:    &product.Category,
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
			Category:    &product.Category,
		})
	}

	return response, nil
}

func (repo *ProductRepo) GetById(ctx context.Context, id uuid.UUID) (models.ProductResponse, error) {
	var product models.Product
	err := repo.DB.WithContext(ctx).Where("id = ?", id).First(&product).Error
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
		Category:    &product.Category,
	}
	return response, nil
}

func (repo *ProductRepo) Update(ctx context.Context, product models.ProductResponse) (models.ProductResponse, error) {
	var oldProduct models.Product
	err := repo.DB.WithContext(ctx).Where("id = ?", product.ID).First(&oldProduct).Error
	if err != nil {
		return models.ProductResponse{}, err
	}
	err = repo.DB.WithContext(ctx).Model(&oldProduct).Updates(product).Error
	if err != nil {
		return models.ProductResponse{}, err
	}
	response := models.ProductResponse{
		ID:          &oldProduct.ID,
		Name:        &oldProduct.Name,
		Description: &oldProduct.Description,
		Price:       &oldProduct.Price,
		Stock:       &oldProduct.Stock,
		ImageUrl:    &oldProduct.ImageUrl,
		Category:    &oldProduct.Category,
	}
	return response, nil
}

func (repo *ProductRepo) Delete(ctx context.Context, id uuid.UUID) error {
	var product models.Product
	err := repo.DB.WithContext(ctx).Where("id = ?", id).First(&product).Error
	if err != nil {
		return err
	}
	err = repo.DB.WithContext(ctx).Delete(&product).Error
	if err != nil {
		return err
	}
	return nil
}
