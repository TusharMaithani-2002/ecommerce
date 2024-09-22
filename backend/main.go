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
	db := internal.InitDB()

	if db == nil {
		log.Fatal("Error connecting database")
	}

	// defining controllers and services
	userServices := services.UserService{}
	userController := controllers.UserController{}

	userServices.InitUserService(db)
	userController.InitUserController(r,userServices)

	
	r.Run(":8000")
}