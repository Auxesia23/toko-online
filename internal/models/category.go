package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name string `json:"name" gorm:"type:varchar(100);not null;unique"`
}

type CategoryResponse struct {
	ID   *uint   `json:"id"`
	Name *string `json:"name"`
}
