package middlewares

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mochivi/real-time-leaderboards/internal/auth"
)

func abortWithError(c *gin.Context, statusCode int, message string) {
    c.JSON(statusCode, gin.H{"error": message})
    c.Abort()
}

func abortWithLoginRedirection(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
	c.Redirect(http.StatusFound, "api/v1/auth/login")
    c.Abort()
}

/* 
Auth middleware:
1. Ensures AccessToken is refreshed on every request as long as the user has a valid AccessToken or RefreshToken
2. RefreshToken itself is not refreshed, this is done to prevent bad actos from gaining control to an user account for longer than the TTL of the
	stolen credentials.

	AccessTokenTTL = 5 minutes
	RefreshTokenTTL = 30 minutes
	
	If only the AccessToken is stolen:
		Account will be compromised for only the TTL of the token
	If only the RefreshToken is stolen:
		1. Authorization header will be missing, completely preventing login from bad actors.
		2. If a valid JWT is provided in the authorization header just to fulfill validation and use the RefreshToken
			the account wll be compromised for the RefreshTokenTTL.

3. Ideally, SSO would be implemented alongside the JWT validation to ensure users can easily log back into the system
	whenever their RefreshTokens expire. The current implementation is secure but users will have to input their credentials every 30 minutes
*/


// JWTService is injected into the middleware by the server
func ValidateAuth(j auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// JWT access_token comes in the header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// abortWithError(c, http.StatusUnauthorized, "Missing Authorization header")
			abortWithLoginRedirection(c, http.StatusUnauthorized, "Missing Authorization header")
			return
		}

		// Get only the tokenString
		tokenString, ok := j.ParseTokenFromHeader(authHeader)
		if !ok {
			log.Println("Failed to parse JWT from header")
			abortWithError(c, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// Retrieve the user claims
		useRefreshToken := false
		userClaims, err := j.VerifyToken(tokenString)
		if err != nil {

			// If access_token is expired, try using the refresh_token
			if errors.Is(err, jwt.ErrTokenExpired) {
				useRefreshToken = true
			} else {
				log.Printf("Failed to verify JWT AccessToken: %v", err)
				abortWithError(c, http.StatusUnauthorized, "Unauthorized")
				return
			}
			
		}

		// We will try to use the refreshToken is the provided accessToken is expired
		if useRefreshToken {
			refreshTokenCookie, err := c.Request.Cookie("refresh_token")
				if err != nil {
					if errors.Is(err, http.ErrNoCookie) {
						abortWithError(c, http.StatusBadRequest, "Missing refresh_token cookie")
					} else {
						abortWithError(c, http.StatusInternalServerError, "Internal server error")
					}
					return
				}

				refreshTokenString := refreshTokenCookie.Value
				userClaims, err = j.VerifyToken(refreshTokenString)
				if err != nil {
					log.Printf("Failed to verify JWT RefreshToken: %v", err)
					abortWithError(c, http.StatusUnauthorized, "Unauthorized")
					return
				}
		}

		// Update user cookies with the validated roles
		tokens, err := j.CreateAccessTokens(userClaims.UserID, userClaims.Role)
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

		// Set the user claims in the request and move on to the next handler
		c.Set("UserClaims", userClaims)
		c.Next()
	}
}

// Admin validation should only be called from authenticated endpoints
func ValidateAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Receive userClaims from context
		claims, ok := c.Get("UserClaims")
		log.Printf("Admin request: %+v", claims)
		if !ok {
			abortWithError(c, http.StatusUnauthorized, "Missing user claims")
			return
		}

		// Ensure userClaims is of the correct type
		userClaims, ok := claims.(*auth.CustomClaims)
		if !ok {
			abortWithError(c, http.StatusInternalServerError, "Invalid user claims type")
			return
		}

		// Check if user role is administrator
		if userClaims.Role != "administrator" {
			abortWithError(c, http.StatusUnauthorized, "Administrator priviliges required")
			return
		}

		c.Next()
	}
}