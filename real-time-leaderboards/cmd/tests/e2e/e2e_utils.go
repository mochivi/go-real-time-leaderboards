package e2e_test

/*

// Later modify this function to support e2e testing

func setupServer() *server.Server {
	serverConfig := config.ServerConfig{
		Host: utils.GetEnvString("SERVER_HOSTNAME", "localhost"),
		Port: utils.GetEnvInt("SERVER_PORT", 8080),
	}

	mockUserRepo := mocks.MockUserRepo{}
	mockLeaderboardsRepo := mocks.MockLeaderboardsRepo{}
	mockJWTService := mocks.MockJWTService{}
	mockRedisService := mocks.MockRedisService{}

	// Setup succesfull authentication
	mockJWTService.On("ParseTokenFromHeader").Return("tokenString", nil)
	mockJWTService.On("VerifyToken").Return(&auth.CustomClaims{
		UserID: "1",
		Role: "administrator",
	}, nil)
	mockJWTService.On("ParseTokenFromHeader").Return("tokenString", nil)

	controllers := struct{
		Leaderboards handlers.LeaderboardController
		Auth handlers.AuthController
		Users handlers.UserController
	}{
		Leaderboards: handlers.NewLeaderboardController(&mockLeaderboardsRepo),
		Auth: handlers.NewAuthController(&mockUserRepo, &mockJWTService),
		Users: handlers.NewUserController(&mockUserRepo),
	}

	services := struct{
		JWTService auth.JWTService
		RedisService redis.RedisService
	}{
		JWTService: &mockJWTService,
		RedisService: &mockRedisService,
	}

	dependencies := server.DependencyContainer{
		Controllers: controllers,
		Services: services,
	}

	return server.NewServer(serverConfig, dependencies)
}

*/