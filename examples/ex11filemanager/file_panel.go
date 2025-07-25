package ex11filemanager

import (
	"fmt"

	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nuiforms/ui"
)

type FilePanel struct {
	ui.Widget

	currentEntry *Entry

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
	c.fileList.SetOnKeyDown(c.fileListKeyDown)
	c.contentPanel.AddWidget(c.fileList)

	c.bottomPanel = ui.NewPanel()
	c.bottomPanel.AddWidgetOnGrid(ui.NewLabel("File panel bottom"), 0, 0)
	c.AddWidgetOnGrid(c.bottomPanel, 0, 2)

	rootEntries := readRootEntries()

	c.currentEntry = rootEntries[0]

	c.loadDirectory(c.currentEntry)

	return &c
}

func (c *FilePanel) fileListKeyDown(key nuikey.Key, mods nuikey.KeyModifiers) bool {
	if key == nuikey.KeyEnter {
		currentRowIndex := c.fileList.CurrentRow()
		if currentRowIndex < 0 || currentRowIndex >= c.fileList.RowCount() {
			return false
		}
		fileName := c.fileList.GetCellText(0, currentRowIndex)
		newEntry := c.currentEntry.CreateChildEntry(fileName)
		c.loadDirectory(newEntry)

		return true
	}

	if key == nuikey.KeyBackspace {
		if len(c.currentEntry.ServicePath) > 1 {
			newEntry := c.currentEntry.CreateParentEntry()
			c.loadDirectory(newEntry)
			return true
		}
	}

	return false
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

func (c *FilePanel) loadDirectory(entry *Entry) {
	if entry == nil {
		return
	}

	entries, err := ReadEntry(entry)
	if err != nil {
		c.fileList.SetRowCount(0)
		c.fileList.SetCellText(0, 0, fmt.Sprintf("Error: %v", err))
		return
	}

	c.fileList.SetRowCount(len(entries))
	for i, entry := range entries {
		c.fileList.SetCellText(0, i, entry.DisplayName())
		if entry.IsDir {
			c.fileList.SetCellText(1, i, "<DIR>")
		} else {
			c.fileList.SetCellText(1, i, fmt.Sprint(entry.Size))
		}
		c.fileList.SetCellText(2, i, entry.Modified.Format("2006-01-02 15:04:05"))
	}

	c.currentEntry = entry
}
