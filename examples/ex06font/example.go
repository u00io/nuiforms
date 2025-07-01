package ex06font

import (
	"github.com/u00io/nuiforms/ui"
)

func Run(form *ui.Form) {
	btn := ui.NewButton("abc23456789def")
	panel := form.Panel()
	panel.AddWidgetOnGrid(ui.NewHSpacer(), 0, 0)
	panel.AddWidgetOnGrid(btn, 1, 0)
	panel.AddWidgetOnGrid(ui.NewHSpacer(), 2, 0)
	panel.AddWidgetOnGrid(ui.NewVSpacer(), 0, 1)
}
