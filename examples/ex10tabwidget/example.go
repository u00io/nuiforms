package ex10tabwidget

import "github.com/u00io/nuiforms/ui"

func Run(form *ui.Form) {
	form.SetTitle("Example 10 - Tab Widget")

	tabWidget := ui.NewTabWidget()
	tabWidget.AddPage("Page 1", ui.NewLabel("Content of Page 1"))
	tabWidget.AddPage("Page 2", ui.NewLabel("Content of Page 2"))
	tabWidget.AddPage("Page 3", ui.NewLabel("Content of Page 3"))
	tabWidget.AddPage("Page 4", ui.NewLabel("Content of Page 4"))
	tabWidget.AddPage("Page 5", ui.NewLabel("Content of Page 5"))
	form.Panel().AddWidgetOnGrid(tabWidget, 0, 0)
}
