package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cart struct {
	ID        uuid.UUID `json:"id" gorm:"type:text;primary_key"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	ProductID uuid.UUID `json:"product_id" gorm:"not null"`
	Quantity  int16     `json:"quantity" gorm:"type:int(8);not null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	User    User    `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE;"`
	Product Product `json:"product" gorm:"foreignKey:ProductID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE;"`
}

func (cart *Cart) BeforeCreate(tx *gorm.DB) (err error) {
	cart.ID = uuid.New()
	return
}

func (p *Product) AfterDelete(tx *gorm.DB) error {
	return tx.Where("product_id = ?", p.ID).Delete(&Cart{}).Error
}


type CartInput struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int16     `json:"quantity"`
}

type CartUpdate struct {
	Quantity  int16     `json:"quantity"`
}

type CartResponse struct {
	ID       *uuid.UUID   `json:"id"`
	Quantity *int16       `json:"quantity"`
	Product  *CartProduct `json:"product"`
}

type CartProduct struct {
	ID       *uuid.UUID `json:"id"`
	Name     *string    `json:"name"`
	Price    *int32     `json:"price" `
	Stock    *int16     `json:"stock" `
	ImageUrl *string    `json:"image_url" `
}
