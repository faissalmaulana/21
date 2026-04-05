package main

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"go.uber.org/fx"
)

type EchoMuxParams struct {
	fx.In
}

func NewEchoMux(p EchoMuxParams) http.Handler {
	e := echo.New()

	return e
}
