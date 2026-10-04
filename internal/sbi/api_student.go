package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getStudentRoute() []Route {
	return []Route{
		{
			Name:    "Student GET",
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: s.HTTPGetStudent,
		},
		{
			Name:    "Student POST",
			Method:  http.MethodPost,
			Pattern: "/",
			APIFunc: s.HTTPPostStudent,
		},
	}
}

func (s *Server) HTTPGetStudent(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello Student",
	})
}

func (s *Server) HTTPPostStudent(c *gin.Context) {
	type RequestBody struct {
		Name string `json:"name"`
	}

	var body RequestBody

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hello " + body.Name,
	})
}
