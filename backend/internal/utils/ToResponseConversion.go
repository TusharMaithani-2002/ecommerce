package utils

import (
	"github.com/my_ecommerce/internal/dto"
	"github.com/my_ecommerce/internal/models"
)

func ConvertToProductResponseDTO(product models.Product) dto.ProductNoSellerResponse {
	return dto.ProductNoSellerResponse{
		ID:          product.ID,
		Name:        product.Name,
		Category:    product.Category,
		Description: product.Description,
		Images:      product.Images,
		Quantity:    product.Quantity,
		SellerId:    product.SellerID,
		Price:       product.Price,
		Rating:      product.Rating,
		RatingCount: product.RatingCount,
	}
}

func ConvertProductsToDTOs(products []models.Product) []dto.ProductNoSellerResponse {
	var productDTOs []dto.ProductNoSellerResponse

	for _, product := range products {
		productDTO := ConvertToProductResponseDTO(product)
		productDTOs = append(productDTOs, productDTO)
	}

	return productDTOs
}
