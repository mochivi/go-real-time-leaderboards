package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mochivi/go-real-time-leaderboards/internal/models"
	"github.com/mochivi/go-real-time-leaderboards/internal/storage"
	redis "github.com/mochivi/go-real-time-leaderboards/internal/storage/redis"
)

type LeaderboardController struct {
	repo  storage.LeaderboardRepo
	redis redis.RedisService
}

func NewLeaderboardController(repo storage.LeaderboardRepo, redisService redis.RedisService) LeaderboardController {
	return LeaderboardController{
		repo:  repo,
		redis: redisService,
	}
}

// Any leaderboard level operation is restricted to admins only

// Returns a JSON encoded leaderboard
func (l LeaderboardController) Get(c *gin.Context) {
	leaderboardID := c.Param("id")
	if leaderboardID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"message": "missing leaderboard id",
		})
		return
	}

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
	if leaderboardID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"message": "missing leaderboard id",
		})
		return
	}

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

	// If leaderboard is created with live status, add to cache
	if leaderboard.Live {
		l.updateCache(c.Request.Context(), leaderboard)
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

	// If leaderboard is created as live, add to cache
	if updatedLeaderboard.Live {
		// Leaderboards will have their expiration updated to 2 hours after creation/updates.
		// This is nice so that recent games are also cached and users can look them up faster
		// The cost of keeping the recent games cached is not high if there is enough RAM available for it
		l.updateCache(c.Request.Context(), updatedLeaderboard)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    updatedLeaderboard,
		"message": "Leaderboard updated",
	})
}

func (l LeaderboardController) Delete(c *gin.Context) {
	leaderboardID := c.Param("id")
	if leaderboardID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"message": "missing leaderboard id",
		})
		return
	}

	if err := l.repo.Delete(c.Request.Context(), leaderboardID); err != nil {
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
			"error":   http.StatusInternalServerError,
			"message": errorMessage,
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (l LeaderboardController) updateCache(ctx context.Context, leaderboard *models.Leaderboard) error {
	if err := l.redis.Set(
		ctx,
		leaderboard.RedisKey(),
		leaderboard,
		7200, // 2 hours is the default
	); err != nil {
		log.Println("Failed to add leaderboard to redis: %v", err)
		return fmt.Errorf("failed to add leaderboard to redis: %w", err)
	}

	return nil
}
