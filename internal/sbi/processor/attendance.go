package processor

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (p *Processor) ReturnAttendance(c *gin.Context) {
	con := p.Context()

	if len(con.AttendanceData) == 0 {
		c.String(http.StatusOK, "No attendance recorded")
		return
	} else {
		names := ""
		for _, name := range con.AttendanceData {
			names += name + ", "
		}
		c.String(http.StatusOK, "Attendance: "+names[:len(names)-2])
		return
	}
}

func (p *Processor) PostAttendance(c *gin.Context, targetName string) {
	con := p.Context()

	for n := range con.AttendanceData {
		if con.AttendanceData[n] == targetName {
			c.String(http.StatusConflict, "Attendance already recorded: "+targetName)
			return
		}
	}

	con.AttendanceData = append(con.AttendanceData, targetName)

	c.String(http.StatusOK, "Attendance recorded: "+targetName)
}
