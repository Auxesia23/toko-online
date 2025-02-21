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
	GetById(ctx context.Context, id uuid.UUID) (models.Product, error)
	Update(ctx context.Context, product models.Product) (models.Product, error)
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

	err = repo.DB.WithContext(ctx).Preload("Category").First(&product, &product.ID).Error
	if err != nil {
		return models.ProductResponse{}, err
	}

	categoryResponse := models.CategoryResponse{
		ID:   product.Category.ID,
		Name: product.Category.Name,
	}

	response := models.ProductResponse{
		ID:          &product.ID,
		Name:        &product.Name,
		Description: &product.Description,
		Price:       &product.Price,
		Stock:       &product.Stock,
		ImageUrl:    &product.ImageUrl,
		Category:    &categoryResponse,
	}
	return response, nil
}

func (repo *ProductRepo) GetList(ctx context.Context) ([]models.ProductResponse, error) {
	var products []models.Product
	err := repo.DB.WithContext(ctx).Preload("Category").Find(&products).Error
	if err != nil {
		return []models.ProductResponse{}, err
	}

	var response []models.ProductResponse

	for _, product := range products {
		categoryResponse := models.CategoryResponse{
			ID:   product.Category.ID,
			Name: product.Category.Name,
		}

		response = append(response, models.ProductResponse{
			ID:          &product.ID,
			Name:        &product.Name,
			Description: &product.Description,
			Price:       &product.Price,
			Stock:       &product.Stock,
			ImageUrl:    &product.ImageUrl,
			Category:    &categoryResponse,
		})
	}

	return response, nil
}

func (repo *ProductRepo) GetById(ctx context.Context, id uuid.UUID) (models.Product, error) {
	var product models.Product
	err := repo.DB.WithContext(ctx).Where("id = ?", id).First(&product).Error
	if err != nil {
		return models.Product{}, err
	}
	var category models.Category
	err = repo.DB.WithContext(ctx).Where("id = ?",product.CategoryID).First(&category).Error
	if err == nil {
		product.Category = category
	}

	return product, nil
}

func (repo *ProductRepo) Update(ctx context.Context, product models.Product) (models.Product, error) {
	var oldProduct models.Product
	err := repo.DB.WithContext(ctx).Where("id = ?", product.ID).First(&oldProduct).Error
	if err != nil {
		return models.Product{}, err
	}
	err = repo.DB.WithContext(ctx).Model(&oldProduct).Updates(product).Error
	if err != nil {
		return models.Product{}, err
	}

	var category models.Category
	err = repo.DB.WithContext(ctx).Where("id = ?",product.CategoryID).First(&category).Error
	if err == nil {
		oldProduct.Category = category
	}

	return oldProduct, nil
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
