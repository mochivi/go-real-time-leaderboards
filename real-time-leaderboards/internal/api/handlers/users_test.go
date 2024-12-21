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
	"github.com/stretchr/testify/mock"
)

func TestUsersGet(t *testing.T) {

	// Setup test cases
	testCases := []struct {
		name string // Name of the test
		mockRepo *mocks.MockUserRepo // Used whenever the handler is expected to reach a mock call
		expectedStatus int // expected resulting status code
		requestOpts requestOpts // Anything to add to the header, body params for the request
	}{
		{
			name: "get user",
			mockRepo: setupMocks("GetByID", []any{"1"}, []any{&models.User{ID: "1"}, nil}),
			expectedStatus: http.StatusOK,
			requestOpts: requestOpts{
				params: map[string]string{
					"id": "1",
				},
			},
		},
		{
			name: "get user missing user id",
			expectedStatus: http.StatusBadRequest,
			requestOpts: requestOpts{
				params: map[string]string{
					"id": "",
				},
			},
		},
	}
	
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			
			// Recreate controller with new mock on every testcase
			uc := NewUserController(testCase.mockRepo)

			// Execute request and received recorded and decoded response
			w := executeRequest(
				"GET", 
				[]gin.HandlerFunc{uc.Get},
				testCase.requestOpts,
			)

			// Assert expectations
			assert.Equal(t, testCase.expectedStatus, w.Code)
			if testCase.mockRepo != nil {
				testCase.mockRepo.AssertExpectations(t)
			}
		})
	}
}

func TestUsersRegister(t *testing.T) {

	// Create request body
	registerUser := models.RegisterUser{
		Username: "test username",
		Email: "test@test.com",
		Password: "password123",
	}

	// Setup test cases
	testCases := []struct {
		name string // Name of the test
		mockRepo *mocks.MockUserRepo // Used whenever the handler is expected to reach a mock call
		expectedStatus int // expected resulting status code
		requestOpts requestOpts // Anything to add to the header, body params for the request
	}{
		{
			name: "register user",
			mockRepo: setupMocks("Create", []any{&registerUser}, []any{&models.User{ID: "1"}, nil}),
			expectedStatus: http.StatusCreated,
			requestOpts: requestOpts{body: registerUser},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			uc := NewUserController(testCase.mockRepo)

			// Execute request and received recorded and decoded response
			w := executeRequest(
				"POST", 
				[]gin.HandlerFunc{uc.Register},
				testCase.requestOpts,
			)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			if testCase.mockRepo != nil {
				testCase.mockRepo.AssertExpectations(t)
			}
		})
	}
}

func TestUsersUpdate(t *testing.T) {

	// Setup test cases
	testCases := []struct {
		name string // Name of the test
		mockRepo *mocks.MockUserRepo // Used whenever the handler is expected to reach a mock call
		userID string // User ID of who is making the request
		userRole string // Provide user role for authentication testing
		expectedStatus int // expected resulting status code
		requestOpts requestOpts // Anything to add to the header, body params for the request
	}{
		{
			name: "succesful update user",
			mockRepo: setupMocks("Update", []any{mock.AnythingOfType("*models.UpdateUser")}, []any{&models.User{ID: "1"}, nil}),
			userID: "1", // The ID of the user making the request
			userRole: "visitor",
			expectedStatus: http.StatusCreated,
			requestOpts: requestOpts{body: models.UpdateUser{
				ID: "1", // User ID to be updated
				Username: "New username",
				Email: "test@test.com",
				Role: "visitor",
			}},
		},
		{
			name: "admin update another user",
			mockRepo: setupMocks("Update", []any{mock.AnythingOfType("*models.UpdateUser")}, []any{&models.User{ID: "3"}, nil}),
			userID: "1",
			userRole: "administrator",
			expectedStatus: http.StatusCreated,
			requestOpts: requestOpts{body: models.UpdateUser{
				ID: "3",
				Username: "New username",
				Email: "test@test.com",
				Role: "visitor",
			}},
		},
		{
			name: "error update another user",
			mockRepo: &mocks.MockUserRepo{},
			userID: "1",
			userRole: "visitor",
			expectedStatus: http.StatusUnauthorized,
			requestOpts: requestOpts{body: models.UpdateUser{
				ID: "3",
				Username: "New username",
				Email: "test@test.com",
				Role: "visitor",
			}},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			uc := NewUserController(testCase.mockRepo)

			// Execute request and received recorded and decoded response
			w := executeRequest(
				"PUT", 
				[]gin.HandlerFunc{
					mocks.MockValidateAuthMiddleware(&auth.CustomClaims{
						UserID: testCase.userID,
						Role: testCase.userRole,
					}),
					uc.Update,
				},
				testCase.requestOpts,
			)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			if testCase.mockRepo != nil {
				testCase.mockRepo.AssertExpectations(t)
			}
		})
	}
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
	w := executeRequest("DELETE", testHandlers)

	// Assert responses
	assert.Equal(t, http.StatusOK, w.Code)
}

func setupMocks(funcName string, args, returns []any) *mocks.MockUserRepo {
	mockUserRepo := mocks.MockUserRepo{}
	mockUserRepo.On(funcName, args...).Return(returns...)
	return &mockUserRepo
}

type requestOpts struct {
	headers map[string]string
	body any
	params map[string]string
}

func (r requestOpts) Body() ([]byte, bool) {
	if r.body != nil {
		body, _ := json.Marshal(r.body)
		return body, true
	}
	return nil, false
}


func (r requestOpts) Headers() (map[string]string, bool) {
	if r.headers != nil {
		return r.headers, true
	}
	return nil, false
}

func (r requestOpts) Params() (map[string]string, bool) {
	if r.params != nil {
		return r.params, true
	}
	return nil, false
}

func executeRequest(
	method string,
	testHandlers []gin.HandlerFunc,
	requestOpts ...requestOpts,
) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)

	// Create recorder and gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(method, "/", nil)

	// Check if received a report body, otherwise leave as nil
	if len(requestOpts) > 0 {
		requestOpts := requestOpts[0]
		
		// Set body
		if body, ok := requestOpts.Body(); ok {
			c.Request.Body = io.NopCloser(bytes.NewReader(body)) 
		}

		// Set headers
		if headers, ok := requestOpts.Headers(); ok {
			for k, v := range headers {
				c.Request.Header.Set(k, v)
			}
		}

		// Set params
		if setParams, ok := requestOpts.Params(); ok {
			params := []gin.Param{}
			for k, v := range setParams {
				params = append(params, gin.Param{
					Key: k,
					Value: v,
				})
			}
			c.Params = params
		}
	}

	// Go one handler by one in the provided list
	for _, handler := range testHandlers {
		handler(c)
	}

	return w
}