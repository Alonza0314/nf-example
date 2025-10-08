package sbi

import (
	"net/http"

	"github.com/Alonza0314/nf-example/internal/logger"
	"github.com/Alonza0314/nf-example/internal/sbi/processor"
	"github.com/gin-gonic/gin"
)

func (s *Server) getMessageRoute() []Route {
	return []Route{
		{
			Name:    "Post Message",
			Method:  http.MethodPost,
			Pattern: "/",
			APIFunc: s.HTTPPostMessage,
			// Use
			// curl -X POST http://127.0.0.163:8000/msg/ \
			//   -H "Content-Type: application/json" \
			//   -d '{"content":"Hello World","author":"Anya"}' -w "\n"
		},
		{
			Name:    "Get All Messages",
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: s.HTTPGetMessages,
			// Use
			// curl -X GET http://127.0.0.163:8000/msg/ -w "\n"
		},
		{
			Name:    "Get Message by ID",
			Method:  http.MethodGet,
			Pattern: "/:id", // ":" is used for dynamic parameter
			APIFunc: s.HTTPGetMessageByID,
			// Use
			// curl -X GET http://127.0.0.163:8000/msg/{message-id} -w "\n"
		},
	}
}

func (s *Server) HTTPPostMessage(c *gin.Context) {
	logger.SBILog.Infof("In HTTPPostMessage")

	var req processor.PostMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.SBILog.Errorf("Invalid request body: %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}
	// if req has redundant fields, it will be ignored
	s.Processor().PostMessage(c, req)
}

func (s *Server) HTTPGetMessages(c *gin.Context) {
	logger.SBILog.Infof("In HTTPGetMessages")

	s.Processor().GetMessages(c)
}

func (s *Server) HTTPGetMessageByID(c *gin.Context) {
	logger.SBILog.Infof("In HTTPGetMessageByID")

	messageID := c.Param("id")

	s.Processor().GetMessageByID(c, messageID)
}
