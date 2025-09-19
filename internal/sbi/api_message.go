package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getMessageRoute() []Route {
	return []Route{
		{
			Name:    "get messages",
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: s.getMessageRecord,
			// Use
			// curl -X GET http://127.0.0.163:8000/message/ -w "\n"
		},
	}
}

func (s *Server) putMessageRoute() []Route {
	return []Route{
		{
			Name:    "add message",
			Method:  http.MethodPut,
			Pattern: "/:Message",
			APIFunc: s.addNewMessage,
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

func (s *Server) addNewMessage(c *gin.Context) {

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

func (s *Server) getMessageRecord(c *gin.Context) {
	s.Processor().GetMessageRecord(c)
}
