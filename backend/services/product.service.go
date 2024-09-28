package services

import (
	"encoding/json"
	"log"
	"time"

	"github.com/my_ecommerce/internal/dto"
	"github.com/my_ecommerce/internal/models"
	"gorm.io/datatypes"
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

	var imagesArray datatypes.JSON
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


	var imagesArray datatypes.JSON
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

func (p *ProductService) DeleteProduct(id int) error {
	if err := p.db.Where("id = ?",id).Delete(&models.Product{}).Error; err != nil {
		return err
	}
	return nil
}

type UpdateProductRequest struct {
	Name        *string   
	Price       *float64   
	Description *string   
	Images       *datatypes.JSON 
	Quantity    *int      
	Category    *string 
}
func (p *ProductService) UpdateProduct(id int, request dto.UpdatedProduct) (*dto.ProductNoSellerResponse, error) {

	var product models.Product
	if err := p.db.First(&product, id).Error; err != nil {
		return nil, err
	}

	if request.Name != nil {
		product.Name = *request.Name
	}
	if request.Description != nil {
		product.Description = *request.Description
	}
	if request.Images != nil {
		product.Images = *request.Images
	}
	if request.Quantity != nil {
		product.Quantity = *request.Quantity
	}
	if request.Category != nil {
		product.Category = *request.Category
	}

	if request.Price != nil {
		product.Price = *request.Price
	}

	err := p.db.Save(&product).Error

	if err != nil {
		return nil, err
	}

	productReponse := &dto.ProductNoSellerResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Category:    product.Category,
		SellerId:    product.SellerID,
		Quantity:    product.Quantity,
		Price:       product.Price,
		Images:      product.Images,
		CreatedAt:   product.CreatedAt,
	}

	return productReponse, nil
}


