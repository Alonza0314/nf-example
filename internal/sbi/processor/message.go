package processor

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (p *Processor) AddNewMessage(c *gin.Context, newMessage string) {
	p.Context().MessageMu.Lock()
	defer p.Context().MessageMu.Unlock()
	//addd message
	p.Context().MessageRecord = append(p.Context().MessageRecord, newMessage)
	c.String(http.StatusOK, "add a new message!")
}
func (p *Processor) GetMessageRecord(c *gin.Context) {
	p.Context().MessageMu.Lock()
	defer p.Context().MessageMu.Unlock()

	//no content
	if len(p.Context().MessageRecord) == 0 {
		c.String(http.StatusOK, "no message now, add some messagess!")
		return
	}
	//get record
	Record := ""
	for _, s := range p.Context().MessageRecord {
		Record += fmt.Sprintf("%s\n", s)
	}
	c.String(http.StatusOK, Record)
}
