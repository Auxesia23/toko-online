package repository

import (
	"context"
	"gorm.io/gorm"
	"github.com/Auxesia23/toko-online/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user models.User) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
	Update(ctx context.Context, user models.User, email string) (models.User, error)
}

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepo{
		DB: db,
	}
}

func (repo *UserRepo) Create(ctx context.Context, user models.User) (models.User, error) {
	if err := repo.DB.WithContext(ctx).Create(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (repo *UserRepo) GetByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	if err := repo.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (repo *UserRepo) Update(ctx context.Context, user models.User, email string) (models.User, error) {
	var existingUser models.User


	if err := repo.DB.WithContext(ctx).Where("email = ?", email).First(&existingUser).Error; err != nil {
		return models.User{}, err
	}

	if err := repo.DB.WithContext(ctx).Model(&existingUser).Updates(user).Error; err != nil {
		return models.User{}, err
	}

	return existingUser, nil
}

