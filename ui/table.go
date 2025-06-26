package ui

import (
	"fmt"
	"image/color"

	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
)

type Table struct {
	Widget

	// Content of the table
	cols map[int]*tableColumn
	rows map[int]*tableRow

	rowHeight          int // Can be changed
	defaultColumnWidth int // Default width for columns if not set

	rowCount    int
	columnCount int

	columnResizingIndex int

	cellBorderWidth int
	cellBorderColor color.RGBA

	cellPadding int

	selectingRow  bool
	selectingCell bool

	// Selection
	currentCellX int
	currentCellY int

	onSelectionChanged func(x int, y int)
}

type tableRow struct {
	cells map[int]*tableCell
}

type tableColumn struct {
	name  string
	width int
}

type tableCell struct {
	text string
}

func NewTable() *Table {
	var c Table
	c.InitWidget()
	c.SetTypeName("Table")
	c.SetXExpandable(true)
	c.SetYExpandable(true)
	c.SetOnPaint(c.draw)
	c.SetAllowScroll(true, true)
	c.SetBackgroundColor(color.RGBA{R: 50, G: 60, B: 70, A: 255})
	c.SetOnMouseDown(c.onMouseDown)
	c.SetOnMouseUp(c.onMouseUp)
	c.SetOnKeyDown(c.onKeyDown)
	c.SetOnKeyUp(c.onKeyUp)
	c.SetOnMouseMove(c.onMouseMove)

	c.rows = make(map[int]*tableRow)
	c.rowHeight = 30
	c.cols = make(map[int]*tableColumn)
	c.defaultColumnWidth = 200

	c.columnResizingIndex = -1

	c.cellBorderWidth = 1
	c.cellBorderColor = color.RGBA{R: 100, G: 100, B: 100, A: 255}
	c.cellPadding = 3

	c.selectingRow = true
	c.selectingCell = true

	return &c
}

func (c *Table) SetOnSelectionChanged(callback func(x int, y int)) {
	c.onSelectionChanged = callback
}

func (c *Table) RowCount() int {
	return c.rowCount
}

func (c *Table) SetRowCount(count int) {
	c.rowCount = count
	c.updateInnerSize()
}

func (c *Table) SetColumnCount(count int) {
	c.columnCount = count
	c.updateInnerSize()
}

func (c *Table) SetColumnWidth(col int, width int) {
	if col < 0 || col >= c.columnCount {
		return
	}
	colInfo, exists := c.cols[col]
	if !exists {
		colInfo = &tableColumn{name: "", width: c.defaultColumnWidth}
		c.cols[col] = colInfo
	}
	colInfo.width = width
	c.updateInnerSize()
}

func (c *Table) SetColumnName(col int, name string) {
	if col < 0 || col >= c.columnCount {
		return
	}
	colInfo, exists := c.cols[col]
	if !exists {
		colInfo = &tableColumn{name: "", width: c.defaultColumnWidth}
		c.cols[col] = colInfo
	}
	colInfo.name = name
	c.updateInnerSize()
}

func (c *Table) SetCellText(col int, row int, text string) {
	rowObj, exists := c.rows[row]
	if !exists {
		rowObj = &tableRow{cells: make(map[int]*tableCell)}
		c.rows[row] = rowObj
	}
	cellObj, exists := rowObj.cells[col]
	if !exists {
		cellObj = &tableCell{}
		rowObj.cells[col] = cellObj
	}
	cellObj.text = text
	UpdateMainForm()
}

func (c *Table) SetCurrentCell(col int, row int) {
	if row < 0 || row >= c.rowCount || col < 0 || col >= c.columnCount {
		return
	}
	if c.selectingCell {
		c.currentCellX = col
	} else {
		c.currentCellX = 0
	}
	c.currentCellY = row
	c.ScrollToCell(row, col)
	UpdateMainForm()
	if c.onSelectionChanged != nil {
		c.onSelectionChanged(c.currentCellX, c.currentCellY)
	}
}

