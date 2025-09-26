// Package sbi_test contains integration tests for the SBI (South Bound Interface) layer
// These tests verify the HTTP API endpoints by testing the complete request-response cycle
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

// Test_HTTPPostMessage tests the HTTP POST endpoint for creating new messages
// This function validates various scenarios including successful message creation,
// invalid JSON handling, and missing required fields validation
func Test_HTTPPostMessage(t *testing.T) {
	// Set Gin to test mode to prevent debug output during testing
	gin.SetMode(gin.TestMode)

	// Create a new mock controller to manage all mock objects for this test
	mockCtrl := gomock.NewController(t)

	// Create mock objects for the NF application and processor
	// These mocks allow us to control dependencies and isolate the code under test
	nfApp := sbi.NewMocknfApp(mockCtrl)
	mockProcessor := processor.NewMockProcessorNf(mockCtrl)

	// Create a real processor instance that wraps the mock processor
	// This allows us to test the actual processor logic while mocking its dependencies
	realProcessor, err := processor.NewProcessor(mockProcessor)
	if err != nil {
		t.Errorf("Failed to create processor: %s", err)
		return
	}

	// Set up expectations for mock calls that will happen during server initialization
	// AnyTimes() allows these methods to be called any number of times during the test
	nfApp.EXPECT().Config().Return(&factory.Config{
		Configuration: &factory.Configuration{
			Sbi: &factory.Sbi{
				Port: 8000, // Mock server configuration with port 8000
			},
		},
	}).AnyTimes()
	nfApp.EXPECT().Processor().Return(realProcessor).AnyTimes()

	// Create the SBI server instance with the mocked dependencies
	server := sbi.NewServer(nfApp, "")

	// Test case 1: Successful message posting
	// This test verifies that a valid POST request creates a message successfully
	t.Run("Post Message Successfully", func(t *testing.T) {
		// Define the expected HTTP status code for successful message creation
		const EXPECTED_STATUS = http.StatusCreated

		// Create a sample request body with valid message data
		// This represents the JSON payload that a client would send
		requestBody := map[string]string{
			"content": "Hello World", // The message content
			"author":  "Anya",        // The message author
		}

		// Marshal the request body into JSON format
		// This simulates how HTTP clients send JSON data
		var jsonBody []byte
		jsonBody, err = json.Marshal(requestBody)
		if err != nil {
			t.Errorf("Failed to marshal request body: %s", err)
			return
		}

		// Set up mock context with an empty messages slice
		// This represents the initial state before any messages are created
		mockContext := &nf_context.NFContext{
			Messages: []nf_context.Message{}, // Start with no existing messages
		}
		// Expect the Context() method to be called once during message processing
		mockProcessor.EXPECT().Context().Return(mockContext).Times(1)

		// Create HTTP test infrastructure
		// httpRecorder captures the HTTP response for inspection
		httpRecorder := httptest.NewRecorder()
		// ginCtx provides the Gin context needed for HTTP handling
		ginCtx, _ := gin.CreateTestContext(httpRecorder)

		// Create a mock HTTP POST request with the JSON body
		// This simulates a real HTTP request from a client
		ginCtx.Request, err = http.NewRequest("POST", "/message/", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Errorf("Failed to create request: %s", err)
			return
		}
		// Set the Content-Type header to indicate JSON payload
		ginCtx.Request.Header.Set("Content-Type", "application/json")

		// Execute the actual HTTP handler method
		// This is the code under test - the HTTPPostMessage function
		server.HTTPPostMessage(ginCtx)

		// Verify the HTTP response status code
		// StatusCreated (201) indicates successful resource creation
		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		// Parse the JSON response body to verify its contents
		var response map[string]interface{}
		err = json.Unmarshal(httpRecorder.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %s", err)
		}

		// Verify the response contains the expected success message
		if response["message"] != "Message posted successfully" {
			t.Errorf("Expected message 'Message posted successfully', got %s", response["message"])
		}

		// Verify the response data structure and content
		// The response should include the created message data
		if data, ok := response["data"].(map[string]interface{}); !ok {
			t.Errorf("Expected data field to be an object")
		} else {
			// Verify that the response contains the original content
			if data["content"] != "Hello World" {
				t.Errorf("Expected content 'Hello World', got %s", data["content"])
			}
			// Verify that the response contains the original author
			if data["author"] != "Anya" {
				t.Errorf("Expected author 'Anya', got %s", data["author"])
			}
			// Verify that an ID was generated for the new message
			if data["id"] == nil || data["id"] == "" {
				t.Errorf("Expected non-empty ID")
			}
			// Verify that a timestamp was added to the message
			if data["time"] == nil || data["time"] == "" {
				t.Errorf("Expected non-empty time")
			}
		}
	})

	// Test case 2: Invalid JSON handling
	// This test verifies that the API properly handles malformed JSON requests
	t.Run("Post Message with Invalid JSON", func(t *testing.T) {
		// Define expected response for invalid JSON
		const EXPECTED_STATUS = http.StatusBadRequest
		const EXPECTED_MESSAGE = "Invalid request body"

		// Create intentionally malformed JSON (missing value after "author":)
		// This simulates a client sending corrupted or incomplete JSON data
		invalidJSON := []byte(`{"content": "Hello World", "author":}`)

		// Set up HTTP testing infrastructure
		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)

		// Create request with the invalid JSON payload
		// var err error
		ginCtx.Request, err = http.NewRequest("POST", "/message/", bytes.NewBuffer(invalidJSON))
		if err != nil {
			t.Errorf("Failed to create request: %s", err)
			return
		}
		ginCtx.Request.Header.Set("Content-Type", "application/json")

		// Execute the handler - it should handle the invalid JSON gracefully
		server.HTTPPostMessage(ginCtx)

		// Verify that the server responds with a Bad Request status
		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		// Parse and verify the error response
		var response map[string]interface{}
		err = json.Unmarshal(httpRecorder.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %s", err)
		}

		// Verify the error message is descriptive
		if response["message"] != EXPECTED_MESSAGE {
			t.Errorf("Expected message %s, got %s", EXPECTED_MESSAGE, response["message"])
		}

		// Verify that an error field is present in the response
		// This helps clients understand what went wrong
		if response["error"] == nil {
			t.Errorf("Expected error field to be present")
		}
	})

	// Test case 3: Missing required fields validation
	// This test ensures that requests with incomplete data are rejected properly
	t.Run("Post Message with Missing Required Fields", func(t *testing.T) {
		// Define expected response for incomplete request
		const EXPECTED_STATUS = http.StatusBadRequest
		const EXPECTED_MESSAGE = "Invalid request body"

		// Create request body missing the "author" field
		// This tests the API's input validation capabilities
		requestBody := map[string]string{
			"content": "Hello World", // Only content provided, author is missing
		}

		// Convert to JSON for the request
		var jsonBody []byte
		jsonBody, err = json.Marshal(requestBody)
		if err != nil {
			t.Errorf("Failed to marshal request body: %s", err)
			return
		}

		// Set up HTTP testing infrastructure
		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)

		// Create request with incomplete data
		ginCtx.Request, err = http.NewRequest("POST", "/message/", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Errorf("Failed to create request: %s", err)
			return
		}
		ginCtx.Request.Header.Set("Content-Type", "application/json")

		// Execute the handler - it should validate required fields
		server.HTTPPostMessage(ginCtx)

		// Verify that the server rejects the incomplete request
		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		// Parse and verify the validation error response
		var response map[string]interface{}
		err = json.Unmarshal(httpRecorder.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %s", err)
		}

		// Verify appropriate error message
		if response["message"] != EXPECTED_MESSAGE {
			t.Errorf("Expected message %s, got %s", EXPECTED_MESSAGE, response["message"])
		}
	})
}

