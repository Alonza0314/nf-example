package sbi

import (
	"net/http"

	"github.com/Alonza0314/nf-example/internal/logger"
	"github.com/gin-gonic/gin"
)

func (s *Server) getNetworkRoute() []Route {
	return []Route{
		{
			Name:    "Hello Network!",
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: func(c *gin.Context) {
				c.JSON(http.StatusOK, "Network Device API")
			},
		},
		{
			Name:    "Get Network Device",
			Method:  http.MethodGet,
			Pattern: "/device/:Name",
			APIFunc: s.HTTPGetNetworkDevice,
		},
		{
			Name:    "Add Network Device",
			Method:  http.MethodPost,
			Pattern: "/device",
			APIFunc: s.HTTPPostNetworkDevice,
		},
	}
}

func (s *Server) HTTPGetNetworkDevice(c *gin.Context) {
	logger.SBILog.Infof("In HTTPGetNetworkDevice")

	targetName := c.Param("Name")
	if targetName == "" {
		c.String(http.StatusBadRequest, "No device name provided")
		return
	}

	s.Processor().FindNetworkDevice(c, targetName)
}

func (s *Server) HTTPPostNetworkDevice(c *gin.Context) {
	logger.SBILog.Infof("In HTTPPostNetworkDevice")

	type RequestBody struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}
	var body RequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if body.Name == "" || body.Type == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and type are required"})
		return
	}

	s.Processor().AddNetworkDevice(c, body.Name, body.Type)
}
