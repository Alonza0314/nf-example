package processor_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	nf_context "github.com/Alonza0314/nf-example/internal/context"
	"github.com/Alonza0314/nf-example/internal/sbi/processor"
	"github.com/gin-gonic/gin"
	gomock "go.uber.org/mock/gomock"
)

func Test_SearchDragonBallCharacter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockCtrl := gomock.NewController(t)
	processorNf := processor.NewMockProcessorNf(mockCtrl)
	processor, err := processor.NewProcessor(processorNf)
	if err != nil {
		t.Errorf("Failed to create processor: %s", err)
		return
	}

	t.Run("Find Character That Exists", func(t *testing.T) {
		const INPUT_NAME = "Goku"
		const EXPECTED_STATUS = http.StatusOK
		const EXPECTED_BODY = "Character: " + INPUT_NAME + ", Powerlevel: 7\n"
		processorNf.EXPECT().Context().Return(&nf_context.NFContext{
			DragonBallData: map[string]int32{
				"Goku": 7,
			},
		})

		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)
		processor.SearchDragonBallCharacter(ginCtx, INPUT_NAME)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		if httpRecorder.Body.String() != EXPECTED_BODY {
			t.Errorf("Expected body %s, got %s", EXPECTED_BODY, httpRecorder.Body.String())
		}
	})

	t.Run("Find Character That Does Not Exist", func(t *testing.T) {
		const INPUT_NAME = "Andy"
		const EXPECTED_STATUS = http.StatusNotFound
		const EXPECTED_BODY = "[" + INPUT_NAME + "] not found in Dragon Ball\n"

		processorNf.EXPECT().Context().Return(&nf_context.NFContext{
			DragonBallData: map[string]int32{
				"Goku": 7,
			},
		})

		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)
		processor.SearchDragonBallCharacter(ginCtx, INPUT_NAME)

		if httpRecorder.Code != EXPECTED_STATUS {
			t.Errorf("Expected status code %d, got %d", EXPECTED_STATUS, httpRecorder.Code)
		}

		if httpRecorder.Body.String() != EXPECTED_BODY {
			t.Errorf("Expected body %s, got %s", EXPECTED_BODY, httpRecorder.Body.String())
		}
	})
}

func Test_FightDragonBall(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockCtrl := gomock.NewController(t)
	processorNf := processor.NewMockProcessorNf(mockCtrl)
	processor, err := processor.NewProcessor(processorNf)
	if err != nil {
		t.Errorf("Failed to create processor: %s", err)
		return
	}

	processorNf.EXPECT().Context().Return(&nf_context.NFContext{
		DragonBallData: map[string]int32{
			"Goku":    7,
			"Vegeta":  6,
			"Krillin": 7,
		},
	}).AnyTimes()

	tests := []struct {
		name           string
		targetName1    string
		targetName2    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "targetName1 not found",
			targetName1:    "Andy",
			targetName2:    "Vegeta",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "[Andy] not found in Dragon Ball\n",
		},
		{
			name:           "targetName2 not found",
			targetName1:    "Goku",
			targetName2:    "Andy",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "[Andy] not found in Dragon Ball\n",
		},
		{
			name:           "Goku defeats Vegeta",
			targetName1:    "Goku",
			targetName2:    "Vegeta",
			expectedStatus: http.StatusOK,
			expectedBody:   "Goku defeats Vegeta\n",
		},
		{
			name:           "Vegeta defeats Krillin",
			targetName1:    "Vegeta",
			targetName2:    "Krillin",
			expectedStatus: http.StatusOK,
			expectedBody:   "Krillin defeats Vegeta\n",
		},
		{
			name:           "Tie",
			targetName1:    "Goku",
			targetName2:    "Krillin",
			expectedStatus: http.StatusOK,
			expectedBody:   "Goku ties with Krillin\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			httpRecorder := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(httpRecorder)
			processor.FightDragonBall(ginCtx, tc.targetName1, tc.targetName2)

			if httpRecorder.Code != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, httpRecorder.Code)
			}
			if httpRecorder.Body.String() != tc.expectedBody {
				t.Errorf("Expected body %q, got %q", tc.expectedBody, httpRecorder.Body.String())
			}
		})
	}
}

//nolint:dupl
func Test_AddDragonBallCharacter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockCtrl := gomock.NewController(t)
	processorNf := processor.NewMockProcessorNf(mockCtrl)

	processorNf.EXPECT().Context().Return(&nf_context.NFContext{
		DragonBallData: map[string]int32{
			"Goku": 7,
		},
	}).AnyTimes()
	processor, err := processor.NewProcessor(processorNf)
	if err != nil {
		t.Errorf("Failed to create processor: %s", err)
		return
	}

	t.Run("Add existing character", func(t *testing.T) {
		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)
		processor.AddDragonBallCharacter(ginCtx, "Goku", 7)

		if httpRecorder.Code != http.StatusConflict {
			t.Errorf("Expected status code %d, got %d", http.StatusConflict, httpRecorder.Code)
		}
		expectedBody := "Character Goku already exists with Powerlevel 7\n"
		if httpRecorder.Body.String() != expectedBody {
			t.Errorf("Expected body %q, got %q", expectedBody, httpRecorder.Body.String())
		}
	})

	t.Run("Add new character", func(t *testing.T) {
		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)
		processor.AddDragonBallCharacter(ginCtx, "Vegeta", 6)

		if httpRecorder.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, httpRecorder.Code)
		}
		expectedBody := "Add Character Vegeta with Powerlevel 6\n"
		if httpRecorder.Body.String() != expectedBody {
			t.Errorf("Expected body %q, got %q", expectedBody, httpRecorder.Body.String())
		}
	})
}

//nolint:dupl
func Test_UpdateDragonBallCharacter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockCtrl := gomock.NewController(t)
	processorNf := processor.NewMockProcessorNf(mockCtrl)

	processorNf.EXPECT().Context().Return(&nf_context.NFContext{
		DragonBallData: map[string]int32{
			"Goku": 7,
		},
	}).AnyTimes()
	processor, err := processor.NewProcessor(processorNf)
	if err != nil {
		t.Errorf("Failed to create processor: %s", err)
		return
	}

	t.Run("Update non-existing character", func(t *testing.T) {
		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)
		processor.UpdateDragonBallCharacter(ginCtx, "Vegeta", 6)

		if httpRecorder.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, httpRecorder.Code)
		}
		expectedBody := "Character Vegeta not found\n"
		if httpRecorder.Body.String() != expectedBody {
			t.Errorf("Expected body %q, got %q", expectedBody, httpRecorder.Body.String())
		}
	})

	t.Run("Update existing character", func(t *testing.T) {
		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)
		processor.UpdateDragonBallCharacter(ginCtx, "Goku", 10)

		if httpRecorder.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, httpRecorder.Code)
		}
		expectedBody := "Update Character Goku with Powerlevel 10\n"
		if httpRecorder.Body.String() != expectedBody {
			t.Errorf("Expected body %q, got %q", expectedBody, httpRecorder.Body.String())
		}
	})
}
