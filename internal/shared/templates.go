package shared

import (
	"time"

	"github.com/timenglesf/bike-checkover-checklist/internal/checklist"
)

const (
	DateLayout     = "January 2, 2006"
	DateTimeLayout = "January 2, 2006 at 3:04 PM"
)

type FlashMessage struct {
	Message string
	Type    string
}

type TemplateData struct {
	// IsAuthenticated bool
	// IsAdmin         bool
	CheckList *checklist.CheckList
	Flash     *FlashMessage
	Date      time.Time
	//	CSRFToken   string
	CurrentYear int
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
