// Package processor_test contains unit tests for the processor layer
// These tests focus on business logic validation, data processing, and core functionality
// The processor layer handles the actual message operations after HTTP parsing is complete
package processor_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	nf_context "github.com/Alonza0314/nf-example/internal/context"
	"github.com/Alonza0314/nf-example/internal/sbi/processor"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// Test_PostMessage tests the core business logic for creating new messages
// This test focuses on the processor layer functionality, including data validation,
// ID generation, timestamp creation, and proper storage in the context
func Test_PostMessage(t *testing.T) {
	// Set Gin to test mode to suppress debug output during testing
	gin.SetMode(gin.TestMode)

	// Create a mock controller to manage mock object lifecycles
	mockCtrl := gomock.NewController(t)

	// Create a mock ProcessorNf interface to isolate the processor logic being tested
	// This allows us to control the Context() method behavior without real dependencies
	processorNf := processor.NewMockProcessorNf(mockCtrl)

	// Create the actual processor instance with the mocked dependencies
	// This is the system under test - we test the real processor logic
	p, err := processor.NewProcessor(processorNf)
	if err != nil {
		t.Errorf("Failed to create processor: %s", err)
		return
	}

	// Test case: Successful message creation with all business logic validation
	// This test verifies ID generation, timestamp creation, and data persistence
	t.Run("Post Message Successfully", func(t *testing.T) {
		// Define expected test outcomes
		const EXPECTED_STATUS = 201                            // HTTP Created status
		const INPUT_CONTENT = "Hello World"                    // Test message content
		const INPUT_AUTHOR = "Anya"                            // Test message author
		const EXPECTED_MESSAGE = "Message posted successfully" // Expected success message

		// Set up mock context with empty messages array
		// This represents the initial state before any messages are created
		mockContext := &nf_context.NFContext{
			Messages: []nf_context.Message{}, // Start with empty message store
		}

		// Set expectation that Context() method will be called exactly once
		// This ensures the processor accesses the storage context as expected
		processorNf.EXPECT().Context().Return(mockContext).Times(1)

		// Set up HTTP testing infrastructure
		// httpRecorder captures the HTTP response for verification
		httpRecorder := httptest.NewRecorder()
		// ginCtx provides the Gin context needed for HTTP response handling
		ginCtx, _ := gin.CreateTestContext(httpRecorder)

		// Create the request object with test data
		// This represents the parsed and validated input from the HTTP layer
		req := processor.PostMessageRequest{
			Content: INPUT_CONTENT, // Message text content
			Author:  INPUT_AUTHOR,  // Message author information
		}

		// Execute the core business logic method
		// This is the main functionality being tested
		p.PostMessage(ginCtx, req)

		// Verify the HTTP response status code
		// 201 Created indicates successful resource creation
		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		// Parse the JSON response to verify its structure and content
		var response processor.PostMessageResponse
		err := json.Unmarshal(httpRecorder.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %s", err)
		}

		// Verify the response message indicates success
		if response.Message != EXPECTED_MESSAGE {
			t.Errorf("Expected message %s, got %s", EXPECTED_MESSAGE, response.Message)
		}

		// Verify that the original content is preserved in the response
		if response.Data.Content != INPUT_CONTENT {
			t.Errorf("Expected content %s, got %s", INPUT_CONTENT, response.Data.Content)
		}

		// Verify that the original author is preserved in the response
		if response.Data.Author != INPUT_AUTHOR {
			t.Errorf("Expected author %s, got %s", INPUT_AUTHOR, response.Data.Author)
		}

		// Verify that a unique ID was generated for the new message
		// IDs are essential for message identification and retrieval
		if response.Data.ID == "" {
			t.Errorf("Expected non-empty ID, got empty string")
		}

		// Verify that the generated ID is a valid UUID format
		// This ensures consistent ID formatting across the system
		_, err = uuid.Parse(response.Data.ID)
		if err != nil {
			t.Errorf("Expected valid UUID, got %s", response.Data.ID)
		}

		// Verify that a timestamp was automatically generated
		// Timestamps are crucial for message ordering and audit trails
		if response.Data.Time == "" {
			t.Errorf("Expected non-empty time, got empty string")
		}

		// Verify that the timestamp follows RFC3339 format (ISO 8601)
		// This ensures consistent time formatting and timezone handling
		_, err = time.Parse(time.RFC3339, response.Data.Time)
		if err != nil {
			t.Errorf("Expected time in RFC3339 format, got %s", response.Data.Time)
		}

		// Verify that the message was properly stored in the context
		// This confirms that the data persistence logic works correctly
		if len(mockContext.Messages) != 1 {
			t.Errorf("Expected 1 message in context, got %d", len(mockContext.Messages))
		}
	})
}

