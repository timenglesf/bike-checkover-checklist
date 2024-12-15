// This file contains the shared template data that is used in the templ templates and handlers.
// It also contains some helper functions to format dates and times.
package shared

import (
	"time"

	"github.com/timenglesf/bike-checkover-checklist/internal/models"
	"github.com/timenglesf/bike-checkover-checklist/internal/validator"
)

const (
	DateLayout     = "January 2, 2006"
	DateTimeLayout = "January 2, 2006 at 3:04 PM"
)

type FlashType string

var (
	FlashSuccess FlashType = "success_alert"
	FlashError   FlashType = "error_alert"
	FlashWarning FlashType = "warning_alert"
)

type FlashMessage struct {
	Message string
	Type    FlashType
}

type ChecklistListEntry struct {
	Id          string
	CreatedAt   time.Time
	Checklist   models.Checklist
	Description models.BikeDescription
}

func ConvertChecklistToChecklistListEntry(clDoc models.ChecklistDocument) ChecklistListEntry {
	return ChecklistListEntry{
		Id:          clDoc.ID.Hex(),
		CreatedAt:   clDoc.CreatedAt.Time(),
		Description: clDoc.Description,
	}
}

type TemplateData struct {
	// IsAdmin         bool
	Flash               *FlashMessage
	PinForm             PinForm
	Date                time.Time
	IsAuthenticated     bool
	User                *models.User
	ChecklistDisplay    *models.ChecklistDisplay
	ChecklistDocumentId string
	ChecklistList       []ChecklistListEntry
	//	CSRFToken   string
	CurrentYear int
}

// Form Data

type PinForm struct {
	validator.Validator `form:"-"`
	Pin                 string `form:"pin"`
}

func HumanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(DateLayout)
}

func HumanDateTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(DateTimeLayout)
}
