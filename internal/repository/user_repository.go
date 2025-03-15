package repository

import (
	"context"
	"errors"

	"github.com/Auxesia23/toko-online/internal/models"
	"github.com/Auxesia23/toko-online/internal/utils"
	"gorm.io/gorm"
	"github.com/golang-jwt/jwt"
)

type UserRepository interface {
	Create(ctx context.Context, user models.User) (models.UserResponse, error)
	Verify(ctx context.Context, token string)(error)
	GetByID(ctx context.Context, id uint) (models.UserResponse, error)
	Update(ctx context.Context, user models.User, id uint) (models.UserResponse, error)
	Login(ctx context.Context, email string, password string) (string, error)
	GoogleLogin(ctx context.Context, code string) (string, error)
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

	var newUser models.User
	err = repo.DB.WithContext(ctx).First(&newUser, user.ID).Error
	if err != nil {
		return models.UserResponse{}, err
	}
	
	token, _ := utils.GenerateToken(&user)
	_ = utils.SendVerificationEmail(user.Email, token)

	response := models.UserResponse{
		Name:  &user.Name,
		Email: &user.Email,
	}

	return response, nil
}

func (repo *UserRepo) Verify(ctx context.Context, token string) error {
	jwtToken, err := utils.VerifyJWT(token)
	if err != nil {
		return err
	}

	claims := jwtToken.Claims.(jwt.MapClaims)
	userIDFloat := claims["user_id"].(float64)
	userID := uint(userIDFloat)

	var user models.User
	err = repo.DB.WithContext(ctx).First(&user, userID).Error
	if err != nil {
		return err
	}

	err = repo.DB.WithContext(ctx).Model(&user).Update("verified", true).Error
	if err != nil {
		return err
	}

	return nil
}


func (repo *UserRepo) GetByID(ctx context.Context, id uint) (models.UserResponse, error) {
	var user models.User
	if err := repo.DB.WithContext(ctx).First(&user, &id).Error; err != nil {
		return models.UserResponse{}, err
	}

	response := models.UserResponse{
		Name:  &user.Name,
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
		Name:  &existingUser.Name,
		Email: &existingUser.Email,
	}

	return response, nil
}

func (repo *UserRepo) Login(ctx context.Context, email string, password string) (string, error) {
	var user models.User
	err := repo.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !user.Verified {
		return "", errors.New("email not verified")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(&user)
	if err != nil {
		return "", err
	}

	return token, nil

}

func (repo *UserRepo) GoogleLogin(ctx context.Context, code string) (string, error) {
	token, err := utils.ExchangeCodeForToken(code)
	if err != nil {
		return "", err
	}

	userInfo, err := utils.FetchGoogleUserInfo(token.AccessToken)
	if err != nil {
		return "", err
	}

	var user models.User
	result := repo.DB.Where("email = ?", userInfo.Email).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			user = models.User{
				Name:  userInfo.Name,
				Email: userInfo.Email,
				Picture: userInfo.Picture,
				Verified: userInfo.VerifiedEmail,
			}
			if err := repo.DB.Create(&user).Error; err != nil {
				return "", err
			}
		} else {
			return "", result.Error
		}
	}

	jwt, err := utils.GenerateToken(&user)
	if err != nil {
		return "", err
	}
	return jwt, nil
}