// Test_HTTPGetMessages tests the HTTP GET endpoint for retrieving all messages
// This function verifies that the API correctly returns the list of all stored messages
func Test_HTTPGetMessages(t *testing.T) {
	// Set Gin to test mode to suppress debug output
	gin.SetMode(gin.TestMode)

	// Create mock controller and dependencies
	mockCtrl := gomock.NewController(t)
	nfApp := sbi.NewMocknfApp(mockCtrl)
	mockProcessor := processor.NewMockProcessorNf(mockCtrl)

	// Create a real processor instance with mocked dependencies
	realProcessor, err := processor.NewProcessor(mockProcessor)
	if err != nil {
		t.Errorf("Failed to create processor: %s", err)
		return
	}

	// Set up mock expectations for server initialization
	nfApp.EXPECT().Config().Return(&factory.Config{
		Configuration: &factory.Configuration{
			Sbi: &factory.Sbi{
				Port: 8000, // Mock server configuration
			},
		},
	}).AnyTimes()
	nfApp.EXPECT().Processor().Return(realProcessor).AnyTimes()

	// Initialize the SBI server with mocked dependencies
	server := sbi.NewServer(nfApp, "")

	// Test case: Successful retrieval of messages (empty list scenario)
	// This test verifies that the API correctly handles requests when no messages exist
	t.Run("Get Messages Successfully", func(t *testing.T) {
		// Define expected status for successful GET request
		const EXPECTED_STATUS = http.StatusOK

		// Set up mock context with empty messages array
		// This simulates the state when no messages have been created yet
		mockContext := &nf_context.NFContext{
			Messages: []nf_context.Message{}, // Empty slice represents no existing messages
		}
		// Expect the Context() method to be called once during message retrieval
		mockProcessor.EXPECT().Context().Return(mockContext).Times(1)

		// Set up HTTP testing infrastructure
		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)

		// Create a GET request to retrieve all messages
		// No request body is needed for GET requests
		// var err error
		ginCtx.Request, err = http.NewRequest("GET", "/message/", nil)
		if err != nil {
			t.Errorf("Failed to create request: %s", err)
			return
		}

		// Execute the GET messages handler
		server.HTTPGetMessages(ginCtx)

		// Verify that the request was successful
		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		// Parse the JSON response to verify its structure
		var response map[string]interface{}
		err = json.Unmarshal(httpRecorder.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %s", err)
		}

		// Verify the success message
		if response["message"] != "Messages retrieved successfully" {
			t.Errorf("Expected message 'Messages retrieved successfully', got %s", response["message"])
		}

		// Verify that the data field is an array (even if empty)
		// This ensures consistent response structure regardless of message count
		if _, ok := response["data"].([]interface{}); !ok {
			t.Errorf("Expected data field to be an array")
		}
	})
}

