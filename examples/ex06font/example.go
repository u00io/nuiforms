package ex06font

import (
	"github.com/u00io/nuiforms/ui"
)

func Run(form *ui.Form) {
	lbl := ui.NewLabel("Hello, World!---------------")
	panel := form.Panel()
	panel.AddWidgetOnGrid(lbl, 0, 0)
	panel.AddWidgetOnGrid(ui.NewHSpacer(), 1, 0)
	panel.AddWidgetOnGrid(ui.NewVSpacer(), 0, 1)
}
