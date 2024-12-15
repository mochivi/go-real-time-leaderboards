package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mochivi/real-time-leaderboards/conf"
	"github.com/mochivi/real-time-leaderboards/internal/api/handlers"
	"github.com/mochivi/real-time-leaderboards/internal/api/middlewares"
	"github.com/mochivi/real-time-leaderboards/internal/auth"
	"github.com/mochivi/real-time-leaderboards/internal/storage/redis"
)

// Repositories are owned by each Controller
type DependencyContainer struct {
	Controllers struct {
		Leaderboards handlers.LeaderboardController
		Auth handlers.AuthController
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
	apiGroup := s.engine.Group("/api")
	v1Group := apiGroup.Group("/v1")

	// Auth endpoints
	authGoup := v1Group.Group("/auth")
	{	
		authGoup.POST("/login", s.dependencies.Controllers.Auth.Login)
		authGoup.POST("/logout", s.dependencies.Controllers.Auth.Logout)
		authGoup.POST(
			"/refresh", 
			middlewares.ValidateAuth(s.dependencies.Services.JWTService),
			s.dependencies.Controllers.Auth.RefreshToken,
		)
	}

	// User endpoints
	publicUsersGroup := v1Group.Group("/users")
	{ // Viewing and registering users does not require authentication
		publicUsersGroup.POST("/register", s.dependencies.Controllers.Users.Register)
		publicUsersGroup.GET("/:id", s.dependencies.Controllers.Users.Get)
	}
	authUsersGroup := v1Group.Group("/users", middlewares.ValidateAuth(s.dependencies.Services.JWTService))
	{ // Updating and deleting users requires authentication
		authUsersGroup.PUT("/", s.dependencies.Controllers.Users.Update)
		authUsersGroup.DELETE("/:id", s.dependencies.Controllers.Users.Delete)
	}

	// Leaderboard endpoints
	publicleaderboardsGroup := v1Group.Group("/leaderboards")
	{
		publicleaderboardsGroup.GET("/:id", s.dependencies.Controllers.Leaderboards.Get)
		publicleaderboardsGroup.GET("/entries/:id", s.dependencies.Controllers.Leaderboards.GetEntries)
	}
	adminleaderboardsGroup := v1Group.Group("/leaderboards")
	{
		adminleaderboardsGroup.POST("/", s.dependencies.Controllers.Leaderboards.Create)
		adminleaderboardsGroup.POST("/entries",	s.dependencies.Controllers.Leaderboards.CreateEntry)
		adminleaderboardsGroup.PUT("/", s.dependencies.Controllers.Leaderboards.Update)
		adminleaderboardsGroup.DELETE("/:id",s.dependencies.Controllers.Leaderboards.Delete)
	}
	
	s.engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})
}

func (s *Server) Run() error {
	return s.engine.Run(s.config.GetPort())
}