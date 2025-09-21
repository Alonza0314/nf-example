package sbi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	nf_context "github.com/Alonza0314/nf-example/internal/context"
	"github.com/Alonza0314/nf-example/internal/sbi"
	"github.com/Alonza0314/nf-example/internal/sbi/processor"
	"github.com/Alonza0314/nf-example/pkg/factory"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

func Test_HTTPPostMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockCtrl := gomock.NewController(t)
	nfApp := sbi.NewMocknfApp(mockCtrl)
	mockProcessor := processor.NewMockProcessorNf(mockCtrl)

	// Create a real processor with mock dependencies
	realProcessor, err := processor.NewProcessor(mockProcessor)
	if err != nil {
		t.Errorf("Failed to create processor: %s", err)
		return
	}

	nfApp.EXPECT().Config().Return(&factory.Config{
		Configuration: &factory.Configuration{
			Sbi: &factory.Sbi{
				Port: 8000,
			},
		},
	}).AnyTimes()
	nfApp.EXPECT().Processor().Return(realProcessor).AnyTimes()

	server := sbi.NewServer(nfApp, "")

	t.Run("Post Message Successfully", func(t *testing.T) {
		const EXPECTED_STATUS = http.StatusCreated

		requestBody := map[string]string{
			"content": "Hello World",
			"author":  "Anya",
		}

		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			t.Errorf("Failed to marshal request body: %s", err)
			return
		}

		// Mock context with initial empty messages
		mockContext := &nf_context.NFContext{
			Messages: []nf_context.Message{},
		}
		mockProcessor.EXPECT().Context().Return(mockContext).Times(1)

		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)

		ginCtx.Request, err = http.NewRequest("POST", "/message/", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Errorf("Failed to create request: %s", err)
			return
		}
		ginCtx.Request.Header.Set("Content-Type", "application/json")

		server.HTTPPostMessage(ginCtx)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		var response map[string]interface{}
		err = json.Unmarshal(httpRecorder.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %s", err)
		}

		if response["message"] != "Message posted successfully" {
			t.Errorf("Expected message 'Message posted successfully', got %s", response["message"])
		}

		// Check if data field exists and has correct structure
		if data, ok := response["data"].(map[string]interface{}); !ok {
			t.Errorf("Expected data field to be an object")
		} else {
			if data["content"] != "Hello World" {
				t.Errorf("Expected content 'Hello World', got %s", data["content"])
			}
			if data["author"] != "Anya" {
				t.Errorf("Expected author 'Anya', got %s", data["author"])
			}
			if data["id"] == nil || data["id"] == "" {
				t.Errorf("Expected non-empty ID")
			}
			if data["time"] == nil || data["time"] == "" {
				t.Errorf("Expected non-empty time")
			}
		}
	})

	t.Run("Post Message with Invalid JSON", func(t *testing.T) {
		const EXPECTED_STATUS = http.StatusBadRequest
		const EXPECTED_MESSAGE = "Invalid request body"

		invalidJSON := []byte(`{"content": "Hello World", "author":}`)

		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)

		var err error
		ginCtx.Request, err = http.NewRequest("POST", "/message/", bytes.NewBuffer(invalidJSON))
		if err != nil {
			t.Errorf("Failed to create request: %s", err)
			return
		}
		ginCtx.Request.Header.Set("Content-Type", "application/json")

		server.HTTPPostMessage(ginCtx)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		var response map[string]interface{}
		err = json.Unmarshal(httpRecorder.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %s", err)
		}

		if response["message"] != EXPECTED_MESSAGE {
			t.Errorf("Expected message %s, got %s", EXPECTED_MESSAGE, response["message"])
		}

		if response["error"] == nil {
			t.Errorf("Expected error field to be present")
		}
	})

	t.Run("Post Message with Missing Required Fields", func(t *testing.T) {
		const EXPECTED_STATUS = http.StatusBadRequest
		const EXPECTED_MESSAGE = "Invalid request body"

		// Missing author field
		requestBody := map[string]string{
			"content": "Hello World",
		}

		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			t.Errorf("Failed to marshal request body: %s", err)
			return
		}

		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)

		ginCtx.Request, err = http.NewRequest("POST", "/message/", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Errorf("Failed to create request: %s", err)
			return
		}
		ginCtx.Request.Header.Set("Content-Type", "application/json")

		server.HTTPPostMessage(ginCtx)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		var response map[string]interface{}
		err = json.Unmarshal(httpRecorder.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %s", err)
		}

		if response["message"] != EXPECTED_MESSAGE {
			t.Errorf("Expected message %s, got %s", EXPECTED_MESSAGE, response["message"])
		}
	})
}

