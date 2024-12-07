package shared

import (
	"time"

	"github.com/timenglesf/bike-checkover-checklist/internal/checklist"
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

type TemplateData struct {
	// IsAuthenticated bool
	// IsAdmin         bool
	CheckList       *checklist.CheckList
	Flash           *FlashMessage
	PinForm         PinForm
	Date            time.Time
	IsAuthenticated bool
	User            models.User
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
