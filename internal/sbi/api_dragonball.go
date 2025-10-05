package sbi

import (
	"net/http"

	"github.com/Alonza0314/nf-example/internal/logger"
	"github.com/gin-gonic/gin"
)

func (s *Server) getDragonBallRoute() []Route {
	return []Route{
		{
			Name:    "Hello Dragon Ball!",
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: func(c *gin.Context) {
				c.JSON(http.StatusOK, "Hello Dragon Ball!")
			},
			// Use
			// curl -X GET http://127.0.0.163:8000/dragonball/
		},
		{
			Name:    "Dragon Ball Search Character",
			Method:  http.MethodGet,
			Pattern: "/character/:name",
			APIFunc: s.HTTPSearchDragonBallCharacter,
			// Use
			// curl -X GET http://127.0.0.163:8000/dragonball/character/Goku
		},
		{
			Name:    "Dragon Ball Fight",
			Method:  http.MethodPost,
			Pattern: "/battle",
			APIFunc: s.HTTPDragonBallFight,
			// Use
			// curl -X POST "http://127.0.0.163:8000/dragonball/battle" -d '{"name1": "Goku", "name2": "Vegeta"}'
		},
		{
			Name:    "Add Dragon Ball Character",
			Method:  http.MethodPost,
			Pattern: "/character",
			APIFunc: s.HTTPAddDragonBallCharacter,
			// Use
			// curl -X POST "http://127.0.0.163:8000/dragonball/character" -d '{"Name": "Saitama", "Powerlevel": 10000}'
		},
		{
			Name:    "Update Dragon Ball Character's Powerlevel",
			Method:  http.MethodPut,
			Pattern: "/character/:name",
			APIFunc: s.HTTPUpdateDragonBallCharacter,
			// Use
			// curl -X PUT "http://127.0.0.163:8000/dragonball/character/Goku" -d '{"Powerlevel":  500}'
		},
	}
}

func (s *Server) HTTPSearchDragonBallCharacter(c *gin.Context) {
	logger.SBILog.Infof("In HTTPSearchDragonBallCharacter")
	targetName := c.Param("name")

	if targetName == "" {
		c.String(http.StatusBadRequest, "No name provided")
		return
	}

	s.Processor().SearchDragonBallCharacter(c, targetName)
}

func (s *Server) HTTPDragonBallFight(c *gin.Context) {
	logger.SBILog.Infof("In HTTPDragonBallFight")

	type RequestBody struct {
		TargetName1 string `json:"name1"`
		TargetName2 string `json:"name2"`
	}
	var requestbody RequestBody
	if err := c.ShouldBindBodyWithJSON(&requestbody); err != nil {
		c.String(http.StatusBadRequest, "error")
		return
	}
	if requestbody.TargetName1 == "" {
		c.String(http.StatusBadRequest, "No name1 provided")
		return
	}
	if requestbody.TargetName2 == "" {
		c.String(http.StatusBadRequest, "No name2 provided")
		return
	}

	s.Processor().FightDragonBall(c, requestbody.TargetName1, requestbody.TargetName2)
}

func (s *Server) HTTPAddDragonBallCharacter(c *gin.Context) {
	logger.SBILog.Infof("In HTTPAddDragonBallCharacter")

	type RequestBody struct {
		Name       string `json:"name"`
		PowerLevel *int32 `json:"powerLevel"`
	}
	var requestbody RequestBody
	if err := c.ShouldBindBodyWithJSON(&requestbody); err != nil {
		c.String(http.StatusBadRequest, "error")
		return
	}
	if requestbody.Name == "" {
		c.String(http.StatusBadRequest, "No name provided")
		return
	}
	if requestbody.PowerLevel == nil {
		c.String(http.StatusBadRequest, "No Powerlevel provided")
		return
	}

	s.Processor().AddDragonBallCharacter(c, requestbody.Name, *requestbody.PowerLevel)
}

func (s *Server) HTTPUpdateDragonBallCharacter(c *gin.Context) {
	logger.SBILog.Infof("In HTTPUpdateDragonBallCharacter")
	targetName := c.Param("name")

	if targetName == "" {
		c.String(http.StatusBadRequest, "No name provided")
		return
	}

	type RequestBody struct {
		PowerLevel *int32 `json:"powerLevel"`
	}
	var requestbody RequestBody
	if err := c.ShouldBindBodyWithJSON(&requestbody); err != nil {
		c.String(http.StatusBadRequest, "error")
		return
	}

	if requestbody.PowerLevel == nil {
		c.String(http.StatusBadRequest, "No Powerlevel provided")
		return
	}

	s.Processor().UpdateDragonBallCharacter(c, targetName, *requestbody.PowerLevel)
}
