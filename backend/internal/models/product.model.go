package models

import (
	"gorm.io/datatypes"
	"time"
)

type Product struct {
	ID          int            `gorm:"unique;primaryKey;autoIncrement" `
	Name        string         `gorm:"not null"`
	Price       float64        `gorm:"not null"`
	Description string         `gorm:"type:text"`
	Images      datatypes.JSON `gorm:"type:jsonb"`
	Quantity    int            `gorm:"not null"`
	Category    string         `gorm:"size:100"`
	Rating		float32
	RatingCount	int
	SellerID    int            `gorm:"not null"`
	Seller      User           `gorm:"foreignKey:SellerID;references:ID;onDelete:CASCADE"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
