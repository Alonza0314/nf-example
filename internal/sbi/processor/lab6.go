package processor

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (p *Processor) GetLab6Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"service":     "lab6",
		"status":      "running",
		"description": "Lab 6 API service for Git and GitHub practice",
	})
}

func (p *Processor) CreateLab6Message(c *gin.Context, sender string, message string) {
	sender = strings.TrimSpace(sender)
	message = strings.TrimSpace(message)

	if sender == "" || message == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "sender and message are required",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"service":    "lab6",
		"sender":     sender,
		"message":    message,
		"receivedAt": time.Now().UTC().Format(time.RFC3339),
	})
}
