package dto

import (
	"time"
)

// while sending user at
type UserResponse struct {
	ID      int
	Email   string
	Address string
	Name    string
	Role    string
}

type ProductReponse struct {
	ID          int
	Name        string
	Price       float64
	Category    string
	Images      []string
	Quantity    int
	SellerId    int
	Description string
	Seller      UserResponse
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
