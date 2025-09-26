package sbi

import (
	"fmt"
	"net/http"

	"github.com/Alonza0314/nf-example/internal/logger"
	"github.com/Alonza0314/nf-example/pkg/app"
	"github.com/gin-gonic/gin"

	"github.com/free5gc/util/httpwrapper"
	logger_util "github.com/free5gc/util/logger"
)

type Route struct {
	Name    string
	Method  string
	Pattern string
	APIFunc gin.HandlerFunc
}

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

func newRouter(s *Server) *gin.Engine {
	router := logger_util.NewGinWithLogrus(logger.GinLog)

	defaultGroup := router.Group("/default")
	applyRoutes(defaultGroup, s.getDefaultRoute())

	myPutGetMessageGroup := router.Group("/message")
	applyRoutes(myPutGetMessageGroup, s.myPutGetMessageRoute())

	spyFamilyGroup := router.Group("/spyfamily")
	applyRoutes(spyFamilyGroup, s.getSpyFamilyRoute())

	taskGroup := router.Group("/task")
	applyRoutes(taskGroup, s.getTaskRoute())

	messageGroup := router.Group("/msg") // add for lab6
	applyRoutes(messageGroup, s.getMessageRoute())

	return router
}

func bindRouter(nf app.App, router *gin.Engine, tlsKeyLogPath string) (*http.Server, error) {
	sbiConfig := nf.Config().Configuration.Sbi
	bindAddr := fmt.Sprintf("%s:%d", sbiConfig.BindingIPv4, sbiConfig.Port)
	return httpwrapper.NewHttp2Server(bindAddr, tlsKeyLogPath, router)
}
