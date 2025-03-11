package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name      string `json:"name" gorm:"type:varchar(100);not null"`
	Email     string `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Password  string `json:"password" gorm:"type:varchar(100);not null"`
	Picture   string `json:"picture" gorm:"type:varchar(255)"`
	Superuser bool   `json:"is_superuser" gorm:"default:false"`
	Verified  bool   `json:"verified" gorm:"default:false"`
}

type UserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdate struct {
	Name string `json:"name"`
}
