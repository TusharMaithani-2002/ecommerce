package models

import "time"

type Comment struct {
	ID        int     `gorm:"unique;primaryKey;autoIncrement"`
	Body      string  `gorm:"type:text;not null"`
	ProductID int     `gorm:"not null"`
	Product   Product `gorm:"foreignKey:ProductID;onDelete:cascade"`
	UserID    int     `gorm:"not null"`
	User      User    `gorm:"foreignKey:UserID;onDelete:cascade"`
	Likes     uint    `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