// Test_GetMessages tests the business logic for retrieving all stored messages
// This test verifies data retrieval functionality and proper response formatting
// It covers both empty state and populated state scenarios
func Test_GetMessages(t *testing.T) {
	// Set Gin to test mode to suppress debug output
	gin.SetMode(gin.TestMode)

	// Create mock controller for managing mock object lifecycles
	mockCtrl := gomock.NewController(t)

	// Create mock ProcessorNf interface to control dependencies
	processorNf := processor.NewMockProcessorNf(mockCtrl)

	// Create the actual processor instance being tested
	p, err := processor.NewProcessor(processorNf)
	if err != nil {
		t.Errorf("Failed to create processor: %s", err)
		return
	}

	// Test case 1: Retrieving messages when the system is empty
	// This test verifies proper handling of the initial state with no messages
	t.Run("Get Messages Successfully - Empty List", func(t *testing.T) {
		// Define expected test outcomes for empty state
		const EXPECTED_STATUS = 200                                // HTTP OK status
		const EXPECTED_MESSAGE = "Messages retrieved successfully" // Success message

		// Set up mock context representing an empty message store
		// This simulates the system state when no messages have been created
		mockContext := &nf_context.NFContext{
			Messages: []nf_context.Message{}, // Empty slice represents no stored messages
		}

		// Set expectation for Context() method call
		// The processor needs to access the context to retrieve messages
		processorNf.EXPECT().Context().Return(mockContext).Times(1)

		// Set up HTTP testing infrastructure
		httpRecorder := httptest.NewRecorder()           // Captures HTTP response
		ginCtx, _ := gin.CreateTestContext(httpRecorder) // Provides Gin context

		// Execute the business logic method for retrieving messages
		p.GetMessages(ginCtx)

		// Verify the HTTP response status code
		// 200 OK indicates successful data retrieval
		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		// Parse and verify the JSON response structure
		var response processor.GetMessagesResponse
		err := json.Unmarshal(httpRecorder.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %s", err)
		}

		// Verify the success message is present
		if response.Message != EXPECTED_MESSAGE {
			t.Errorf("Expected message %s, got %s", EXPECTED_MESSAGE, response.Message)
		}

		// Verify that the data array is empty but properly initialized
		// This ensures consistent response structure even with no data
		if len(response.Data) != 0 {
			t.Errorf("Expected 0 messages, got %d", len(response.Data))
		}
	})

	// Test case 2: Retrieving messages when data exists in the system
	// This test verifies proper data serialization and response formatting with actual data
	t.Run("Get Messages Successfully - With Data", func(t *testing.T) {
		// Define expected outcomes for populated state
		const EXPECTED_STATUS = 200                                // HTTP OK status
		const EXPECTED_MESSAGE = "Messages retrieved successfully" // Success message

		// Create test data representing existing messages in the system
		// This simulates the state after messages have been created and stored
		testMessages := []nf_context.Message{
			{
				ID:      "test-id-1",            // First message with unique ID
				Content: "Test message 1",       // First message content
				Author:  "Author 1",             // First message author
				Time:    "2023-01-01T12:00:00Z", // First message timestamp (ISO 8601)
			},
			{
				ID:      "test-id-2",            // Second message with unique ID
				Content: "Test message 2",       // Second message content
				Author:  "Author 2",             // Second message author
				Time:    "2023-01-01T12:01:00Z", // Second message timestamp (1 minute later)
			},
		}

		// Set up mock context with the test data
		// This represents a populated message store with existing messages
		mockContext := &nf_context.NFContext{
			Messages: testMessages, // Pre-populated message slice
		}

		// Set expectation for Context() method access
		processorNf.EXPECT().Context().Return(mockContext).Times(1)

		// Set up HTTP testing infrastructure
		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)

		// Execute the message retrieval business logic
		p.GetMessages(ginCtx)

		// Verify successful HTTP response
		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		// Parse and verify the JSON response with actual data
		var response processor.GetMessagesResponse
		err := json.Unmarshal(httpRecorder.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %s", err)
		}

		// Verify success message
		if response.Message != EXPECTED_MESSAGE {
			t.Errorf("Expected message %s, got %s", EXPECTED_MESSAGE, response.Message)
		}

		// Verify that the correct number of messages are returned
		// This ensures all stored messages are properly retrieved
		if len(response.Data) != 2 {
			t.Errorf("Expected 2 messages, got %d", len(response.Data))
		}

		// Verify the first message data integrity
		// This ensures the data is correctly serialized and ordered
		if response.Data[0].ID != "test-id-1" {
			t.Errorf("Expected first message ID test-id-1, got %s", response.Data[0].ID)
		}

		// Verify the second message content integrity
		// This confirms that all message fields are properly preserved
		if response.Data[1].Content != "Test message 2" {
			t.Errorf("Expected second message content 'Test message 2', got %s", response.Data[1].Content)
		}
	})
}

