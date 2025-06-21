package examples

import (
	"github.com/u00io/nuiforms/examples/ex01base"
	"github.com/u00io/nuiforms/examples/ex02widget"
	"github.com/u00io/nuiforms/ui"
)

func Run() {
	{
		form := ui.NewForm()
		form.SetTitle("Examples")
		form.SetSize(800, 600)
		{
			btnEx01 := ui.NewButton()
			btnEx01.SetText("Example 01")
			btnEx01.SetOnButtonClick(func(btn *ui.Button) {
				form.Panel().RemoveAllWidgets()
				ex01base.Run(form)
			})
			form.Panel().AddWidgetOnGrid(btnEx01, 0, 0)
		}
		{
			btnEx02 := ui.NewButton()
			btnEx02.SetText("Example 02")
			btnEx02.SetOnButtonClick(func(btn *ui.Button) {
				form.Panel().RemoveAllWidgets()
				ex02widget.Run(form)
			})
			form.Panel().AddWidgetOnGrid(btnEx02, 0, 1)
		}

		form.Panel().AddWidgetOnGrid(ui.NewVSpacer(), 0, 10)
		form.Exec()
	}
}
