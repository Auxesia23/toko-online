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
	DeletedAt gorm.DeletedAt `gorm:"index"`

	User    User    `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE;"`
	Product Product `json:"product" gorm:"foreignKey:ProductID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE;"`
}

func (cart *Cart) BeforeCreate(tx *gorm.DB) (err error) {
	cart.ID = uuid.New()
	return
}

type CartInput struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int16     `json:"quantity"`
}

type CartResponse struct {
	ProductName  *string `json:"product"`
	ProductPrice *int32  `json:"product_price"`
	Quantity     *int16  `json:"quantity"`
}
