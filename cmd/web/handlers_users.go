package main

import (
	"net/http"

	"github.com/timenglesf/bike-checkover-checklist/internal/models"
	"github.com/timenglesf/bike-checkover-checklist/internal/shared"
)

func (app application) handlePostUserCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	if !data.IsAdmin {
		app.clientError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var form models.UserForm

	if err := app.decodeForm(r, &form); err != nil {
		app.clientError(w, http.StatusBadRequest, "error decoding form", err)
		return
	}

	// if the pins do not match send a flash message
	if form.Pin != form.PinConfirm {
		form.Pin = ""
		form.PinConfirm = ""
		data.UserCreationFormData = form
		data.Flash.Message = "Pins must match"
		data.Flash.Type = shared.FlashError
		app.renderPage(w, r, app.pageTemplates.UserCreation, "Create User", &data)
		return
	}
	// if the pin already exists send a flash message

	count, err := app.users.GetDocumentCountByPin(r.Context(), form.Pin)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if count > 0 {
		data.UserCreationFormData = form
		data.Flash.Message = "Pin already exists"
		data.Flash.Type = shared.FlashError
		app.renderPage(w, r, app.pageTemplates.UserCreation, "Create User", &data)
		return
	}

	u := models.CreateUser(
		form.FirstName,
		form.LastName,
		form.Pin,
		form.StoreId,
	)

	if err := app.users.Insert(r.Context(), u); err != nil {
		app.serverError(w, r, err)
		return
	}

	data.Flash.Message = "User created successfully"
	data.Flash.Type = shared.FlashSuccess
	app.renderPage(w, r, app.pageTemplates.UserCreation, "Create User", &data)
}
