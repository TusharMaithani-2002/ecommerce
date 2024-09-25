package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	// "github.com/golang-jwt/jwt/v5"
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
		Name     string `json:"name" form:"user" binding:"required"`
		Address  string `json:"address" form:"address"`
		Email    string `json:"email" form:"email" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
		Role     string `json:"role" form:"role" binding:"required"`
	}
	return func(c *gin.Context) {
		var userBody User

		if err := c.BindJSON(&userBody); err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		user, err := u.userSevices.CreateUser(userBody.Name, userBody.Password, userBody.Email, userBody.Address, userBody.Role)

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
