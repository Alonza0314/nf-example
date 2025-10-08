package processor_test

import (
	"net/http/httptest"
	"testing"

	nf_context "github.com/Alonza0314/nf-example/internal/context"
	"github.com/Alonza0314/nf-example/internal/sbi/processor"
	"github.com/gin-gonic/gin"
	gomock "go.uber.org/mock/gomock"
)

func Test_ReturnAttendence(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockCtrl := gomock.NewController(t)
	processorNf := processor.NewMockProcessorNf(mockCtrl)
	processor, err := processor.NewProcessor(processorNf)
	if err != nil {
		t.Errorf("Failed to create processor: %s", err)
		return
	}
	t.Run("No Attendence Recorded", func(t *testing.T) {
		const EXPECTED_STATUS = 200
		const EXPECTED_BODY = "No attendence recorded"
		processorNf.EXPECT().Context().Return(&nf_context.NFContext{
			AttendenceData: []string{},
		})

		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)
		processor.ReturnAttendence(ginCtx)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		if httpRecorder.Body.String() != EXPECTED_BODY {
			t.Errorf("Expected body %s, got %s", EXPECTED_BODY, httpRecorder.Body.String())
		}
	})

	t.Run("Some Attendence Recorded", func(t *testing.T) {
		const EXPECTED_STATUS = 200
		const EXPECTED_BODY = `{"Attendence":"Alice, Bob, Charlie"}`
		processorNf.EXPECT().Context().Return(&nf_context.NFContext{
			AttendenceData: []string{"Alice", "Bob", "Charlie"},
		})
		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)
		processor.ReturnAttendence(ginCtx)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		if httpRecorder.Body.String() != EXPECTED_BODY {
			t.Errorf("Expected body %s, got %s", EXPECTED_BODY, httpRecorder.Body.String())
		}

	})
}

func Test_PostAttendence(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockCtrl := gomock.NewController(t)
	processorNf := processor.NewMockProcessorNf(mockCtrl)
	processor, err := processor.NewProcessor(processorNf)

	if err != nil {
		t.Errorf("Failed to create processor: %s", err)
		return
	}

	t.Run("Post New Attendence", func(t *testing.T) {
		const INPUT_NAME = "David"
		const EXPECTED_STATUS = 200
		const EXPECTED_BODY = "Attendence recorded: "+ INPUT_NAME
		processorNf.EXPECT().Context().Return(&nf_context.NFContext{
			AttendenceData: []string{"Alice", "Bob", "Charlie"},
		})
		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)
		processor.PostAttendence(ginCtx, INPUT_NAME)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		if httpRecorder.Body.String() != EXPECTED_BODY {
			t.Errorf("Expected body %s, got %s", EXPECTED_BODY, httpRecorder.Body.String())
		}
	})

	t.Run("Post Duplicate Attendence", func(t *testing.T) {
		const INPUT_NAME = "Alice"
		const EXPECTED_STATUS = 409
		const EXPECTED_BODY = "Attendence already recorded: " + INPUT_NAME
		processorNf.EXPECT().Context().Return(&nf_context.NFContext{
			AttendenceData: []string{"Alice", "Bob", "Charlie"},
		})
		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)
		processor.PostAttendence(ginCtx, INPUT_NAME)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		if httpRecorder.Body.String() != EXPECTED_BODY {
			t.Errorf("Expected body %s, got %s", EXPECTED_BODY, httpRecorder.Body.String())
		}
	})
}
