package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mochivi/real-time-leaderboards/internal/auth"
	"github.com/mochivi/real-time-leaderboards/internal/models"
	"github.com/mochivi/real-time-leaderboards/internal/storage"
)

type AuthController struct {
	repo storage.UserRepo
	jwtService auth.JWTService 
}

func NewAuthController(repo storage.UserRepo, jwtService auth.JWTService) AuthController {
	return AuthController{
		repo: repo,
		jwtService: jwtService,
	}
}

/* 
Auth handler is responsible for Login, Logout and TokenRefresh operations
	It differs from auth middleware in the sense that the handler will retrieve
	user information from the non-cached DB (redis), therefore, it is slower.
Middleware implementation will only verify if the tokens to maintain user session.
*/

func (a AuthController) Login(c *gin.Context) {

	// Bind request body with AuthRequest model
	authRequest := models.AuthRequest{}
	if err := c.ShouldBindBodyWithJSON(&authRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to parse login request from request body",
			"message": "Authentication failed",
		})
		return
	}

	// Check username and password in the repo
	user, err := a.repo.GetByUsername(c.Request.Context(), authRequest.Username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User does not exist",
		})
		return
	}

	// Validate provided password with the password hash
	if !user.ValidatePasswordHash(authRequest.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Username or password incorrect",
		})
		return
	}

	// Generate JWT token to send with response
	tokens, err := a.jwtService.CreateAccessTokens(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate JWT",
		})
		return
	}

	// Set tokens in the user cookies
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"access_token",
		tokens.AccessToken,
		3600,
		"",
		"",
		false,
		true,
	)
	c.SetCookie(
		"refresh_token",
		tokens.RefreshToken,
		3600,
		"",
		"",
		false,
		true,
	)
	c.JSON(http.StatusOK, gin.H{})
}

// Should revoke user tokens, for now, we will assume the front-end will clear the user local storage
func (a AuthController) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func (a AuthController) RefreshToken(c *gin.Context) {

	// Receive refresh token in the cookies from the request
	cookies := c.Request.Cookies()
	var refreshToken string
	for _, cookie := range cookies {
		if cookie.Name == "refresh_token" {
			refreshToken = cookie.Value
		}
	}

	// Validate if token was received
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Refresh token cookie not sent",
			"message": "Failed to refresh access token",
		})
		return
	}

	// Get userID & Role to generate the custom claims
	user := models.User{}
	if err := c.ShouldBindBodyWithJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing required user information",
			"message": "Failed to refresh access token",
		})
		return
	}

	// Create new access and refresh tokens for user
	tokens, err := a.jwtService.CreateAccessTokens(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate JWT",
		})
		return
	}

	// Set tokens in the user cookies
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"access_token",
		tokens.AccessToken,
		3600,
		"",
		"",
		false,
		true,
	)
	c.SetCookie(
		"refresh_token",
		tokens.RefreshToken,
		3600,
		"",
		"",
		false,
		true,
	)
	c.JSON(http.StatusOK, gin.H{})
}