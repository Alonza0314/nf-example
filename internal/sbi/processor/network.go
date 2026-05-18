package processor

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (p *Processor) FindNetworkDevice(c *gin.Context, targetName string) {
	if deviceType, ok := p.Context().NetworkData[targetName]; ok {
		c.JSON(http.StatusOK, gin.H{
			"name": targetName,
			"type": deviceType,
		})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": fmt.Sprintf("device [%s] not found", targetName),
	})
}

func (p *Processor) AddNetworkDevice(c *gin.Context, name string, deviceType string) {
	if _, ok := p.Context().NetworkData[name]; ok {
		c.JSON(http.StatusConflict, gin.H{
			"error": fmt.Sprintf("device [%s] already exists", name),
		})
		return
	}
	p.Context().NetworkData[name] = deviceType
	c.JSON(http.StatusCreated, gin.H{
		"name": name,
		"type": deviceType,
	})
}
