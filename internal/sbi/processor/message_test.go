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

func Test_PostMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockCtrl := gomock.NewController(t)
	processorNf := processor.NewMockProcessorNf(mockCtrl)
	p, err := processor.NewProcessor(processorNf)
	if err != nil {
		t.Errorf("Failed to create processor: %s", err)
		return
	}

	t.Run("Post Message Successfully", func(t *testing.T) {
		const EXPECTED_STATUS = 201
		const INPUT_CONTENT = "Hello World"
		const INPUT_AUTHOR = "Anya"
		const EXPECTED_MESSAGE = "Message posted successfully"

		// Mock context with initial empty messages
		mockContext := &nf_context.NFContext{
			Messages: []nf_context.Message{},
		}

		processorNf.EXPECT().Context().Return(mockContext).Times(1)

		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)

		req := processor.PostMessageRequest{
			Content: INPUT_CONTENT,
			Author:  INPUT_AUTHOR,
		}

		p.PostMessage(ginCtx, req)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		var response processor.PostMessageResponse
		err := json.Unmarshal(httpRecorder.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %s", err)
		}

		if response.Message != EXPECTED_MESSAGE {
			t.Errorf("Expected message %s, got %s", EXPECTED_MESSAGE, response.Message)
		}

		if response.Data.Content != INPUT_CONTENT {
			t.Errorf("Expected content %s, got %s", INPUT_CONTENT, response.Data.Content)
		}

		if response.Data.Author != INPUT_AUTHOR {
			t.Errorf("Expected author %s, got %s", INPUT_AUTHOR, response.Data.Author)
		}

		// Verify that ID is generated and not empty
		if response.Data.ID == "" {
			t.Errorf("Expected non-empty ID, got empty string")
		}

		// Verify that ID is a valid UUID
		_, err = uuid.Parse(response.Data.ID)
		if err != nil {
			t.Errorf("Expected valid UUID, got %s", response.Data.ID)
		}

		// Verify that time is set and in correct format
		if response.Data.Time == "" {
			t.Errorf("Expected non-empty time, got empty string")
		}

		_, err = time.Parse(time.RFC3339, response.Data.Time)
		if err != nil {
			t.Errorf("Expected time in RFC3339 format, got %s", response.Data.Time)
		}

		// Verify message was added to context
		if len(mockContext.Messages) != 1 {
			t.Errorf("Expected 1 message in context, got %d", len(mockContext.Messages))
		}
	})
}

func Test_GetMessages(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockCtrl := gomock.NewController(t)
	processorNf := processor.NewMockProcessorNf(mockCtrl)
	p, err := processor.NewProcessor(processorNf)
	if err != nil {
		t.Errorf("Failed to create processor: %s", err)
		return
	}

	t.Run("Get Messages Successfully - Empty List", func(t *testing.T) {
		const EXPECTED_STATUS = 200
		const EXPECTED_MESSAGE = "Messages retrieved successfully"

		mockContext := &nf_context.NFContext{
			Messages: []nf_context.Message{},
		}

		processorNf.EXPECT().Context().Return(mockContext).Times(1)

		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)

		p.GetMessages(ginCtx)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		var response processor.GetMessagesResponse
		err := json.Unmarshal(httpRecorder.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %s", err)
		}

		if response.Message != EXPECTED_MESSAGE {
			t.Errorf("Expected message %s, got %s", EXPECTED_MESSAGE, response.Message)
		}

		if len(response.Data) != 0 {
			t.Errorf("Expected 0 messages, got %d", len(response.Data))
		}
	})

	t.Run("Get Messages Successfully - With Data", func(t *testing.T) {
		const EXPECTED_STATUS = 200
		const EXPECTED_MESSAGE = "Messages retrieved successfully"

		testMessages := []nf_context.Message{
			{
				ID:      "test-id-1",
				Content: "Test message 1",
				Author:  "Author 1",
				Time:    "2023-01-01T12:00:00Z",
			},
			{
				ID:      "test-id-2",
				Content: "Test message 2",
				Author:  "Author 2",
				Time:    "2023-01-01T12:01:00Z",
			},
		}

		mockContext := &nf_context.NFContext{
			Messages: testMessages,
		}

		processorNf.EXPECT().Context().Return(mockContext).Times(1)

		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)

		p.GetMessages(ginCtx)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		var response processor.GetMessagesResponse
		err := json.Unmarshal(httpRecorder.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %s", err)
		}

		if response.Message != EXPECTED_MESSAGE {
			t.Errorf("Expected message %s, got %s", EXPECTED_MESSAGE, response.Message)
		}

		if len(response.Data) != 2 {
			t.Errorf("Expected 2 messages, got %d", len(response.Data))
		}

		if response.Data[0].ID != "test-id-1" {
			t.Errorf("Expected first message ID test-id-1, got %s", response.Data[0].ID)
		}

		if response.Data[1].Content != "Test message 2" {
			t.Errorf("Expected second message content 'Test message 2', got %s", response.Data[1].Content)
		}
	})
}

func Test_GetMessageByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockCtrl := gomock.NewController(t)
	processorNf := processor.NewMockProcessorNf(mockCtrl)
	p, err := processor.NewProcessor(processorNf)
	if err != nil {
		t.Errorf("Failed to create processor: %s", err)
		return
	}

	testMessages := []nf_context.Message{
		{
			ID:      "existing-id",
			Content: "Existing message",
			Author:  "Test Author",
			Time:    "2023-01-01T12:00:00Z",
		},
		{
			ID:      "another-id",
			Content: "Another message",
			Author:  "Another Author",
			Time:    "2023-01-01T12:01:00Z",
		},
	}

	t.Run("Find Message That Exists", func(t *testing.T) {
		const INPUT_ID = "existing-id"
		const EXPECTED_STATUS = 200
		const EXPECTED_MESSAGE = "Message found"

		mockContext := &nf_context.NFContext{
			Messages: testMessages,
		}

		processorNf.EXPECT().Context().Return(mockContext).Times(1)

		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)

		p.GetMessageByID(ginCtx, INPUT_ID)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		var response processor.PostMessageResponse
		err := json.Unmarshal(httpRecorder.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %s", err)
		}

		if response.Message != EXPECTED_MESSAGE {
			t.Errorf("Expected message %s, got %s", EXPECTED_MESSAGE, response.Message)
		}

		if response.Data.ID != INPUT_ID {
			t.Errorf("Expected ID %s, got %s", INPUT_ID, response.Data.ID)
		}

		if response.Data.Content != "Existing message" {
			t.Errorf("Expected content 'Existing message', got %s", response.Data.Content)
		}

		if response.Data.Author != "Test Author" {
			t.Errorf("Expected author 'Test Author', got %s", response.Data.Author)
		}
	})

	t.Run("Find Message That Does Not Exist", func(t *testing.T) {
		const INPUT_ID = "non-existing-id"
		const EXPECTED_STATUS = 404
		const EXPECTED_MESSAGE = "Message not found"
		const EXPECTED_ERROR = "No message found with the specified ID"

		mockContext := &nf_context.NFContext{
			Messages: testMessages,
		}

		processorNf.EXPECT().Context().Return(mockContext).Times(1)

		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)

		p.GetMessageByID(ginCtx, INPUT_ID)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		var response map[string]interface{}
		err := json.Unmarshal(httpRecorder.Body.Bytes(), &response)
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
