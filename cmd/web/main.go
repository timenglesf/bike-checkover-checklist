package main

import (
	"log/slog"
	"os"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/timenglesf/bike-checkover-checklist/internal/db"
	"github.com/timenglesf/bike-checkover-checklist/internal/models"
	"github.com/timenglesf/bike-checkover-checklist/ui/template"
)

const (
	DB_CONNECTION_ATTEMPTS  = 10
	DB_NAME                 = "bike-check-in"
	DB_USER_COLLECTION      = "users"
	DB_CHECKLIST_COLLECTION = "checklists"
)

type application struct {
	logger *slog.Logger
	cfg    *config
	// UI Templates
	pageTemplates    *template.Pages
	partialTemplates *template.Partials
	// Database & Models
	db             *mongo.Client
	sessionManager *scs.SessionManager
	users          *models.UserModel
	checklist      *models.ChecklistModel
	// models data.Models

	formDecoder *form.Decoder
}

type config struct {
	port                 string
	defaultAdminPassword string
}

func main() {
	// Create logger
	logger := createLogger()

	// Connect to the database.
	dbClient, err := db.ConnectWithRetries(
		logger,
		DB_CONNECTION_ATTEMPTS,
	)
	if err != nil {
		logger.Error("Failed to connect to database", "err", err.Error())
		os.Exit(1)
	}
	defer db.Close(dbClient)

	// Initialize session manager
	sessionManager := initializeSessionManager(dbClient)
	if err != nil {
		logger.Error("unable to initialize session manager", "error", err)
		panic("failed to initialize session manager")
	}

	// Create the application state
	app, err := createApplication(dbClient, logger, sessionManager)
	if err != nil {
		app.logger.Error("Error creating application", "err", err)
		return
	}

	// Start the HTTP server.
	svr := app.intializeServer()

	app.logger.Info("starting server", "host", svr.Addr, "port", app.cfg.port)

	err = svr.ListenAndServe()
	if err != nil {
		app.logger.Error("server error", "err", err.Error())
		os.Exit(1)
	}
}
