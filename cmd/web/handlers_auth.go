package main

import (
	"errors"
	"net/http"

	"github.com/timenglesf/bike-checkover-checklist/internal/shared"
	"github.com/timenglesf/bike-checkover-checklist/internal/validator"
	"golang.org/x/crypto/bcrypt"

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

/////////////////////////////////
////////// ADMIN HANDLERS ///////
/////////////////////////////////

func (app *application) handlePostAdminLogin(w http.ResponseWriter, r *http.Request) {
	var form shared.AdminLoginForm

	err := app.decodeForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, "error decoding form", err)
		return
	}

	app.logger.Info("login attempt", "username", form.Username, "password", form.Password)

	// Get admin from database
	admin, err := app.admin.FindByUsername(r.Context(), form.Username)
	data := app.newTemplateData(r)
	data.AdminFormData = form
	if err != nil {
		app.logger.Error("error retrieving admin", "err", err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			data.Flash.Message = "Invalid username or password"
			data.Flash.Type = shared.FlashError
			app.renderPage(w, r, app.pageTemplates.AdminLogin, "Admin Login", &data)
			return
		}
		app.serverError(w, r, err)
	}

	// Validate password
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(form.Password))
	if err != nil {
		app.logger.Warn("invalid password", "username", admin.Username, "ip", r.RemoteAddr)
		app.logger.Error("error comparing password", "err", err)
		data.Flash.Message = "Invalid username or password"
		data.Flash.Type = shared.FlashError
		app.renderPage(w, r, app.pageTemplates.AdminLogin, "Admin Login", &data)
		return
	}

	app.sessionManager.Put(r.Context(), SessionIsAdmin, true)
	app.logger.Info("admin logged in", "username", admin.Username, "ip", r.RemoteAddr)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
