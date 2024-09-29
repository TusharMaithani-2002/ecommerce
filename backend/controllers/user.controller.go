package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	// "github.com/golang-jwt/jwt/v5"
	"github.com/my_ecommerce/internal/dto"
	"github.com/my_ecommerce/internal/middleware"
	"github.com/my_ecommerce/internal/utils"
	"github.com/my_ecommerce/services"
)

type UserController struct {
	userSevices services.UserService
}

func (u *UserController) InitUserController(router *gin.Engine, userService services.UserService) {

	userRouter := router.Group("/user")
	userRouter.GET("/:id", u.getUserById())
	userRouter.POST("/register", u.registerUser())
	userRouter.POST("/login", u.loginUser())
	userRouter.DELETE("/delete/:id", middleware.VerifyUser(), u.deleteUser())
	userRouter.POST("/logout", u.logOut())
	userRouter.PATCH("/update/:id", middleware.VerifyUser(), u.updateUser())
	u.userSevices = userService
}

func (u *UserController) getUserById() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		// converting my string id to int
		numId, err := strconv.Atoi(id)

		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		user, err := u.userSevices.GetUser(numId)

		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
		}

		c.JSON(200, gin.H{
			"data": user,
		})
	}
}

func (u *UserController) registerUser() gin.HandlerFunc {

	type User struct {
		Name        string `json:"name" form:"user" binding:"required"`
		Address     string `json:"address" form:"address"`
		Email       string `json:"email" form:"email" binding:"required"`
		Password    string `json:"password" form:"password" binding:"required"`
		Role        string `json:"role" form:"role" binding:"required"`
		PhoneNumber string `json:"phoneNumber" form:"phoneNumber"`
	}
	return func(c *gin.Context) {
		var userBody User

		if err := c.BindJSON(&userBody); err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		user, err := 
		u.userSevices.
		CreateUser(userBody.Name, userBody.Password, userBody.Email, userBody.Address, userBody.Role,userBody.PhoneNumber)

		if err != nil {
			c.JSON(400, gin.H{
				"error": err,
			})
			return
		}

		jwtToken, err := utils.GenerateJWT(user.Email, user.Role, user.ID)

		if err != nil {
			c.JSON(http.StatusMethodNotAllowed, gin.H{
				"error":   err.Error(),
				"message": "user created! but jwt creation failed",
			})
			return
		}

		c.SetCookie("go_ecommerce", jwtToken, 259200, "/", "localhost", false, true)

		c.JSON(201, gin.H{
			"data": user,
		})
	}
}

func (u *UserController) loginUser() gin.HandlerFunc {

	type Login struct {
		Email    string `json:"email" form:"email" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
	}

	return func(c *gin.Context) {
		var loginData Login

		if err := c.BindJSON(&loginData); err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		loggedInUser, err := u.userSevices.LoginUser(loginData.Email, loginData.Password)
		if err != nil {
			c.JSON(404, gin.H{
				"error": err.Error(),
			})
			return
		}

		jwtToken, err := utils.GenerateJWT(loggedInUser.Email, loggedInUser.Role, loggedInUser.ID)

		if err != nil {
			c.JSON(http.StatusMethodNotAllowed, gin.H{
				"error":   err.Error(),
				"message": "unable to login! jwt creation failed",
			})
			return
		}

		c.SetCookie("go_ecommerce", jwtToken, 259200, "/", "localhost", false, true)

		c.JSON(200, gin.H{
			"data": loggedInUser,
		})

	}
}

func (u *UserController) deleteUser() gin.HandlerFunc {

	return func(c *gin.Context) {
		id := c.Param("id")
		numId, err := strconv.Atoi(id)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "userId not valid",
			})
			return
		}

		if err := u.userSevices.DeleteUser(numId); err != nil {
			c.JSON(http.StatusExpectationFailed, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "user deleted successfully",
		})

	}
}

func (u *UserController) logOut() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("go_ecommerce", "", -1, "/", "localhost", false, true)

		c.JSON(http.StatusOK, gin.H{
			"message": "user logged out successfully",
		})
	}
}

func (u *UserController) updateUser() gin.HandlerFunc {

	type UpdateUserRequest struct {
		Name        *string `json:"name,omitempty" form:"name,omitempty"`
		Address     *string `json:"address,omitempty" form:"address,omitempty"`
		Role        *string `json:"role,omitempty" form:"role,omitempty"`
		PhoneNumber *string `json:"phoneNumber,omitempty" form:"phoneNumber,omitempty"`
	}

	return func(c *gin.Context) {
		id := c.Param("id")
		numId, err := strconv.Atoi(id)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "user Id not valid number",
			})
			return
		}

		cookieId, ok := c.Get("cookieId")
		userId := cookieId.(int)

		if userId != numId {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "you are not authorized to the update this user",
			})
			return
		}

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "cookie is invalid",
			})
			return
		}

		var userRequest UpdateUserRequest
		if err := c.BindJSON(&userRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		request := dto.UpdatedUser{
			Name:        userRequest.Name,
			Role:        userRequest.Role,
			Address:     userRequest.Address,
			PhoneNumber: userRequest.PhoneNumber,
		}
		userResponse, err := u.userSevices.UpdateUser(numId, request)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"message": "error while updating user",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": userResponse,
		})
	}
}
