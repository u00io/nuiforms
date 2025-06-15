package ex02widget

import (
	"image/color"

	"github.com/u00io/nuiforms/ui"
)

func Run() {
	form := ui.NewForm()
	form.SetTitle("Hello, World!")
	form.SetSize(800, 600)

	{
		myWidget := NewMyWidget(color.RGBA{0, 100, 50, 255})
		myWidget.SetPosition(8, 8)
		myWidget.SetSize(300, 480)
		myWidget.SetAnchors(true, true, false, true)
		form.Panel().AddWidget(myWidget)
	}
	{
		myWidget := NewMyWidget(color.RGBA{0, 50, 100, 255})
		myWidget.SetPosition(320, 8)
		myWidget.SetSize(470, 480)
		myWidget.SetAnchors(true, true, true, true)
		form.Panel().AddWidget(myWidget)
	}

	{
		myWidget := NewMyWidget(color.RGBA{100, 0, 50, 255})
		myWidget.SetPosition(8, 500)
		myWidget.SetSize(780, 90)
		myWidget.SetAnchors(true, false, true, true)
		form.Panel().AddWidget(myWidget)

		txt := ui.NewTextBlock()
		txt.SetPosition(150, 0)
		txt.SetSize(100, 80)
		myWidget.AddWidget(txt)

		btn := ui.NewButton()
		btn.SetPosition(450, 30)
		myWidget.AddWidget(btn)
	}

	form.Exec()
}
