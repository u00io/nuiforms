package ex09textbox

import "github.com/u00io/nuiforms/ui"

func Run(form *ui.Form) {
	form.SetTitle("Example 01 - Base Form")
	panel := form.Panel()
	panel1 := form.Panel()
	lbl := ui.NewLabel("This is a label")
	panel1.AddWidgetOnGrid(lbl, 0, 0)
	txt1 := ui.NewTextBox()
	panel1.AddWidgetOnGrid(txt1, 0, 1)
	panel.AddWidgetOnGrid(panel1, 0, 0)
}
