package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEchoMux(t *testing.T) {
	t.Run("NewEchoMux", func(t *testing.T) {
		p := EchoMuxParams{}

		handler := NewEchoMux(p)

		assert.NotNil(t, handler)
	})
}
