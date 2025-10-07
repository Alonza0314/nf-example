package sbi

import (
	"net/http"

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
		{
			Name:    "Post Attendence",
			Method:  http.MethodPost,
			Pattern: "/",
			APIFunc: s.PostAttendence,
		},
	}
}

func (s *Server) GetAttendence(c *gin.Context) {
	logger.SBILog.Infof("In HTTPGetAttendence")

	s.Processor().ReturnAttendence(c)
}

func (s *Server) PostAttendence(c *gin.Context) {
	logger.SBILog.Infof("In HTTPPostAttendence")

	targetName := c.Param("Name")
	s.Processor().PostAttendence(c, targetName)
}
