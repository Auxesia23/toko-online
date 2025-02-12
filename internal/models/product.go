package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string `json:"name" gorm:"type:varchar(100);not null"`
	Description string `json:"description" gorm:"type:varchar(255);not null"`
	Price       int32  `json:"price" gorm:"type:int(12);not null"`
	Stock       int16  `json:"stock" gorm:"type:int(8);not null"`
	ImageUrl    string `json:"image_url" gorm:"type:varchar(100);not null"`
}

type ProductInput struct {
	Name string `json:"name"`
	Description string `json:"description"`
	Price       int32  `json:"price"`
	Stock       int16  `json:"stock"`
	ImageUrl    string `json:"image_url"`

}
