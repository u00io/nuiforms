package ex03masterdetail

import (
	"github.com/u00io/nuiforms/ui"
)

func Run() {
	form := ui.NewForm()
	form.SetTitle("Hello, World!")
	form.SetSize(800, 600)

	leftPanelWidth := 300

	var tabWidget *ui.Widget

	closeCurrentTab := func() {
		if tabWidget != nil {
			form.Panel().RemoveWidget(tabWidget)
			tabWidget = nil
		}
	}

	openTab1 := func() {
		closeCurrentTab()
		tabWidget = NewTab1Widget()
		tabWidget.SetPosition(300, 0)
		tabWidget.SetSize(form.Panel().Width()-leftPanelWidth, form.Panel().Height())
		tabWidget.SetAnchors(true, true, true, true)
		form.Panel().AddWidget(tabWidget)
	}

	openTab2 := func() {
		closeCurrentTab()
		tabWidget = NewTab2Widget()
		tabWidget.SetPosition(300, 0)
		tabWidget.SetSize(form.Panel().Width()-leftPanelWidth, form.Panel().Height())
		tabWidget.SetAnchors(true, true, true, true)
		form.Panel().AddWidget(tabWidget)
	}

	masterWidget := NewMasterWidget()
	masterWidget.SetPosition(0, 0)
	masterWidget.SetSize(leftPanelWidth, form.Panel().Height())
	masterWidget.SetAnchors(true, true, false, true)
	masterWidget.SetProp("open1", openTab1)
	masterWidget.SetProp("open2", openTab2)
	form.Panel().AddWidget(masterWidget)

	openTab1()

	form.Exec()
}
