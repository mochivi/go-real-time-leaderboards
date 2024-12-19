package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mochivi/go-real-time-leaderboards/internal/auth"
	"github.com/mochivi/go-real-time-leaderboards/internal/mocks"
	"github.com/mochivi/go-real-time-leaderboards/internal/models"
	"github.com/stretchr/testify/assert"
)

// Represents the default JSON response from the API
// will always return at least one field, used to decode
// JSON sent from the handlers
type defaultUsersAPIResponse struct {
	Data *models.User `json:"data"`
	Message string `json:"message"`
	Error string `json:"error"`
}

func TestUsersGet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create mock and user controller
	mockUserRepo := setupMocks(
		"GetByID", 
		[]any{"1"},
		[]any{&models.User{	ID: "1"}, nil},
	)
	uc := NewUserController(mockUserRepo)

	// Execute request and received recorded and decoded response
	w, wResponse := executeRequest(
		"GET",
		"/api/v1/users/:id",
		"/api/v1/users/1",
		[]gin.HandlerFunc{uc.Get},
	)		

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "1", wResponse.Data.ID)
}

func TestUsersRegister(t *testing.T) {

	// Create request body
	registerUser := models.RegisterUser{
		Username: "test username",
		Email: "test@test.com",
		Password: "password123",
	}
	bodyBuf, _ := json.Marshal(registerUser)

	// Create mock and user controller
	mockUserRepo := setupMocks(
		"Create",
		[]any{&registerUser},
		[]any{&models.User{ID: "1"}, nil})
	uc := NewUserController(mockUserRepo)

	// Execute request and received recorded and decoded response
	w, wResponse := executeRequest(
		"POST",
		"/api/v1/users/register",
		"/api/v1/users/register",
		[]gin.HandlerFunc{uc.Register},
		bodyBuf,
	)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "1", wResponse.Data.ID)
	assert.Equal(t, "User created", wResponse.Message)
}

func TestUsersUpdate(t *testing.T) {
	// Prepare user to send in request body
	updateUser := models.UpdateUser{
		ID: "1",
		Username: "New username",
		Email: "test@test.com",
		Role: "administrator",
	}
	bodyBuf, _ := json.Marshal(updateUser)

	// Prepare user storage mock
	mockUserRepo := setupMocks(
		"Update", 
		[]any{&updateUser}, 
		[]any{&models.User{ID: "1"}, nil},
	)
	uc := NewUserController(mockUserRepo)

	// Setup handlers request will flow through
	testHandlers := []gin.HandlerFunc{
		mocks.MockValidateAuthMiddleware(&auth.CustomClaims{
			UserID: "1",
			Role: "administrator",
		}),
		uc.Update,
	}

	// Execute request and received recorded and decoded response
	w, wResponse := executeRequest(
		"PUT",
		"/api/v1/users",
		"/api/v1/users",
		testHandlers,
		bodyBuf,
	)

	// Assert responses
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "1", wResponse.Data.ID)
	assert.Equal(t, "User updated", wResponse.Message)
}

func TestUsersDelete(t *testing.T) {

	// Setup mock repo, assign functions to call and handlers to call in order
	mockUserRepo := setupMocks("Delete", []any{"1"}, []any{nil})
	uc := NewUserController(mockUserRepo)
	testHandlers := []gin.HandlerFunc{
		mocks.MockValidateAuthMiddleware(&auth.CustomClaims{
			UserID: "1",
			Role: "administrator",
		}),
		uc.Delete,
	}
	w, wResponse := executeRequest(
		"DELETE",
		"/api/v1/users/:id",
		"/api/v1/users/1",
		testHandlers,
	)

	// Assert responses
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "User deleted", wResponse.Message)
}

func setupMocks(funcName string, args, returns []any) *mocks.MockUserRepo {
	mockUserRepo := mocks.MockUserRepo{}
	mockUserRepo.On(funcName, args...).Return(returns...)
	return &mockUserRepo
}

func executeRequest(
	method string,
	requestPath string,
	requestedPath string,
	handlers []gin.HandlerFunc,
	bodyBuf ...[]byte,
) (*httptest.ResponseRecorder, defaultUsersAPIResponse) {
	gin.SetMode(gin.TestMode)
	
	// Setup router and routes
	router := gin.Default()
	router.Handle(method, requestPath, handlers...)
	
	// Check if received a report body, otherwise leave as nil
	var body io.Reader = nil
	if len(bodyBuf) > 0 {
		body = bytes.NewReader(bodyBuf[0])
	}

	// Create recorder, execute request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, requestedPath, body)
	router.ServeHTTP(w, req)

	// Decode response body
	wResponse := defaultUsersAPIResponse{}
	json.NewDecoder(w.Body).Decode(&wResponse)

	return w, wResponse
}