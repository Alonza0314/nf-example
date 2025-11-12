package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// GET /myapi/hello
func HelloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello from My API!"})
}

// POST /myapi/echo
func EchoHandler(c *gin.Context) {
	var data map[string]interface{}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"received": data})
}
