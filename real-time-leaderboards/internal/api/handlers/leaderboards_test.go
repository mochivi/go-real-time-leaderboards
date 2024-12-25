package handlers

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mochivi/go-real-time-leaderboards/internal/mocks"
	"github.com/mochivi/go-real-time-leaderboards/internal/models"
	"github.com/mochivi/go-real-time-leaderboards/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLeaderboardsGet(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		name           string
		mockRepo       *mocks.MockLeaderboardsRepo
		expectedStatus int
		requestOpts    requestOpts
	}{
		{
			name:           "get leaderboard",
			mockRepo:       setupLeaderboardRepoMock("Get", []any{"1"}, []any{&models.Leaderboard{ID: "1"}, nil}),
			expectedStatus: http.StatusOK,
			requestOpts: requestOpts{
				params: map[string]string{
					"id": "1",
				},
			},
		},
		{
			name:           "get leaderboard missing id",
			expectedStatus: http.StatusBadRequest,
			requestOpts: requestOpts{
				params: map[string]string{
					"id": "",
				},
			},
		},
		{
			name:           "get leaderboard db not found",
			mockRepo:       setupLeaderboardRepoMock("Get", []any{"1"}, []any{&models.Leaderboard{}, storage.ErrNotFound}),
			expectedStatus: http.StatusNotFound,
			requestOpts: requestOpts{
				params: map[string]string{
					"id": "1",
				},
			},
		},
		{
			name:           "get leaderboard db error",
			mockRepo:       setupLeaderboardRepoMock("Get", []any{"1"}, []any{&models.Leaderboard{}, ErrRepoOperation}),
			expectedStatus: http.StatusInternalServerError,
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
			uc := NewLeaderboardController(testCase.mockRepo)

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

func TestLeaderboardsGetEntries(t *testing.T) {

	// Setup test cases
	testCases := []struct {
		name           string
		mockRepo       *mocks.MockLeaderboardsRepo
		expectedStatus int
		requestOpts    requestOpts
	}{
		{
			name:           "get leaderboard entry",
			mockRepo:       setupLeaderboardRepoMock("GetEntries", []any{"1"}, []any{[]models.LeaderboardEntry{{ID: "1"}}, nil}),
			expectedStatus: http.StatusOK,
			requestOpts: requestOpts{
				params: map[string]string{
					"id": "1",
				},
			},
		},
		{
			name:           "get leaderboard entry missing id",
			expectedStatus: http.StatusBadRequest,
			requestOpts: requestOpts{
				params: map[string]string{
					"id": "",
				},
			},
		},
		{
			name:           "get leaderboard entry db error",
			mockRepo:       setupLeaderboardRepoMock("GetEntries", []any{"1"}, []any{[]models.LeaderboardEntry{}, ErrRepoOperation}),
			expectedStatus: http.StatusInternalServerError,
			requestOpts: requestOpts{
				params: map[string]string{
					"id": "1",
				},
			},
		},
		{
			name:           "get leaderboard entry db not found",
			mockRepo:       setupLeaderboardRepoMock("GetEntries", []any{"1"}, []any{[]models.LeaderboardEntry{}, storage.ErrNotFound}),
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
			uc := NewLeaderboardController(testCase.mockRepo)

			// Execute request and received recorded and decoded response
			w := executeRequest(
				[]gin.HandlerFunc{uc.GetEntries},
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

func TestLeaderboardsCreate(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		name           string
		mockRepo       *mocks.MockLeaderboardsRepo
		expectedStatus int
		requestOpts    requestOpts
	}{
		{
			name: "create leaderboard",
			mockRepo: setupLeaderboardRepoMock(
				"Create",
				[]any{mock.AnythingOfType("*models.LeaderboardRequest")},
				[]any{&models.Leaderboard{ID: "1"}, nil}),
			expectedStatus: http.StatusCreated,
			requestOpts: requestOpts{
				body: models.LeaderboardRequest{Name: "test-leaderboard"},
			},
		},
		{
			name: "create leaderboard db error",
			mockRepo: setupLeaderboardRepoMock(
				"Create",
				[]any{mock.AnythingOfType("*models.LeaderboardRequest")},
				[]any{&models.Leaderboard{ID: "1"}, errors.New("db error")}),
			expectedStatus: http.StatusInternalServerError,
			requestOpts: requestOpts{
				body: models.LeaderboardRequest{Name: "test-leaderboard"},
			},
		},
		{
			name: "create leaderboard db conflict",
			mockRepo: setupLeaderboardRepoMock(
				"Create",
				[]any{mock.AnythingOfType("*models.LeaderboardRequest")},
				[]any{&models.Leaderboard{}, storage.ErrConflict}),
			expectedStatus: http.StatusInternalServerError,
			requestOpts: requestOpts{
				body: models.LeaderboardRequest{Name: "test-leaderboard"},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			// Recreate controller with new mock on every testcase
			uc := NewLeaderboardController(testCase.mockRepo)

			// Execute request and received recorded and decoded response
			w := executeRequest(
				[]gin.HandlerFunc{uc.Create},
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

func TestLeaderboardsCreateEntry(t *testing.T) {

	exampleEntry := models.LeaderboardEntryRequest{
		LeaderboardID: "1",
		UserID:        "1",
		Score:         10,
	}

	// Setup test cases
	testCases := []struct {
		name           string
		mockRepo       *mocks.MockLeaderboardsRepo
		expectedStatus int
		requestOpts    requestOpts
	}{
		{
			name: "create leaderboard entry",
			mockRepo: setupLeaderboardRepoMock(
				"CreateEntry",
				[]any{mock.AnythingOfType("*models.LeaderboardEntryRequest")},
				[]any{&models.LeaderboardEntry{ID: "1"}, nil}),
			expectedStatus: http.StatusCreated,
			requestOpts:    requestOpts{body: exampleEntry},
		},
		{
			name: "create leaderboard entry db error",
			mockRepo: setupLeaderboardRepoMock(
				"CreateEntry",
				[]any{mock.AnythingOfType("*models.LeaderboardEntryRequest")},
				[]any{&models.LeaderboardEntry{ID: "1"}, errors.New("db error")}),
			expectedStatus: http.StatusInternalServerError,
			requestOpts:    requestOpts{body: exampleEntry},
		},
		{
			name: "create leaderboard entry db conflict",
			mockRepo: setupLeaderboardRepoMock(
				"CreateEntry",
				[]any{mock.AnythingOfType("*models.LeaderboardEntryRequest")},
				[]any{&models.LeaderboardEntry{}, storage.ErrConflict}),
			expectedStatus: http.StatusInternalServerError,
			requestOpts:    requestOpts{body: exampleEntry},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			// Recreate controller with new mock on every testcase
			uc := NewLeaderboardController(testCase.mockRepo)

			// Execute request and received recorded and decoded response
			w := executeRequest(
				[]gin.HandlerFunc{uc.CreateEntry},
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

func TestLeaderboardsUpdate(t *testing.T) {

	exampleLeaderboardUpdate := models.UpdateLeaderboardRequest{
		ID:   "1",
		Name: "test-updated-name",
	}

	// Setup test cases
	testCases := []struct {
		name           string
		mockRepo       *mocks.MockLeaderboardsRepo
		expectedStatus int
		requestOpts    requestOpts
	}{
		{
			name: "update leaderboard",
			mockRepo: setupLeaderboardRepoMock(
				"Update",
				[]any{mock.AnythingOfType("*models.UpdateLeaderboardRequest")},
				[]any{&models.Leaderboard{ID: "1"}, nil},
			),
			expectedStatus: http.StatusOK,
			requestOpts:    requestOpts{body: exampleLeaderboardUpdate},
		},
		{
			name: "update leaderboard db error",
			mockRepo: setupLeaderboardRepoMock(
				"Update",
				[]any{mock.AnythingOfType("*models.UpdateLeaderboardRequest")},
				[]any{&models.Leaderboard{ID: "1"}, errors.New("db error")},
			),
			expectedStatus: http.StatusInternalServerError,
			requestOpts:    requestOpts{body: exampleLeaderboardUpdate},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			// Recreate controller with new mock on every testcase
			uc := NewLeaderboardController(testCase.mockRepo)

			// Execute request and received recorded and decoded response
			w := executeRequest(
				[]gin.HandlerFunc{uc.Update},
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

func TestLeaderboardsDelete(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		name           string
		mockRepo       *mocks.MockLeaderboardsRepo
		expectedStatus int
		requestOpts    requestOpts
	}{
		{
			name:           "delete leaderboard",
			mockRepo:       setupLeaderboardRepoMock("Delete", []any{"1"}, []any{nil}),
			expectedStatus: http.StatusNoContent,
			requestOpts: requestOpts{
				params: map[string]string{
					"id": "1",
				},
			},
		},
		{
			name:           "delete leaderboard missing id",
			expectedStatus: http.StatusBadRequest,
			requestOpts: requestOpts{
				params: map[string]string{
					"id": "",
				},
			},
		},
		{
			name:           "delete leaderboard db not found",
			mockRepo:       setupLeaderboardRepoMock("Delete", []any{"1"}, []any{storage.ErrNotFound}),
			expectedStatus: http.StatusNotFound,
			requestOpts: requestOpts{
				params: map[string]string{
					"id": "1",
				},
			},
		},
		{
			name:           "delete leaderboard db error",
			mockRepo:       setupLeaderboardRepoMock("Delete", []any{"1"}, []any{ErrRepoOperation}),
			expectedStatus: http.StatusInternalServerError,
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
			uc := NewLeaderboardController(testCase.mockRepo)

			// Execute request and received recorded and decoded response
			w := executeRequest(
				[]gin.HandlerFunc{uc.Delete},
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
