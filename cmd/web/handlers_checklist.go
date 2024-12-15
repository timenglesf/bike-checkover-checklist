package main

import (
	"net/http"

	"github.com/timenglesf/bike-checkover-checklist/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (app *application) handleDisplayChecklist(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.renderPage(w, r, app.pageTemplates.CheckList, "Check List", &data)
}

func (app *application) postChecklist(w http.ResponseWriter, r *http.Request) {
	var form models.ChecklistForm
	if err := app.decodeForm(r, &form); err != nil {
		app.clientError(w, http.StatusBadRequest, "Error decoding form", err)
		return
	}

	userIdStr := app.sessionManager.GetString(r.Context(), SessionUserID)
	userId, err := primitive.ObjectIDFromHex(userIdStr)
	if err != nil {
		app.logger.Error("error converting userId", "userId", userIdStr, "err", err)
		app.clientError(w, http.StatusUnprocessableEntity, "error", err)
	}

	clDoc, err := app.checklist.GetRecentActiveChecklist(r.Context(), userId)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	cl := form.ConvertFormToChecklist()
	desc := form.ConvertFormToBikeDescription()

	if err := app.checklist.SubmitChecklist(r.Context(), clDoc.ID, cl, desc); err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/bike/"+clDoc.ID.Hex(), http.StatusSeeOther)
}

func (app *application) putChecklist(w http.ResponseWriter, r *http.Request) {
	var form models.ChecklistForm
	if err := app.decodeForm(r, &form); err != nil {
		app.clientError(w, http.StatusBadRequest, "Error decoding form", err)
		return
	}
	app.logger.Info("form", "form", form)

	userIdStr := app.sessionManager.GetString(r.Context(), SessionUserID)
	userId, err := primitive.ObjectIDFromHex(userIdStr)
	if err != nil {
		app.logger.Error("error converting userId", "userId", userIdStr, "err", err)
		app.clientError(w, http.StatusUnprocessableEntity, "error", err)
	}

	clDoc, err := app.checklist.GetRecentActiveChecklist(r.Context(), userId)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	cl := form.ConvertFormToChecklist()
	if err := app.checklist.Update(r.Context(), clDoc.ID, cl); err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *application) handleChecklistReset(w http.ResponseWriter, r *http.Request) {
	userIdStr := app.sessionManager.GetString(r.Context(), SessionUserID)
	if userIdStr == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	userId, err := primitive.ObjectIDFromHex(userIdStr)
	if err != nil {
		app.logger.Error("error converting userId", "userId", userIdStr, "err", err)
		app.clientError(w, http.StatusUnprocessableEntity, "error", err)
	}
	clDoc, err := app.checklist.GetRecentActiveChecklist(r.Context(), userId)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	if err = app.checklist.Reset(r.Context(), clDoc.ID); err != nil {
		app.serverError(w, r, err)
		return
	}
	data := app.newTemplateData(r)
	app.renderPage(w, r, app.pageTemplates.CheckList, "Check List", &data)
}
