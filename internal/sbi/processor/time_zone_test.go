package processor_test

import (
	"net/http/httptest"
	"testing"

	nf_context "github.com/Alonza0314/nf-example/internal/context"
	"github.com/Alonza0314/nf-example/internal/sbi/processor"
	"github.com/gin-gonic/gin"
	gomock "go.uber.org/mock/gomock"
)

func setupTimeZoneProcessor(t *testing.T) (*processor.Processor, *processor.MockProcessorNf) {
	gin.SetMode(gin.TestMode)

	mockCtrl := gomock.NewController(t)
	processorNf := processor.NewMockProcessorNf(mockCtrl)
	proc, err := processor.NewProcessor(processorNf)
	if err != nil {
		t.Fatalf("Failed to create processor: %s", err)
	}
	return proc, processorNf
}

func Test_HandleGetTimeZone(t *testing.T) {
	proc, processorNf := setupTimeZoneProcessor(t)

	t.Run("Get TimeZone for Existing City", func(t *testing.T) {
		const INPUT_CITY = "Taipei"
		const EXPECTED_STATUS = 200
		const EXPECTED_BODY = "UTC+8"

		processorNf.EXPECT().Context().Return(&nf_context.NFContext{
			TimeZoneData: map[string]string{
				"Taipei": "UTC+8",
				"Tokyo":  "UTC+9",
			},
		})

		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)
		proc.HandleGetTimeZone(ginCtx, INPUT_CITY)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		if httpRecorder.Body.String() != EXPECTED_BODY {
			t.Errorf("Expected body %s, got %s", EXPECTED_BODY, httpRecorder.Body.String())
		}
	})

	t.Run("Get TimeZone for Non-Existing City", func(t *testing.T) {
		const INPUT_CITY = "Unknown"
		const EXPECTED_STATUS = 404
		const EXPECTED_BODY = "[Unknown] not found"

		processorNf.EXPECT().Context().Return(&nf_context.NFContext{
			TimeZoneData: map[string]string{
				"Taipei": "UTC+8",
				"Tokyo":  "UTC+9",
			},
		})

		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)
		proc.HandleGetTimeZone(ginCtx, INPUT_CITY)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		if httpRecorder.Body.String() != EXPECTED_BODY {
			t.Errorf("Expected body %s, got %s", EXPECTED_BODY, httpRecorder.Body.String())
		}
	})
}

func Test_HandleAddNewCityTimeZone(t *testing.T) {
	proc, processorNf := setupTimeZoneProcessor(t)

	t.Run("Add New City Successfully", func(t *testing.T) {
		const EXPECTED_STATUS = 200
		const EXPECTED_BODY = "Time zone of Chicago is set to UTC-6"

		req := processor.TimeZoneRequest{
			City:     "Chicago",
			TimeZone: "UTC-6",
		}

		timeZoneData := map[string]string{
			"Taipei": "UTC+8",
			"Tokyo":  "UTC+9",
		}

		processorNf.EXPECT().Context().Return(&nf_context.NFContext{
			TimeZoneData: timeZoneData,
		}).Times(2) // Called twice: once to check if city exists, once to set timezone

		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)
		proc.HandleAddNewCityTimeZone(ginCtx, req)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		if httpRecorder.Body.String() != EXPECTED_BODY {
			t.Errorf("Expected body %s, got %s", EXPECTED_BODY, httpRecorder.Body.String())
		}

		// Verify city was added to the data
		if timeZoneData["Chicago"] != "UTC-6" {
			t.Errorf("Expected Chicago to be added with timezone UTC-6, but was not found")
		}
	})

	t.Run("Add Existing City - Conflict", func(t *testing.T) {
		const EXPECTED_STATUS = 409
		const EXPECTED_BODY = "City 'Taipei' already exists"

		req := processor.TimeZoneRequest{
			City:     "Taipei",
			TimeZone: "UTC+7",
		}

		processorNf.EXPECT().Context().Return(&nf_context.NFContext{
			TimeZoneData: map[string]string{
				"Taipei": "UTC+8",
				"Tokyo":  "UTC+9",
			},
		}).Times(1) // Called once to check if city exists (conflict case)

		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)
		proc.HandleAddNewCityTimeZone(ginCtx, req)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		if httpRecorder.Body.String() != EXPECTED_BODY {
			t.Errorf("Expected body %s, got %s", EXPECTED_BODY, httpRecorder.Body.String())
		}
	})
}

