package ex11filemanager

import (
	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nuiforms/ui"
)

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

	currentFilePanelIndex int
	columnResizing        bool // Prevent recursive column resizing calls

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
	//c.topPanel.SetMinHeight(50)
	c.topPanel.SetYExpandable(false)
	c.topPanel.SetAllowScroll(false, false)
	c.AddWidgetOnGrid(c.topPanel, 0, 0)

	c.contentPanel = ui.NewPanel()
	c.AddWidgetOnGrid(c.contentPanel, 1, 0)
	c.filePanels = make([]*FilePanel, 0)
	panel1 := NewFilePanel()
	panel1.SetName("Panel 1")
	panel1.SetOnFocused(c.panel1Focused)
	panel1.SetOnColumnResize(c.onPanel1ColumnResize)
	c.filePanels = append(c.filePanels, panel1)
	panel2 := NewFilePanel()
	panel2.SetName("Panel 2")
	panel2.SetOnFocused(c.panel2Focused)
	panel2.SetOnColumnResize(c.onPanel2ColumnResize)
	c.filePanels = append(c.filePanels, panel2)
	c.contentPanel.AddWidgetOnGrid(panel1, 0, 0)
	c.contentPanel.AddWidgetOnGrid(panel2, 0, 1)

	c.bottomPanel = ui.NewPanel()
	c.cmdLine = ui.NewTextBox()
	c.cmdLine.SetName("CommandLine")
	c.cmdLine.SetEmptyText("Enter command here...")
	c.cmdLine.SetOnTextBoxKeyDown(func() {
		key := ui.CurrentEvent().Parameter.(*ui.EventTextboxKeyDown).Key
		mods := ui.CurrentEvent().Parameter.(*ui.EventTextboxKeyDown).Mods
		if key == nuikey.KeyArrowUp || key == nuikey.KeyArrowDown {
			c.filePanels[c.currentFilePanelIndex].Focus()
			c.filePanels[c.currentFilePanelIndex].ProcessKeyDown(key, mods)
			ui.CurrentEvent().Parameter.(*ui.EventTextboxKeyDown).Processed = true
			return
		}
		return
	})
	c.bottomPanel.AddWidgetOnGrid(c.cmdLine, 0, 0)
	c.AddWidgetOnGrid(c.bottomPanel, 2, 0)

	ui.MainForm.SetOnGlobalKeyDown(c.onKeyDown)

	c.currentFilePanelIndex = -1
	c.selectFilePanel(0)
	c.filePanels[c.currentFilePanelIndex].Focus()

	return &c
}

func (c *MainWidget) onPanel1ColumnResize(col int, newWidth int) {
	if c.columnResizing {
		return
	}
	c.columnResizing = true
	c.filePanels[1].SetColumnWidth(col, newWidth)
	c.columnResizing = false
}

func (c *MainWidget) onPanel2ColumnResize(col int, newWidth int) {
	if c.columnResizing {
		return
	}
	c.columnResizing = true
	c.filePanels[0].SetColumnWidth(col, newWidth)
	c.columnResizing = false
}

func (c *MainWidget) onKeyDown(key nuikey.Key, mods nuikey.KeyModifiers) bool {
	if key == nuikey.KeyTab {
		if c.currentFilePanelIndex == 0 {
			c.selectFilePanel(1)
		} else {
			c.selectFilePanel(0)
		}
		c.filePanels[c.currentFilePanelIndex].Focus()
		return true
	}

	if key == nuikey.KeyArrowRight && !c.cmdLine.IsFocused() {
		c.cmdLine.Focus()
		c.cmdLine.MoveCursorToEnd()
		c.cmdLine.SelectAllText()
		return true
	}

	if key == nuikey.KeyEsc {
		c.cmdLine.SetText("")
		if c.cmdLine.IsFocused() {
			c.filePanels[c.currentFilePanelIndex].Select()
		}
	}
	return false
}

func (c *MainWidget) panel1Focused() {
	c.selectFilePanel(0)
}

func (c *MainWidget) panel2Focused() {
	c.selectFilePanel(1)
}

func (c *MainWidget) selectFilePanel(index int) {
	if index < 0 || index >= len(c.filePanels) {
		return
	}
	if c.currentFilePanelIndex == index {
		return
	}
	c.currentFilePanelIndex = index
	if c.currentFilePanelIndex == 0 {
		c.filePanels[0].Select()
		c.filePanels[1].Unselect()
	} else {
		c.filePanels[0].Unselect()
		c.filePanels[1].Select()
	}
}
