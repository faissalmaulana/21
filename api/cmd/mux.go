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

	PingHandler        handler.Handle `name:"pingHandler"`
	GetProjectsHandler handler.Handle `name:"getProjectsHandler"`
	PostProjectHandler handler.Handle `name:"postProjectHandler"`
}

func NewEchoMux(p EchoMuxParams) http.Handler {
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.CORS("http://localhost:5173"))

	api := e.Group("/api")

	api.GET("/ping", p.PingHandler.HandleFunc)

	projects := api.Group("/projects")
	projects.POST("", p.PostProjectHandler.HandleFunc)
	projects.GET("/", p.GetProjectsHandler.HandleFunc)

	return e
}
