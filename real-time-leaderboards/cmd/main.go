package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/mochivi/real-time-leaderboards/conf"
	"github.com/mochivi/real-time-leaderboards/internal/api/handlers"
	"github.com/mochivi/real-time-leaderboards/internal/auth"
	"github.com/mochivi/real-time-leaderboards/internal/server"
	"github.com/mochivi/real-time-leaderboards/internal/storage"
	"github.com/mochivi/real-time-leaderboards/internal/storage/redis"
	"github.com/mochivi/real-time-leaderboards/utils"
)


func main() {

	// Load environment variables, not needed in containers
	godotenv.Load()

	// Create config
	cfg := conf.Config{
		ServerConfig: conf.ServerConfig{
			Host: utils.GetEnvString("SERVER_HOSTNAME", "localhost"),
			Port: utils.GetEnvInt("SERVER_PORT", 8080),
		},
		DBConfig: conf.DBConfig{
			Host: utils.GetEnvString("POSTGRES_HOST", "localhost"),
			Port: utils.GetEnvInt("POSTGRES_PORT", 5432),
			MaxOpenConns: utils.GetEnvInt("POSTGRES_MAX_OPEN_CONNS", 10),
			MaxIdleConns: utils.GetEnvInt("POSTGRES_MAX_IDLE_CONNS", 10),
			MaxIdleTime:  utils.GetEnvString("POSTGRES_MAX_IDLE_TIME", "30s"),
		},
		RedisConfig: conf.RedisConfig{
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

	// Initialize dependencies
	dependencies := initDependencies(pgDB, jwtService, redisService)

	// Initialize server
	server := server.NewServer(cfg.ServerConfig, dependencies)

	// Run server
	log.Println("Starting server...")
	if err := server.Run(); err != nil {
		log.Fatalf("Server shutting down: %v", err)
	}

}

// Initializes database
func initDB(dbConfig conf.DBConfig) *sql.DB {

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

func initServices(redisConfig conf.RedisConfig) (auth.JWTService, redis.RedisService) {
	
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

func initDependencies(pgDB *sql.DB, jwtService auth.JWTService, redisService redis.RedisService) server.DependencyContainer {
	
	// Initialize dependencies, supporting dependency injection
	// Provide concrete implementations
	dependencies := server.DependencyContainer{
		
		Controllers: struct{
			Leaderboards handlers.LeaderboardController
			Auth handlers.AuthController
			Users handlers.UserController
		}{
			Leaderboards: handlers.NewLeaderboardController(
				storage.NewLeaderboardRepoPG(pgDB),
			),
			Auth: handlers.NewAuthController(
				storage.NewUserRepoPG(pgDB),
				jwtService,
			),
			Users: handlers.NewUserController(
				storage.NewUserRepoPG(pgDB),
			),
		},

		Services: struct{
			JWTService auth.JWTService
			RedisService redis.RedisService
		}{
			JWTService: jwtService,
			RedisService: redisService,
		},
	}


	return dependencies
}