package main

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/timenglesf/bike-checkover-checklist/internal/shared"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (app *application) handleDisplayMainPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	if !data.IsAuthenticated {
		app.renderPage(w, r, app.pageTemplates.UserLogin, "Bike Intake Checklist: Login", &data)
		return
	}

	userIdStr := app.sessionManager.GetString(r.Context(), SessionUserID)
	app.logger.Info("user access", "userId", userIdStr)

	userId, err := primitive.ObjectIDFromHex(userIdStr)
	if err != nil {
		app.logger.Error("error converting userId", "userId", userIdStr, "err", err)
		app.clientError(w, http.StatusUnprocessableEntity, "error", err)
	}

	clDoc, err := app.checklist.GetRecentActiveChecklist(r.Context(), userId)
	if err != nil {
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
	app.renderPage(w, r, app.pageTemplates.CheckList, "Bike Intake Checklist", &data)
}

func (app *application) handleNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	data := app.newTemplateData(r)
	app.logger.Info("Page Not Found")
	app.renderPage(w, r, app.pageTemplates.NotFound, "Page Not Found", &data)
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

func (app *application) handleDisplayUserHistory(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	if !data.IsAuthenticated {
		app.renderPage(w, r, app.pageTemplates.UserLogin, "Bike Intake Checklist: Login", &data)
		return
	}
	userIdStr := app.sessionManager.GetString(r.Context(), SessionUserID)

	userObjId, err := primitive.ObjectIDFromHex(userIdStr)
	if err != nil {
		app.logger.Error("error converting userId", "userId", userIdStr, "err", err)
		app.clientError(w, http.StatusUnprocessableEntity, "error", err)
	}

	clDocs, err := app.checklist.GetUserChecklists(r.Context(), userObjId)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	clList := make([]shared.ChecklistListEntry, len(clDocs))
	for i, clDoc := range clDocs {
		clList[i] = shared.ConvertChecklistToChecklistListEntry(clDoc)
	}
	data.ChecklistList = clList

	app.renderPage(w, r, app.pageTemplates.UserHistory, "User History", &data)
}
