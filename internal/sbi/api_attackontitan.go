package sbi

import (
	"net/http"

	"github.com/Alonza0314/nf-example/internal/logger"
	"github.com/gin-gonic/gin"
)

func (s *Server) getAttackOnTitanRoute() []Route {
	return []Route{
		{
			Name:    "Hello Attack on Titan!",
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: func(c *gin.Context) {
				c.JSON(http.StatusOK, "Hello Attack on Titan!")
			},
			// Use
			// curl -X GET http://127.0.0.163:8000/attackontitan/ -w "\n"
		},
		{
			Name:    "Attack on Titan Character",
			Method:  http.MethodGet,
			Pattern: "/character/:Name",
			APIFunc: s.HTTPSearchAttackOnTitanCharacter,
			// Use
			// curl -X GET http://127.0.0.163:8000/attackontitan/character/Eren -w "\n"
		},
		{
			Name:    "Add Attack on Titan Character",
			Method:  http.MethodPost,
			Pattern: "/character",
			APIFunc: s.HTTPAddAttackOnTitanCharacter,
			// Use
			// curl -X POST http://127.0.0.163:8000/attackontitan/character -H "Content-Type: application/json" -d '{"name":"Zeke","lastName":"Yeager"}' -w "\n"
		},
	}
}

func (s *Server) HTTPSearchAttackOnTitanCharacter(c *gin.Context) {
	logger.SBILog.Infof("In HTTPSearchAttackOnTitanCharacter")
	targetName := c.Param("Name")
	if targetName == "" {
		c.String(http.StatusBadRequest, "No name provided")
		return
	}
	s.Processor().FindAttackOnTitanCharacterName(c, targetName)
}

func (s *Server) HTTPAddAttackOnTitanCharacter(c *gin.Context) {
	logger.SBILog.Infof("In HTTPAddAttackOnTitanCharacter")
	s.Processor().AddAttackOnTitanCharacter(c)
}