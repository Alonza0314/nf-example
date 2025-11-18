package custom

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func AddCustomService(r *gin.Engine) {
    r.GET("/custom/hello", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "Hello from Evelin’s custom API!",
        })
    })

    r.POST("/custom/echo", func(c *gin.Context) {
        var data map[string]interface{}
        c.BindJSON(&data)
        c.JSON(http.StatusOK, gin.H{
            "received": data,
        })
    })
}

