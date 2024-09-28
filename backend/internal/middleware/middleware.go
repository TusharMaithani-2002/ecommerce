package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/my_ecommerce/internal/utils"
)

func VerifyUser() gin.HandlerFunc {

	return func(c *gin.Context) {
		cookie,err := c.Cookie("go_ecommerce")
		
		if err != nil {
			c.JSON(401,gin.H{
				"error":"cookie not found! err: "+err.Error(),
			})
			c.Abort() 
			return
		}

		decoded, err := utils.DecodeJWT(cookie)
		if err != nil || decoded.Email == "" {
			c.JSON(401,gin.H{
				"error":"token invalid! err: "+err.Error(),
			})
			c.Abort() 
			return
		}
		expiryTime := decoded.ExpiresAt

		currentTime := time.Now().Unix()

		if currentTime > expiryTime {
			c.JSON(401,gin.H{
				"error": "token expired",
			})
			c.Abort() 
			return
		}

		c.Set("cookieId",decoded.ID)
		c.Set("role",decoded.Role)
		c.Next()

	}
}

/*
checking whether user is seller are same
*/
func verifyWriter() gin.HandlerFunc {

	return func(c *gin.Context) {
		
	}
}

/*
checking whether user is admin or writer
*/
func verrifyPrivilege() gin.HandlerFunc {
	return func(c *gin.Context) {}
}