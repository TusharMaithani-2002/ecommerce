package models

type Cart struct {
	ID     int        `gorm:"primaryKey;autoIncrement"`
	UserID int        `gorm:"not null"`
	User   User       `gorm:"foreignKey:UserID"`
	Items  []CartItem `gorm:"foreignKey:CartID"`
}
