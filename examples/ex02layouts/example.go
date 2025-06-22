package ex02layouts

import (
	"github.com/u00io/nuiforms/ui"
)

func Run(form *ui.Form) {
	form.SetTitle("Layouts")
	panel := form.Panel()

	{
		panelTop := ui.NewPanel()
		panel.AddWidgetOnGrid(panelTop, 0, 0)
		myWidget := NewMyWidget("TOP")
		panelTop.AddWidgetOnGrid(myWidget, 0, 0)
		panelTop.SetMaxHeight(50)
		panelTop.SetPanelPadding(0)
	}

	{
		panelMiddle := ui.NewPanel()
		panelMiddle.SetPanelPadding(0)
		panel.AddWidgetOnGrid(panelMiddle, 0, 1)

		panelMiddleLeft := ui.NewPanel()
		panelMiddleLeft.SetPanelPadding(0)
		panelMiddle.AddWidgetOnGrid(panelMiddleLeft, 0, 0)
		leftWidget := NewMyWidget("LEFT")
		panelMiddleLeft.AddWidgetOnGrid(leftWidget, 0, 0)
		panelMiddleLeft.SetMaxWidth(200)

		panelMiddleCenter := ui.NewPanel()
		panelMiddleCenter.SetPanelPadding(0)
		panelMiddle.AddWidgetOnGrid(panelMiddleCenter, 1, 0)
		centerWidget := NewMyWidget("CENTER")
		panelMiddleCenter.AddWidgetOnGrid(centerWidget, 0, 0)

		panelMiddleRight := ui.NewPanel()
		panelMiddleRight.SetPanelPadding(0)
		panelMiddle.AddWidgetOnGrid(panelMiddleRight, 2, 0)
		panelMiddleRight.SetMaxWidth(100)
		rightWidget := NewMyWidget("RIGHT")
		panelMiddleRight.AddWidgetOnGrid(rightWidget, 0, 0)
	}

	{
		panelBottom := ui.NewPanel()
		panel.AddWidgetOnGrid(panelBottom, 0, 2)
		panelBottom.SetPanelPadding(0)
		panelBottom.SetMaxHeight(50)
		myWidget := NewMyWidget("BOTTOM")
		panelBottom.AddWidgetOnGrid(myWidget, 0, 0)
	}
}
