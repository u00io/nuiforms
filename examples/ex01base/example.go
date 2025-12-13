package ex01base

import "github.com/u00io/nuiforms/ui"

func Run(form *ui.Form) {
	form.SetTitle("Example 01 - Base Form")
	panel := form.Panel()
	lbl := ui.NewLabel("init text")
	panel.AddWidgetOnGrid(lbl, 0, 0)

	btn := ui.NewButton("Click Me")
	btn.SetOnButtonClick(func() {
		lbl.SetText("Button clicked!")
		btn.SetText("Clicked")
	})
	panel.AddWidgetOnGrid(btn, 1, 0)
}
