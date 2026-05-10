package sbi

import (
	"net/http"

	"github.com/Alonza0314/nf-example/internal/logger"
	"github.com/gin-gonic/gin"
)

type Lab6MessageRequest struct {
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

func (s *Server) getLab6Route() []Route {
	return []Route{
		{
			Name:    "Lab6 Status",
			Method:  http.MethodGet,
			Pattern: "/status",
			APIFunc: s.HTTPGetLab6Status,
		},
		{
			Name:    "Lab6 Message",
			Method:  http.MethodPost,
			Pattern: "/message",
			APIFunc: s.HTTPCreateLab6Message,
		},
	}
}

func (s *Server) HTTPGetLab6Status(c *gin.Context) {
	logger.SBILog.Infof("In HTTPGetLab6Status")
	s.Processor().GetLab6Status(c)
}

func (s *Server) HTTPCreateLab6Message(c *gin.Context) {
	logger.SBILog.Infof("In HTTPCreateLab6Message")

	var request Lab6MessageRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid JSON body",
		})
		return
	}

	s.Processor().CreateLab6Message(c, request.Sender, request.Message)
}
