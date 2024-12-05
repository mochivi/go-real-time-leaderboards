package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mochivi/real-time-leaderboards/internal/models"
	"github.com/mochivi/real-time-leaderboards/internal/storage"
)

type LeaderboardController struct {
	repo storage.LeaderboardRepo
}

func NewLeaderboardController(repo storage.LeaderboardRepo) LeaderboardController {
	return LeaderboardController{
		repo: repo,
	}
}

// Any leaderboard level operation is restricted to admins only

// Returns a JSON encoded leaderboard
func (l *LeaderboardController) Get(c *gin.Context) {
	leaderboardID := c.Param("id")

	leaderboard, err := l.repo.Get(c.Request.Context(), leaderboardID)
	if err != nil {
		var errorMessage string
		switch err {
			case storage.ErrNotFound:
				errorMessage = "Leaderboard not found"
			default:
				errorMessage = "Something went wrong"
		}
		c.JSON(http.StatusNotFound, gin.H{
			"error":   http.StatusInternalServerError,
			"message": errorMessage,
		})
		return
	}

	// leaderboard exists
	c.JSON(http.StatusOK, gin.H{
		"data": leaderboard,
	})
}


func (l *LeaderboardController) Create(c *gin.Context) {

	newLeaderboard := models.Leaderboard{}
	if err := c.ShouldBindBodyWithJSON(&newLeaderboard); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   http.StatusBadRequest,
			"message": "Invalid request body",
		})
		return
	}

	leaderboard, err := l.repo.Create(c.Request.Context(), &newLeaderboard)
	if err != nil {
		var errorMessage string
		switch err {
			case storage.ErrConflict:
				errorMessage = "Leaderboard already exists"
			default:
				errorMessage = "Something went wrong"
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   http.StatusInternalServerError,
			"message": errorMessage,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    leaderboard,
		"message": "Leaderboard created",
	})
}



func (l *LeaderboardController) Update(c *gin.Context) {
	
	leaderboard := models.Leaderboard{}
	if err := c.ShouldBindBodyWithJSON(&leaderboard); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   http.StatusBadRequest,
			"message": "Invalid request body",
		})
		return
	}

	updatedLeaderboard, err := l.repo.Update(c.Request.Context(), &leaderboard)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    updatedLeaderboard,
		"message": "Leaderboard updated",
	})
}

func (l *LeaderboardController) Delete(c *gin.Context) {
	leaderboardID := c.Param("id")

	if err := l.repo.Delete(c.Request.Context(), leaderboardID); err != nil {
		var errorMessage string
		switch err {
			case storage.ErrNotFound:
				errorMessage = "Leaderboard not found"
			default:
				errorMessage = "Something went wrong"
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   http.StatusInternalServerError,
			"message": errorMessage,
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
