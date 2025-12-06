package ex09textbox

import "github.com/u00io/nuiforms/ui"

func Run(form *ui.Form) {
	form.SetTitle("Example 01 - Base Form")
	panel := form.Panel()
	txt1 := ui.NewTextBox()
	panel.AddWidgetOnGrid(txt1, 2, 0)
}
