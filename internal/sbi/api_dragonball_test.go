package sbi_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Alonza0314/nf-example/internal/sbi"
	"github.com/Alonza0314/nf-example/pkg/factory"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

// go test will automatically run any function with a name beginning with "Test"
//
//nolint:dupl
func Test_HTTPSearchDragonBallCharacter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockCtrl := gomock.NewController(t)
	// Create a mock
	nfApp := sbi.NewMocknfApp(mockCtrl)
	// Set up expected return values, and can be called any number of times
	nfApp.EXPECT().Config().Return(&factory.Config{
		Configuration: &factory.Configuration{
			Sbi: &factory.Sbi{
				Port: 8000,
			},
		},
	}).AnyTimes()
	server := sbi.NewServer(nfApp, "")

	t.Run("No name provided", func(t *testing.T) {
		const EXPECTED_STATUS = http.StatusBadRequest
		const EXPECTED_BODY = "No name provided"

		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)

		var err error
		ginCtx.Request, err = http.NewRequest("GET", "/dragonball", nil)
		if err != nil {
			t.Errorf("Failed to create request: %s", err)
			return
		}

		server.HTTPSearchDragonBallCharacter(ginCtx)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		if httpRecorder.Body.String() != EXPECTED_BODY {
			t.Errorf("Expected body %s, got %s", EXPECTED_BODY, httpRecorder.Body.String())
		}
	})
}

//nolint:dupl
func Test_HTTPDragonBallFight(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockCtrl := gomock.NewController(t)
	nfApp := sbi.NewMocknfApp(mockCtrl)
	nfApp.EXPECT().Config().Return(&factory.Config{
		Configuration: &factory.Configuration{
			Sbi: &factory.Sbi{Port: 8000},
		},
	}).AnyTimes()
	server := sbi.NewServer(nfApp, "")

	// define test cases
	tests := []struct {
		name           string
		jsonBody       string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "error",
			jsonBody:       `{`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "error",
		},
		{
			name:           "No name1 provided",
			jsonBody:       `{"name2": "Vegeta"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "No name1 provided",
		},
		{
			name:           "No name2 provided",
			jsonBody:       `{"name1":"Goku"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "No name2 provided",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			httpRecorder := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(httpRecorder)

			req, err := http.NewRequest("POST", "/dragonball/battle", bytes.NewBufferString(tc.jsonBody))
			if err != nil {
				t.Errorf("Failed to create request: %s", err)
				return
			}
			ginCtx.Request = req

			server.HTTPDragonBallFight(ginCtx)

			if httpRecorder.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d body=%s",
					tc.expectedStatus, httpRecorder.Code, httpRecorder.Body.String())
			}

			if httpRecorder.Body.String() != tc.expectedBody {
				t.Errorf("expected body %q, got %q",
					tc.expectedBody, httpRecorder.Body.String())
			}
		})
	}
}

//nolint:dupl
func Test_HTTPAddDragonBallCharacter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockCtrl := gomock.NewController(t)
	nfApp := sbi.NewMocknfApp(mockCtrl)
	nfApp.EXPECT().Config().Return(&factory.Config{
		Configuration: &factory.Configuration{
			Sbi: &factory.Sbi{Port: 8000},
		},
	}).AnyTimes()
	server := sbi.NewServer(nfApp, "")

	// define test cases
	tests := []struct {
		name           string
		jsonBody       string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "error",
			jsonBody:       `{`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "error",
		},
		{
			name:           "No name provided",
			jsonBody:       `{"powerLevel":100}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "No name provided",
		},
		{
			name:           "No Powerlevel provided",
			jsonBody:       `{"name":"Goku"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "No Powerlevel provided",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			httpRecorder := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(httpRecorder)

			req, err := http.NewRequest("POST", "/dragonball", bytes.NewBufferString(tc.jsonBody))
			if err != nil {
				t.Errorf("Failed to create request: %s", err)
				return
			}
			ginCtx.Request = req

			server.HTTPAddDragonBallCharacter(ginCtx)

			if httpRecorder.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d body=%s",
					tc.expectedStatus, httpRecorder.Code, httpRecorder.Body.String())
			}

			if httpRecorder.Body.String() != tc.expectedBody {
				t.Errorf("expected body %q, got %q",
					tc.expectedBody, httpRecorder.Body.String())
			}
		})
	}
}

func Test_HTTPUpdateDragonBallCharacter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockCtrl := gomock.NewController(t)
	nfApp := sbi.NewMocknfApp(mockCtrl)
	nfApp.EXPECT().Config().Return(&factory.Config{
		Configuration: &factory.Configuration{
			Sbi: &factory.Sbi{Port: 8000},
		},
	}).AnyTimes()
	server := sbi.NewServer(nfApp, "")

	// define test cases
	tests := []struct {
		name           string
		jsonBody       string
		expectedStatus int
		expectedBody   string
		url_param      string
	}{
		{
			name:           "No name provided",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "No name provided",
			url_param:      "",
		},
		{
			name:           "error",
			jsonBody:       `{`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "error",
			url_param:      "/Character",
		},
		{
			name:           "No Powerlevel provided",
			jsonBody:       `{}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "No Powerlevel provided",
			url_param:      "/Character",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			httpRecorder := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(httpRecorder)
			// check if :name is provided
			ginCtx.Params = gin.Params{gin.Param{Key: "name", Value: tc.url_param}}

			req, err := http.NewRequest("PUT", "/dragonball", bytes.NewBufferString(tc.jsonBody))
			if err != nil {
				t.Errorf("Failed to create request: %s", err)
				return
			}
			ginCtx.Request = req

			server.HTTPUpdateDragonBallCharacter(ginCtx)

			if httpRecorder.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d body=%s",
					tc.expectedStatus, httpRecorder.Code, httpRecorder.Body.String())
			}

			if httpRecorder.Body.String() != tc.expectedBody {
				t.Errorf("expected body %q, got %q",
					tc.expectedBody, httpRecorder.Body.String())
			}
		})
	}
}
