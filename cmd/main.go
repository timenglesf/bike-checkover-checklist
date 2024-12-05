package main

import (
	"log/slog"
	"os"

	"github.com/timenglesf/bike-checkover-checklist/ui/template"
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
	// db *mongo.Client
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

	svr := app.intializeServer()

	app.logger.Info("starting server", "host", svr.Addr, "port", app.cfg.port)
	err = svr.ListenAndServe()
	if err != nil {
		app.logger.Error("server error", "err", err.Error())
		os.Exit(1)
	}
}
