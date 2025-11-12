package sbi

import (
	"fmt"
	"net/http"

	api "github.com/Alonza0314/nf-example/service/api"
	"github.com/Alonza0314/nf-example/internal/logger"
	"github.com/Alonza0314/nf-example/pkg/app"
	"github.com/gin-gonic/gin"

	"github.com/free5gc/util/httpwrapper"
	logger_util "github.com/free5gc/util/logger"
)

// Estructura base de rutas
type Route struct {
	Name    string
	Method  string
	Pattern string
	APIFunc gin.HandlerFunc
}

// Aplica rutas a un grupo dado
func applyRoutes(group *gin.RouterGroup, routes []Route) {
	for _, route := range routes {
		switch route.Method {
		case "GET":
			group.GET(route.Pattern, route.APIFunc)
		case "POST":
			group.POST(route.Pattern, route.APIFunc)
		case "PUT":
			group.PUT(route.Pattern, route.APIFunc)
		case "PATCH":
			group.PATCH(route.Pattern, route.APIFunc)
		case "DELETE":
			group.DELETE(route.Pattern, route.APIFunc)
		}
	}
}

// Definición del router principal
func newRouter(s *Server) *gin.Engine {
	router := logger_util.NewGinWithLogrus(logger.GinLog)

	// Grupo /default
	defaultGroup := router.Group("/default")
	applyRoutes(defaultGroup, s.getDefaultRoute())

	// Grupo /spyfamily
	spyFamilyGroup := router.Group("/spyfamily")
	applyRoutes(spyFamilyGroup, s.getSpyFamilyRoute())

	// 🚀 NUEVO GRUPO /myapi
	myapiGroup := router.Group("/myapi")
	{
		myapiGroup.GET("/hello", api.HelloHandler)
		myapiGroup.POST("/echo", api.EchoHandler)
	}

	return router
}

// Enlaza el router al servidor HTTP
func bindRouter(nf app.App, router *gin.Engine, tlsKeyLogPath string) (*http.Server, error) {
	sbiConfig := nf.Config().Configuration.Sbi
	bindAddr := fmt.Sprintf("%s:%d", sbiConfig.BindingIPv4, sbiConfig.Port)
	return httpwrapper.NewHttp2Server(bindAddr, tlsKeyLogPath, router)
}
