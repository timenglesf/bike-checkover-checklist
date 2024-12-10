package main

import (
	"errors"
	"net/http"

	"github.com/timenglesf/bike-checkover-checklist/internal/shared"
	"github.com/timenglesf/bike-checkover-checklist/internal/validator"
	"go.mongodb.org/mongo-driver/mongo"
)

func (app *application) handleDisplayMainPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	if !data.IsAuthenticated {
		app.renderPage(w, r, app.pageTemplates.CheckList, "Check List", &data)
	}
	// Tmp
	u := app.sessionManager.GetString(r.Context(), SessionUserName)
	app.logger.Info("User", "user", u)
	_, err := w.Write([]byte(u))
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) handleDisplayChecklist(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.renderPage(w, r, app.pageTemplates.CheckList, "Check List", &data)
}

func (app *application) handlePostUserLogin(w http.ResponseWriter, r *http.Request) {
	var form shared.PinForm

	err := app.decodeForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, "Error decoding form", err)
		return
	}

	form.CheckField(
		validator.NotBlank(form.Pin),
		"pin",
		"pin must not be blank",
	)

	app.logger.Info("login attempt", "pin", form.Pin)

	data := app.newTemplateData(r)
	data.PinForm = form

	if !form.Valid() {
		data.Flash.Message = "Pin must not be blank"
		data.Flash.Type = shared.FlashError
		app.renderPage(w, r, app.pageTemplates.UserLogin, "User Login", &data)
		return
	}

	user, err := app.users.FindByPin(form.Pin, r.Context())
	if err != nil {
		app.logger.Error("Error retrieving user", "err", err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			data.Flash.Message = "Invalid pin"
			data.Flash.Type = shared.FlashError
			app.renderPage(w, r, app.pageTemplates.UserLogin, "User Login", &data)
			return
		}
	}
	data.User = &user
	app.sessionManager.Put(r.Context(), SessionUserID, user.ID.Hex())
	app.sessionManager.Put(r.Context(), SessionUserName, user.FirstName)
}

func (app *application) handleNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	data := app.newTemplateData(r)
	app.logger.Info("Page Not Found")
	app.renderPage(w, r, app.pageTemplates.NotFound, "Page Not Found", &data)
}

// TMP
func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.renderPage(w, r, app.pageTemplates.UserLogin, "User Login", &data)
}
