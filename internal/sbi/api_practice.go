package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PracticeRequest struct {
	Name string `json:"name"`
}

func (s *Server) getPracticeRoute() []Route {
	return []Route{
		{
			Name:    "GetPractice",
			Method:  "GET",
			Pattern: "/list",
			APIFunc: GetPractice,
		},
		{
			Name:    "CreatePractice",
			Method:  "POST",
			Pattern: "/create",
			APIFunc: CreatePractice,
		},
	}
}

func GetPractice(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Practice API is working",
	})
}

func CreatePractice(c *gin.Context) {
	var body PracticeRequest

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task created",
		"name":    body.Name,
	})
}
