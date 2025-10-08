package processor

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (p *Processor) ReturnAttendence(c *gin.Context) {
	con := p.Context()

	if len(con.AttendenceData) == 0 {
		c.String(http.StatusOK, "No attendence recorded")
		return
	} else {
		names := ""
		for _, name := range con.AttendenceData {
			names += name + ", "
		}
		c.String(http.StatusOK, `{"Attendence":"`+names[:len(names)-2]+`"}`)
		return
	}
}

func (p *Processor) PostAttendence(c *gin.Context, targetName string) {
	con := p.Context()

	for n := range con.AttendenceData {
		if con.AttendenceData[n] == targetName {
			c.String(http.StatusConflict, "Attendence already recorded: "+targetName)
			return
		}
	}

	con.AttendenceData = append(con.AttendenceData, targetName)

	c.String(http.StatusOK, "Attendence recorded: "+targetName)
}
