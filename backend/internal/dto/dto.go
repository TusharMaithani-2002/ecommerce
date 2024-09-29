package dto

import (
	"time"

	"gorm.io/datatypes"
)

// while sending user at
type UserResponse struct {
	ID          int
	Email       string
	Address     string
	Name        string
	Role        string
	PhoneNumber string
}

type ProductReponse struct {
	ID          int
	Name        string
	Price       float64
	Category    string
	Images      datatypes.JSON
	Quantity    int
	SellerId    int
	Description string
	Seller      UserResponse
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ProductNoSellerResponse struct {
	ID          int
	Name        string
	Price       float64
	Category    string
	Images      datatypes.JSON
	Quantity    int
	SellerId    int
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UpdatedProduct struct {
	Name        *string
	Price       *float64
	Description *string
	Images      *datatypes.JSON
	Quantity    *int
	Category    *string
}

type UpdatedUser struct {
	Name        *string
	Address     *string
	Role        *string
	PhoneNumber *string
}
