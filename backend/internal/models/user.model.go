package models

type User struct {
	ID       int    `gorm:"unique;primaryKey;autoIncrement" `
	Name     string 
	Email    string `gorm:"unique"`
	Password string 
	Address  string 
}
