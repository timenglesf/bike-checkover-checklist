package main

import (
	"net/http"
)

func (app *application) handleDisplayMainPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.renderPage(w, r, app.pageTemplates.CheckList, "Check List", &data)
}

func (app *application) handleNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	data := app.newTemplateData(r)
	app.logger.Info("Page Not Found")
	app.renderPage(w, r, app.pageTemplates.NotFound, "Page Not Found", &data)
}
