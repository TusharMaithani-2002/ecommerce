package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/my_ecommerce/internal/dto"
	"github.com/my_ecommerce/internal/middleware"
	"github.com/my_ecommerce/services"
	"gorm.io/datatypes"
)

type ProductController struct {
	productServices services.ProductService
}

func (p *ProductController) InitProductController(router *gin.Engine, productServices services.ProductService) {

	productRouter := router.Group("/product")

	productRouter.GET("/:id", p.getProduct())
	productRouter.GET("/all/:page",p.getAllProducts())
	productRouter.POST("",middleware.VerifyUser(),p.createProduct())
	productRouter.DELETE("/delete/:id",middleware.VerifyUser(),p.deleteProduct())
	productRouter.PATCH("/update/:id",middleware.VerifyUser(),p.updateProduct())
	productRouter.GET("/filter",p.getFilteredProducts())
	p.productServices = productServices
}

func (p *ProductController) getAllProducts() gin.HandlerFunc {
	return func(c *gin.Context) {

		page := c.Param("page")
		pageNumber, err := strconv.Atoi(page)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":err,
				"message":"Page number not valid",
			})
			return
		}

		products, err := p.productServices.GetAllProducts(pageNumber)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":products,
		})
	}
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

		// checking whether email in cookie and seller email is same
		cookieId,ok := c.Get("cookieId")
		userId := cookieId.(int)

		if !ok || userId != productBody.SellerID {
			c.JSON(http.StatusUnauthorized,gin.H{
				"error":"cookie is invalid",
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

func (p *ProductController) deleteProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")
		numId, err := strconv.Atoi(id)

		if err != nil {
			c.JSON(http.StatusBadRequest,gin.H{
				"error":"product Id not valid number",
			})
			return
		}

		product, err := p.productServices.GetProduct(numId)

		if err != nil {
			c.JSON(http.StatusBadRequest,gin.H{
				"error":err.Error(),
				"message":"failed to get product with current Id",
			})
			return
		}
		// first we have to check whether the deleting request is made by the seller or not
		role := c.MustGet("role").(string)

		if role != "admin" {
			// you are not the admin then you must be the user to delete the product
			cookieId := c.MustGet("cookieId").(int)
			if product.SellerId != cookieId {
				c.JSON(http.StatusForbidden,gin.H{
					"error":"you are not authorized for the deletion of the product",
				})
				return
			}
		}


		if err := p.productServices.DeleteProduct(numId); err != nil {
			c.JSON(http.StatusBadRequest,gin.H{
				"error":err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK,gin.H{
			"message":fmt.Sprintf("product with id = %s deleted",id),
		})
	}
}

func (p *ProductController) updateProduct() gin.HandlerFunc {

	type UpdateProductRequest struct {
		Name        *string   `json:"name,omitempty" form:"name,omitempty"`
		Price       *float64   `json:"price,omitempty" form:"price,omitempty"`
		Description *string   `json:"description,omitempty" form:"description,omitempty"`
		Images       *datatypes.JSON `json:"images,omitempty" form:"images,omitempty"`
		Quantity    *int      `json:"quantity,omitempty" form:"quantity,omitempty"`
		Category    *string   `json:"category,omitempty" form:"category,omitempty"`
	}



	return func(c *gin.Context) {

		id := c.Param("id")
		numId, err := strconv.Atoi(id)

		if err != nil {
			c.JSON(http.StatusBadRequest,gin.H{
				"error":"product Id not valid number",
			})
			return
		}

		product, err := p.productServices.GetProduct(numId)

		if err != nil {
			c.JSON(http.StatusBadRequest,gin.H{
				"error":err.Error(),
				"message":"failed to get product with current Id",
			})
			return
		}

		cookieId,ok := c.Get("cookieId")
		userId := cookieId.(int)
		
		if !ok {
			c.JSON(http.StatusUnauthorized,gin.H{
				"error":"cookie is invalid",
			})
			return
		}

		if userId != product.SellerId {
			c.JSON(http.StatusUnauthorized,gin.H{
				"error":"you are not authorized to the product",
			})
			return
		}

		var request UpdateProductRequest
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest,gin.H{		
				"error":err.Error(),
			})
			return
		}
		var updateRequest dto.UpdatedProduct
		updateRequest.Name = request.Name
		updateRequest.Category = request.Category
		updateRequest.Description = request.Description
		updateRequest.Images = request.Images
		updateRequest.Price = request.Price
		updateRequest.Quantity = request.Quantity

		updatedProduct, err := p.productServices.UpdateProduct(numId, updateRequest)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":err.Error(),
			})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{
			"data": updatedProduct,
			"messaage":"product updated successfully",
		})

	}
}

func (p *ProductController) getFilteredProducts() gin.HandlerFunc {

	return func(c *gin.Context) {
		category := c.Query("category")
		name := c.Query("name")
		description := c.Query("description")
		minPriceStr := c.Query("minPrice")
		maxPriceStr := c.Query("maxPrice")
		pageNumberStr := c.Query("pageNumber")
		sorting := c.Query("sorting")

		var minPrice* float64 = nil
		var maxPrice* float64 = nil
		var pageNumber int = 1

		// applying checks for price and page number

		if minPriceStr != "" {
			min, err := strconv.ParseFloat(minPriceStr, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":err.Error(),
				})
				return
			}

			minPrice = &min
		}

		if maxPriceStr != "" {
			max, err := strconv.ParseFloat(maxPriceStr, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":err.Error(),
				})
				return
			}
			maxPrice = &max
		}

		if pageNumberStr != "" {
			page, err := strconv.Atoi(pageNumberStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":err.Error(),
				})
				return
			}
			pageNumber = page
		}


		products, err := p.productServices.GetFilteredProducts(category, name, description, sorting, minPrice, maxPrice,pageNumber)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":products,
		})
	}
}