package main

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/timenglesf/bike-checkover-checklist/internal/models"
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

	clDoc, err := app.checklist.GetRecentActiveChecklist(r.Context(), userId)
	if err != nil {
		// if the error is that there is no documument use the one that is created with app.newTemplateData
		if errors.Is(err, mongo.ErrNoDocuments) {
			app.logger.Info("No active checklist found")
			user, err := app.users.FindById(r.Context(), userId)
			if err != nil {
				app.clientError(w, http.StatusNotFound, "no user found in database", err)
			}

			cl := data.ChecklistDisplay.ExtractChecklist()

			if err = app.checklist.Insert(r.Context(), *cl, user); err != nil {
				app.serverError(w, r, err)
			}

			app.renderPage(w, r, app.pageTemplates.CheckList, "Check List", &data)
			return
		}
		app.serverError(w, r, err)
		return
	}
	app.logger.Info("checklist found with id", "id", clDoc.ID.Hex())
	data.ChecklistDisplay.UpdateStatusFromChecklist(clDoc.Checklist)
	data.ChecklistDocumentId = clDoc.ID.Hex()
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
	app.logger.Info("User logged in", "user", user)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) handleNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	data := app.newTemplateData(r)
	app.logger.Info("Page Not Found")
	app.renderPage(w, r, app.pageTemplates.NotFound, "Page Not Found", &data)
}

func (app *application) postChecklist(w http.ResponseWriter, r *http.Request) {
	// init form struct
	var form models.ChecklistForm
	// decode form
	if err := app.decodeForm(r, &form); err != nil {
		app.clientError(w, http.StatusBadRequest, "Error decoding form", err)
		return
	}

	// logged in user's id
	userIdStr := app.sessionManager.GetString(r.Context(), SessionUserID)

	// convert to ObjectID
	userId, err := primitive.ObjectIDFromHex(userIdStr)
	if err != nil {
		app.logger.Error("error converting userId", "userId", userIdStr, "err", err)
		app.clientError(w, http.StatusUnprocessableEntity, "error", err)
	}

	// get most recent active checklist document
	clDoc, err := app.checklist.GetRecentActiveChecklist(r.Context(), userId)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// convert form to checklist and bike description
	cl := form.ConvertFormToChecklist()
	desc := form.ConvertFormToBikeDescription()

	// submit and complete checklist
	if err := app.checklist.SubmitChecklist(r.Context(), clDoc.ID, cl, desc); err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/bike/"+clDoc.ID.Hex(), http.StatusSeeOther)
}

func (app *application) putChecklist(w http.ResponseWriter, r *http.Request) {
	// init form struct
	var form models.ChecklistForm
	// decode form
	if err := app.decodeForm(r, &form); err != nil {
		app.clientError(w, http.StatusBadRequest, "Error decoding form", err)
		return
	}
	app.logger.Info("form", "form", form)
	// logged in user's id
	userIdStr := app.sessionManager.GetString(r.Context(), SessionUserID)
	// convert to ObjectID
	userId, err := primitive.ObjectIDFromHex(userIdStr)
	if err != nil {
		app.logger.Error("error converting userId", "userId", userIdStr, "err", err)
		app.clientError(w, http.StatusUnprocessableEntity, "error", err)
	}
	// get most recent active checklist document
	clDoc, err := app.checklist.GetRecentActiveChecklist(r.Context(), userId)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	// convert form to checklist
	cl := form.ConvertFormToChecklist()

	// submit and complete checklist
	if err := app.checklist.Update(r.Context(), clDoc.ID, cl); err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *application) getBikeDisplay(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	docId := chi.URLParam(r, "slug")
	app.logger.Info("docId", "docId", docId)
	docObjId, _ := primitive.ObjectIDFromHex(docId)
	clDoc, _ := app.checklist.Get(r.Context(), docObjId)

	data.ChecklistDisplay.BikeDescription = clDoc.Description
	data.ChecklistDisplay.UpdateStatusFromChecklist(clDoc.Checklist)
	app.renderPage(w, r, app.pageTemplates.BikeDisplay, "Bike Display", &data)
}
