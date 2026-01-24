package examples

import (
	"github.com/u00io/nuiforms/examples/ex00layout"
	"github.com/u00io/nuiforms/examples/ex01base"
	"github.com/u00io/nuiforms/examples/ex02label"
	"github.com/u00io/nuiforms/examples/ex03masterdetail"
	"github.com/u00io/nuiforms/examples/ex06font"
	"github.com/u00io/nuiforms/examples/ex09textbox"
	"github.com/u00io/nuiforms/examples/ex10tabwidget"
	"github.com/u00io/nuiforms/examples/ex11filemanager"
	"github.com/u00io/nuiforms/examples/ex12cards"
	"github.com/u00io/nuiforms/examples/ex13table"
	"github.com/u00io/nuiforms/examples/ex14popup"
	"github.com/u00io/nuiforms/examples/ex15dialog"
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

			form.Panel().AddWidgetOnGrid(btn, form.Panel().NextGridRow(), 0)
		}

		addButton("Example 00 - Layouts", func() {
			form.Panel().RemoveAllWidgets()
			ex00layout.Run(form)
		})

		addButton("Example 01 - Base Form", func() {
			form.Panel().RemoveAllWidgets()
			ex01base.Run(form)
		})

		addButton("Example 02 - Label", func() {
			form.Panel().RemoveAllWidgets()
			ex02label.Run(form)
		})

		addButton("Example 03 - Master Detail", func() {
			form.Panel().RemoveAllWidgets()
			ex03masterdetail.Run(form)
		})

		addButton("Example 06 - Font", func() {
			form.Panel().RemoveAllWidgets()
			ex06font.Run(form)
		})

		addButton("Example 09 - TextBox", func() {
			form.Panel().RemoveAllWidgets()
			ex09textbox.Run(form)
		})

		addButton("Example 10 - Tab Widget", func() {
			form.Panel().RemoveAllWidgets()
			ex10tabwidget.Run(form)
		})

		addButton("File Manager Example", func() {
			form.Panel().RemoveAllWidgets()
			ex11filemanager.Run(form)
		})

		addButton("Example 12 - Cards", func() {
			form.Panel().RemoveAllWidgets()
			ex12cards.Run(form)
		})

		addButton("Example 13 - Table Widget", func() {
			form.Panel().RemoveAllWidgets()
			ex13table.Run(form)
		})

		addButton("Example 14 - Popup", func() {
			form.Panel().RemoveAllWidgets()
			ex14popup.Run(form)
		})

		addButton("Example 15 - Dialog", func() {
			form.Panel().RemoveAllWidgets()
			ex15dialog.Run(form)
		})

		addButton("Example 16 - Light Theme", func() {
			ui.ApplyLightTheme()
		})

		form.Panel().AddWidgetOnGrid(ui.NewVSpacer(), form.Panel().NextGridRow(), 0)
		form.Exec()
	}
}
