package repository

import (
	"context"
	"errors"

	"github.com/Auxesia23/toko-online/internal/models"
	"github.com/Auxesia23/toko-online/internal/utils"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user models.User) (models.UserResponse, error)
	GetByID(ctx context.Context, id uint) (models.UserResponse, error)
	Update(ctx context.Context, user models.User, id uint) (models.UserResponse, error)
	Login(ctx context.Context, email string, password string)(string, error)
}

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepo{
		DB: db,
	}
}

func (repo *UserRepo) Create(ctx context.Context, user models.User) (models.UserResponse, error) {
	var existingUser models.User
	err := repo.DB.WithContext(ctx).Where("email = ?", user.Email).First(&existingUser).Error
	if err == nil {
		return models.UserResponse{}, errors.New("email already exist")
	}

	if err := repo.DB.WithContext(ctx).Create(&user).Error; err != nil {
		return models.UserResponse{}, err
	}

	response := models.UserResponse{
		Name:  &user.Name,
		Email: &user.Email,
	}

	return response, nil
}

func (repo *UserRepo) GetByID(ctx context.Context, id uint) (models.UserResponse, error) {
	var user models.User
	if err := repo.DB.WithContext(ctx).First(&user, &id).Error; err != nil {
		return models.UserResponse{}, err
	}

	response := models.UserResponse{
		Name: &user.Name,
		Email: &user.Email,
	}
	return response, nil
}

func (repo *UserRepo) Update(ctx context.Context, user models.User, id uint) (models.UserResponse, error) {
	var existingUser models.User

	if err := repo.DB.WithContext(ctx).First(&existingUser, &id).Error; err != nil {
		return models.UserResponse{}, err
	}

	if err := repo.DB.WithContext(ctx).Model(&existingUser).Updates(&user).Error; err != nil {
		return models.UserResponse{}, err
	}

	response := models.UserResponse{
		Name: &existingUser.Name,
		Email: &existingUser.Email,
	}

	return response, nil
}

func (repo *UserRepo) Login(ctx context.Context, email string, password string)(string, error){
	var user models.User
	err := repo.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.Password){
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(&user)
	if err != nil {
		return "", err
	}

	return token, nil

}