package secondapi

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func GetStatus(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status": "second api running",
    })
}

type MessageRequest struct {
    Message string `json:"message"`
}

func PostMessage(c *gin.Context) {
    var req MessageRequest
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "received": req.Message,
    })
}
