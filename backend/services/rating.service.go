package services

import (
	"errors"
	"log"
	"time"

	"github.com/my_ecommerce/internal/dto"
	"github.com/my_ecommerce/internal/models"
	"gorm.io/gorm"
)

type RatingService struct {
	db *gorm.DB
}

func (r *RatingService) InitRatingService(database *gorm.DB) {
	r.db = database
	err := r.db.AutoMigrate(&models.Rating{})
	if err != nil {
		log.Fatal("error while migrating rating model")
	}
}

func (r *RatingService) AddRating(userId, productId int, value float32) (*dto.RatingResponse, error) {

	// starting a transaction for atomic operations

	tx := r.db.Begin()

	var product models.Product
	err := tx.Where("id = ?", productId).First(&product).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// checking wether user has already rated, if yes then update rating
	var productRating models.Rating
	if err = tx.Where("user_id = ? and product_id = ?", userId, productId).
		First(&productRating).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {

		tx.Rollback()
		return nil, err

	}

	if productRating.ID != 0 {
		// rating already exists
		currentRating := product.Rating * float32(product.RatingCount)
		product.Rating = (currentRating - productRating.Value + value) / float32(product.RatingCount)

		productRating.Value = value
		if err = tx.Save(&productRating).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

	} else {
		rating := models.Rating{
			UserID:    userId,
			ProductID: productId,
			Value:     value,
			CreatedAt: time.Now(),
		}
		// adding entry in rating table
		if err = tx.Create(&rating).Error; err != nil {
			return nil, err
		}

		// update rating in product
		product.Rating = (product.Rating*float32(product.RatingCount) + value) / (float32(product.RatingCount + 1))
		product.RatingCount += 1
	}

	if err = tx.Save(&product).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	var ratingResponse = &dto.RatingResponse{
		UserID:    userId,
		ProductID: productId,
		Value:     value,
	}

	return ratingResponse, nil
}
