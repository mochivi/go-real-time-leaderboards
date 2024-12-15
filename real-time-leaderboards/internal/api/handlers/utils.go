package handlers

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mochivi/real-time-leaderboards/internal/auth"
)

// Parses the user claims from the request
// This will only work with endpoints that use the validateAuth middleware
func parseUserClaims(c *gin.Context) (*auth.CustomClaims, error) {
	
	// Get user claims
	claims, ok := c.Get("UserClaims")
	log.Printf("Admin request: %+v", claims)
	if !ok {
		return nil, errors.New("not enough privileges")
	}

	// Ensure userClaims is of the correct type
	userClaims, ok := claims.(*auth.CustomClaims)
	if !ok {
		return nil, errors.New("invalid user claims type")
	}

	return userClaims, nil
}