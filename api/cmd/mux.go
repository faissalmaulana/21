package main

import (
	"net/http"

	"github.com/faissalmaulana/21/api/cmd/handler"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"go.uber.org/fx"
)

type EchoMuxParams struct {
	fx.In

	PingHandler handler.Handle `name:"pingHandler"`
}

func NewEchoMux(p EchoMuxParams) http.Handler {
	e := echo.New()
	e.Use(middleware.RequestLogger())

	e.GET("/ping", p.PingHandler.HandleFunc)

	return e
}