// Test_GetMessageByID tests the business logic for retrieving a specific message by its ID
// This test validates message lookup functionality, including both successful retrieval
// and proper error handling for non-existent messages
func Test_GetMessageByID(t *testing.T) {
	// Set Gin to test mode to suppress debug output
	gin.SetMode(gin.TestMode)

	// Create mock controller for managing mock objects
	mockCtrl := gomock.NewController(t)

	// Create mock ProcessorNf interface for dependency control
	processorNf := processor.NewMockProcessorNf(mockCtrl)

	// Create the processor instance being tested
	p, err := processor.NewProcessor(processorNf)
	if err != nil {
		t.Errorf("Failed to create processor: %s", err)
		return
	}

	// Set up test data representing existing messages in the system
	// This data will be used across multiple test scenarios
	testMessages := []nf_context.Message{
		{
			ID:      "existing-id",          // First test message with known ID
			Content: "Existing message",     // Content for verification
			Author:  "Test Author",          // Author for verification
			Time:    "2023-01-01T12:00:00Z", // Timestamp in ISO 8601 format
		},
		{
			ID:      "another-id",           // Second test message with different ID
			Content: "Another message",      // Different content for testing
			Author:  "Another Author",       // Different author for testing
			Time:    "2023-01-01T12:01:00Z", // Later timestamp
		},
	}

	// Test case 1: Successfully finding and retrieving an existing message
	// This test verifies the happy path where the requested message exists
	t.Run("Find Message That Exists", func(t *testing.T) {
		// Define test parameters for successful lookup
		const INPUT_ID = "existing-id"           // ID that exists in our test data
		const EXPECTED_STATUS = 200              // HTTP OK status for successful retrieval
		const EXPECTED_MESSAGE = "Message found" // Success message

		// Set up mock context with the test messages
		// This simulates a populated message store containing our target message
		mockContext := &nf_context.NFContext{
			Messages: testMessages, // Use the predefined test data
		}

		// Set expectation for Context() method call during message lookup
		processorNf.EXPECT().Context().Return(mockContext).Times(1)

		// Set up HTTP testing infrastructure
		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)

		// Execute the message lookup business logic
		// This tests the core ID-based search functionality
		p.GetMessageByID(ginCtx, INPUT_ID)

		// Verify successful HTTP response status
		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		// Parse the JSON response to verify message data integrity
		var response processor.PostMessageResponse
		err := json.Unmarshal(httpRecorder.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %s", err)
		}

		// Verify the success message
		if response.Message != EXPECTED_MESSAGE {
			t.Errorf("Expected message %s, got %s", EXPECTED_MESSAGE, response.Message)
		}

		// Verify that the correct message was retrieved by checking the ID
		if response.Data.ID != INPUT_ID {
			t.Errorf("Expected ID %s, got %s", INPUT_ID, response.Data.ID)
		}

		// Verify that the message content matches the expected data
		// This ensures the entire message object was correctly retrieved
		if response.Data.Content != "Existing message" {
			t.Errorf("Expected content 'Existing message', got %s", response.Data.Content)
		}

		// Verify that the message author matches the expected data
		// This confirms all message fields are properly preserved
		if response.Data.Author != "Test Author" {
			t.Errorf("Expected author 'Test Author', got %s", response.Data.Author)
		}
	})

	// Test case 2: Handling requests for non-existent messages
	// This test verifies proper error handling when the requested message ID doesn't exist
	t.Run("Find Message That Does Not Exist", func(t *testing.T) {
		// Define test parameters for failed lookup scenario
		const INPUT_ID = "non-existing-id"                              // ID that doesn't exist in test data
		const EXPECTED_STATUS = 404                                     // HTTP Not Found status
		const EXPECTED_MESSAGE = "Message not found"                    // Error message for user
		const EXPECTED_ERROR = "No message found with the specified ID" // Detailed error description

		// Set up mock context with the existing test messages
		// The requested ID will not be found in this dataset
		mockContext := &nf_context.NFContext{
			Messages: testMessages, // Same test data, but won't contain the requested ID
		}

		// Set expectation for Context() method call during failed lookup
		processorNf.EXPECT().Context().Return(mockContext).Times(1)

		// Set up HTTP testing infrastructure
		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)

		// Execute the message lookup with non-existent ID
		// This tests the error handling path in the business logic
		p.GetMessageByID(ginCtx, INPUT_ID)

		// Verify that the response indicates resource not found
		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		// Parse the error response to verify proper error formatting
		var response map[string]interface{}
		err := json.Unmarshal(httpRecorder.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %s", err)
		}

		// Verify that the error message is user-friendly and informative
		if response["message"] != EXPECTED_MESSAGE {
			t.Errorf("Expected message %s, got %s", EXPECTED_MESSAGE, response["message"])
		}

		// Verify that detailed error information is provided
		// This helps with debugging and provides context to API consumers
		if response["error"] != EXPECTED_ERROR {
			t.Errorf("Expected error %s, got %s", EXPECTED_ERROR, response["error"])
		}
	})
}
