package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTemplateData(t *testing.T) {
	app := &application{}
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	data := app.newTemplateData(r)

	assert.WithinDuration(t, time.Now(), data.Date, time.Second)
	assert.Equal(t, time.Now().Year(), data.CurrentYear)
	assert.NotNil(t, data.Flash)
}
