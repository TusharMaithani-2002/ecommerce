package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/my_ecommerce/services"
)

type ProductController struct {
	productServices services.ProductService
}

func (p *ProductController) InitProductController(router *gin.Engine, productServices services.ProductService) {

	productRouter := router.Group("/product")

	productRouter.GET("/:id", p.getProduct())
	productRouter.POST("",p.createProduct())
	p.productServices = productServices
}

func (p *ProductController) getProduct() gin.HandlerFunc {

	return func(c *gin.Context) {
		id := c.Param("id")
		numId, err := strconv.Atoi(id)

		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		product, err := p.productServices.GetProduct(numId)
		if err != nil {
			c.JSON(http.StatusBadRequest,gin.H{
				"error":err.Error(),
			})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{
			"data": product,
		})

	}
}

func (p *ProductController) createProduct() gin.HandlerFunc {

	type Product struct {
		Name        string   `json:"name" form:"name" binding:"required"`
		Price       float64   `json:"price" form:"price" binding:"required"`
		Description string   `json:"description" form:"description" binding:"required"`
		Images       []string `json:"images" form:"images"`
		Qunatity    int      `json:"quantity" form:"quantity" binding:"required"`
		Category    string   `json:"category" form:"category"`
		SellerID    int     `json:"sellerId" form:"sellerId" binding:"required"`
	}

	return func(c *gin.Context) {

		var productBody Product
		if err := c.BindJSON(&productBody); err != nil {
			c.JSON(http.StatusBadRequest,gin.H{		
				"error":err.Error(),
			})
			return
		}

		product, err := p.productServices.CreateProduct(productBody.Name,productBody.Description,productBody.Category,int(productBody.SellerID),productBody.Qunatity,productBody.Price,productBody.Images)

		if err != nil {
			c.JSON(http.StatusBadRequest,gin.H{
				"error":err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"data": product,
		})
	}
}
