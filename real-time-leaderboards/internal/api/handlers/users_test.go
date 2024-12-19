package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mochivi/go-real-time-leaderboards/internal/mocks"
	"github.com/mochivi/go-real-time-leaderboards/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestUsersGet(t *testing.T) {
	// Create mock and user controller
	mockUserRepo := mocks.MockUserRepo{}
	mockUserRepo.On("GetByID", "1").Return(&models.User{
		ID: "1",
	}, nil)
	uc := NewUserController(&mockUserRepo)

	// Setup router
	router := gin.Default()
	router.GET("/api/v1/users/:id", uc.Get)
	
	// Setup recorder and make request, assert response
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users/1", nil)
	router.ServeHTTP(w, req)

	// Decode response JSON body
	wResponse := struct{
		Data models.User `json:"data"`
	}{}
	json.NewDecoder(w.Body).Decode(&wResponse)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "1", wResponse.Data.ID)
}

func TestUsersRegister(t *testing.T) {

	// Create mock and user controller
	mockUserRepo := mocks.MockUserRepo{}
	
	uc := NewUserController(&mockUserRepo)

	// Setup router
	router := gin.Default()
	router.POST("/api/v1/users/register", uc.Register)
	
	// Create request body
	registerUser := models.RegisterUser{
		Username: "test username",
		Email: "test@test.com",
		Password: "password123",
	}
	bufBody, _ := json.Marshal(registerUser)

	// Define mock functionality
	mockUserRepo.On(
		"Create", 
		&registerUser,
	).Return(&models.User{
		ID: "1",
	}, nil)

	// Setup recorder and make request, assert response
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/users/register", bytes.NewReader(bufBody))
	router.ServeHTTP(w, req)

	// Decode response JSON body
	wResponse := struct{
		Data models.User `json:"data"`
		Message string `json:"message"`
	}{}
	json.NewDecoder(w.Body).Decode(&wResponse)

	log.Printf("%+v\n", wResponse)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "1", wResponse.Data.ID)
	assert.Equal(t, "User created", wResponse.Message)
}