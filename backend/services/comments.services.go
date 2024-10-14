package services

import (
	"fmt"
	"time"

	"github.com/my_ecommerce/internal/dto"
	"github.com/my_ecommerce/internal/models"
	"gorm.io/gorm"
)

type CommentServices struct {
	db *gorm.DB
}

func (c* CommentServices) InitCommentService(database *gorm.DB) {
	c.db = database
	c.db.AutoMigrate(&models.Comment{})
}

func (c* CommentServices) AddComment(body string, userId, productId int) (*dto.CommentResponse, error) {

	comment := models.Comment{
		UserID: userId,
		Body: body,
		ProductID: productId,
		CreatedAt: time.Now(),
	}
	if err := c.db.Create(&comment).Error; err != nil {
		return nil, err
	}

	commentResponse := &dto.CommentResponse{
		ID: comment.ID,
		UserID: comment.UserID,
		ProductID: comment.ProductID,
		Likes: comment.Likes,
		Body: comment.Body,
	}
	return commentResponse, nil
}

func (c *CommentServices) AddLikeToComment(userId, commentId int) (*dto.CommentResponse, error) {

	comment := models.Comment{}
	if err := c.db.Where("id = ? and user_id = ?",commentId,userId).First(&comment).Error; err != nil {
		return nil, err
	}

	comment.Likes = comment.Likes + 1
	
	if err := c.db.Save(&comment).Error; err != nil {
		return nil, err
	}
	
	commentResponse := &dto.CommentResponse{
		ID: comment.ID,
		UserID: comment.UserID,
		ProductID: comment.ProductID,
		Likes: comment.Likes,
		Body: comment.Body,
	}
	return commentResponse, nil
}

func (c* CommentServices) UpdateComment(userId, commentId int, body string) (*dto.CommentResponse, error) {

	comment := models.Comment{}
	if err := c.db.Where("id = ? and user_id = ?",commentId,userId).First(&comment).Error; err != nil {
		return nil, err
	}

	comment.Body = body

	if err := c.db.Save(&comment).Error; err != nil {
		return nil, err
	}

	commentResponse := &dto.CommentResponse{
		ID: comment.ID,
		UserID: comment.UserID,
		ProductID: comment.ProductID,
		Likes: comment.Likes,
		Body: comment.Body,
	}
	return commentResponse, nil
}

func (c *CommentServices) DeleteComment(userId, commentId int, role string) error {

	comment := models.Comment{}
	if err := c.db.First(&comment,commentId).Error; err != nil {
		return err
	}

	if role != "admin" && comment.UserID != userId {
		return fmt.Errorf("not authorized to delete comment")
	}

	if err := c.db.Delete(&comment).Error; err != nil {
		return err
	}

	return nil
}