func (c *Table) CurrentRow() int {
	if c.currentCellY < 0 || c.currentCellY >= c.rowCount {
		return -1
	}
	return c.currentCellY
}

func (c *Table) CurrentColumn() int {
	if c.currentCellX < 0 || c.currentCellX >= c.columnCount {
		return -1
	}
	return c.currentCellX
}

func (c *Table) GetCellText(col int, row int) string {
	if row < 0 || row >= c.rowCount || col < 0 || col >= c.columnCount {
		return ""
	}
	rowObj, exists := c.rows[row]
	if !exists {
		return ""
	}
	cellObj, exists := rowObj.cells[col]
	if !exists {
		return ""
	}
	return cellObj.text
}

func (c *Table) ScrollToCell(row, col int) {
	if row < 0 || row >= c.rowCount || col < 0 || col >= c.columnCount {
		return
	}

	leftTopX := c.columnOffset(col)
	leftTopY := c.headerHeight() + row*c.rowHeight
	c.ScrollEnsureVisible(leftTopX, leftTopY-c.rowHeight)

	rightBottomX := leftTopX + c.columnWidth(col)
	rightBottomY := leftTopY + c.rowHeight
	c.ScrollEnsureVisible(rightBottomX, rightBottomY)
	UpdateMainForm()
}

func (c *Table) onMouseDown(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) {
	headerColumnBorder := c.headerColumnBorderByPosition(x, y)
	if headerColumnBorder >= 0 {
		fmt.Println("Header column border clicked:", headerColumnBorder)
		c.columnResizingIndex = headerColumnBorder
		return
	}

	headerColumn := c.headerColumnByPosition(x, y)
	if headerColumn >= 0 {
		fmt.Println("Header column clicked:", headerColumn)
		return
	}

	col, row := c.cellByPosition(x, y)
	if row >= 0 && col >= 0 {
		c.SetCurrentCell(col, row)
	}
}

func (c *Table) onMouseUp(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) {
	c.columnResizingIndex = -1
}

func (c *Table) onKeyDown(key nuikey.Key, mods nuikey.KeyModifiers) {
	if key == nuikey.KeyArrowLeft {
		if c.currentCellX > 0 {
			c.SetCurrentCell(c.currentCellX-1, c.currentCellY)
			UpdateMainForm()
		}
	}

	if key == nuikey.KeyArrowRight {
		if c.currentCellX < c.columnCount-1 {
			c.SetCurrentCell(c.currentCellX+1, c.currentCellY)
			UpdateMainForm()
		}
	}

	if key == nuikey.KeyArrowUp {
		if c.currentCellY > 0 {
			c.SetCurrentCell(c.currentCellX, c.currentCellY-1)
			UpdateMainForm()
		}
	}

	if key == nuikey.KeyArrowDown {
		if c.currentCellY < c.rowCount-1 {
			c.SetCurrentCell(c.currentCellX, c.currentCellY+1)
			UpdateMainForm()
		}
	}

	if key == nuikey.KeyHome {
		c.SetCurrentCell(c.currentCellX, 0)
		UpdateMainForm()
	}

	if key == nuikey.KeyEnd {
		c.SetCurrentCell(c.currentCellX, c.rowCount-1)
		UpdateMainForm()
	}

	if key == nuikey.KeyPageUp {
		pageSizeInRows := c.Height() / c.rowHeight
		targetRow := c.currentCellY - pageSizeInRows
		if targetRow < 0 {
			targetRow = 0
		}
		if targetRow != c.currentCellY {
			c.SetCurrentCell(c.currentCellX, targetRow)
			UpdateMainForm()
		}
	}

	if key == nuikey.KeyPageDown {
		pageSizeInRows := c.Height() / c.rowHeight
		targetRow := c.currentCellY + pageSizeInRows
		if targetRow >= c.rowCount {
			targetRow = c.rowCount - 1
		}
		if targetRow != c.currentCellY {
			c.SetCurrentCell(c.currentCellX, targetRow)
			UpdateMainForm()
		}
	}
}

func (c *Table) onKeyUp(key nuikey.Key, mods nuikey.KeyModifiers) {
}

