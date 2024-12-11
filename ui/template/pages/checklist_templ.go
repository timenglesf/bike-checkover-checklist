// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.793
package pages

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"github.com/timenglesf/bike-checkover-checklist/internal/shared"
	"github.com/timenglesf/bike-checkover-checklist/ui/template/component"
)

func CheckList(d *shared.TemplateData) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<h1 class=\"mb-8 text-3xl text-center\">Bike Checklist</h1><form><ul class=\"flex flex-col gap-4 mx-2\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = component.ChecklistItem(d.Checklist.BrakePad.Id, d.Checklist.BrakePad.Status, d.Checklist.BrakePad.Name, d.Checklist.BrakePad.Description).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = component.ChecklistItem(d.Checklist.Chain.Id, d.Checklist.Chain.Status, d.Checklist.Chain.Name, d.Checklist.Chain.Description).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = component.ChecklistItem(d.Checklist.Tires.Id, d.Checklist.Tires.Status, d.Checklist.Tires.Name, d.Checklist.Tires.Description).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = component.ChecklistItem(d.Checklist.Cassette.Id, d.Checklist.Cassette.Status, d.Checklist.Cassette.Name, d.Checklist.Cassette.Description).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = component.ChecklistItem(d.Checklist.CablesHousing.Id, d.Checklist.CablesHousing.Status, d.Checklist.CablesHousing.Name, d.Checklist.CablesHousing.Description).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = component.ChecklistItem(d.Checklist.Tubes.Id, d.Checklist.Tubes.Status, d.Checklist.Tubes.Name, d.Checklist.Tubes.Description).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = component.ChecklistItem(d.Checklist.ChainRing.Id, d.Checklist.ChainRing.Status, d.Checklist.ChainRing.Name, d.Checklist.ChainRing.Description).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = component.ChecklistItem(d.Checklist.FrontWheel.Id, d.Checklist.FrontWheel.Status, d.Checklist.FrontWheel.Name, d.Checklist.FrontWheel.Description).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = component.ChecklistItem(d.Checklist.PadFunction.Id, d.Checklist.PadFunction.Status, d.Checklist.PadFunction.Name, d.Checklist.PadFunction.Description).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = component.ChecklistItem(d.Checklist.Derailleur.Id, d.Checklist.Derailleur.Status, d.Checklist.Derailleur.Name, d.Checklist.Derailleur.Description).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = component.ChecklistItem(d.Checklist.RearWheel.Id, d.Checklist.RearWheel.Status, d.Checklist.RearWheel.Name, d.Checklist.RearWheel.Description).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = component.ChecklistItem(d.Checklist.RotorRim.Id, d.Checklist.RotorRim.Status, d.Checklist.RotorRim.Name, d.Checklist.RotorRim.Description).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = component.ChecklistItem(d.Checklist.Hanger.Id, d.Checklist.Hanger.Status, d.Checklist.Hanger.Name, d.Checklist.Hanger.Description).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = component.ChecklistItem(d.Checklist.Shifting.Id, d.Checklist.Shifting.Status, d.Checklist.Shifting.Name, d.Checklist.Shifting.Description).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</ul></form>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
