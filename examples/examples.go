package examples

import (
	"github.com/u00io/nuiforms/examples/ex01base"
	"github.com/u00io/nuiforms/examples/ex02messagebox"
	"github.com/u00io/nuiforms/examples/ex03customwidget"
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

		addButton("Example 01 - Base Form", ex01base.NewExampleForm)
		addButton("Example 02 - MessageBox", ex02messagebox.NewExample)
		addButton("Example 03 - Custom Widget", ex03customwidget.NewExample)

		form.Panel().AddWidget(form.Panel().NextGridRow(), 0, ui.NewVSpacer())
		form.Show()
		form.Exec()
	}
}

func RunSample() {
	form := ex03customwidget.NewExample()
	form.Show()
	form.Exec()
}
