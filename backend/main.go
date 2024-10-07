package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/my_ecommerce/controllers"
	"github.com/my_ecommerce/internal/database"
	"github.com/my_ecommerce/services"
)

func main() {

	r := gin.Default()

	// initializing db
	db := database.InitDB()

	if db == nil {
		log.Fatal("Error connecting database")
	}

	// defining controllers and services
	userServices := services.UserService{}
	userController := controllers.UserController{}

	userServices.InitUserService(db)
	userController.InitUserController(r,userServices)

	productServices := services.ProductService{}
	productController := controllers.ProductController{}

	productServices.InitProductService(db)
	productController.InitProductController(r,productServices)

	ratingServices := services.RatingService{}
	ratingController := controllers.RatingController{}

	ratingServices.InitRatingService(db)
	ratingController.InitRatingController(r, ratingServices)
	
	r.Run(":8000")
}