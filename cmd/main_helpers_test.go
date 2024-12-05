package main

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInitializeServer(t *testing.T) {
	app := &application{
		cfg: &config{
			port: "8080",
		},
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
	server := app.intializeServer()

	assert.Equal(t, ":8080", server.Addr)
	assert.NotNil(t, server.Handler)
	assert.Equal(t, time.Minute, server.IdleTimeout)
	assert.Equal(t, 5*time.Second, server.ReadTimeout)
	assert.Equal(t, 10*time.Second, server.WriteTimeout)
	assert.NotNil(t, server.ErrorLog)
}

func TestCreateLogger(t *testing.T) {
	logger := createLogger()
	assert.NotNil(t, logger)
}

func TestCreateApplication(t *testing.T) {
	app, err := createApplication()
	assert.NoError(t, err)
	assert.NotNil(t, app)
	assert.NotNil(t, app.logger)
	assert.NotNil(t, app.cfg)
	assert.NotNil(t, app.pageTemplates)
	assert.NotNil(t, app.partialTemplates)
}

func TestCreateConfig(t *testing.T) {
	os.Setenv("PORT", "8080")
	os.Setenv("DEFAULT_ADMIN_PASSWORD", "password")
	cfg := createConfig()

	assert.Equal(t, "8080", cfg.port)
	assert.Equal(t, "password", cfg.defaultAdminPassword)
}

func TestGetEnv(t *testing.T) {
	os.Setenv("TEST_ENV", "test_value")
	value := getEnv("TEST_ENV", "default_value")
	assert.Equal(t, "test_value", value)

	value = getEnv("NON_EXISTENT_ENV", "default_value")
	assert.Equal(t, "default_value", value)
}
