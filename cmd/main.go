package main

import (
	"log/slog"
	"os"

	"github.com/alexedwards/scs/v2"
	"github.com/timenglesf/bike-checkover-checklist/internal/db"
	"github.com/timenglesf/bike-checkover-checklist/ui/template"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DB_CONNECTION_ATTEMPTS = 10
	DB_NAME                = "bike-check-in"
)

type application struct {
	logger           *slog.Logger
	cfg              *config
	pageTemplates    *template.Pages
	partialTemplates *template.Partials
	db               *mongo.Client
	sessionManager   *scs.SessionManager
	// models data.Models
}

type config struct {
	port                 string
	defaultAdminPassword string
}

func main() {
	app, err := createApplication()
	if err != nil {
		app.logger.Error("Error creating application", "err", err)
		return
	}

	// Connect to the database.
	dbClient, err := db.ConnectWithRetries(app.logger, DB_CONNECTION_ATTEMPTS)
	if err != nil {
		app.logger.Error("Failed to connect to database", "err", err.Error())
		os.Exit(1)
	}
	defer db.Close(dbClient)

	app.db = dbClient

	// Initialize session manager
	sessionManager := initializeSessionManager(dbClient)
	if err != nil {
		app.logger.Error("unable to initialize session manager", "error", err)
		panic("failed to initialize session manager")
	}
	app.sessionManager = sessionManager

	// Start the HTTP server.
	svr := app.intializeServer()

	app.logger.Info("starting server", "host", svr.Addr, "port", app.cfg.port)

	err = svr.ListenAndServe()
	if err != nil {
		app.logger.Error("server error", "err", err.Error())
		os.Exit(1)
	}
}
