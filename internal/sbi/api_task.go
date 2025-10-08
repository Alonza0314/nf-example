package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) HTTPCreateNewTask(c *gin.Context) {
	s.Processor().CreateNewTask(c)
}

func (s *Server) HTTPGetAllTasks(c *gin.Context) {
	s.Processor().GetAllTasks(c)
}

func (s *Server) getTaskRoute() []Route {
	return []Route{
		{
			Name:    "Get All Tasks",
			Method:  http.MethodGet,
			Pattern: "/tasks",
			APIFunc: s.HTTPGetAllTasks,
		},
		{
			Name:    "Create New Task",
			Method:  http.MethodPost,
			Pattern: "/tasks",
			APIFunc: s.HTTPCreateNewTask,
		},
	}
}
