package processor

import (
	"net/http"
	"time"

	nf_context "github.com/Alonza0314/nf-example/internal/context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostMessageRequest struct {
	Content string `json:"content" binding:"required"`
	Author  string `json:"author" binding:"required"`
}

type PostMessageResponse struct {
	Message string             `json:"message"`
	Data    nf_context.Message `json:"data"`
}

type GetMessagesResponse struct {
	Message string               `json:"message"`
	Data    []nf_context.Message `json:"data"`
}

func (p *Processor) PostMessage(c *gin.Context, req PostMessageRequest) {
	newMessage := nf_context.Message{
		ID:      uuid.New().String(),
		Content: req.Content,
		Author:  req.Author,
		Time:    time.Now().Format(time.RFC3339),
	}

	// add message to context
	ctx := p.Context()
	ctx.Messages = append(ctx.Messages, newMessage)

	// return success response
	response := PostMessageResponse{
		Message: "Message posted successfully",
		Data:    newMessage,
	}

	c.JSON(http.StatusCreated, response)
}

func (p *Processor) GetMessages(c *gin.Context) {
	ctx := p.Context()

	response := GetMessagesResponse{
		Message: "Messages retrieved successfully",
		Data:    ctx.Messages,
	}

	c.JSON(http.StatusOK, response)
}

func (p *Processor) GetMessageByID(c *gin.Context, messageID string) {
	ctx := p.Context()

	// find message with specified ID
	for _, message := range ctx.Messages {
		if message.ID == messageID {
			response := PostMessageResponse{
				Message: "Message found",
				Data:    message,
			}
			c.JSON(http.StatusOK, response)
			return
		}
	}

	// if message not found
	c.JSON(http.StatusNotFound, gin.H{
		"message": "Message not found",
		"error":   "No message found with the specified ID",
	})
}
