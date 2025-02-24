package repository

import (
	"context"
	"errors"

	"github.com/Auxesia23/toko-online/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartRepository interface {
	Create(ctx context.Context, input models.CartInput, userID uint) (error)
	GetList(ctx context.Context, userID uint) ([]models.CartResponse, error)
	Delete(ctx context.Context, cartID uuid.UUID, userID uint)(error)
	Increment(ctx context.Context, cartID uuid.UUID, userID uint)(error)
	Decrement(ctx context.Context, cartID uuid.UUID, userID uint)(error)
}

type CartRepo struct {
	DB *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &CartRepo{
		DB: db,
	}
}

func (repo *CartRepo) Create(ctx context.Context, input models.CartInput, userID uint) (error) {
	var existingCart models.Cart

	err := repo.DB.WithContext(ctx).Where("product_id = ? AND user_id = ?", input.ProductID, userID).First(&existingCart).Error
	if err == nil {
		existingCart.Quantity += 1
		err = repo.DB.WithContext(ctx).Save(&existingCart).Error
		if err != nil {
			return err
		}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		newCart := models.Cart{
			UserID:    userID,
			ProductID: input.ProductID,
			Quantity:  1,
		}
		err = repo.DB.WithContext(ctx).Create(&newCart).Error
		if err != nil {
			return err
		}
		existingCart = newCart
	} else {
		return err
	}

	return nil
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

func (repo *CartRepo) Delete(ctx context.Context, cartID uuid.UUID, userID uint) error {
	var cart models.Cart 
	err := repo.DB.WithContext(ctx).Preload("User").First(&cart,&cartID).Error
	if err != nil {
		return err
	}

	if cart.User.ID != userID {
		return errors.New("forbiden")
	}

	err = repo.DB.WithContext(ctx).Delete(&cart).Error
	if err != nil {
		return err
	}

	return nil

}

func (repo *CartRepo) Increment(ctx context.Context, cartID uuid.UUID, userID uint) (error){
	var cart models.Cart
	err := repo.DB.WithContext(ctx).Preload("User").First(&cart,&cartID).Error
	if err != nil {
		return err
	}


	if cart.User.ID != userID {
		return errors.New("forbiden")
	}

	cart.Quantity += 1
	err = repo.DB.WithContext(ctx).Save(&cart).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *CartRepo) Decrement(ctx context.Context, cartID uuid.UUID, userID uint) (error){
	var cart models.Cart
	err := repo.DB.WithContext(ctx).Preload("User").First(&cart,&cartID).Error
	if err != nil {
		return err
	}


	if cart.User.ID != userID {
		return errors.New("forbiden")
	}

	if cart.Quantity == 1 {
		err = repo.DB.WithContext(ctx).Delete(&cart).Error
		if err != nil {
			return err
		}
	} else {
		cart.Quantity -= 1
		err = repo.DB.WithContext(ctx).Save(&cart).Error
		if err != nil {
			return err
		}
	}
	return nil
}