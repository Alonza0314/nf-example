package sbi

import (
	"fmt"
	"net/http"

	"github.com/andy89923/nf-example/internal/logger"
	"github.com/andy89923/nf-example/pkg/app"
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

	spyFamilyGroup := router.Group("/spyfamily")
	applyRoutes(spyFamilyGroup, s.getSpyFamilyRoute())

	stu113Group := router.Group("/stu113")
	applyRoutes(stu113Group, s.getStu113Route())

	saoGroup := router.Group("/sao")
	applyRoutes(saoGroup, s.getSAORoute())

	kuroumiGroup := router.Group("/kuromi")
	applyRoutes(kuroumiGroup, s.getKuromiRoute())
	applyRoutes(kuroumiGroup, s.postKuromiRoute())

	return router
}

func bindRouter(nf app.App, router *gin.Engine, tlsKeyLogPath string) (*http.Server, error) {
	sbiConfig := nf.Config().Configuration.Sbi
	bindAddr := fmt.Sprintf("%s:%d", sbiConfig.BindingIPv4, sbiConfig.Port)
	return httpwrapper.NewHttp2Server(bindAddr, tlsKeyLogPath, router)
}
