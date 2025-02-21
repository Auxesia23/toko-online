package repository

import (
	"context"
	"errors"

	"github.com/Auxesia23/toko-online/internal/models"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(context.Context, models.Category) (models.CategoryResponse, error)
	GetList(context.Context) ([]models.CategoryResponse, error)
	GetByID(context.Context, uint) (models.CategoryResponse, error)
	Delete(context.Context, uint) error
}

type CategoryRepo struct {
	DB *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &CategoryRepo{
		DB: db,
	}
}

func (repo *CategoryRepo) Create(ctx context.Context, category models.Category) (models.CategoryResponse, error) {
	var existingCategory models.Category

	err := repo.DB.WithContext(ctx).Where("name = ?", category.Name).First(&existingCategory).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.CategoryResponse{}, err
	}

	if err == nil {
		return models.CategoryResponse{}, errors.New("category already exists")
	}

	err = repo.DB.WithContext(ctx).Create(&category).Error
	if err != nil {
		return models.CategoryResponse{}, err
	}

	response := models.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}

	return response, nil
}

func (repo *CategoryRepo) GetList(ctx context.Context) ([]models.CategoryResponse, error) {
	var categories []models.Category
	err := repo.DB.WithContext(ctx).Find(&categories).Error
	if err != nil {
		return []models.CategoryResponse{}, err
	}

	var response []models.CategoryResponse
	for _, category := range categories {
		response = append(response, models.CategoryResponse{
			ID:   category.ID,
			Name: category.Name,
		})
	}
	return response, nil
}

func (repo *CategoryRepo) Delete(ctx context.Context, id uint) error {
	var category models.Category
	err := repo.DB.WithContext(ctx).Where("id = ?", id).First(&category).Error
	if err != nil {
		return err
	}

	err = repo.DB.WithContext(ctx).Delete(&category).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *CategoryRepo) GetByID(ctx context.Context, id uint) (models.CategoryResponse, error) {
	var category models.Category
	err := repo.DB.WithContext(ctx).Where("id = ?", id).First(&category).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.CategoryResponse{}, err
	} else if err != nil {
		return models.CategoryResponse{}, err
	}

	return models.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}, nil
}

