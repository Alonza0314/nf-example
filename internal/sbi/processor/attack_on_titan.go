package processor

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (p *Processor) FindAttackOnTitanCharacterName(c *gin.Context, targetName string) {
	if lastName, ok := p.Context().AttackOnTitanData[targetName]; ok {
		c.String(http.StatusOK, fmt.Sprintf("Character: %s %s", targetName, lastName))
		return
	}
	c.String(http.StatusNotFound, fmt.Sprintf("[%s] not found in Attack on Titan", targetName))
}

func (p *Processor) AddAttackOnTitanCharacter(c *gin.Context) {
	var body struct {
		Name     string `json:"name"`
		LastName string `json:"lastName"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.String(http.StatusBadRequest, "Invalid request body")
		return
	}
	if body.Name == "" || body.LastName == "" {
		c.String(http.StatusBadRequest, "Name and LastName are required")
		return
	}
	p.Context().AttackOnTitanData[body.Name] = body.LastName
	c.String(http.StatusOK, fmt.Sprintf("Character %s %s added successfully", body.Name, body.LastName))
}