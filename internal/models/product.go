package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string    `json:"name" gorm:"type:varchar(100);not null"`
	CategoryID  uint     `json:"category_id" gorm:"not null"`
	Description string    `json:"description" gorm:"type:varchar(255);not null"`
	Price       int32       `json:"price" gorm:"not null"` 
	Stock       int16     `json:"stock" gorm:"not null"`
	ImageUrl    string    `json:"image_url" gorm:"type:varchar(100);not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Category Category `json:"category" gorm:"foreignKey:CategoryID;constraint:OnUpdate:RESTRICT,OnDelete:SET NULL;"`
}

func (product *Product) BeforeCreate(tx *gorm.DB) (err error) {
	product.ID = uuid.New()
	return
}

func (p *Product) AfterDelete(tx *gorm.DB) error {
	return tx.Where("product_id = ?", p.ID).Delete(&Cart{}).Error
}

type ProductResponse struct {
	ID          *uuid.UUID        `json:"id" `
	Name        *string           `json:"name"`
	Description *string           `json:"description"`
	Price       *int32            `json:"price" `
	Stock       *int16            `json:"stock" `
	ImageUrl    *string           `json:"image_url" `
	Category    *CategoryResponse `json:"category"`
}
