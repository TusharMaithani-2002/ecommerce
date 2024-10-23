package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/my_ecommerce/internal/dto"
	"github.com/my_ecommerce/internal/models"
	"gorm.io/gorm"
)

type CartServices struct {
	db *gorm.DB
}

func (c *CartServices) InitCartServices(database *gorm.DB) {
	c.db = database
	c.db.AutoMigrate(&models.Cart{}, &models.CartItem{})
}

func (c *CartServices) GetCart(userId int) (*dto.CartResponsePopulated, error) {

	// checking whether cart exists -> createa  -> show cart
	var cart models.Cart
	if err := c.db.Preload("Items").Where("user_id = ?", userId).
		First(&cart).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var cartResponse dto.CartResponsePopulated

	if cart.ID == 0 {
		// cart doesn't exists
		// create a cart

		cartBody := models.Cart{
			UserID: userId,
		}

		if err := c.db.Create(&cartBody).Error; err != nil {
			return nil, err
		}

		cartResponse.ID = cartBody.ID
		cartResponse.UserID = cartBody.UserID
	} else {
		// return cart
		cartResponse.ID = cart.ID
		cartResponse.UserID = cart.UserID

		var cartItemResponse []dto.CartProductResponse
		for _, item := range cart.Items {
			var cartItem dto.CartProductResponse
			var product models.Product
			// come up with something better
			if err := c.db.Where("id = ?", item.ProductID).First(&product).Error; err != nil {
				continue
			}

			cartItem.ID = item.ID
			cartItem.Quantity = item.Quantity
			cartItem.Name = product.Name
			cartItem.Price = product.Price

			if len(product.Images) > 0 {
				var images []string

				if err := json.Unmarshal(product.Images, &images); err != nil {
					fmt.Println("Error unmarshalling images: ", err.Error())
				}

				if len(images) > 0 {
					cartItem.Image = images[0]
				}
			}

			cartItemResponse = append(cartItemResponse, cartItem)
		}
		cartResponse.Products = cartItemResponse
	}

	return &cartResponse, nil
}

func (c *CartServices) AddItem(cartId, userId, productId, quantity int) (*dto.CartResponse, error) {

	// add new item
	// increase quantity

	// first check whether the cart is created or not
	var cart models.Cart
	tx := c.db.Begin()

	if err := tx.Preload("Items").Where("id = ? and user_id = ?", cartId, userId).
		First(&cart).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return nil, err
	}

	if cart.ID == 0 {
		c.GetCart(userId)
		// fetching cart after its created
		if err := tx.Where("id = ? and user_id = ?", cartId, userId).
			First(&cart).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return nil, err
		}
	}

	itemExists := false

	// for changes to persist outside loop we have to access loop using index or by using pointers

	for i := range cart.Items {
		if cart.Items[i].ProductID == productId {
			// item present in cart
			cart.Items[i].Quantity = cart.Items[i].Quantity + 1

			if err := tx.Save(&cart.Items[i]).Error; err != nil {
				tx.Rollback()
				return nil, err
			}

			itemExists = true
			break
		}
	}

	if !itemExists {
		// add item to cart
		newItem := models.CartItem{
			ProductID: productId,
			CartID:    cartId,
			Quantity:  1,
		}

		if err := tx.Create(&newItem).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		items := cart.Items
		items = append(items, newItem)
		cart.Items = items
	}

	if err := tx.Save(&cart).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	var cartResponse dto.CartResponse
	cartResponse.ID = cart.ID
	cartResponse.UserID = cart.UserID

	var cartItemResponse []dto.CartItemResponse
	for _, item := range cart.Items {

		var cartItem dto.CartItemResponse

		cartItem.Quantity = item.Quantity
		cartItem.ProductID = item.ProductID

		cartItemResponse = append(cartItemResponse, cartItem)
	}
	cartResponse.Products = cartItemResponse

	return &cartResponse, nil
}

func (c *CartServices) DecrementItem(cartId, productId int) (*dto.CartResponse, error) {
	// find cart item and remove it
	var cartItem models.CartItem

	if err := c.db.Where("cart_id = ? and product_id = ?", cartId, productId).First(&cartItem).Error; err != nil {
		return nil, err
	}

	if cartItem.Quantity == 1 {
		if err := c.RemoveItem(cartId, productId); err != nil {
			return nil, err
		}
		cartItem.Quantity -= 1
	} else {
		cartItem.Quantity -= 1
		if err := c.db.Save(&cartItem).Error; err != nil {
			return nil, err
		}
	}

	var cartResponse dto.CartResponse
	var items []dto.CartItemResponse
	items = append(items, dto.CartItemResponse{
		ProductID: cartItem.ProductID,
		Quantity:  cartItem.Quantity,
	})

	cartResponse.ID = cartId
	cartResponse.Products = items

	return &cartResponse, nil

}

func (c *CartServices) RemoveItem(cartId, productId int) error {

	if err := c.db.Delete(models.CartItem{}, "product_id = ? and cart_id = ?", productId, cartId).Error; err != nil {
		return err
	}
	return nil
}
