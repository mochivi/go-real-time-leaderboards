package handlers

import (
	"net/http"
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
		name           string              // Name of the test
		mockRepo       *mocks.MockUserRepo // Used whenever the handler is expected to reach a mock call
		expectedStatus int                 // expected resulting status code
		requestOpts    requestOpts         // Anything to add to the header, body params for the request
	}{
		{
			name:           "get user",
			mockRepo:       setupUserRepoMock("GetByID", []any{"1"}, []any{&models.User{ID: "1"}, nil}),
			expectedStatus: http.StatusOK,
			requestOpts: requestOpts{
				params: map[string]string{
					"id": "1",
				},
			},
		},
		{
			name:           "get user missing user id",
			expectedStatus: http.StatusBadRequest,
			requestOpts: requestOpts{
				params: map[string]string{
					"id": "",
				},
			},
		},
		{
			name:           "get user db error",
			mockRepo:       setupUserRepoMock("GetByID", []any{"1"}, []any{&models.User{}, ErrRepoOperation}),
			expectedStatus: http.StatusNotFound,
			requestOpts: requestOpts{
				params: map[string]string{
					"id": "1",
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
		Email:    "test@test.com",
		Password: "password123",
	}

	// Setup test cases
	testCases := []struct {
		name           string              // Name of the test
		mockRepo       *mocks.MockUserRepo // Used whenever the handler is expected to reach a mock call
		expectedStatus int                 // expected resulting status code
		requestOpts    requestOpts         // Anything to add to the header, body params for the request
	}{
		{
			name:           "register user",
			mockRepo:       setupUserRepoMock("Create", []any{&registerUser}, []any{&models.User{ID: "1"}, nil}),
			expectedStatus: http.StatusCreated,
			requestOpts:    requestOpts{body: registerUser},
		},
		{
			name:           "register user db error",
			mockRepo:       setupUserRepoMock("Create", []any{&registerUser}, []any{&models.User{}, ErrRepoOperation}),
			expectedStatus: http.StatusInternalServerError,
			requestOpts:    requestOpts{body: registerUser},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			uc := NewUserController(testCase.mockRepo)

			// Execute request and received recorded and decoded response
			w := executeRequest(
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
		name           string              // Name of the test
		mockRepo       *mocks.MockUserRepo // Used whenever the handler is expected to reach a mock call
		userID         string              // User ID of who is making the request
		userRole       string              // Provide user role for authentication testing
		expectedStatus int                 // expected resulting status code
		requestOpts    requestOpts         // Anything to add to the header, body params for the request
	}{
		{
			name:           "succesful update user",
			mockRepo:       setupUserRepoMock("Update", []any{mock.AnythingOfType("*models.UpdateUser")}, []any{&models.User{ID: "1"}, nil}),
			userID:         "1", // The ID of the user making the request
			userRole:       "visitor",
			expectedStatus: http.StatusCreated,
			requestOpts: requestOpts{body: models.UpdateUser{
				ID:       "1", // User ID to be updated
				Username: "New username",
				Email:    "test@test.com",
				Role:     "visitor",
			}},
		},
		{
			name:           "admin update another user",
			mockRepo:       setupUserRepoMock("Update", []any{mock.AnythingOfType("*models.UpdateUser")}, []any{&models.User{ID: "3"}, nil}),
			userID:         "1",
			userRole:       "administrator",
			expectedStatus: http.StatusCreated,
			requestOpts: requestOpts{body: models.UpdateUser{
				ID:       "3",
				Username: "New username",
				Email:    "test@test.com",
				Role:     "visitor",
			}},
		},
		{
			name:           "error update another user",
			mockRepo:       &mocks.MockUserRepo{},
			userID:         "1",
			userRole:       "visitor",
			expectedStatus: http.StatusUnauthorized,
			requestOpts: requestOpts{body: models.UpdateUser{
				ID:       "3",
				Username: "New username",
				Email:    "test@test.com",
				Role:     "visitor",
			}},
		},
		{
			name: "update user db error",
			mockRepo: setupUserRepoMock(
				"Update",
				[]any{mock.AnythingOfType("*models.UpdateUser")},
				[]any{new(models.User), ErrRepoOperation},
			),
			userID:         "1",
			userRole:       "visitor",
			expectedStatus: http.StatusInternalServerError,
			requestOpts: requestOpts{body: models.UpdateUser{
				ID:       "1",
				Username: "New username",
				Email:    "test@test.com",
				Role:     "visitor",
			}},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			uc := NewUserController(testCase.mockRepo)

			// Execute request and received recorded and decoded response
			w := executeRequest(
				[]gin.HandlerFunc{
					mocks.MockValidateAuthMiddleware(&auth.CustomClaims{
						UserID: testCase.userID,
						Role:   testCase.userRole,
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

	testCases := []struct {
		name           string
		mockRepo       *mocks.MockUserRepo
		userID         string
		userRole       string
		expectedStatus int
		requestOpts    requestOpts
	}{
		{
			name:           "succesful delete user",
			mockRepo:       setupUserRepoMock("Delete", []any{"1"}, []any{nil}),
			userID:         "1", // The ID of the user making the request
			userRole:       "visitor",
			expectedStatus: http.StatusOK,
			requestOpts:    requestOpts{params: map[string]string{"id": "1"}},
		},
		{
			name:           "admin delete another user",
			mockRepo:       setupUserRepoMock("Delete", []any{"3"}, []any{nil}),
			userID:         "1",
			userRole:       "administrator",
			expectedStatus: http.StatusOK,
			requestOpts:    requestOpts{params: map[string]string{"id": "3"}},
		},
		{
			name:           "error delete another user",
			mockRepo:       &mocks.MockUserRepo{},
			userID:         "1",
			userRole:       "visitor",
			expectedStatus: http.StatusUnauthorized,
			requestOpts:    requestOpts{params: map[string]string{"id": "3"}},
		},
		{
			name:           "update user db error",
			mockRepo:       setupUserRepoMock("Delete", []any{"1"}, []any{ErrRepoOperation}),
			userID:         "1",
			userRole:       "visitor",
			expectedStatus: http.StatusInternalServerError,
			requestOpts:    requestOpts{params: map[string]string{"id": "1"}},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			uc := NewUserController(testCase.mockRepo)

			// Execute request and received recorded and decoded response
			w := executeRequest(
				[]gin.HandlerFunc{
					mocks.MockValidateAuthMiddleware(&auth.CustomClaims{
						UserID: testCase.userID,
						Role:   testCase.userRole,
					}),
					uc.Delete,
				},
				testCase.requestOpts,
			)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			if testCase.mockRepo != nil {
				testCase.mockRepo.AssertExpectations(t)
			}
		})
	}

	// Setup mock repo, assign functions to call and handlers to call in order
	mockUserRepo := setupUserRepoMock("Delete", []any{"1"}, []any{nil})
	uc := NewUserController(mockUserRepo)
	testHandlers := []gin.HandlerFunc{
		mocks.MockValidateAuthMiddleware(&auth.CustomClaims{
			UserID: "1",
			Role:   "administrator",
		}),
		uc.Delete,
	}
	w := executeRequest(testHandlers, requestOpts{params: map[string]string{"id": "1"}})

	// Assert responses
	assert.Equal(t, http.StatusOK, w.Code)
}
