package ui

import (
	"fmt"
	"image/color"

	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
)

type Table struct {
	widget Widget

	cols map[int]*tableColumn
	rows map[int]*tableRow

	rowHeight int

	rowCount    int
	columnCount int

	currentCellX int
	currentCellY int

	defaultColumnWidth int
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
	c.widget.InitWidget()
	c.widget.SetXExpandable(true)
	c.widget.SetYExpandable(true)
	c.widget.SetOnPaint(c.draw)
	c.widget.SetAllowScroll(true, true)
	c.widget.SetBackgroundColor(color.RGBA{R: 50, G: 60, B: 70, A: 255})
	c.widget.SetOnMouseDown(c.onMouseDown)

	c.rows = make(map[int]*tableRow)
	c.rowHeight = 30
	c.cols = make(map[int]*tableColumn)
	c.defaultColumnWidth = 200

	return &c
}

func (c *Table) Widgeter() any {
	return &c.widget
}

func (c *Table) onMouseDown(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) {
	headerColumnBorder := c.headerColumnBorderByPosition(x, y)
	if headerColumnBorder >= 0 {
		fmt.Println("Header column border clicked:", headerColumnBorder)
		return
	}

	headerColumn := c.headerColumnByPosition(x, y)
	if headerColumn >= 0 {
		fmt.Println("Header column clicked:", headerColumn)
		return
	}

	row, col := c.cellByPosition(x, y)
	if row >= 0 && col >= 0 {
		c.currentCellX = col
		c.currentCellY = row
	}
}

func (c *Table) headerColumnBorderByPosition(x int, y int) int {
	if y < c.widget.scrollY || y >= c.widget.scrollY+c.headerHeight() {
		return -1
	}
	if x < 0 {
		return -1
	}
	if x >= c.widget.innerWidth {
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
	if y < c.widget.scrollY || y >= c.widget.scrollY+c.headerHeight() {
		return -1
	}
	if x < 0 {
		return -1
	}
	if x >= c.widget.innerWidth {
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

func (c *Table) cellByPosition(x, y int) (int, int) {
	if y < c.widget.scrollY || y >= c.widget.scrollY+c.rowCount*c.rowHeight {
		return -1, -1
	}
	col := 0
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
	row := (y - c.headerHeight()) / c.rowHeight
	if row < 0 || row >= c.rowCount {
		return -1, -1
	}
	return row, col
}

func (c *Table) headerHeight() int {
	return c.rowHeight
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

func (c *Table) SetCellText(row, col int, text string) {
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

func (c *Table) draw(cnv *Canvas) {
	yOffset := 0

	yOffset += c.rowHeight

	for rowIndex := 0; rowIndex < c.rowCount; rowIndex++ {
		rowObj, exists := c.rows[rowIndex]
		if exists {
			for colIndex := 0; colIndex < c.columnCount; colIndex++ {
				cellObj, exists := rowObj.cells[colIndex]
				if exists {
					x := c.columnOffset(colIndex)
					y := yOffset

					columnWidth := c.columnWidth(colIndex)
					selected := c.currentCellX == colIndex && c.currentCellY == rowIndex
					backColor := color.RGBA{R: 50, G: 60, B: 70, A: 255}
					if selected {
						backColor = color.RGBA{R: 80, G: 90, B: 100, A: 255}
					}
					cnv.FillRect(x, y, columnWidth, c.rowHeight, backColor)
					cnv.DrawTextMultiline(x, y, columnWidth, c.rowHeight, HAlignLeft, VAlignCenter, cellObj.text, color.RGBA{R: 200, G: 200, B: 200, A: 255}, "robotomono", 16, false)
				}
			}
		}

		yOffset += c.rowHeight
	}

	for colIndex := 0; colIndex < c.columnCount; colIndex++ {
		colObj, exists := c.cols[colIndex]
		if exists {
			x := c.columnOffset(colIndex)
			cnv.FillRect(x, c.widget.scrollY, colObj.width, c.rowHeight, color.RGBA{R: 70, G: 80, B: 90, A: 255})
			cnv.DrawTextMultiline(x, c.widget.scrollY, colObj.width, c.rowHeight, HAlignLeft, VAlignCenter, colObj.name, color.RGBA{R: 200, G: 200, B: 200, A: 255}, "robotomono", 16, false)
		}
	}

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
	c.widget.SetInnerSize(width, c.rowCount*c.rowHeight)
}
