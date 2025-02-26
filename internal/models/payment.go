package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Payment struct {
	ID      uuid.UUID `json:"id"`
	OrderID uuid.UUID `json:"order_id" gorm:"not null"`
	Status  string    `json:"status" gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (payment *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	payment.ID = uuid.New()
	return
}
