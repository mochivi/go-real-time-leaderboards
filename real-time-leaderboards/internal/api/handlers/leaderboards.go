package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

func (l *LeaderboardController) CreateLeaderboard(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": "Leaderboard created",
	})
}