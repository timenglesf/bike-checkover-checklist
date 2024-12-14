package main

import (
	"errors"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/a-h/templ"
	"github.com/go-playground/form/v4"

	"github.com/timenglesf/bike-checkover-checklist/internal/models"
	"github.com/timenglesf/bike-checkover-checklist/internal/shared"
)

// ENUMS
const (
	SessionUserID   = "authenticatedUserId"
	SessionIsAdmin  = "isAdmin"
	SessionUserName = "firstName"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)
	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int, msg string, err error) {
	app.logger.Error(msg, "err", err.Error())
	http.Error(w, http.StatusText(status), status)
}

func (app *application) isAuthenticated(r *http.Request) bool {
	return app.sessionManager.Exists(r.Context(), SessionUserID)
}

func (app *application) newTemplateData(r *http.Request) shared.TemplateData {
	return shared.TemplateData{
		Date:             time.Now(),
		CurrentYear:      time.Now().Year(),
		Flash:            &shared.FlashMessage{},
		ChecklistDisplay: models.CreateChecklistDisplay(),
		IsAuthenticated:  app.isAuthenticated(r),
		// IsAdmin:         app.isAdmin(r),
		// CSRFToken:       nosurf.Token(r),
	}
}

func (app *application) renderPage(
	w http.ResponseWriter,
	r *http.Request,
	templateFunc func(data *shared.TemplateData) templ.Component,
	title string,
	data *shared.TemplateData,
) {
	page := templateFunc(data)
	base := app.pageTemplates.Base(title, page, data)
	err := base.Render(r.Context(), w)
	if err != nil {
		app.serverError(w, r, err)
	}
}

// a helper function to decode form data into a struct using the automatic form decoder
func (app *application) decodeForm(r *http.Request, dst any) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	err := app.formDecoder.Decode(&dst, r.PostForm)

	var inValidDecoderError *form.InvalidDecoderError

	if err != nil {
		if errors.As(err, &inValidDecoderError) {
			panic(err)
		}
		return err
	}

	return nil
}
