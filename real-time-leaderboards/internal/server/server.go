package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mochivi/real-time-leaderboards/conf"
	"github.com/mochivi/real-time-leaderboards/internal/api/handlers"
	"github.com/mochivi/real-time-leaderboards/internal/auth"
	"github.com/mochivi/real-time-leaderboards/internal/storage/redis"
)

// Repositories are owned by each Controller
type DependencyContainer struct {
	Controllers struct {
		Leaderboards handlers.LeaderboardController
		Users handlers.UserController
	}
	Services struct {
		JWTService auth.JWTService
		RedisService redis.RedisService
	}
}

type Server struct {
	config conf.ServerConfig
	engine *gin.Engine
	dependencies DependencyContainer
}

func NewServer(cfg conf.ServerConfig, dependencies DependencyContainer) *Server {
	server := &Server{
		config: cfg,
		engine: gin.Default(),
		dependencies: dependencies,
	}
	server.mount()
	return server
}

func (s *Server) mount() {
	s.engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})
}

func (s *Server) Run() error {
	return s.engine.Run(s.config.Addr())
}