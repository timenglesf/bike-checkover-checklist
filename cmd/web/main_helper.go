package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mongodbstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/timenglesf/bike-checkover-checklist/internal/models"
	"github.com/timenglesf/bike-checkover-checklist/ui/template"
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

func createApplication(
	dbClient *mongo.Client,
	logger *slog.Logger,
	sessionManager *scs.SessionManager,
) (*application, error) {
	app := &application{
		logger:           logger,
		sessionManager:   sessionManager,
		cfg:              createConfig(),
		pageTemplates:    template.CreatePages(),
		partialTemplates: template.CreatePartials(),
		db:               dbClient,
		users: &models.UserModel{
			DB:             dbClient,
			DBName:         DB_NAME,
			CollectionName: DB_USER_COLLECTION,
		},
		checklist: &models.ChecklistModel{
			DB:             dbClient,
			DBName:         DB_NAME,
			CollectionName: DB_CHECKLIST_COLLECTION,
		},
		formDecoder: form.NewDecoder(),
	}
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