func (c *Table) onMouseMove(x int, y int, mods nuikey.KeyModifiers) {
	if c.columnResizingIndex >= 0 {
		if c.columnResizingIndex < 0 || c.columnResizingIndex >= c.columnCount {
			return
		}
		colInfo, exists := c.cols[c.columnResizingIndex]
		if !exists {
			colInfo = &tableColumn{name: "", width: c.defaultColumnWidth}
			c.cols[c.columnResizingIndex] = colInfo
		}
		newWidth := x - c.columnOffset(c.columnResizingIndex)
		if newWidth < 50 {
			newWidth = 50
		}
		if newWidth != colInfo.width {
			colInfo.width = newWidth
			c.updateInnerSize()
			c.SetMouseCursor(nuimouse.MouseCursorResizeHor)
			UpdateMainForm()
		}
		return
	}

	headerColumnBorder := c.headerColumnBorderByPosition(x, y)
	if headerColumnBorder >= 0 {
		c.SetMouseCursor(nuimouse.MouseCursorResizeHor)
		return
	}
	c.SetMouseCursor(nuimouse.MouseCursorArrow)
}

func (c *Table) draw(cnv *Canvas) {
	yOffset := 0

	yOffset += c.rowHeight

	visibleRow1, visibleRow2 := c.visibleRows()
	//fmt.Println("Visible rows:", visibleRow1, visibleRow2)

	yOffset += visibleRow1 * c.rowHeight

	for rowIndex := visibleRow1; rowIndex < visibleRow2; rowIndex++ {
		rowObj1, rowExists := c.rows[rowIndex]
		{
			_ = rowExists
			for colIndex := 0; colIndex < c.columnCount; colIndex++ {
				var cellObj *tableCell
				var cellExists bool
				if rowObj1 != nil {
					cellObj, cellExists = rowObj1.cells[colIndex]
				}
				{
					_ = cellExists
					x := c.columnOffset(colIndex)
					y := yOffset

					columnWidth := c.columnWidth(colIndex)

					/*selected := c.currentCellX == colIndex && c.currentCellY == rowIndex
					if c.selectingRow {
						selected = c.currentCellY == rowIndex
					}*/

					rowIsSelected := c.currentCellY == rowIndex
					if !c.selectingRow {
						rowIsSelected = false
					}

					cellIsSelected := c.currentCellX == colIndex && c.currentCellY == rowIndex
					if !c.selectingCell {
						cellIsSelected = false
					}

					backColor := color.RGBA{R: 50, G: 60, B: 70, A: 255}
					if rowIsSelected {
						backColor = color.RGBA{R: 80, G: 90, B: 100, A: 255}
					}
					if cellIsSelected {
						backColor = color.RGBA{R: 100, G: 110, B: 120, A: 255}
					}
					cnv.FillRect(x, y, columnWidth, c.rowHeight, backColor)

					cellText := ""
					if cellObj != nil {
						cellText = cellObj.text
					}

					cnv.DrawTextMultiline(x+c.cellPadding, y+c.cellPadding, columnWidth-c.cellPadding*2, c.rowHeight-c.cellPadding*2, HAlignLeft, VAlignCenter, cellText, color.RGBA{R: 200, G: 200, B: 200, A: 255}, "robotomono", 16, false)
				}
			}
		}

		yOffset += c.rowHeight
	}

	for colIndex := 0; colIndex < c.columnCount; colIndex++ {
		colObj, exists := c.cols[colIndex]
		if exists {
			x := c.columnOffset(colIndex)
			cnv.FillRect(x, c.scrollY, colObj.width, c.rowHeight, color.RGBA{R: 70, G: 80, B: 90, A: 255})
			cnv.DrawTextMultiline(x+c.cellPadding, c.scrollY+c.cellPadding, colObj.width-c.cellPadding*2, c.rowHeight-c.cellPadding*2, HAlignLeft, VAlignCenter, colObj.name, color.RGBA{R: 200, G: 200, B: 200, A: 255}, "robotomono", 16, false)
		}
	}

	// Draw cell borders
	cnv.Save()
	cnv.SetDirectTranslateAndClip(cnv.state.translateX+c.scrollX, cnv.state.translateY+c.scrollY+c.headerHeight(), c.Width(), c.Height()-c.headerHeight())
	for rowIndex := visibleRow1; rowIndex < visibleRow2+1; rowIndex++ {
		x1 := 0
		y1 := rowIndex*c.rowHeight - c.scrollY
		x2 := c.innerWidth
		y2 := y1
		cnv.DrawLine(x1, y1, x2, y2, c.cellBorderWidth, c.cellBorderColor)
	}
	cnv.Restore()

	for colIndex := 0; colIndex < c.columnCount+1; colIndex++ {
		x1 := c.columnOffset(colIndex)
		y1 := visibleRow1*c.rowHeight - c.headerHeight()
		x2 := x1
		y2 := (visibleRow2 + 1) * c.rowHeight
		cnv.DrawLine(x1, y1, x2, y2, c.cellBorderWidth, c.cellBorderColor)
	}
}

