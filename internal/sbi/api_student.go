package sbi

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

// Rutas para el nuevo servicio /student
func (s *Server) getStudentRoute() []Route {
    return []Route{
        {
            Name:   "GetStudent",
            Method: http.MethodGet,
            Pattern: "/:id",
            APIFunc: s.HTTPGetStudent,
        },
        {
            Name:   "CreateStudent",
            Method: http.MethodPost,
            Pattern: "/",
            APIFunc: s.HTTPCreateStudent,
        },
    }
}

// Estructura para el cuerpo del POST
type Student struct {
    ID   string `json:"id" binding:"required"`
    Name string `json:"name" binding:"required"`
}

// GET /student/:id
func (s *Server) HTTPGetStudent(c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "id":   id,
        "name": "Student " + id,
        "msg":  "Student info from nf-example",
    })
}

// POST /student/
func (s *Server) HTTPCreateStudent(c *gin.Context) {
    var student Student

    if err := c.ShouldBindJSON(&student); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "msg":     "Student created successfully",
        "student": student,
    })
}
