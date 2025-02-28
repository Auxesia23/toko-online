package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID     uint      `json:"user_id" gorm:"not null"`
	TotalPrice int32     `json:"total_price" gorm:"not null"` // Hapus `type:int(12)`
	Status     string    `json:"status" gorm:"type:varchar(20);not null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	User       User        `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE;"`
	OrderItems []OrderItem `json:"order_items" gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Payment    Payment     `json:"payment" gorm:"foreignKey:OrderID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE;"`
}

type OrderItem struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID   uuid.UUID `json:"order_id" gorm:"not null"`
	ProductID uuid.UUID `json:"product_id" gorm:"not null"`
	Quantity  int16     `json:"quantity" gorm:"not null"` 
	Price     int32     `json:"price" gorm:"not null"`    

	Product Product `json:"product" gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

type OrderInput struct {
	Carts []uuid.UUID `json:"carts"`
}

type OrderItemResponse struct {
	Quantity        *int16  `json:"quantity"`
	ProductName     *string `json:"product_name"`
	ProductPrice    *int32  `json:"product_price"`
	ProductImageUrl *string `json:"image_url"`
}

type OrderResponse struct {
	ID         *uuid.UUID `json:"id" `
	TotalPrice *int32     `json:"total_price"`
	Status     *string    `json:"status"`

	CreatedAt  *time.Time
	OrderItems *[]OrderItemResponse `json:"items"`
	Payment    *Payment             `json:"payment"`
}

type OrderPreview struct {
	TotalPrice *int32               `json:"total_price"`
	OrderItems *[]OrderItemResponse `json:"items"`
}

type OrderID struct {
	ID *uuid.UUID `json:"id"`
}
