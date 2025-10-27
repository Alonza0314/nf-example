// internal/sbi/processor/time_zone.go
package processor

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleGetTimeZone 查詢時區
func (p *Processor) HandleGetTimeZone(c *gin.Context, city string) {
	tzData := p.Context().TimeZoneData
	if tz, ok := tzData[city]; ok {
		c.String(http.StatusOK, tz)
		return
	}
	c.String(http.StatusNotFound, fmt.Sprintf("[%s] not found", city))
}

// TimeZoneRequest used for POST /city
type TimeZoneRequest struct {
	City     string `json:"City"`
	TimeZone string `json:"TimeZone"`
}

// HandleAddNewCityTimeZone 新增城市時區
func (p *Processor) HandleAddNewCityTimeZone(c *gin.Context, req TimeZoneRequest) {
	if _, ok := p.Context().TimeZoneData[req.City]; ok {
		c.String(http.StatusConflict, fmt.Sprintf("City '%s' already exists", req.City))
		return
	}
	p.Context().TimeZoneData[req.City] = req.TimeZone
	c.String(http.StatusOK, fmt.Sprintf("Time zone of %s is set to %s", req.City, req.TimeZone))
}

// HTTPResetCityTimeZone 重設時區
func (p *Processor) HandleResetCityTimeZone(c *gin.Context, city string, newTZ string) {
	if _, ok := p.Context().TimeZoneData[city]; !ok {
		c.String(http.StatusNotFound, fmt.Sprintf("City '%s' not found", city))
		return
	}
	p.Context().TimeZoneData[city] = newTZ
	c.String(http.StatusOK, fmt.Sprintf("Time zone of %s is reset to %s", city, newTZ))
}

// HTTPDeleteCityTimeZone 刪除城市時區
func (p *Processor) HandleDeleteCityTimeZone(c *gin.Context, city string) {
	if _, ok := p.Context().TimeZoneData[city]; !ok {
		c.String(http.StatusNotFound, fmt.Sprintf("City '%s' not found", city))
		return
	}
	delete(p.Context().TimeZoneData, city)
	c.String(http.StatusOK, fmt.Sprintf("City '%s' has been removed", city))
}
