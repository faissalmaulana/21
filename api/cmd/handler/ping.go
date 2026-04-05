package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type PingHandler struct{}

func NewPingHandler() Handle {
	return &PingHandler{}
}

func (p *PingHandler) HandleFunc(c *echo.Context) error {

	return c.JSON(http.StatusOK, map[string]string{"message": "pong"})
}
