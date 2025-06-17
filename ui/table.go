package ui

import "image/color"

type Table struct {
	widget Widget

	cols map[int]*tableColumn
	rows map[int]*tableRow

	rowHeight int

	rowCount    int
	columnCount int

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

	c.rows = make(map[int]*tableRow)
	c.rowHeight = 30
	c.cols = make(map[int]*tableColumn)
	c.defaultColumnWidth = 200

	return &c
}

func (c *Table) Widgeter() any {
	return &c.widget
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
	for rowIndex := 0; rowIndex < c.rowCount; rowIndex++ {
		rowObj, exists := c.rows[rowIndex]
		if !exists {
			continue
		}
		for colIndex := 0; colIndex < c.columnCount; colIndex++ {
			cellObj, exists := rowObj.cells[colIndex]
			if !exists {
				continue
			}
			x := c.columnOffset(colIndex)
			y := rowIndex * c.rowHeight
			//cnv.FillRect(x, y, 200, 30, color.RGBA{R: 255, G: 255, B: 255, A: 255})
			cnv.DrawTextMultiline(x, y, 200, 30, HAlignLeft, VAlignCenter, cellObj.text, color.RGBA{R: 200, G: 200, B: 200, A: 255}, "robotomono", 16, false)
		}
	}

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
