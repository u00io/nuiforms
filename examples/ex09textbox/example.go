package ex09textbox

import "github.com/u00io/nuiforms/ui"

func Run(form *ui.Form) {
	form.SetTitle("Example 01 - Base Form")
	panel := form.Panel()
	btn1 := ui.NewButton()
	btn1.SetText("Click Me")
	panel.AddWidgetOnGrid(btn1, 0, 0)
	btn2 := ui.NewButton()
	btn2.SetText("Click Me Too")
	panel.AddWidgetOnGrid(btn2, 0, 1)
	txt1 := ui.NewTextBox()
	panel.AddWidgetOnGrid(txt1, 0, 2)
}
