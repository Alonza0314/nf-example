package sbi

import (
    "github.com/gabbyrz024/nf-example/internal/secondapi"
)

func (s *Server) getSecondAPIRoute() []Route {
    return []Route{
        {
            Name:    "secondapi-status",
            Method:  "GET",
            Pattern: "/status",
            APIFunc: secondapi.GetStatus,
        },
        {
            Name:    "secondapi-message",
            Method:  "POST",
            Pattern: "/message",
            APIFunc: secondapi.PostMessage,
        },
    }
}
