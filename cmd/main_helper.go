package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mongodbstore"
	"github.com/alexedwards/scs/v2"
	"github.com/timenglesf/bike-checkover-checklist/ui/template"
	"go.mongodb.org/mongo-driver/mongo"
)

// Initialize the HTTP server with configuration settings
func (app *application) intializeServer() *http.Server {
	fmt.Println("port", app.cfg.port)
	return &http.Server{
		Addr:         ":" + app.cfg.port,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}
}

func createLogger() *slog.Logger {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return logger
}

func createApplication() (*application, error) {
	app := &application{
		logger:           createLogger(),
		cfg:              createConfig(),
		pageTemplates:    template.CreatePages(),
		partialTemplates: template.CreatePartials(),
	}
	app.logger = createLogger()
	app.cfg = createConfig()
	return app, nil
}

func createConfig() *config {
	cfg := &config{}
	cfg.port = getEnv("PORT", "8081")
	cfg.defaultAdminPassword = getEnv("DEFUALT_ADMIN_PASSWORD", "password")

	return cfg
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func initializeSessionManager(dbClient *mongo.Client) *scs.SessionManager {
	sessionManager := scs.New()
	sessionManager.Store = mongodbstore.New(dbClient.Database(DB_NAME))
	sessionManager.Lifetime = 24 * 7 * time.Hour
	return sessionManager
}
