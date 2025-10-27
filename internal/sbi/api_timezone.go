package sbi

import (
	"net/http"

	"github.com/Alonza0314/nf-example/internal/logger"
	"github.com/Alonza0314/nf-example/internal/sbi/processor"
	"github.com/gin-gonic/gin"
)

func (s *Server) getTimeZoneRoute() []Route {
	return []Route{
		{
			Name:    "Welcome to timezone service",
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: func(c *gin.Context) {
				c.String(http.StatusOK, "Welcome to time zone query service")
			},
			// Use
			// curl -X GET http://127.0.0.163:8000/timezone/ -w "\n"
		},
		{
			Name:    "Query city time zone",
			Method:  http.MethodGet,
			Pattern: "/city/:City",
			APIFunc: s.HTTPGetTimeZoneByCity,
			// Use
			// curl -X GET http://127.0.0.163:8000/timezone/city/Taipei -w "\n"
		},
		{
			Name:    "Add new city time zone",
			Method:  http.MethodPost,
			Pattern: "/city",
			APIFunc: s.HTTPAddNewCityTimeZone,
			// Use
			// curl -X POST http://127.0.0.163:8000/timezone/city -d '{"City": "Chicago", "TimeZone": "UTC-5"}' -w "\n"
		},
		{
			Name:    "Reset city time zone",
			Method:  http.MethodPost,
			Pattern: "/city/:City",
			APIFunc: s.HTTPResetCityTimeZone,
			// Use
			// curl -X POST http://127.0.0.163:8000/timezone/city/Chicago -d '{"TimeZone": "UTC-6"}' -w "\n"
		},
		{
			Name:    "Delete city time zone",
			Method:  http.MethodDelete,
			Pattern: "/city/:City",
			APIFunc: s.HTTPDeleteCityTimeZone,
			// Usage:
			// curl -X DELETE http://127.0.0.163:8000/timezone/city/Chicago -w "\n"
		},
	}
}

func (s *Server) HTTPGetTimeZoneByCity(c *gin.Context) {
	logger.SBILog.Infof("In HTTPGetTimeZoneByCity")

	city := c.Param("City")
	if city == "" {
		c.String(http.StatusBadRequest, "No city provided")
		return
	}
	s.Processor().HandleGetTimeZone(c, city)
}

func (s *Server) HTTPAddNewCityTimeZone(c *gin.Context) {
	logger.SBILog.Infof("In HTTPAddNewCityTimeZone")

	var req processor.TimeZoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.City == "" || req.TimeZone == "" {
		c.String(http.StatusBadRequest, "City and TimeZone fields are required")
		return
	}
	s.Processor().HandleAddNewCityTimeZone(c, req)
}

func (s *Server) HTTPResetCityTimeZone(c *gin.Context) {
	logger.SBILog.Infof("In HTTPResetCityTimeZone")

	city := c.Param("City")
	if city == "" {
		c.String(http.StatusBadRequest, "No city provided")
		return
	}

	var req struct {
		TZ string `json:"TimeZone"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "Invalid JSON format, expected object with TimeZone field")
		return
	}

	if req.TZ == "" {
		c.String(http.StatusBadRequest, "TimeZone field is required")
		return
	}

	s.Processor().HandleResetCityTimeZone(c, city, req.TZ)
}

func (s *Server) HTTPDeleteCityTimeZone(c *gin.Context) {
	logger.SBILog.Infof("In HTTPDeleteCityTimeZone")

	city := c.Param("City")
	if city == "" {
		c.String(http.StatusBadRequest, "No city provided")
		return
	}
	s.Processor().HandleDeleteCityTimeZone(c, city)
}
