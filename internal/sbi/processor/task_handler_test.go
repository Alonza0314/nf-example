package processor_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/Alonza0314/nf-example/internal/context"
	"github.com/Alonza0314/nf-example/internal/sbi/processor"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_TaskHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockNfApp := processor.NewMockProcessorNf(mockCtrl)

	nfContext := &context.NFContext{
		Tasks:      make([]context.Task, 0),
		TaskMutex:  sync.RWMutex{},
		NextTaskID: 0,
	}
	mockNfApp.EXPECT().Context().Return(nfContext).AnyTimes()

	proc, err := processor.NewProcessor(mockNfApp)
	assert.NoError(t, err)

	t.Run("Create Task", func(t *testing.T) {
		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)

		taskData := map[string]string{"name": "Test Task"}
		var body []byte
		body, err = json.Marshal(taskData)
		assert.NoError(t, err)

		ginCtx.Request, err = http.NewRequest(http.MethodPost, "/task/tasks", bytes.NewReader(body))
		assert.NoError(t, err)

		ginCtx.Request.Header.Set("Content-Type", "application/json")

		proc.CreateNewTask(ginCtx)

		assert.Equal(t, http.StatusCreated, httpRecorder.Code)

		var createdTask context.Task
		err = json.Unmarshal(httpRecorder.Body.Bytes(), &createdTask)
		assert.NoError(t, err)
		assert.Equal(t, "Test Task", createdTask.Name)
		assert.Equal(t, 1, createdTask.ID)
	})

	t.Run("Get All Tasks", func(t *testing.T) {
		httpRecorder := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(httpRecorder)

		ginCtx.Request, err = http.NewRequest(http.MethodGet, "/task/tasks", nil)
		assert.NoError(t, err)

		proc.GetAllTasks(ginCtx)

		assert.Equal(t, http.StatusOK, httpRecorder.Code)

		var tasks []context.Task
		err = json.Unmarshal(httpRecorder.Body.Bytes(), &tasks)
		assert.NoError(t, err)
		assert.Len(t, tasks, 1)
		assert.Equal(t, "Test Task", tasks[0].Name)
	})
}
