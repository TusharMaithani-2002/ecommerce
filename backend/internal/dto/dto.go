package dto

import (
	"gorm.io/datatypes"
	"time"
)

// while sending user at
type UserResponse struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	Name        string `json:"name"`
	Role        string `json:"role"`
	PhoneNumber string `json:"phoneNumber"`
}

type ProductReponse struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Price       float64        `json:"price"`
	Category    string         `json:"category"`
	Images      datatypes.JSON `json:"images"`
	Quantity    int            `json:"quantity"`
	SellerId    int            `json:"sellerId"`
	Description string         `json:"description"`
	Seller      UserResponse   `json:"seller"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	Rating      float32        `json:"rating"`
	RatingCount int            `json:"ratingCount"`
}

type ProductNoSellerResponse struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Price       float64        `json:"price"`
	Category    string         `json:"category"`
	Images      datatypes.JSON `json:"images"`
	Quantity    int            `json:"quantity"`
	SellerId    int            `json:"sellerId"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	Rating      float32        `json:"rating"`
	RatingCount int            `json:"ratingCount"`
}

type UpdatedProduct struct {
	Name        *string         `json:"name"`
	Price       *float64        `json:"price"`
	Description *string         `json:"description"`
	Images      *datatypes.JSON `json:"images"`
	Quantity    *int            `json:"quantity"`
	Category    *string         `json:"category"`
	Rating      *float32        `json:"rating"`
	RatingCount *int            `json:"ratingCount"`
}

type UpdatedUser struct {
	Name        *string `json:"name"`
	Address     *string `json:"address"`
	Role        *string `json:"role"`
	PhoneNumber *string `json:"phoneNumber"`
}

type PaginatedProductsResponse struct {
	Products    []ProductNoSellerResponse `json:"products"`
	CurrentPage int                       `json:"currentPage"`
	PageSize    int                       `json:"pageSize"`
	NextPage    bool                      `json:"nextPage"`
}

type RatingResponse struct {
	UserID    int     `json:"userId"`
	ProductID int     `json:"productId"`
	Value     float32 `json:"value"`
}
