package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/faissalmaulana/21/api/cmd/handler"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPingHandler(t *testing.T) {
	pingHandler := handler.PingHandler{}
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	err := pingHandler.HandleFunc(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)

	var result map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &result)
	require.NoError(t, err)

	assert.Equal(t, map[string]string{"message": "pong"}, result)
}
