package models

import (
	"gorm.io/datatypes"
	"time"
)

type Product struct {
	ID          int            `gorm:"unique;primaryKey;autoIncrement" `
	Name        string         `gorm:"size:255;not null"`
	Price       float64        `gorm:"not null"`
	Description string         `gorm:"type:text"`
	Images      datatypes.JSON `gorm:"type:jsonb"`
	Quantity    int            `gorm:"not null"`
	Category    string         `gorm:"size:100"`
	SellerID    int            `gorm:"not null"`
	Seller      User           `gorm:"foreignKey:SellerID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
