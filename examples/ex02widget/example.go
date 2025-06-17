package ex02widget

import (
	"github.com/u00io/nuiforms/ui"
)

func Run() {
	form := ui.NewForm()
	form.SetTitle("Hello, World!")

	panel := form.Panel()
	panel.SetPanelPadding(50)
	myWidget := NewMyWidget(ui.DefaultBackground)
	panel.AddWidgetOnGrid(myWidget, 0, 0)

	form.Exec()
}
