package services

import (
	"encoding/json"
	"log"
	"time"

	"github.com/my_ecommerce/internal/dto"
	"github.com/my_ecommerce/internal/models"
	"gorm.io/gorm"
)

type ProductService struct {
	db *gorm.DB
}

func (p *ProductService) InitProductService(database *gorm.DB) {
	p.db = database
	if err := p.db.AutoMigrate(&models.Product{}); err != nil {
		log.Fatal("Failed to auto migrate product")
	}
}

func (p *ProductService) GetProduct(id int) (*dto.ProductReponse, error) {

	var product models.Product
	if err := p.db.Where("id = ?", id).Preload("Seller").First(&product).Error; err != nil {
		return nil, err
	}

	var imagesArray []string
	json.Unmarshal(product.Images, &imagesArray)

	sellerResponse := dto.UserResponse{
		Name:    product.Seller.Name,
		Email:   product.Seller.Email,
		ID:      product.Seller.ID,
		Address: product.Seller.Address,
		Role:    product.Seller.Role,
	}

	productResponse := &dto.ProductReponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Category:    product.Category,
		SellerId:    product.SellerID,
		Quantity:    product.Quantity,
		Price:       product.Price,
		Images:      imagesArray,
		CreatedAt:   product.CreatedAt,
		Seller:      sellerResponse,
	}

	return productResponse, nil
}

func (p *ProductService) CreateProduct(name, description, category string, sellerId, quantity int, price float64, images []string) (*dto.ProductNoSellerResponse, error) {

	imagesJSON, _ := json.Marshal(images)
	product := &models.Product{
		Name:        name,
		Description: description,
		Category:    category,
		SellerID:    sellerId,
		Quantity:    quantity,
		Price:       price,
		Images:      imagesJSON,
		CreatedAt:   time.Now(),
	}

	if err := p.db.Create(&product).Error; err != nil {
		return nil, err
	}


	var imagesArray []string
	json.Unmarshal(product.Images, &imagesArray)

	productReponse := &dto.ProductNoSellerResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Category:    product.Category,
		SellerId:    product.SellerID,
		Quantity:    product.Quantity,
		Price:       product.Price,
		Images:      imagesArray,
		CreatedAt:   product.CreatedAt,
	}

	return productReponse, nil
}
