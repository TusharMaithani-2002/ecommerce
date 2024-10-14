package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/my_ecommerce/internal/middleware"
	"github.com/my_ecommerce/services"
)

type CommentController struct {
	commentServices *services.CommentServices
}

func (c *CommentController) InitCommentController(router *gin.Engine, commentServices *services.CommentServices) {

	commentRouter := router.Group("/comment")
	commentRouter.POST("/", middleware.VerifyUser(), c.addComment())
	commentRouter.PATCH("/like", middleware.VerifyUser(), c.addLikeToComment())
	commentRouter.DELETE("/delete",middleware.VerifyUser(), c.deleteComment())
	commentRouter.PATCH("/update", middleware.VerifyUser(), c.updateComment())
	c.commentServices = commentServices
}

func (cc *CommentController) addComment() gin.HandlerFunc {
	type Comment struct {
		UserID    int    `json:"userId" form:"userId" binding:"required"`
		ProductID int    `json:"productId" form:"productId" binding:"required"`
		Body      string `json:"body" form:"body" binding:"required"`
	}

	return func(c *gin.Context) {
		var commentBody Comment
		if err := c.BindJSON(&commentBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		cookieId, exists := c.Get("cookieId")
		numId := cookieId.(int)

		if !exists || numId != commentBody.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "cookie invalid",
			})
			return
		}

		commentResponse, err := cc.commentServices.AddComment(commentBody.Body,commentBody.UserID, commentBody.ProductID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"data":commentResponse,
		})
	}
}

func (cc *CommentController) addLikeToComment() gin.HandlerFunc {

	type CommentLikeRequest struct {
		CommentID int `json:"commentId" form:"commentId" binding:"required"`
		UserID int `json:"userId" form:"userId" binding:"required"`
	}
	return func(c *gin.Context) {
		var commentLikeRequest CommentLikeRequest

		if err := c.BindJSON(&commentLikeRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":err.Error(),
			})
			return
		}

		cookieId, exists := c.Get("cookieId")
		numId := cookieId.(int)

		if !exists || numId != commentLikeRequest.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "cookie invalid",
			})
			return
		}

		commentReponse, err := cc.commentServices.AddLikeToComment(commentLikeRequest.UserID, commentLikeRequest.CommentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":commentReponse,
		})
	}
}

func (cc* CommentController) updateComment() gin.HandlerFunc {

	type UpdateCommentRequest struct {
		UserID int `json:"userId" form:"userId" binding:"required"`
		CommentID int `json:"commentId" form:"commentId" binding:"required"`
		Body string `json:"body" form:"body" binding:"required"`
	}

	return func(c *gin.Context) {

		var updateCommentRequest UpdateCommentRequest
		if err := c.BindJSON(&updateCommentRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":err,
			})
			return
		}

		cookieId, exists := c.Get("cookieId")
		numId := cookieId.(int)

		if !exists || numId != updateCommentRequest.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "cookie invalid",
			})
			return
		}

		commentReponse, err := cc.commentServices.UpdateComment(updateCommentRequest.UserID, updateCommentRequest.CommentID, updateCommentRequest.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H {
			"data":commentReponse,
		})
	}
}

func (cc *CommentController) deleteComment() gin.HandlerFunc {
	
	type DeleteCommentRequest struct {
		UserID int `json:"userId" form:"userId" binding:"required"`
		CommentID int `json:"commentId" form:"commentId" binding:"required"`
	}
	return func(c *gin.Context) {

		var deleteCommentRequest DeleteCommentRequest
		if err := c.BindJSON(&deleteCommentRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":err.Error(),
			})
			return
		}

		cookieId, existsId := c.Get("cookieId")
		numId := cookieId.(int)
		cookieRole, existsRole := c.Get("cookieRole")
		role := cookieRole.(string)
		if !existsId || !existsRole || numId != deleteCommentRequest.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "cookie invalid",
			})
			return
		}

		if err := cc.commentServices.DeleteComment(deleteCommentRequest.UserID, deleteCommentRequest.CommentID, role); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H {
				"error":err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":"comment deleted successfuly!",
		})
	}
}