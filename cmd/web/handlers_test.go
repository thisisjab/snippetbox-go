package main

import (
	"github.com/go-playground/assert"
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	status, _, body := ts.get(t, "/ping")

	assert.Equal(t, status, http.StatusOK)
	assert.Equal(t, body, "pong")
}
