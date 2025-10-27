package sbi

import (
	"io"
	"net/http"

	"github.com/Alonza0314/nf-example/internal/logger"
	"github.com/gin-gonic/gin"
)

func (s *Server) getAttendanceRoute() []Route {
	return []Route{
		{
			Name:    "Get Attendance",
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: s.GetAttendance,
		},
		// curl -X GET http://127.0.0.163:8000/attendance/ -w "\n"
		{
			Name:    "Post Attendance",
			Method:  http.MethodPost,
			Pattern: "/",
			APIFunc: s.PostAttendance,
		},
		// curl -X POST http://127.0.0.163:8000/attendance/ -d 'John' -w "\n"
	}
}

func (s *Server) GetAttendance(c *gin.Context) {
	logger.SBILog.Infof("In HTTPGetAttandence")

	s.Processor().ReturnAttendance(c)
}

func (s *Server) PostAttendance(c *gin.Context) {
	logger.SBILog.Infof("In HTTPPostAttendance")

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	targetName := string(body)
	if targetName == "" {
		c.String(http.StatusBadRequest, "error: no name provided")
		return
	}
	s.Processor().PostAttendance(c, targetName)
}
