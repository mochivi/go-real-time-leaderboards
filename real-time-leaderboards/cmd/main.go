package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/mochivi/go-real-time-leaderboards/config"
	"github.com/mochivi/go-real-time-leaderboards/internal/api/handlers"
	"github.com/mochivi/go-real-time-leaderboards/internal/auth"
	"github.com/mochivi/go-real-time-leaderboards/internal/server"
	"github.com/mochivi/go-real-time-leaderboards/internal/storage"
	"github.com/mochivi/go-real-time-leaderboards/internal/storage/redis"
	"github.com/mochivi/go-real-time-leaderboards/utils"
)

func main() {

	// Load environment variables, not needed in containers
	godotenv.Load()

	// Create config
	cfg := config.Config{
		ServerConfig: config.ServerConfig{
			Host: utils.GetEnvString("SERVER_HOSTNAME", "localhost"),
			Port: utils.GetEnvInt("SERVER_PORT", 8080),
		},
		DBConfig: config.DBConfig{
			Host: utils.GetEnvString("POSTGRES_HOST", "localhost"),
			Port: utils.GetEnvInt("POSTGRES_PORT", 5432),
			MaxOpenConns: utils.GetEnvInt("POSTGRES_MAX_OPEN_CONNS", 10),
			MaxIdleConns: utils.GetEnvInt("POSTGRES_MAX_IDLE_CONNS", 10),
			MaxIdleTime:  utils.GetEnvString("POSTGRES_MAX_IDLE_TIME", "30s"),
		},
		RedisConfig: config.RedisConfig{
			Host: utils.GetEnvString("REDIS_HOST", "localhost"),
			Port: utils.GetEnvInt("REDIS_PORT", 6379),
			Password: utils.GetEnvString("REDIS_PASSWORD", "redis"),
		},
	}

	// Initialize postgres database
	pgDB := initDB(cfg.DBConfig)
	defer pgDB.Close()

	// Initialize services
	jwtService, redisService := initServices(cfg.RedisConfig)

	// Initialize dependencies with concrete types
	dependencies := initDependencies(
		storage.NewLeaderboardRepoPG(pgDB),
		storage.NewUserRepoPG(pgDB),
		jwtService,
		redisService,
	)

	// Initialize server
	server := server.NewServer(cfg.ServerConfig, dependencies)

	// Run server
	log.Println("Starting server...")
	if err := server.Run(); err != nil {
		log.Fatalf("Server shutting down: %v", err)
	}

}

// Accept interfaces for creating any kind of configuration
func initDependencies(
	leaderboardRepo storage.LeaderboardRepo,
	userRepo storage.UserRepo,
	jwtService auth.JWTService, 
	redisService redis.RedisService,
) server.DependencyContainer {
	
	controllers := struct{
		Leaderboards handlers.LeaderboardController
		Auth handlers.AuthController
		Users handlers.UserController
	}{
		Leaderboards: handlers.NewLeaderboardController(leaderboardRepo),
		Auth: handlers.NewAuthController(userRepo, jwtService),
		Users: handlers.NewUserController(userRepo),
	}

	services := struct{
		JWTService auth.JWTService
		RedisService redis.RedisService
	}{
		JWTService: jwtService,
		RedisService: redisService,
	}

	dependencies := server.DependencyContainer{
		Controllers: controllers, 
		Services: services,
	}

	return dependencies
}

// Initialize services
func initServices(redisConfig config.RedisConfig) (auth.JWTService, redis.RedisService) {
	
	// Redis service
	log.Println("Connecting to redis database...")
	redisService := redis.NewRedisService(redisConfig)
	log.Println("Connected to redis db")

	// JWT service
	accessTokenTTL := utils.GetEnvInt("JWT_ACCESS_TOKEN_TTL", 5)
	refreshTokenTTL := utils.GetEnvInt("JWT_REFRESH_TOKEN_TTL", 5)
	jtwtService := auth.NewJWTService(
		utils.GetEnvString("JWT_SECRET", ""),
		time.Duration(accessTokenTTL) * time.Minute,
		time.Duration(refreshTokenTTL) * time.Minute,
	)

	return jtwtService, redisService
}

// Initializes database
func initDB(dbConfig config.DBConfig) *sql.DB {

	log.Println("Connecting to postgres database...")
	pgDB, err := storage.NewPostgres(
		dbConfig.DSN(),
		dbConfig.MaxOpenConns,
		dbConfig.MaxIdleConns,
		dbConfig.MaxIdleTime,
	)
	if err != nil {
		log.Fatalf("Failed to connect to postgres database: %v", err)
	}

	log.Println("Connected to postgres db")

	return pgDB
}