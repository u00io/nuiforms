package ex11filemanager

import (
	"fmt"
	"os"

	"github.com/u00io/nuiforms/ui"
)

type FilePanel struct {
	ui.Widget

	topPanel      *ui.Panel
	topPanelLabel *ui.Label

	contentPanel *ui.Panel
	fileList     *ui.Table
	bottomPanel  *ui.Panel
}

func NewFilePanel() *FilePanel {
	var c FilePanel
	c.InitWidget()
	c.SetTypeName("FilePanel")
	c.SetXExpandable(true)
	c.SetYExpandable(true)

	c.topPanel = ui.NewPanel()
	c.topPanelLabel = ui.NewLabel("File Panel Top")

	c.topPanel.AddWidgetOnGrid(c.topPanelLabel, 0, 0)
	c.AddWidgetOnGrid(c.topPanel, 0, 0)

	c.contentPanel = ui.NewPanel()
	c.AddWidgetOnGrid(c.contentPanel, 0, 1)
	c.fileList = ui.NewTable()
	c.fileList.SetName("FileList")
	c.fileList.SetColumnCount(3)
	c.fileList.SetColumnName(0, "Name")
	c.fileList.SetColumnName(1, "Size")
	c.fileList.SetColumnName(2, "Modified")
	c.fileList.SetAllowScroll(false, true)
	c.fileList.SetSelectingCell(false)
	c.contentPanel.AddWidget(c.fileList)

	c.bottomPanel = ui.NewPanel()
	c.bottomPanel.AddWidgetOnGrid(ui.NewLabel("File panel bottom"), 0, 0)
	c.AddWidgetOnGrid(c.bottomPanel, 0, 2)

	c.loadDirectory("D:/")

	return &c
}

func (c *FilePanel) Select() {
	c.topPanelLabel.SetText("Selected")
	c.fileList.SetShowSelection(true)
	c.fileList.Focus()
}

func (c *FilePanel) Unselect() {
	c.topPanelLabel.SetText("Unselected")
	c.fileList.SetShowSelection(false)
}

func (c *FilePanel) SetOnFocused(onFocused func()) {
	c.fileList.SetOnFocused(onFocused)
}

func (c *FilePanel) loadDirectory(path string) {
	// Get file list from the specified path
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return
	}

	c.fileList.SetRowCount(len(dirEntries))
	for i, entry := range dirEntries {
		fileInfo, err := entry.Info()
		if err != nil {
			continue
		}
		c.fileList.SetCellText(0, i, fileInfo.Name())
		if fileInfo.IsDir() {
			c.fileList.SetCellText(1, i, "<DIR>")
		} else {
			c.fileList.SetCellText(1, i, fmt.Sprint(fileInfo.Size()))
		}
		c.fileList.SetCellText(2, i, fileInfo.ModTime().Format("2006-01-02 15:04:05"))
	}
}
