package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	utils "github.com/nade-harlow/ecom-api/internal/app/utils/auth"
	"net/http"
	"strings"
)

func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header["Authorization"]

		if header == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "No token in header"})
			return
		}

		if header[0] == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "No token in header"})
			return
		}

		// get Token
		token := strings.Split(header[0], " ")[1]

		// verify token
		claims, err := utils.ValidateJwtAuthenticity(token)
		if err != nil {
			// update user record here
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		c.Set("user", claims)

		c.Next()

	}
}

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userD utils.DecodedUser
		user, found := c.Get("user")
		if !found {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized to perform this action"})
			return
		}

		err := mapstructure.Decode(user, &userD)
		if err != nil {
			fmt.Println("error decoding token", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong. Try again"})
			return
		}

		if userD.Role != "admin" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized to perform this action"})
			return
		}

		c.Next()
	}
}
