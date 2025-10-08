package processor

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type PostFortuneRequest struct {
	Fortune string `json:"fortune" binding:"required"`
}

func (p *Processor) GetFortune(c *gin.Context) {
	ctx := p.Context()

	ctx.FortuneMutex.RLock()
	defer ctx.FortuneMutex.RUnlock()

	if len(ctx.Fortunes) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "No fortunes available.",
		})
		return
	}

	// Seed the random number generator
	rand.New(rand.NewSource(time.Now().UnixNano()))
	// Get a random fortune
	fortune := ctx.Fortunes[rand.Intn(len(ctx.Fortunes))]

	c.JSON(http.StatusOK, gin.H{
		"message": "Here is your fortune for today!",
		"fortune": fortune,
	})
}

func (p *Processor) PostFortune(c *gin.Context, req PostFortuneRequest) {
	ctx := p.Context()

	ctx.FortuneMutex.Lock()
	defer ctx.FortuneMutex.Unlock()

	ctx.Fortunes = append(ctx.Fortunes, req.Fortune)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Fortune added successfully",
		"fortune": req.Fortune,
	})
}
