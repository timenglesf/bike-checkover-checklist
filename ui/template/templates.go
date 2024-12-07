package template

import (
	"github.com/a-h/templ"
	"github.com/timenglesf/bike-checkover-checklist/internal/shared"
	"github.com/timenglesf/bike-checkover-checklist/ui/template/pages"
	"github.com/timenglesf/bike-checkover-checklist/ui/template/partials"
)

type Pages struct {
	Base      func(title string, page templ.Component, data *shared.TemplateData) templ.Component
	NotFound  func(data *shared.TemplateData) templ.Component
	CheckList func(data *shared.TemplateData) templ.Component
	UserLogin func(data *shared.TemplateData) templ.Component
	// LogIn     func(data *shared.TemplateData) templ.Component
}

type Partials struct {
	Header func(data *shared.TemplateData) templ.Component
	Footer func(data *shared.TemplateData) templ.Component
}

func CreatePages() *Pages {
	return &Pages{
		Base:      Base,
		NotFound:  pages.NotFound,
		CheckList: pages.CheckList,
		UserLogin: pages.UserLogin,
	}
}

func CreatePartials() *Partials {
	return &Partials{
		Header: partials.PageHeader,
		Footer: partials.PageFooter,
	}
}