func Test_HTTPGetMessages(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockCtrl := gomock.NewController(t)
	nfApp := sbi.NewMocknfApp(mockCtrl)
	mockProcessor := processor.NewMockProcessorNf(mockCtrl)

	// Create a real processor with mock dependencies
	realProcessor, err := processor.NewProcessor(mockProcessor)
	if err != nil {
		t.Errorf("Failed to create processor: %s", err)
		return
	}

	nfApp.EXPECT().Config().Return(&factory.Config{
		Configuration: &factory.Configuration{
			Sbi: &factory.Sbi{
				Port: 8000,
			},
		},
	}).AnyTimes()
	nfApp.EXPECT().Processor().Return(realProcessor).AnyTimes()

	server := sbi.NewServer(nfApp, "")

	t.Run("Get Messages Successfully", func(t *testing.T) {
		const EXPECTED_STATUS = http.StatusOK

		// Mock context with empty messages
		mockContext := &nf_context.NFContext{
			Messages: []nf_context.Message{},
		}
		mockProcessor.EXPECT().Context().Return(mockContext).Times(1)

		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)

		var err error
		ginCtx.Request, err = http.NewRequest("GET", "/message/", nil)
		if err != nil {
			t.Errorf("Failed to create request: %s", err)
			return
		}

		server.HTTPGetMessages(ginCtx)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		var response map[string]interface{}
		err = json.Unmarshal(httpRecorder.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %s", err)
		}

		if response["message"] != "Messages retrieved successfully" {
			t.Errorf("Expected message 'Messages retrieved successfully', got %s", response["message"])
		}

		// Check if data field exists (should be an array)
		if _, ok := response["data"].([]interface{}); !ok {
			t.Errorf("Expected data field to be an array")
		}
	})
}

func Test_HTTPGetMessageByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockCtrl := gomock.NewController(t)
	nfApp := sbi.NewMocknfApp(mockCtrl)
	mockProcessor := processor.NewMockProcessorNf(mockCtrl)

	// Create a real processor with mock dependencies
	realProcessor, err := processor.NewProcessor(mockProcessor)
	if err != nil {
		t.Errorf("Failed to create processor: %s", err)
		return
	}

	nfApp.EXPECT().Config().Return(&factory.Config{
		Configuration: &factory.Configuration{
			Sbi: &factory.Sbi{
				Port: 8000,
			},
		},
	}).AnyTimes()
	nfApp.EXPECT().Processor().Return(realProcessor).AnyTimes()

	server := sbi.NewServer(nfApp, "")

	t.Run("Get Message by Valid ID - Not Found", func(t *testing.T) {
		const MESSAGE_ID = "test-message-id"
		const EXPECTED_STATUS = http.StatusNotFound

		// Mock context with no messages
		mockContext := &nf_context.NFContext{
			Messages: []nf_context.Message{},
		}
		mockProcessor.EXPECT().Context().Return(mockContext).Times(1)

		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)

		var err error
		ginCtx.Request, err = http.NewRequest("GET", "/message/"+MESSAGE_ID, nil)
		if err != nil {
			t.Errorf("Failed to create request: %s", err)
			return
		}

		// Set the param manually for testing
		ginCtx.Params = gin.Params{
			{Key: "id", Value: MESSAGE_ID},
		}

		server.HTTPGetMessageByID(ginCtx)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		var response map[string]interface{}
		err = json.Unmarshal(httpRecorder.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %s", err)
		}

		if response["message"] != "Message not found" {
			t.Errorf("Expected message 'Message not found', got %s", response["message"])
		}
	})

	t.Run("Get Message with Empty ID", func(t *testing.T) {
		const EXPECTED_STATUS = http.StatusBadRequest
		const EXPECTED_MESSAGE = "Message ID is required"
		const EXPECTED_ERROR = "No message ID provided in URL path"

		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)

		var err error
		ginCtx.Request, err = http.NewRequest("GET", "/message/", nil)
		if err != nil {
			t.Errorf("Failed to create request: %s", err)
			return
		}

		// Set empty param to simulate missing ID
		ginCtx.Params = gin.Params{
			{Key: "id", Value: ""},
		}

		server.HTTPGetMessageByID(ginCtx)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		var response map[string]interface{}
		err = json.Unmarshal(httpRecorder.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %s", err)
		}

		if response["message"] != EXPECTED_MESSAGE {
			t.Errorf("Expected message %s, got %s", EXPECTED_MESSAGE, response["message"])
		}

		if response["error"] != EXPECTED_ERROR {
			t.Errorf("Expected error %s, got %s", EXPECTED_ERROR, response["error"])
		}
	})
}
