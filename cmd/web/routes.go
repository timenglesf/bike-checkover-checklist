package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/timenglesf/bike-checkover-checklist/ui"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.NotFound(http.HandlerFunc(app.handleNotFound))

	r.Use(app.sessionManager.LoadAndSave)

	// r.MethodNotAllowed(http.HandlerFunc(app.methodNotAllowedResponse))
	r.Use(middleware.Heartbeat("/ping"))

	// Serve static files
	r.Handle("/static/*", http.FileServerFS(ui.Files))

	r.Get("/tmp-checklist", app.handleDisplayChecklist)

	r.Get("/", app.handleDisplayMainPage)

	// r.Get("/login", app.userLogin)
	r.Post("/user/login", app.handlePostUserLogin)
	// r.Get("/checklist", app.getChecklistDisplayHandler)

	r.Post("/checklist", app.postChecklist)
	r.Put("/checklist", app.putChecklist)
	r.Put("/checklist/reset", app.handleChecklistReset)

	r.Get("/bike/{slug}", app.getBikeDisplay)

	r.Get("/user/logout", app.handleLogoutUser)
	r.Get("/history", app.handleDisplayUserHistory)
	return r
}
