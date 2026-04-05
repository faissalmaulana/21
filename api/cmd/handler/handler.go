package handler

import "github.com/labstack/echo/v5"

type Handle interface {
	HandleFunc(c *echo.Context) error
}
