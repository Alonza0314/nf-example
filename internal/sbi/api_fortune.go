package sbi

import (
	"net/http"

	"github.com/Alonza0314/nf-example/internal/logger"
	"github.com/Alonza0314/nf-example/internal/sbi/processor"
	"github.com/gin-gonic/gin"
)

func (s *Server) getFortuneRoute() []Route {
	return []Route{
		{
			Name:    "Get Today's Fortune",
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: s.HTTPGetFortune,
			// Use
			// curl -X GET http://127.0.0.163:8000/fortune/ -w "\n"
		},
		{
			Name:    "Add a new Fortune",
			Method:  http.MethodPost,
			Pattern: "/",
			APIFunc: s.HTTPPostFortune,
			// Use
			// curl -X POST http://127.0.0.163:8000/fortune/ \
			//   -H "Content-Type: application/json" \
			//   -d '{"fortune":"New fortune text"}' -w "\n"
		},
	}
}

func (s *Server) HTTPGetFortune(c *gin.Context) {
	logger.SBILog.Infof("In HTTPGetFortune")

	s.Processor().GetFortune(c)
}

func (s *Server) HTTPPostFortune(c *gin.Context) {
	logger.SBILog.Infof("In HTTPPostFortune")

	var req processor.PostFortuneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.SBILog.Errorf("Invalid request body: %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	s.Processor().PostFortune(c, req)
}

func (s *Server) GetFortuneRoute() []Route {
	return s.getFortuneRoute()
}
