package sbi

import (
	"net/http"

	"github.com/Alonza0314/nf-example/internal/logger"
	"github.com/gin-gonic/gin"
)

func (s *Server) getSpyFamilyRoute() []Route {
	return []Route{
		{
			Name:    "Hello SPYxFAMILY!",
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: func(c *gin.Context) {
				c.JSON(http.StatusOK, "Hello SPYxFAMILY!")
			},
			// Use
			// curl -X GET http://127.0.0.163:8000/spyfamily/ -w "\n"
		},
		{
			Name:    "SPYxFAMILY Character",
			Method:  http.MethodGet,
			Pattern: "/character/:Name",
			APIFunc: s.HTTPSerchSpyFamilyCharacter,
			// Use
			// curl -X GET http://127.0.0.163:8000/spyfamily/Anya -w "\n"
			// "Character: Anya Forger"
		},
	}
}

func (s *Server) HTTPSerchSpyFamilyCharacter(c *gin.Context) {
	logger.SBILog.Infof("In HTTPSerchCharacter")

	targetName := c.Param("Name")
	if targetName == "" {
		c.String(http.StatusBadRequest, "No name provided")
		return
	}

	s.Processor().FindSpyFamilyCharacterName(c, targetName)
}
func HTTPGetStudent(c *gin.Context) {
    name := c.Param("name")

    c.JSON(200, gin.H{
        "student": name,
        "status":  "active",
    })
}

func HTTPCreateStudent(c *gin.Context) {

    type Student struct {
        Name string `json:"name"`
        Age  int    `json:"age"`
    }

    var body Student

    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(400, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(201, gin.H{
        "message": "Student created successfully",
        "student": body,
    })
}
