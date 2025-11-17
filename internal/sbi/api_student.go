package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetStudentData(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"student": "Gilbert",
		"semester": 7,
		"university": "UDC",
		"message": "Datos obtenidos correctamente",
	})
}

type StudentRequest struct {
	Name      string `json:"name" binding:"required"`
	Semester  int    `json:"semester"`
	Career    string `json:"career"`
}

func RegisterStudent(c *gin.Context) {
	var req StudentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "registrado correctamente",
		"name": req.Name,
		"semester": req.Semester,
		"career": req.Career,
	})
}
