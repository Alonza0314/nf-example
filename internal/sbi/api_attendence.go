package sbi

import (
	"net/http"

	"io"

	"github.com/Alonza0314/nf-example/internal/logger"
	"github.com/gin-gonic/gin"
)

func (s *Server) getAttendenceRoute() []Route {
	return []Route{
		{
			Name:    "Get Attendence",
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: s.GetAttendence,
		},
		// curl -X GET http://127.0.0.163:8000/attendence/ -w "\n"
		{
			Name:    "Post Attendence",
			Method:  http.MethodPost,
			Pattern: "/",
			APIFunc: s.PostAttendence,
		},
		// curl -X POST http://127.0.0.163:8000/attendence/ -d 'John' -w "\n"
	}
}

func (s *Server) GetAttendence(c *gin.Context) {
	logger.SBILog.Infof("In HTTPGetAttendence")

	s.Processor().ReturnAttendence(c)
}

func (s *Server) PostAttendence(c *gin.Context) {
	logger.SBILog.Infof("In HTTPPostAttendence")

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	targetName := string(body)
	if targetName == "" {
		c.String(http.StatusBadRequest, "no name provided")
		return
	}
	s.Processor().PostAttendence(c, targetName)
}