func Test_HandleResetCityTimeZone(t *testing.T) {
	proc, processorNf := setupTimeZoneProcessor(t)

	t.Run("Reset Existing City TimeZone", func(t *testing.T) {
		const INPUT_CITY = "Taipei"
		const NEW_TIMEZONE = "UTC+7"
		const EXPECTED_STATUS = 200
		const EXPECTED_BODY = "Time zone of Taipei is reset to UTC+7"

		timeZoneData := map[string]string{
			"Taipei": "UTC+8",
			"Tokyo":  "UTC+9",
		}

		processorNf.EXPECT().Context().Return(&nf_context.NFContext{
			TimeZoneData: timeZoneData,
		}).Times(2) // Called twice: once to check if city exists, once to update timezone

		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)
		proc.HandleResetCityTimeZone(ginCtx, INPUT_CITY, NEW_TIMEZONE)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		if httpRecorder.Body.String() != EXPECTED_BODY {
			t.Errorf("Expected body %s, got %s", EXPECTED_BODY, httpRecorder.Body.String())
		}

		// Verify timezone was updated
		if timeZoneData["Taipei"] != NEW_TIMEZONE {
			t.Errorf("Expected Taipei timezone to be updated to %s, got %s", NEW_TIMEZONE, timeZoneData["Taipei"])
		}
	})

	t.Run("Reset Non-Existing City TimeZone", func(t *testing.T) {
		const INPUT_CITY = "Unknown"
		const NEW_TIMEZONE = "UTC+0"
		const EXPECTED_STATUS = 404
		const EXPECTED_BODY = "City 'Unknown' not found"

		processorNf.EXPECT().Context().Return(&nf_context.NFContext{
			TimeZoneData: map[string]string{
				"Taipei": "UTC+8",
				"Tokyo":  "UTC+9",
			},
		}).Times(1) // Called once to check if city exists (not found case)

		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)
		proc.HandleResetCityTimeZone(ginCtx, INPUT_CITY, NEW_TIMEZONE)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		if httpRecorder.Body.String() != EXPECTED_BODY {
			t.Errorf("Expected body %s, got %s", EXPECTED_BODY, httpRecorder.Body.String())
		}
	})
}

func Test_HandleDeleteCityTimeZone(t *testing.T) {
	proc, processorNf := setupTimeZoneProcessor(t)

	t.Run("Delete Existing City", func(t *testing.T) {
		const INPUT_CITY = "Tokyo"
		const EXPECTED_STATUS = 200
		const EXPECTED_BODY = "City 'Tokyo' has been removed"

		timeZoneData := map[string]string{
			"Taipei": "UTC+8",
			"Tokyo":  "UTC+9",
		}

		processorNf.EXPECT().Context().Return(&nf_context.NFContext{
			TimeZoneData: timeZoneData,
		}).Times(2) // Called twice: once to check if city exists, once to delete from map

		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)
		proc.HandleDeleteCityTimeZone(ginCtx, INPUT_CITY)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		if httpRecorder.Body.String() != EXPECTED_BODY {
			t.Errorf("Expected body %s, got %s", EXPECTED_BODY, httpRecorder.Body.String())
		}

		// Verify city was deleted
		if _, exists := timeZoneData["Tokyo"]; exists {
			t.Errorf("Expected Tokyo to be deleted, but it still exists")
		}
	})

	t.Run("Delete Non-Existing City", func(t *testing.T) {
		const INPUT_CITY = "Unknown"
		const EXPECTED_STATUS = 404
		const EXPECTED_BODY = "City 'Unknown' not found"

		processorNf.EXPECT().Context().Return(&nf_context.NFContext{
			TimeZoneData: map[string]string{
				"Taipei": "UTC+8",
				"Tokyo":  "UTC+9",
			},
		}).Times(1) // Called once to check if city exists (not found case)

		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)
		proc.HandleDeleteCityTimeZone(ginCtx, INPUT_CITY)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		if httpRecorder.Body.String() != EXPECTED_BODY {
			t.Errorf("Expected body %s, got %s", EXPECTED_BODY, httpRecorder.Body.String())
		}
	})
}
