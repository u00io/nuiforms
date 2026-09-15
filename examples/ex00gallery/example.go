package ex00gallery

import "github.com/u00io/nuiforms/ui"

type Example struct {
	ui.Widget

	tabWidget *ui.TabWidget
}

func NewExample() *Example {
	var c Example
	c.InitWidget()

	c.tabWidget = ui.NewTabWidget()
	c.AddWidget(0, 0, c.tabWidget)

	c.tabWidget.AddPage("Label", NewExamplePageLabel())
	c.tabWidget.AddPage("TextEdit", NewExamplePageTextBox())
	c.tabWidget.AddPage("ContextMenu", NewExamplePageContextMenu())
	c.tabWidget.AddPage("Button", NewExamplePageButton())
	c.tabWidget.AddPage("CheckBox", NewExamplePageCheckbox())

	return &c
}

func NewExampleForm() *ui.Form {
	form := ui.NewForm()
	form.SetSize(1100, 900)
	form.Panel().AddWidget(0, 0, NewExample())
	return form
}
