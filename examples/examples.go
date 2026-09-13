package examples

import (
	"github.com/u00io/nuiforms/examples/ex01base"
	"github.com/u00io/nuiforms/examples/ex02messagebox"
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
			form.Panel().AddWidget(btn, form.Panel().NextGridRow(), 0)
		}

		addButton("Example 01 - Base Form", ex01base.NewExampleForm)
		addButton("Example 02 - MessageBox", ex02messagebox.NewExample)

		form.Panel().AddWidget(ui.NewVSpacer(), form.Panel().NextGridRow(), 0)
		form.Show()
		form.Exec()
	}
}
