package main

import (
	"net/http"

	"github.com/faissalmaulana/21/api/cmd/handler"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.Provide(
			NewEchoMux,
			NewHttpServer,
			fx.Annotate(handler.NewPingHandler, fx.As(new(handler.Handle)), fx.ResultTags(`name:"pingHandler"`)),
			zap.NewDevelopment,
		),
		fx.Invoke(func(*http.Server) {}),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
	).Run()
}
