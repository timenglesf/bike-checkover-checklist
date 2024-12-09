package main

import (
	"errors"
	"net/http"

	"github.com/timenglesf/bike-checkover-checklist/internal/shared"
	"github.com/timenglesf/bike-checkover-checklist/internal/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (app *application) handleDisplayMainPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	if !data.IsAuthenticated {
		app.renderPage(w, r, app.pageTemplates.UserLogin, "Bike Intake Checklist: Login", &data)
		return
	}

	// Get user's active checklist using UserId
	userIdStr := app.sessionManager.GetString(r.Context(), SessionUserID)
	app.logger.Info("user access", "userId", userIdStr)

	userId, err := primitive.ObjectIDFromHex(userIdStr)
	if err != nil {
		app.logger.Error("error converting userId", "userId", userIdStr, "err", err)
		app.clientError(w, http.StatusUnprocessableEntity, "error", err)
	}

	cl, err := app.checklist.GetRecentActiveChecklist(r.Context(), userId)
	if err != nil {
		// if the error is that there is no documument use the one that is created with app.newTemplateData
		if errors.Is(err, mongo.ErrNoDocuments) {
			app.logger.Info("No active checklist found")
			user, err := app.users.FindById(r.Context(), userId)
			if err != nil {
				app.clientError(w, http.StatusNotFound, "no user found in database", err)
			}

			if err = app.checklist.Insert(r.Context(), *data.Checklist, user); err != nil {
				app.serverError(w, r, err)
			}

			app.renderPage(w, r, app.pageTemplates.CheckList, "Check List", &data)
			return
		}
		app.serverError(w, r, err)
		return
	}
	app.logger.Info("checklist found with id", "id", cl.ID.Hex())
	data.Checklist = &cl.Checklist
	data.ChecklistDocumentId = cl.ID.Hex()
	// save the new checklist to the db

	app.renderPage(w, r, app.pageTemplates.CheckList, "Bike Intake Checklist", &data)
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
