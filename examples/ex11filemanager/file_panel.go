package ex11filemanager

import (
	"fmt"

	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nuiforms/ui"
)

type FilePanel struct {
	ui.Widget

	folderStack []*Entry

	topPanel      *ui.Panel
	topPanelLabel *ui.Label

	contentPanel *ui.Panel
	fileList     *ui.Table
	bottomPanel  *ui.Panel

	onColumnResize func(col int, newWidth int)
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
	c.AddWidgetOnGrid(c.contentPanel, 1, 0)
	c.fileList = ui.NewTable()
	c.fileList.SetName("FileList")
	c.fileList.SetColumnCount(3)
	c.fileList.SetColumnName(0, "Name")
	c.fileList.SetColumnName(1, "Size")
	c.fileList.SetColumnName(2, "Modified")
	c.fileList.SetAllowScroll(false, true)
	c.fileList.SetSelectingCell(false)
	c.fileList.SetOnKeyDown(c.fileListKeyDown)
	c.fileList.SetOnColumnResize(c.columnResized)
	c.contentPanel.AddWidget(c.fileList)

	c.bottomPanel = ui.NewPanel()
	c.bottomPanel.AddWidgetOnGrid(ui.NewLabel("File panel bottom"), 0, 0)
	c.AddWidgetOnGrid(c.bottomPanel, 2, 0)

	rootEntries := readRootEntries()

	c.gotoFolder(rootEntries[0])

	return &c
}

func (c *FilePanel) gotoFolder(entry *Entry) {
	if entry.isLinkToParentDirectory {
		c.gotoParentFolder()
		return
	}
	if entry == nil || !entry.IsDir {
		return
	}
	if len(c.folderStack) > 0 {
		currentEntry := c.folderStack[len(c.folderStack)-1]
		currentEntry.selectedChildIndex = c.fileList.CurrentRow()
	}
	c.folderStack = append(c.folderStack, entry)
	c.loadDirectory(entry)
}

func (c *FilePanel) gotoParentFolder() {
	if len(c.folderStack) < 2 {
		return
	}
	c.folderStack = c.folderStack[:len(c.folderStack)-1]
	c.loadDirectory(c.folderStack[len(c.folderStack)-1])
}

func (c *FilePanel) fileListKeyDown(key nuikey.Key, mods nuikey.KeyModifiers) bool {
	if key == nuikey.KeyEnter {
		currentRowIndex := c.fileList.CurrentRow()
		if currentRowIndex < 0 || currentRowIndex >= c.fileList.RowCount() {
			return false
		}
		entry, ok := c.fileList.GetCellData(0, currentRowIndex).(*Entry)
		if !ok {
			return false
		}
		c.gotoFolder(entry)

		return true
	}

	if key == nuikey.KeyBackspace {
		c.gotoParentFolder()
		return true
	}

	return false
}

func (c *FilePanel) SetColumnWidth(col int, width int) {
	if c.fileList == nil {
		return
	}
	c.fileList.SetColumnWidth(col, width)
}

func (c *FilePanel) SetOnColumnResize(onColumnResize func(col int, newWidth int)) {
	c.onColumnResize = onColumnResize
}

func (c *FilePanel) columnResized(col int, newWidth int) {
	if c.onColumnResize != nil {
		c.onColumnResize(col, newWidth)
	}
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
	for i, en := range entries {
		c.fileList.SetCellText(0, i, en.DisplayName())
		c.fileList.SetCellData(0, i, en)
		if en.IsDir {
			c.fileList.SetCellText(1, i, "<DIR>")
		} else {
			c.fileList.SetCellText(1, i, fmt.Sprint(en.Size))
		}
		c.fileList.SetCellText(2, i, en.Modified.Format("2006-01-02 15:04:05"))
	}

	c.fileList.SetCurrentCell(0, entry.selectedChildIndex)
}
