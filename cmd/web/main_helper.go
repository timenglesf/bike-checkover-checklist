package main

import (
	"context"
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
		admin: &models.AdminModel{
			DB:             dbClient,
			DBName:         DB_NAME,
			CollectionName: DB_ADMIN_COLLECTION,
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

func (app *application) adminVerification() error {
	if adminExists := app.adminExists(); !adminExists {
		// Check if the admin env vars are set
		if !app.adminEnvsValid() {
			return fmt.Errorf("admin env vars not set or too short")
		}

		// Create the admin
		adminUsername := os.Getenv(MB_ADMIN)
		adminPassword := os.Getenv(MB_PASS)

		app.logger.Info("Creating admin", "username", adminUsername)
		newAdmin, err := models.CreateAdmin(
			adminUsername,
			adminPassword,
			"admin",
			"admin",
		)
		if err != nil {
			return err
		}

		app.logger.Info("Inserting admin into database", "username", adminUsername)
		if err = app.admin.Insert(context.Background(), newAdmin); err != nil {
			return err
		}

		// Check if the admin was created
		err = app.adminVerification()
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *application) adminExists() bool {
	count, err := app.admin.GetDocumentCount(context.Background())
	if err != nil {
		app.logger.Error("Error getting admin count", "err", err)
		return false
	}
	return count > 0
}

func (app *application) adminEnvVarsExist() bool {
	return os.Getenv(MB_ADMIN) != "" && os.Getenv(MB_PASS) != ""
}

func (app *application) adminEnvsValid() bool {
	if !app.adminEnvVarsExist() {
		return false
	}
	name := os.Getenv(MB_ADMIN)
	pass := os.Getenv(MB_PASS)
	if len(name) < 5 || len(pass) < 8 {
		return false
	}

	return true
}
