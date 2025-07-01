package ex11filemanager

import "github.com/u00io/nuiforms/ui"

func Run(form *ui.Form) {
	mainWidget := NewMainWidget()
	form.Panel().AddWidgetOnGrid(mainWidget, 0, 0)
}

type MainWidget struct {
	ui.Widget
	topPanel     *ui.Panel
	contentPanel *ui.Panel
	filePanels   []*FilePanel
	bottomPanel  *ui.Panel

	cmdLine *ui.TextBox
}

func NewMainWidget() *MainWidget {
	var c MainWidget
	c.InitWidget()
	c.SetTypeName("MainWidget")
	c.SetXExpandable(true)
	c.SetYExpandable(true)

	c.topPanel = ui.NewPanel()
	c.topPanel.AddWidgetOnGrid(ui.NewLabel("File Manager Example"), 0, 0)
	c.AddWidgetOnGrid(c.topPanel, 0, 0)

	c.contentPanel = ui.NewPanel()
	c.AddWidgetOnGrid(c.contentPanel, 0, 1)
	c.filePanels = make([]*FilePanel, 0)
	panel1 := NewFilePanel()
	panel1.SetName("Panel 1")
	c.filePanels = append(c.filePanels, panel1)
	panel2 := NewFilePanel()
	panel2.SetName("Panel 2")
	c.filePanels = append(c.filePanels, panel2)
	c.contentPanel.AddWidgetOnGrid(panel1, 0, 0)
	c.contentPanel.AddWidgetOnGrid(panel2, 1, 0)

	c.bottomPanel = ui.NewPanel()
	c.cmdLine = ui.NewTextBox()
	c.cmdLine.SetName("CommandLine")
	c.cmdLine.SetEmptyText("Enter command here...")
	c.bottomPanel.AddWidgetOnGrid(c.cmdLine, 0, 0)
	c.AddWidgetOnGrid(c.bottomPanel, 0, 2)

	return &c
}
