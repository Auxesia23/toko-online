package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID         uuid.UUID `json:"id" gorm:"type:text;primary_key"`
	UserID     uint      `json:"user_id" gorm:"not null"`
	TotalPrice int32     `json:"total_price" gorm:"type:int(12);not null"`
	Status     string    `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`

	CreatedAt time.Time
	UpdatedAt time.Time

	User       User        `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE;"`
	OrderItems []OrderItem `json:"order_items" gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Payment    Payment     `json:"payment" gorm:"foreignKey:OrderID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE;"`
}

func (order *Order) BeforeCreate(tx *gorm.DB) (err error) {
	order.ID = uuid.New()
	return
}

type OrderItem struct {
	ID        uuid.UUID `json:"id" gorm:"type:text;primary_key"`
	OrderID   uuid.UUID `json:"order_id" gorm:"not null"`
	ProductID uuid.UUID `json:"product_id" gorm:"not null"`
	Quantity  int16     `json:"quantity" gorm:"type:int(8);not null"`
	Price     int32     `json:"price" gorm:"type:int(12);not null"`

	Product Product `json:"product" gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

func (orderItem *OrderItem) BeforeCreate(tx *gorm.DB) (err error) {
	orderItem.ID = uuid.New()
	return
}

type OrderInput struct {
	Carts []struct {
		CartID uuid.UUID `json:"cart_id"`
	} `json:"carts"`
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
}
