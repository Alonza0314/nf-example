package sbi

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getOnePieceRoute() []Route {
	return []Route{
		{
			Name:    "Hello Straw Hats",
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: s.HTTPOnePieceGreeting,
		},
		{
			Name:    "Recruit Straw Hat",
			Method:  http.MethodPost,
			Pattern: "/crew",
			APIFunc: s.HTTPOnePieceRecruit,
		},
	}
}

func (s *Server) HTTPOnePieceGreeting(c *gin.Context) {
	c.JSON(http.StatusOK, "Hello Straw Hat Pirates!")
}

func (s *Server) HTTPOnePieceRecruit(c *gin.Context) {
	var request struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	message := fmt.Sprintf("%s has joined the Straw Hat crew!", request.Name)
	c.JSON(http.StatusCreated, message)
}
