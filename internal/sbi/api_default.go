package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getDefaultRoute() []Route {
	return []Route{
		{
			Name:    "Hello free5GC!",
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: func(c *gin.Context) {
				c.JSON(http.StatusOK, "Hello free5GC!")
			},
			// Use
			// curl -X GET http://127.0.0.163:8000/default/ -w "\n"
		},
		{
			Name:    "Exercise GET",
			Method:  http.MethodGet,
			Pattern: "/exercise",
			APIFunc: func(c *gin.Context) {
				c.JSON(http.StatusOK, "This is get")
			},
		},
		{
			Name:    "Exercise POST",
			Method:  http.MethodPost,
			Pattern: "/exercise",
			APIFunc: func(c *gin.Context) {
				c.JSON(http.StatusOK, "This is post")
			},
		},
	}
}