func (c *Table) visibleRows() (min int, max int) {
	min = c.scrollY / c.rowHeight
	max = min + (c.Height()-c.headerHeight())/c.rowHeight
	min = min - 1
	max = max + 1
	if min < 0 {
		min = 0
	}
	if max > c.rowCount {
		max = c.rowCount
	}
	return
}

func (c *Table) columnWidth(col int) int {
	if col < 0 || col >= c.columnCount {
		return c.defaultColumnWidth
	}
	colInfo, exists := c.cols[col]
	if !exists {
		return c.defaultColumnWidth
	}
	return colInfo.width
}

func (c *Table) columnOffset(col int) int {
	result := 0
	for i := 0; i < col; i++ {
		colWidth := c.defaultColumnWidth
		if column, exists := c.cols[i]; exists {
			colWidth = column.width
		}
		result += colWidth
	}
	return result
}

func (c *Table) updateInnerSize() {
	width := 0
	for i := 0; i < c.columnCount; i++ {
		colWidth := c.defaultColumnWidth
		if column, exists := c.cols[i]; exists {
			colWidth = column.width
		}
		width += colWidth
	}
	c.SetInnerSize(width, c.headerHeight()+c.rowCount*c.rowHeight)
	c.checkScrolls()
}

func (c *Table) headerColumnBorderByPosition(x int, y int) int {
	if y < c.scrollY || y >= c.scrollY+c.headerHeight() {
		return -1
	}
	if x < 0 {
		return -1
	}
	if x >= c.innerWidth {
		return -1
	}
	for col := 0; col < c.columnCount; col++ {
		colOffset := c.columnOffset(col)
		colWidth := c.columnWidth(col)
		rigthBorder := colOffset + colWidth
		if x >= rigthBorder-5 && x < rigthBorder+5 {
			return col
		}
	}
	return -1
}

func (c *Table) headerColumnByPosition(x int, y int) int {
	if y < c.scrollY || y >= c.scrollY+c.headerHeight() {
		return -1
	}
	if x < 0 {
		return -1
	}
	if x >= c.innerWidth {
		return -1
	}
	for col := 0; col < c.columnCount; col++ {
		colOffset := c.columnOffset(col)
		colWidth := c.columnWidth(col)
		if x >= colOffset && x < colOffset+colWidth {
			return col
		}
	}
	return -1
}

func (c *Table) cellByPosition(x, y int) (col int, row int) {
	col = 0
	for col < c.columnCount {
		colOffset := c.columnOffset(col)
		colWidth := c.columnWidth(col)
		if x >= colOffset && x < colOffset+colWidth {
			break
		}
		col++
	}
	if col >= c.columnCount {
		return -1, -1
	}
	row = (y - c.headerHeight()) / c.rowHeight
	if row < 0 || row >= c.rowCount {
		return -1, -1
	}
	return col, row
}

func (c *Table) headerHeight() int {
	return c.rowHeight
}
