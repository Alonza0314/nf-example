package newapi

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func GetMessage(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "Hola desde el nuevo GET!",
    })
}

type CreateRequest struct {
    Name string `json:"name"`
}

func CreateItem(c *gin.Context) {
    var req CreateRequest
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "status": "created",
        "name": req.Name,
    })
}
