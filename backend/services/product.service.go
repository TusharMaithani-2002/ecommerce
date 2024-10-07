package services

import (
	"encoding/json"
	"log"
	"time"

	"github.com/my_ecommerce/internal/dto"
	"github.com/my_ecommerce/internal/models"
	"github.com/my_ecommerce/internal/pagination"
	"github.com/my_ecommerce/internal/utils"
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

func (p *ProductService) GetAllProducts(pageNumber int) (*dto.PaginatedProductsResponse, error) {

	const pageSize = 15 // max products shown on a single page
	pagination := pagination.NewPaginate(pageSize, pageNumber).PaginatedResult

	var products []models.Product
	var totalProducts int64

	err := p.db.Model(&models.Product{}).Count(&totalProducts).Error
	if err != nil {
		return nil, err
	}

	nextPageAvailable := totalProducts-int64(pageSize*pageNumber) >= 1

	err = p.db.Scopes(pagination).Find(&products).Error
	if err != nil {
		return nil, err
	}

	productDTOS := utils.ConvertProductsToDTOs(products)

	paginatedResponse := &dto.PaginatedProductsResponse{
		Products:    productDTOS,
		CurrentPage: pageNumber,
		PageSize:    pageSize,
		NextPage:    nextPageAvailable,
	}

	return paginatedResponse, nil
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
		Rating:      product.Rating,
		RatingCount: product.RatingCount,
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
	if err := p.db.Where("id = ?", id).Delete(&models.Product{}).Error; err != nil {
		return err
	}
	return nil
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

	product.UpdatedAt = time.Now()

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
		UpdatedAt:   product.UpdatedAt,
	}

	return productReponse, nil
}

func (p *ProductService) GetFilteredProducts(category, name, description, sorting string, minPrice, maxPrice *float64, pageNumber int) (*dto.PaginatedProductsResponse, error) {

	const pageSize = 15 // max products shown on a single page
	pagination := pagination.NewPaginate(pageSize, pageNumber).PaginatedResult

	var products []models.Product

	query := p.db.Model(&models.Product{})

	if category != "" {
		query = query.Where("category ILIKE ?", "%"+category+"%")
	}
	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if description != "" {
		query = query.Where("to_tsvector(description) @@ plainto_tsquery(?)", description)
	}
	if minPrice != nil && maxPrice != nil {
		query = query.Where("price BETWEEN ? AND ?", minPrice, maxPrice)
	} else if minPrice != nil {
		query = query.Where("price >= ?", *minPrice)
	} else if maxPrice != nil {
		query = query.Where("price <= ?", *maxPrice)
	}

	if sorting == "asc" {
		query = query.Order("price asc")
	} else if sorting == "desc" {
		query = query.Order("price desc")
	}

	var totalProducts int64
	if err := query.Count(&totalProducts).Error; err != nil {
		return nil, err
	}

	if err := query.Scopes(pagination).Find(&products).Error; err != nil {
		return nil, err
	}

	nextPageAvailable := totalProducts-int64(pageSize*pageNumber) >= 1
	productDTOS := utils.ConvertProductsToDTOs(products)

	paginatedResponse := &dto.PaginatedProductsResponse{
		Products:    productDTOS,
		CurrentPage: pageNumber,
		PageSize:    pageSize,
		NextPage:    nextPageAvailable,
	}

	return paginatedResponse, nil

}
