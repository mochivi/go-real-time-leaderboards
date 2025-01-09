package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/mochivi/go-real-time-leaderboards/internal/mocks"
)

// This package is not for testign the utils.go file
// but for holding utils used for other test files

var ErrRepoOperation = errors.New("failed db operation")

// Functions to prepare the repository mocks for each test case

func setupLeaderboardRepoMock(funcName string, args, returns []any) *mocks.MockLeaderboardsRepo {
	mockRepo := mocks.MockLeaderboardsRepo{}
	mockRepo.On(funcName, args...).Return(returns...)
	return &mockRepo
}

func setupUserRepoMock(funcName string, args, returns []any) *mocks.MockUserRepo {
	mockUserRepo := mocks.MockUserRepo{}
	mockUserRepo.On(funcName, args...).Return(returns...)
	return &mockUserRepo
}

func setupRedisServiceMock(funcName string, args, returns []any) *mocks.MockRedisService {
	mockRedisService := mocks.MockRedisService{}
	mockRedisService.On(funcName, args...).Return(returns...)
	return &mockRedisService
}

// Functions to help making the request to the handler below

type requestOpts struct {
	headers map[string]string
	body    any
	params  map[string]string
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

func executeRequest(testHandlers []gin.HandlerFunc, requestOpts ...requestOpts) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)

	// Create recorder and gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// The HTTP method is irrelevant since it is only used for routing the request to the appropriate handler
	// But we are directly calling the handler with the context in this case
	c.Request, _ = http.NewRequest("GET", "/", nil)

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
					Key:   k,
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
