package http

import (
	"github.com/labstack/echo"
)

type Server struct {
	E       *echo.Echo
	Address string
}

func NewHTTPServer(echo *echo.Echo, address string) *Server {
	return &Server{
		E:       echo,
		Address: address,
	}
}

func (s *Server) ServeForever(errorChan chan<- error) {
	if err := s.E.Start(s.Address); err != nil && errorChan != nil {
		errorChan <- err
	}
}
