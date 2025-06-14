package ex04scrolling

import (
	"github.com/u00io/nuiforms/ui"
)

func Run() {
	form := ui.NewForm()
	form.SetTitle("Scrolling Example")
	form.SetSize(800, 600)

	scrollWidget := NewBigWidget()
	scrollWidget.SetPosition(0, 0)
	scrollWidget.SetSize(800, 600)
	scrollWidget.SetAnchors(true, true, true, true)
	form.Panel().AddWidget(scrollWidget)

	form.Exec()
}
