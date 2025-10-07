package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getOnePunchManRoute() []Route {
	return []Route{
		{
			Name:    "I'm bald, and stronger",
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: func(c *gin.Context) {
				c.JSON(http.StatusOK, "I'm bald, and stronger")
			},
			// Use
			// curl -X GET http://127.0.0.163:8000/onepunchman/ -w "\n"
		},
		{
			Name:    "Echo POST",
			Method:  http.MethodPost,
			Pattern: "/echo",
			APIFunc: func(c *gin.Context) {
				// 定義一個 map 接收 JSON
				var requestData map[string]interface{}
				if err := c.ShouldBindJSON(&requestData); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"message": "Received your data!",
					"data":    requestData,
				})
			},
			// curl -X POST http://127.0.0.163:8000/onepunchman/echo -H "Content-Type: application/json" -d '{"name":"Saitama","power":100}'
		},
	}
}
