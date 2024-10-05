package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getKuromiRoute() []Route {
	return []Route{
		{
			Name:    "Hello KUROMI",
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: s.HTTPHelloKuromi,
			// Use
			// curl -X GET http://127.0.0.163:8000/kuromi/ -w "\n"
		},
		{
			Name:    "Big KUROMI",
			Method:  http.MethodPost,
			Pattern: "/",
			APIFunc: s.HTTPBigKuromi,
			// Use
			// curl -X POST http://127.0.0.163:8000/kuromi/ -w "\n"
		},
	}
}

func (s *Server) HTTPHelloKuromi(c *gin.Context) {
	c.JSON(http.StatusOK, "Hello Kuromi!")
}

func (s *Server) HTTPBigKuromi(c *gin.Context) {
	c.String(http.StatusOK, "⣴⣦⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣴⡷\n"+
		"⠈⣿⣷⣦⣄⡀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣤⣶⣿⣿\n"+
		"⠀⢸⣿⣿⣿⣿⣷⣆⣀⣀⣀⣀⣀⣾⣿⣿⣿⣿⡇\n"+
		"⠀⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡇\n"+
		"⠀⠀⠿⢿⣿⣿⣿⣿⡏⡀⠀⡙⣿⣿⣿⣿⣿⠛\n"+
		"⠀⠀⠀⣿⣿⣿⡿⠟⠷⣅⣀⠵⠟⢿⣿⣿⣿⡆\n"+
		"⠀⠀⠀⣿⣿⠏⢲⣤⠀⠀⠀⠀⢠⣶⠙⣿⣿⠃\n"+
		"⠀⠀⠀⠘⢿⡄⠈⠃⠀⢐⢔⠀⠈⠋⢀⡿⠋\n"+
		"⠀⠀⠀⢀⢀⣼⣷⣶⣤⣤⣭⣤⣴⣶⣍\n"+
		"⠀⠀⠀⠈⠈⣈⢰⠿⠛⠉⠉⢻⢇⠆⣁⠁\n"+
		"⠀⠀⠀⠀⠀⠑⢸⠉⠀⠀⠀⠀⠁⡄⢘⣽⣿\n"+
		"⠀⠀⠀⠀⠀⠀⡜⠀⠀⢰⡆⠀⠀⠻⠛⠋\n"+
		"⠀⠀⠀⠀⠀⠀⠑⠒⠒⠈⠈⠒⠒⠊")
}
