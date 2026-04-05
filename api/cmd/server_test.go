package main

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func NewForTest(tb testing.TB, opts ...fx.Option) *fx.App {
	testOpts := []fx.Option{
		// Provide both: Logger and WithLogger so that if the test
		// WithLogger fails, we don't pollute stderr.
		fx.Logger(fxtest.NewTestPrinter(tb)),
		fxtest.WithTestLogger(tb),
		fx.Provide(func() *zap.Logger {
			return zaptest.NewLogger(tb)
		}),
	}
	opts = append(testOpts, opts...)

	return fx.New(opts...)
}

func TestHttpServer(t *testing.T) {
	t.Parallel()

	t.Run("StartNewHttpServer", func(t *testing.T) {
		var srv *http.Server

		app := NewForTest(t,
			fx.Provide(func() http.Handler { return http.DefaultServeMux }),
			fx.Provide(NewHttpServer),
			fx.Invoke(func(s *http.Server) {
				srv = s
			}),
		)

		ctx := context.Background()

		require.NoError(t, app.Start(ctx))
		defer app.Stop(ctx)

		assert.NotNil(t, srv)
		assert.Equal(t, ":8080", srv.Addr)
	})

	t.Run("RequestHttpServer", func(t *testing.T) {

		e := echo.New()
		e.GET("/ping", func(c *echo.Context) error {
			return c.String(http.StatusOK, "pong")
		})

		app := NewForTest(t,
			fx.Provide(func() http.Handler { return e }),
			fx.Provide(NewHttpServer),
			fx.Invoke(func(s *http.Server) {}),
		)

		ctx := context.Background()

		require.NoError(t, app.Start(ctx))
		defer app.Stop(ctx)

		resp, err := http.Get("http://localhost:8080/ping")
		require.NoError(t, err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "pong", string(body))
	})
}
