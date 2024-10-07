package models

import "time"

type Rating struct {
	ID        int `gorm:"primaryKey;autoIncrement"`
	UserID    int
	ProductID int
	Value     float32
	User      User    `gorm:"foreignKey:UserID"`
	Product   Product `gorm:"foreignKey:ProductID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
