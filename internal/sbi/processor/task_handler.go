package processor

import (
	"net/http"
	"sync/atomic"

	"github.com/Alonza0314/nf-example/internal/context"
	"github.com/gin-gonic/gin"
)

func (p *Processor) CreateNewTask(c *gin.Context) {
	var newTask context.Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// next ID
	newID := atomic.AddUint64(&p.Context().NextTaskID, 1)
	newTask.ID = int(newID)

	p.Context().TaskMutex.Lock()
	p.Context().Tasks = append(p.Context().Tasks, newTask)
	p.Context().TaskMutex.Unlock()

	c.JSON(http.StatusCreated, newTask)
}

func (p *Processor) GetAllTasks(c *gin.Context) {
	p.Context().TaskMutex.RLock()
	tasksCopy := make([]context.Task, len(p.Context().Tasks))
	copy(tasksCopy, p.Context().Tasks)
	p.Context().TaskMutex.RUnlock()
	c.JSON(http.StatusOK, tasksCopy)
}
