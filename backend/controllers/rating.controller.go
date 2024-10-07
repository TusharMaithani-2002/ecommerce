package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/my_ecommerce/internal/middleware"
	"github.com/my_ecommerce/services"
)

type RatingController struct {
	ratingServices services.RatingService
}

func (r *RatingController) InitRatingController(router *gin.Engine, ratingServices services.RatingService) {

	ratingRouter := router.Group("rating")
	ratingRouter.POST("",middleware.VerifyUser(),r.addRating())
	r.ratingServices = ratingServices
}

func (r* RatingController) addRating() gin.HandlerFunc {

	type RatingRequest struct {
		ProductID int `json:"productId" form:"productId" binding:"required"`
		UserID int `json:"userId" form:"userId" binding:"required"`
		Value float32 `json:"value" form:"value" binding:"required"`
	}
	return func(c *gin.Context) {
		var requestBody RatingRequest
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":err.Error(),
			})
			return
		}

		cookieId, exists := c.Get("cookieId")
		userId := cookieId.(int)
		if !exists ||  userId != requestBody.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":"cookie invalid",
			})
			return
		}

		ratingResponse, err := r.ratingServices.AddRating(requestBody.UserID, requestBody.ProductID,requestBody.Value)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"data":ratingResponse,
		})

	}
}