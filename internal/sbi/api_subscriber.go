package sbi

import (
	"net/http"

	"github.com/Alonza0314/nf-example/internal/logger"
	"github.com/gin-gonic/gin"
)

func (s *Server) getSubscriberRoute() []Route {
	return []Route{
		{
			Name:    "Get All Subscribers",
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: s.HTTPGetAllSubscribers,
			// Use: curl -X GET http://127.0.0.163:8000/subscriber/ -w "\n"
		},
		{
			Name:    "Add Subscriber",
			Method:  http.MethodPost,
			Pattern: "/",
			APIFunc: s.HTTPAddSubscriber,
			// Use: curl -X POST http://127.0.0.163:8000/subscriber/ \
			//      -H "Content-Type: application/json" \
			//      -d '{"imsi":"208930000000001","name":"UE-Test"}' -w "\n"
		},
	}
}

func (s *Server) HTTPGetAllSubscribers(c *gin.Context) {
	logger.SBILog.Infof("HTTPGetAllSubscribers")
	s.Processor().GetAllSubscribers(c)
}

func (s *Server) HTTPAddSubscriber(c *gin.Context) {
	logger.SBILog.Infof("HTTPAddSubscriber")
	s.Processor().AddSubscriber(c)
}
