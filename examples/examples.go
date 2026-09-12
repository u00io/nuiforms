package examples

import (
	"github.com/u00io/nuiforms/examples/ex01base"
	"github.com/u00io/nuiforms/ui"
)

func Run() {
	{
		form := ui.NewForm()
		form.SetTitle("Examples")
		form.SetSize(800, 600)

		addButton := func(text string, onClick func()) {
			btn := ui.NewButton(text)
			btn.SetOnClick(onClick)

			form.Panel().AddWidget(btn, form.Panel().NextGridRow(), 0)
		}

		addButton("Example 01 - Base Form", func() {
			form.Panel().RemoveAllWidgets()
			ex01base.Run(form)
		})

		form.Panel().AddWidget(ui.NewVSpacer(), form.Panel().NextGridRow(), 0)
		form.Show()
		form.Exec()
	}
}
