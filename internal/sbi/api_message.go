package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) myPutGetMessageRoute() []Route {
	return []Route{
		{
			Name:    "get messages",
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: s.HTTPGetMessageRecord,
			// Use
			// curl -X GET http://127.0.0.163:8000/message/ -w "\n"
			// return all added message
		},
		{
			Name:    "add message",
			Method:  http.MethodPut,
			Pattern: "/:Message",
			APIFunc: s.HTTPAddNewMessage,
			// Use
			// curl -X PUT http://127.0.0.163:8000/message/yourmessage -w "\n"
			// add "yourmessage" to message record
		},
		{
			// empty input handle, will not accept
			Name:    "empty input",
			Method:  http.MethodPut,
			Pattern: "/",
			APIFunc: s.noMessageHandler,
		},
	}
}

func (s *Server) HTTPAddNewMessage(c *gin.Context) {

	newMessage := c.Param("Message")
	if newMessage == "" {
		s.noMessageHandler(c)
		return
	}
	s.Processor().AddNewMessage(c, newMessage)
}
func (s *Server) noMessageHandler(c *gin.Context) {
	c.String(http.StatusBadRequest, "No message provided")
}

func (s *Server) HTTPGetMessageRecord(c *gin.Context) {
	s.Processor().GetMessageRecord(c)
}
