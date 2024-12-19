package mocks

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mochivi/go-real-time-leaderboards/internal/auth"
)

func MockValidateAuthMiddleware(userClaims *auth.CustomClaims) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("UserClaims", userClaims)
		c.Next()
	}
}

func MockValidateAdminMiddleware(userClaims *auth.CustomClaims) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("UserClaims", userClaims)

		// Check if user role is administrator
		if userClaims.Role != "administrator" {
			c.Abort()
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Not enough priviliges",
			})
			return
		}

		c.Next()
	}
}