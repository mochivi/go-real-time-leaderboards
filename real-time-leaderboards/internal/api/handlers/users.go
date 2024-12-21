package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mochivi/go-real-time-leaderboards/internal/models"
	"github.com/mochivi/go-real-time-leaderboards/internal/storage"
)

type UserController struct {
	repo storage.UserRepo
}

func NewUserController(repo storage.UserRepo) UserController {
	return UserController{
		repo: repo,
	}
}


// Create a new user
func (u UserController) Register(c *gin.Context) {

	var registerUser models.RegisterUser
	if err := c.ShouldBindBodyWithJSON(&registerUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	// hash the provided password
	passwordHash, err := registerUser.HashPassword()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	// add user to database
	user, err := u.repo.Create(c.Request.Context(), &registerUser, passwordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": user,
		"message": "User created",
	})
}


func (u UserController) Get(c *gin.Context) {

	// Get id from request
	userID, ok := c.Params.Get("id")
	if !ok || userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID not provided in the request",
			"message": "Bad request",
		})
		return
	}

	// Request user from db
	user, err := u.repo.GetByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found by ID",
			"message": "User not found",
		})
		return
	}

	// Return user
	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}


func (u UserController) Update(c *gin.Context) {

	// Receive data
	var updateUser models.UpdateUser
	if err := c.ShouldBindBodyWithJSON(&updateUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	// Get user claims
	userClaims, err := parseUserClaims(c)
	if err != nil {
		log.Printf("Failed to parse user claims from context: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Validate user can update the provided user
	if userClaims.UserID != updateUser.ID && userClaims.Role != "administrator" {
		log.Printf(
			"User with claimed ID '%s' and role '%s' tried updating user of ID '%s'", 
			userClaims.UserID,
			userClaims.Role,
			updateUser.ID,
		)
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Not enough priviliges for this operation",
		})
		return
	}

	// Update user to database
	user, err := u.repo.Update(c.Request.Context(), &updateUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"message": "Database error",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": user,
		"message": "User updated",
	})
}


func (u UserController) Delete(c *gin.Context) {
	
	// Get id from request
	userID, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID not provided in the request",
			"message": "Bad request",
		})
		return
	}

	// Get user claims
	userClaims, err := parseUserClaims(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Validate user can delete the provided user
	if userClaims.UserID != userID || userClaims.Role != "administrator" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Not enough priviliges for this operation",
		})
	}

	if err := u.repo.Delete(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"message": "Could not delete user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted",
	})
}
