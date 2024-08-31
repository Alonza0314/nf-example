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
			APIFunc: func(c *gin.Context) {
				c.JSON(http.StatusOK, "Hello Kuromi!")
			},
			// Use
			// curl -X GET http://127.0.0.163:8000/kuromi/ -w "\n"
		},
	}
}

func (s *Server) postKuromiRoute() []Route {
	return []Route{
		{
			Name:    "Big KUROMI",
			Method:  http.MethodPost,
			Pattern: "/",
			APIFunc: func(c *gin.Context) {
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
			},
			// Use
			// curl -X POST http://127.0.0.163:8000/kuromi/ -w "\n"
		},
	}
}