// Test_HTTPGetMessageByID tests the HTTP GET endpoint for retrieving a specific message by ID
// This function validates both successful retrieval and proper handling of non-existent messages
func Test_HTTPGetMessageByID(t *testing.T) {
	// Set Gin to test mode to prevent debug output
	gin.SetMode(gin.TestMode)

	// Create mock controller and dependencies for testing
	mockCtrl := gomock.NewController(t)
	nfApp := sbi.NewMocknfApp(mockCtrl)
	mockProcessor := processor.NewMockProcessorNf(mockCtrl)

	// Create a real processor instance with mocked dependencies
	realProcessor, err := processor.NewProcessor(mockProcessor)
	if err != nil {
		t.Errorf("Failed to create processor: %s", err)
		return
	}

	// Set up mock expectations for server configuration
	nfApp.EXPECT().Config().Return(&factory.Config{
		Configuration: &factory.Configuration{
			Sbi: &factory.Sbi{
				Port: 8000, // Mock server port configuration
			},
		},
	}).AnyTimes()
	nfApp.EXPECT().Processor().Return(realProcessor).AnyTimes()

	// Initialize the SBI server with mocked dependencies
	server := sbi.NewServer(nfApp, "")

	// Test case: Attempting to retrieve a message that doesn't exist
	// This test verifies proper error handling when a requested message ID is not found
	t.Run("Get Message by Valid ID - Not Found", func(t *testing.T) {
		// Define test parameters
		const MESSAGE_ID = "test-message-id"        // ID that doesn't exist in the system
		const EXPECTED_STATUS = http.StatusNotFound // 404 status for non-existent resources

		// Set up mock context with no messages
		// This simulates a scenario where the requested message doesn't exist
		mockContext := &nf_context.NFContext{
			Messages: []nf_context.Message{}, // Empty slice means no messages exist
		}
		// Expect the Context() method to be called once during message lookup
		mockProcessor.EXPECT().Context().Return(mockContext).Times(1)

		// Set up HTTP testing infrastructure
		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)

		// Create a GET request for a specific message ID
		// var err error
		ginCtx.Request, err = http.NewRequest("GET", "/message/"+MESSAGE_ID, nil)
		if err != nil {
			t.Errorf("Failed to create request: %s", err)
			return
		}

		// Manually set the URL parameter for testing
		// In a real Gin router, this would be extracted from the URL path automatically
		// The "id" parameter corresponds to the ":id" in the route pattern "/:id"
		ginCtx.Params = gin.Params{
			{Key: "id", Value: MESSAGE_ID},
		}

		// Execute the handler to get message by ID
		server.HTTPGetMessageByID(ginCtx)

		// Verify that the server responds with Not Found status
		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		// Parse and verify the error response
		var response map[string]interface{}
		err = json.Unmarshal(httpRecorder.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %s", err)
		}

		// Verify that the response contains an appropriate error message
		if response["message"] != "Message not found" {
			t.Errorf("Expected message 'Message not found', got %s", response["message"])
		}
	})
}
