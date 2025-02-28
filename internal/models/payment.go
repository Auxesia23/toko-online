package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Payment struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID       uuid.UUID `json:"order_id" gorm:"not null"`
	Status        string    `json:"status" gorm:"not null"`
	MidtransToken string    `json:"token"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
