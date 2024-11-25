package handlers

import "github.com/mochivi/real-time-leaderboards/internal/storage"

type UserController struct {
	repo storage.UserRepo
}

func NewUserController(repo storage.UserRepo) UserController {
	return UserController{
		repo: repo,
	}
}