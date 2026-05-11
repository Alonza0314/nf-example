package processor

import (
	"net/http"

	context "github.com/Alonza0314/nf-example/internal/context"
	"github.com/gin-gonic/gin"
)

func (p *Processor) GetAllSubscribers(c *gin.Context) {
	subscribers := p.Context().Subscribers
	if len(subscribers) == 0 {
		c.JSON(http.StatusOK, []context.Subscriber{})
		return
	}
	c.JSON(http.StatusOK, subscribers)
}

func (p *Processor) AddSubscriber(c *gin.Context) {
	var sub context.Subscriber
	if err := c.ShouldBindJSON(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if sub.IMSI == "" || sub.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "imsi and name are required"})
		return
	}
	p.Context().Subscribers = append(p.Context().Subscribers, sub)
	c.JSON(http.StatusCreated, sub)
}
