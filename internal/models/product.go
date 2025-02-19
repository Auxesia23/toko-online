package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID `json:"id" gorm:"type:text;primary_key"`
	Name        string    `json:"name" gorm:"type:varchar(100);not null"`
	Description string    `json:"description" gorm:"type:varchar(255);not null"`
	Price       int32     `json:"price" gorm:"type:int(12);not null"`
	Stock       int16     `json:"stock" gorm:"type:int(8);not null"`
	ImageUrl    string    `json:"image_url" gorm:"type:varchar(100);not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (product *Product) BeforeCreate(tx *gorm.DB) (err error) {
	product.ID = uuid.New()
	return
}

type ProductResponse struct {
	ID          *uuid.UUID `json:"id" gorm:"type:text;primary_key"`
	Name        *string    `json:"name" gorm:"type:varchar(100);not null"`
	Description *string    `json:"description" gorm:"type:varchar(255);not null"`
	Price       *int32     `json:"price" gorm:"type:int(12);not null"`
	Stock       *int16     `json:"stock" gorm:"type:int(8);not null"`
	ImageUrl    *string    `json:"image_url" gorm:"type:varchar(100);not null"`
}
