package models

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	ProductID uuid.UUID `json:"product_id" gorm:"not null"`
	Quantity  int16     `json:"quantity" gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	User    User    `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE;"`
	Product Product `json:"product" gorm:"foreignKey:ProductID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE;"`
}

type CartInput struct {
	ProductID uuid.UUID `json:"product_id"`
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
