package examples

import (
	"github.com/u00io/nuiforms/examples/ex01base"
	"github.com/u00io/nuiforms/examples/ex02layouts"
	"github.com/u00io/nuiforms/examples/ex03masterdetail"
	"github.com/u00io/nuiforms/examples/ex06font"
	"github.com/u00io/nuiforms/examples/ex09textbox"
	"github.com/u00io/nuiforms/examples/ex10tabwidget"
	"github.com/u00io/nuiforms/ui"
)

func Run() {
	{
		form := ui.NewForm()
		form.SetTitle("Examples")
		form.SetSize(800, 600)

		addButton := func(text string, onClick func(btn *ui.Button)) {
			btn := ui.NewButton(text)
			btn.SetOnButtonClick(onClick)
			form.Panel().AddWidgetOnGrid(btn, 0, form.Panel().NextGridY())
		}

		addButton("Example 01 - Base Form", func(btn *ui.Button) {
			form.Panel().RemoveAllWidgets()
			ex01base.Run(form)
		})

		addButton("Example 02 - Layouts", func(btn *ui.Button) {
			form.Panel().RemoveAllWidgets()
			ex02layouts.Run(form)
		})

		addButton("Example 03 - Master Detail", func(btn *ui.Button) {
			form.Panel().RemoveAllWidgets()
			ex03masterdetail.Run(form)
		})

		addButton("Example 06 - Font", func(btn *ui.Button) {
			form.Panel().RemoveAllWidgets()
			ex06font.Run(form)
		})

		addButton("Example 09 - TextBox", func(btn *ui.Button) {
			form.Panel().RemoveAllWidgets()
			ex09textbox.Run(form)
		})

		addButton("Example 10 - Tab Widget", func(btn *ui.Button) {
			form.Panel().RemoveAllWidgets()
			ex10tabwidget.Run(form)
		})

		form.Panel().AddWidgetOnGrid(ui.NewVSpacer(), 0, form.Panel().NextGridY())
		form.Exec()
	}
}
