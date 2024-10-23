package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/my_ecommerce/internal/middleware"
	"github.com/my_ecommerce/services"
)

type CartController struct {
	cartServices *services.CartServices
}

func (c *CartController) InitCartController(router *gin.Engine, cartServices *services.CartServices) {

	cartRouter := router.Group("/cart")
	cartRouter.GET("/", middleware.VerifyUser(), c.getCart())
	cartRouter.POST("/", middleware.VerifyUser(), c.addItem())
	cartRouter.PATCH("/decrement", middleware.VerifyUser(), c.decrementItem())
	cartRouter.DELETE("/remove", middleware.VerifyUser(), c.removeItem())
	c.cartServices = cartServices
}

func (cc *CartController) getCart() gin.HandlerFunc {
	type CartRequest struct {
		UserID int `json:"userId" form:"userId" binding:"required"`
	}

	return func(c *gin.Context) {
		var cartRequest CartRequest
		if err := c.BindJSON(&cartRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		cookieId, exists := c.Get("cookieId")
		numId := cookieId.(int)

		if !exists || numId != cartRequest.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "cookie invalid",
			})
			return
		}

		cartResponse, err := cc.cartServices.GetCart(cartRequest.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": cartResponse,
		})

	}
}

func (cc *CartController) addItem() gin.HandlerFunc {

	type AddItemRequest struct {
		CartID    int `json:"cartId" form:"cartId" binding:"required"`
		ProductID int `json:"productId" form:"productId" binding:"required"`
		UserID    int `json:"userId" form:"userId" binding:"required"`
	}

	return func(c *gin.Context) {

		var cartItemRequest AddItemRequest
		if err := c.BindJSON(&cartItemRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		cookieId, exists := c.Get("cookieId")
		numId := cookieId.(int)

		if !exists || numId != cartItemRequest.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "cookie invalid",
			})
			return
		}

		cartItemAddReponse, err := cc.cartServices.AddItem(cartItemRequest.CartID, cartItemRequest.UserID, cartItemRequest.ProductID, 1)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": cartItemAddReponse,
		})

	}
}

func (cc *CartController) decrementItem() gin.HandlerFunc {
	type DecrementItemRequest struct {
		CartID    int `json:"cartId" form:"cartId" binding:"required"`
		ProductID int `json:"productId" form:"productId" binding:"required"`
		UserID    int `json:"userId" form:"userId" binding:"required"`
	}
	return func(c *gin.Context) {
		var cartItemRequest DecrementItemRequest
		if err := c.BindJSON(&cartItemRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		cookieId, exists := c.Get("cookieId")
		numId := cookieId.(int)

		if !exists || numId != cartItemRequest.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "cookie invalid",
			})
			return
		}

		cartResponse, err := cc.cartServices.DecrementItem(cartItemRequest.CartID, cartItemRequest.ProductID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": cartResponse,
		})

	}
}

func (cc *CartController) removeItem() gin.HandlerFunc {
	type RemoveItemRequest struct {
		CartID    int `json:"cartId" form:"cartId" binding:"required"`
		ProductID int `json:"productId" form:"productId" binding:"required"`
		UserID    int `json:"userId" form:"userId" binding:"required"`
	}
	return func(c *gin.Context) {
		var cartItemRequest RemoveItemRequest
		if err := c.BindJSON(&cartItemRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		cookieId, exists := c.Get("cookieId")
		numId := cookieId.(int)

		if !exists || numId != cartItemRequest.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "cookie invalid",
			})
			return
		}

		err := cc.cartServices.RemoveItem(cartItemRequest.CartID, cartItemRequest.ProductID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "item deleted from cart",
		})
	}
}
