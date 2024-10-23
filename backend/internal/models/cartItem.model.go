package models

import "time"

type CartItem struct {
	ID        int     `gorm:"primaryKey;autoIncrement"`
	CartID    int     `gorm:"not null"`
	Cart      Cart    `gorm:"ForeignKey:CartID"`
	ProductID int     `gorm:"not null"`
	Product   Product `gorm:"foreignKey:ProductID"`
	Quantity  int     `gorm:"not null;default:1"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
