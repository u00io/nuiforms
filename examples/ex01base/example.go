package ex01base

import "github.com/u00io/nuiforms/ui"

func Run(form *ui.Form) {
	form.SetTitle("Example 01 - Base Form")
	panel := form.Panel()
	lbl := ui.NewLabel("Hello, World!")
	lbl.SetMinSize(300, 500)
	lbl.SetMaxSize(500, 500)
	lbl.SetName("LBL")
	panel.AddWidgetOnGrid(lbl, 0, 0)
}
