package main

import (
	"errors"
	"net/http"

	"github.com/timenglesf/bike-checkover-checklist/internal/shared"
	"github.com/timenglesf/bike-checkover-checklist/internal/validator"

	"go.mongodb.org/mongo-driver/mongo"
)

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

	user, err := app.users.FindByPin(r.Context(), form.Pin)
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
	app.logger.Info("User logged in", "user", user)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) handleLogoutUser(w http.ResponseWriter, r *http.Request) {
	if err := app.sessionManager.Destroy(r.Context()); err != nil {
		app.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
