package sbi_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Alonza0314/nf-example/internal/sbi"
	"github.com/Alonza0314/nf-example/pkg/factory"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

func setupTimeZoneTestServer(t *testing.T) *sbi.Server {
	gin.SetMode(gin.TestMode)

	mockCtrl := gomock.NewController(t)
	nfApp := sbi.NewMocknfApp(mockCtrl)
	nfApp.EXPECT().Config().Return(&factory.Config{
		Configuration: &factory.Configuration{
			Sbi: &factory.Sbi{
				Port: 8000,
			},
		},
	}).AnyTimes()
	return sbi.NewServer(nfApp, "")
}

func createJSONRequest(
	t *testing.T,
	method, url string,
	jsonBody string,
	params gin.Params,
) (*httptest.ResponseRecorder, *gin.Context) {
	httpRecorder := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(httpRecorder)
	if params != nil {
		ginCtx.Params = params
	}

	var body io.Reader
	if jsonBody != "" {
		body = bytes.NewBufferString(jsonBody)
	}

	var err error
	ginCtx.Request, err = http.NewRequest(method, url, body)
	if err != nil {
		t.Errorf("Failed to create request: %s", err)
		return nil, nil
	}
	if jsonBody != "" {
		ginCtx.Request.Header.Set("Content-Type", "application/json")
	}

	return httpRecorder, ginCtx
}

func Test_HTTPGetTimeZoneByCity(t *testing.T) {
	server := setupTimeZoneTestServer(t)

	t.Run("No city provided", func(t *testing.T) {
		const EXPECTED_STATUS = http.StatusBadRequest
		const EXPECTED_BODY = "No city provided"

		httpRecorder, ginCtx := createJSONRequest(t, "GET", "/timezone/city/", "", nil)
		if ginCtx == nil {
			return
		}

		server.HTTPGetTimeZoneByCity(ginCtx)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		if httpRecorder.Body.String() != EXPECTED_BODY {
			t.Errorf("Expected body %s, got %s", EXPECTED_BODY, httpRecorder.Body.String())
		}
	})
}

func Test_HTTPAddNewCityTimeZone(t *testing.T) {
	server := setupTimeZoneTestServer(t)

	t.Run("Invalid JSON", func(t *testing.T) {
		const EXPECTED_STATUS = http.StatusBadRequest
		const EXPECTED_BODY = "Invalid JSON"

		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)

		invalidJSON := bytes.NewBufferString("{invalid json}")
		var err error
		ginCtx.Request, err = http.NewRequest("POST", "/timezone/city", invalidJSON)
		if err != nil {
			t.Errorf("Failed to create request: %s", err)
			return
		}
		ginCtx.Request.Header.Set("Content-Type", "application/json")

		server.HTTPAddNewCityTimeZone(ginCtx)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		if httpRecorder.Body.String() != EXPECTED_BODY {
			t.Errorf("Expected body %s, got %s", EXPECTED_BODY, httpRecorder.Body.String())
		}
	})

	t.Run("Missing required fields", func(t *testing.T) {
		const EXPECTED_STATUS = http.StatusBadRequest
		const EXPECTED_BODY = "City and TimeZone fields are required"

		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)

		incompleteJSON := bytes.NewBufferString(`{"City": "", "TimeZone": "UTC+8"}`)
		var err error
		ginCtx.Request, err = http.NewRequest("POST", "/timezone/city", incompleteJSON)
		if err != nil {
			t.Errorf("Failed to create request: %s", err)
			return
		}
		ginCtx.Request.Header.Set("Content-Type", "application/json")

		server.HTTPAddNewCityTimeZone(ginCtx)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		if httpRecorder.Body.String() != EXPECTED_BODY {
			t.Errorf("Expected body %s, got %s", EXPECTED_BODY, httpRecorder.Body.String())
		}
	})
}

func Test_HTTPResetCityTimeZone(t *testing.T) {
	server := setupTimeZoneTestServer(t)

	t.Run("No city provided", func(t *testing.T) {
		const EXPECTED_STATUS = http.StatusBadRequest
		const EXPECTED_BODY = "No city provided"

		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)

		validJSON := bytes.NewBufferString(`{"TimeZone": "UTC+8"}`)
		var err error
		ginCtx.Request, err = http.NewRequest("POST", "/timezone/city/", validJSON)
		if err != nil {
			t.Errorf("Failed to create request: %s", err)
			return
		}
		ginCtx.Request.Header.Set("Content-Type", "application/json")

		server.HTTPResetCityTimeZone(ginCtx)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		if httpRecorder.Body.String() != EXPECTED_BODY {
			t.Errorf("Expected body %s, got %s", EXPECTED_BODY, httpRecorder.Body.String())
		}
	})

	t.Run("Invalid JSON format", func(t *testing.T) {
		const EXPECTED_STATUS = http.StatusBadRequest
		const EXPECTED_BODY = "Invalid JSON format, expected object with TimeZone field"

		httpRecorder, ginCtx := createJSONRequest(t, "POST", "/timezone/city/Taipei",
			"{invalid json}", gin.Params{gin.Param{Key: "City", Value: "Taipei"}})
		if ginCtx == nil {
			return
		}

		server.HTTPResetCityTimeZone(ginCtx)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		if httpRecorder.Body.String() != EXPECTED_BODY {
			t.Errorf("Expected body %s, got %s", EXPECTED_BODY, httpRecorder.Body.String())
		}
	})

	t.Run("Missing TimeZone field", func(t *testing.T) {
		const EXPECTED_STATUS = http.StatusBadRequest
		const EXPECTED_BODY = "TimeZone field is required"

		httpRecorder, ginCtx := createJSONRequest(t, "POST", "/timezone/city/Taipei",
			`{"TimeZone": ""}`, gin.Params{gin.Param{Key: "City", Value: "Taipei"}})
		if ginCtx == nil {
			return
		}

		server.HTTPResetCityTimeZone(ginCtx)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		if httpRecorder.Body.String() != EXPECTED_BODY {
			t.Errorf("Expected body %s, got %s", EXPECTED_BODY, httpRecorder.Body.String())
		}
	})
}

func Test_HTTPDeleteCityTimeZone(t *testing.T) {
	server := setupTimeZoneTestServer(t)

	t.Run("No city provided", func(t *testing.T) {
		const EXPECTED_STATUS = http.StatusBadRequest
		const EXPECTED_BODY = "No city provided"

		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)

		var err error
		ginCtx.Request, err = http.NewRequest("DELETE", "/timezone/city/", nil)
		if err != nil {
			t.Errorf("Failed to create request: %s", err)
			return
		}

		server.HTTPDeleteCityTimeZone(ginCtx)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		if httpRecorder.Body.String() != EXPECTED_BODY {
			t.Errorf("Expected body %s, got %s", EXPECTED_BODY, httpRecorder.Body.String())
		}
	})
}
