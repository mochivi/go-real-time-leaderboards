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
func (l LeaderboardController) Get(c *gin.Context) {
	leaderboardID := c.Param("id")
	leaderboard, err := l.repo.Get(c.Request.Context(), leaderboardID)
	if err != nil {
		var errorMessage string
		var statusCode int
		switch err {
			case storage.ErrNotFound:
				statusCode = http.StatusNotFound
				errorMessage = "Leaderboard not found"
			default:
				statusCode = http.StatusInternalServerError
				errorMessage = "Something went wrong"
		}
		c.JSON(statusCode, gin.H{
			"error":   "Failed to get data",
			"message": errorMessage,
		})
		return
	}

	// leaderboard exists
	c.JSON(http.StatusOK, gin.H{
		"data": leaderboard,
	})
}

func (l LeaderboardController) GetEntries(c *gin.Context) {
	leaderboardID := c.Param("id")

	leaderboardEntries, err := l.repo.GetEntries(c.Request.Context(), leaderboardID)
	if err != nil {
		var errorMessage string
		var statusCode int
		switch err {
			case storage.ErrNotFound:
				statusCode = http.StatusNotFound
				errorMessage = "Leaderboard not found"
			default:
				statusCode = http.StatusInternalServerError
				errorMessage = "Something went wrong"
		}
		c.JSON(statusCode, gin.H{
			"error":   "Failed to get data",
			"message": errorMessage,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": leaderboardEntries,
	})
}

func (l LeaderboardController) Create(c *gin.Context) {

	newLeaderboardRequest := models.LeaderboardRequest{}
	if err := c.ShouldBindBodyWithJSON(&newLeaderboardRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   http.StatusBadRequest,
			"message": "Invalid request body",
		})
		return
	}

	leaderboard, err := l.repo.Create(c.Request.Context(), &newLeaderboardRequest)
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
		"data":    *leaderboard,
		"message": "Leaderboard created",
	})
}

func (l LeaderboardController) CreateEntry(c *gin.Context) {

	// Bind entry request with fields that can be provided upon creation
	leaderboardEntryRequest := models.LeaderboardEntryRequest{}
	if err := c.ShouldBindBodyWithJSON(&leaderboardEntryRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   http.StatusBadRequest,
			"message": "Invalid request body",
		})
		return
	}
	leaderboardEntryRequest.AddUpdatedAt()

	// Create entry in the database
	leaderboardEntry, err := l.repo.CreateEntry(c.Request.Context(), &leaderboardEntryRequest)
	if err != nil {
		var errorMessage string
		switch err {
			case storage.ErrConflict:
				errorMessage = "Leaderboard entry already exists"
			default:
				errorMessage = "Something went wrong"
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   http.StatusInternalServerError,
			"message": errorMessage,
		})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"data":    leaderboardEntry,
		"message": "Leaderboard entry added",
	})
}

func (l LeaderboardController) Update(c *gin.Context) {
	
	leaderboard := models.UpdateLeaderboardRequest{}
	if err := c.ShouldBindBodyWithJSON(&leaderboard); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   http.StatusBadRequest,
			"message": "Invalid request body",
		})
		return
	}
	leaderboard.AddUpdatedAt()

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

func (l LeaderboardController) Delete(c *gin.Context) {
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
