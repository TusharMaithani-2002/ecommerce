package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/my_ecommerce/internal/database"
)

func main() {

	r := gin.Default()

	// initializing db
	db := internal.InitDB()

	if db == nil {
		log.Fatal("Error connecting database")
	}

	r.GET("/",func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "successfully connected to db",
		})
	})

	r.Run(":8000")
}