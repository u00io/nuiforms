package examples

import (
	"github.com/u00io/nuiforms/examples/ex00gallery"
	"github.com/u00io/nuiforms/examples/ex01base"
	"github.com/u00io/nuiforms/examples/ex02messagebox"
	"github.com/u00io/nuiforms/examples/ex03customwidget"
	"github.com/u00io/nuiforms/examples/ex04dialog"
	"github.com/u00io/nuiforms/examples/ex05chart"
	"github.com/u00io/nuiforms/examples/ex06timechart"
	"github.com/u00io/nuiforms/ui"
)

func Run() {
	{
		form := ui.NewForm()
		form.SetTitle("Examples")
		form.SetSize(1000, 800)

		addButton := func(text string, newFormFunc func() *ui.Form) {
			btn := ui.NewButton(text)
			btn.SetOnClick(func() {
				newForm := newFormFunc()
				newForm.ShowModal(form)
			})
			form.Panel().AddWidget(form.Panel().NextGridRow(), 0, btn)
		}

		addButton("Example 00 - Gallery", ex00gallery.NewExampleForm)
		addButton("Example 01 - Base Form", ex01base.NewExampleForm)
		addButton("Example 02 - MessageBox", ex02messagebox.NewExampleForm)
		addButton("Example 03 - Custom Widget", ex03customwidget.NewExampleForm)
		addButton("Example 04 - Dialog", ex04dialog.NewExampleForm)
		addButton("Example 05 - Chart", ex05chart.NewExampleForm)
		addButton("Example 06 - TimeChart", ex06timechart.NewExampleForm)

		form.Panel().AddWidget(form.Panel().NextGridRow(), 0, ui.NewVSpacer())
		form.Panel().AddButton(form.Panel().NextGridRow(), 0, "Light Theme", func() {
			ui.ApplyLightTheme()
		})
		form.Panel().AddButton(form.Panel().NextGridRow(), 0, "Dark Theme", func() {
			ui.ApplyDarkTheme()
		})
		form.Show()
		form.Exec()
	}
}
