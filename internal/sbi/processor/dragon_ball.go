package processor

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (p *Processor) SearchDragonBallCharacter(c *gin.Context, targetName string) {
	pl, ok := p.Context().DragonBallData[targetName]
	if !ok {
		c.String(http.StatusNotFound, fmt.Sprintf("[%s] not found in Dragon Ball\n", targetName))
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("Character: %s, Powerlevel: %d\n", targetName, pl))
}

func (p *Processor) FightDragonBall(c *gin.Context, targetName1 string, targetName2 string) {
	pl1, ok1 := p.Context().DragonBallData[targetName1]
	pl2, ok2 := p.Context().DragonBallData[targetName2]

	if !ok1 {
		c.String(http.StatusNotFound, fmt.Sprintf("[%s] not found in Dragon Ball\n", targetName1))
		return
	}
	if !ok2 {
		c.String(http.StatusNotFound, fmt.Sprintf("[%s] not found in Dragon Ball\n", targetName2))
		return
	}

	if pl1 > pl2 {
		c.String(http.StatusOK, fmt.Sprintf("%s defeats %s\n", targetName1, targetName2))
	} else if pl1 < pl2 {
		c.String(http.StatusOK, fmt.Sprintf("%s defeats %s\n", targetName2, targetName1))
	} else {
		c.String(http.StatusOK, fmt.Sprintf("%s ties with %s\n", targetName1, targetName2))
	}
}

func (p *Processor) AddDragonBallCharacter(c *gin.Context, targetName string, powerlevel int32) {
	pl, ok := p.Context().DragonBallData[targetName]
	if ok {
		c.String(http.StatusConflict, fmt.Sprintf("Character %s already exists with Powerlevel %d\n", targetName, pl))
		return
	}
	p.Context().DragonBallData[targetName] = powerlevel
	c.String(http.StatusCreated, fmt.Sprintf("Add Character %s with Powerlevel %d\n", targetName, powerlevel))
}

func (p *Processor) UpdateDragonBallCharacter(c *gin.Context, targetName string, powerlevel int32) {
	if _, ok := p.Context().DragonBallData[targetName]; !ok {
		c.String(http.StatusNotFound, fmt.Sprintf("Character %s not found\n", targetName))
		return
	}
	p.Context().DragonBallData[targetName] = powerlevel
	c.String(http.StatusOK, fmt.Sprintf("Update Character %s with Powerlevel %d\n", targetName, powerlevel))
}
