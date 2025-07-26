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
	headerRowsCount  int
	headerRows       map[int]*tableHeaderRow
	headerRowHeights map[int]int
	columnsWidths    map[int]int
	rows             map[int]*tableRow

	rowHeight1         int // Can be changed
	defaultColumnWidth int // Default width for columns if not set

	rowCount    int
	columnCount int

	columnResizingIndex int

	cellBorderWidth int
	cellBorderColor color.RGBA

	cellPadding int

	showSelection bool
	selectingRow  bool
	selectingCell bool

	// Selection
	currentCellX int
	currentCellY int

	onSelectionChanged func(x int, y int)

	headerWidget  *tableHeader
	editorTextBox *TextBox
	innerWidgets  []*innerWidget

	editTriggerDoubleClick bool
	editTriggerEnter       bool
	editTriggerF2          bool
	editTriggerKeyDown     bool
}

type innerWidget struct {
	widget        Widgeter
	posCellX      int
	posCellY      int
	widthInCells  int
	heightInCells int
}

type tableRow struct {
	cells map[int]*tableCell
}

type tableHeaderRow struct {
	cells map[int]*tableHeaderCell
}

type tableHeaderCell struct {
	name  string
	spanX int
	spanY int
}

func (c tableHeaderCell) SpanX() int {
	if c.spanX <= 0 {
		return 1
	}
	return c.spanX
}

func (c tableHeaderCell) SpanY() int {
	if c.spanY <= 0 {
		return 1
	}
	return c.spanY
}

type tableCell struct {
	text string
}

func NewTable() *Table {
	var c Table
	c.InitWidget()
	c.SetAbsolutePositioning(true)
	c.SetTypeName("Table")
	c.SetXExpandable(true)
	c.SetYExpandable(true)
	c.SetAllowScroll(true, true)
	c.SetBackgroundColor(c.BackgroundColor())

	// Events
	c.SetOnPaint(c.draw)
	c.SetOnPostPaint(c.drawPost)
	c.SetOnMouseDown(c.onMouseDown)
	c.SetOnMouseUp(c.onMouseUp)
	c.SetOnMouseDblClick(c.onMouseDblClick)
	//c.SetOnKeyDown(c.onKeyDown)
	c.SetOnKeyUp(c.onKeyUp)
	c.SetOnMouseMove(c.onMouseMove)
	c.SetOnFocused(c.onFocused)

	c.SetOnScrollChanged(func(scrollX, scrollY int) {
		c.updateInnerWidgetsLayout()
	})

	// Init runtime
	c.SetCanBeFocused(true)
	c.rows = make(map[int]*tableRow)
	c.rowHeight1 = 30
	c.headerRows = make(map[int]*tableHeaderRow)
	c.columnsWidths = make(map[int]int)
	c.headerRowHeights = make(map[int]int)
	c.defaultColumnWidth = 200

	c.headerRowsCount = 2

	c.columnResizingIndex = -1

	c.cellBorderWidth = 1
	c.cellBorderColor = color.RGBA{R: 100, G: 100, B: 100, A: 255}
	c.cellPadding = 3

	c.showSelection = true
	c.selectingRow = true
	c.selectingCell = true

	c.headerWidget = newTableHeader()
	c.headerWidget.OnHeaderMouseDown = func(button nuimouse.MouseButton, x, y int, mods nuikey.KeyModifiers) bool {
		return c.onMouseDown(button, x+c.scrollX, y+c.scrollY, mods)
	}
	c.headerWidget.OnHeaderMouseUp = func(button nuimouse.MouseButton, x, y int, mods nuikey.KeyModifiers) bool {
		return c.onMouseUp(button, x+c.scrollX, y+c.scrollY, mods)
	}
	c.headerWidget.OnHeaderMouseMove = func(x, y int, mods nuikey.KeyModifiers) nuimouse.MouseCursor {
		return c.onMouseMoveHeader(x+c.scrollX, y+c.scrollY, mods)
	}
	c.AddWidget(c.headerWidget)

	c.innerWidgets = make([]*innerWidget, 0)

	c.updateInnerWidgetsLayout()

	return &c
}

func (c *Table) SetEditTriggerDoubleClick(enabled bool) {
	c.editTriggerDoubleClick = enabled
}

func (c *Table) SetEditTriggerEnter(enabled bool) {
	c.editTriggerEnter = enabled
}

func (c *Table) SetEditTriggerF2(enabled bool) {
	c.editTriggerF2 = enabled
}

func (c *Table) SetEditTriggerKeyDown(enabled bool) {
	c.editTriggerKeyDown = enabled
}

func (c *Table) SetOnSelectionChanged(callback func(x int, y int)) {
	c.onSelectionChanged = callback
}

func (c *Table) AddWidgetOnTable(widget Widgeter, posCellX int, posCellY int, widthInCells int, heightInCells int) {
	if posCellX < 0 || posCellY < 0 || posCellX >= c.columnCount || posCellY >= c.rowCount {
		return
	}
	if posCellX+widthInCells > c.columnCount || posCellY+heightInCells > c.rowCount {
		return
	}

	var inWidget innerWidget
	inWidget.widget = widget
	inWidget.posCellX = posCellX
	inWidget.posCellY = posCellY
	inWidget.widthInCells = widthInCells
	inWidget.heightInCells = heightInCells
	c.innerWidgets = append(c.innerWidgets, &inWidget)
	c.AddWidget(widget)
	c.updateInnerWidgetsLayout()
}

func (c *Table) rowOffset(row int) int {
	if row < 0 || row >= c.rowCount {
		return 0
	}
	return c.headerHeight() + row*c.rowHeight1
}

func (c *Table) headerCell(colIndex int, rowIndex int) *tableHeaderCell {
	if colIndex < 0 || colIndex >= c.columnCount {
		return &tableHeaderCell{name: ""}
	}

	if rowIndex < 0 || rowIndex >= c.headerRowsCount {
		return &tableHeaderCell{name: ""}
	}

	headerRow, exists := c.headerRows[rowIndex]
	if !exists {
		headerRow = &tableHeaderRow{cells: make(map[int]*tableHeaderCell)}
		c.headerRows[rowIndex] = headerRow
	}

	cell, exists := headerRow.cells[colIndex]
	if !exists {
		cell = &tableHeaderCell{name: ""}
		headerRow.cells[colIndex] = cell
	}

	return cell
}

func (c *Table) headerCellShadowed(colIndex int, rowIndex int) bool {
	result := false

	for headerRowIndex := 0; headerRowIndex < c.headerRowsCount; headerRowIndex++ {
		for headerColIndex := 0; headerColIndex < c.columnCount; headerColIndex++ {
			if headerRowIndex == rowIndex && headerColIndex == colIndex {
				continue
			}
			headerCell := c.headerCell(headerColIndex, headerRowIndex)
			spanX := headerCell.SpanX()
			spanY := headerCell.SpanY()
			if spanX > 1 || spanY > 1 {
				cellSpanX1 := headerColIndex
				cellSpanX2 := headerColIndex + spanX - 1
				cellSpanY1 := headerRowIndex
				cellSpanY2 := headerRowIndex + spanY - 1
				if colIndex >= cellSpanX1 && colIndex <= cellSpanX2 &&
					rowIndex >= cellSpanY1 && rowIndex <= cellSpanY2 {
					result = true
					break
				}
			}
		}
	}

	return result
}

/*func (c *Table) columnWidth(col int) int {
	if col < 0 || col >= c.columnCount {
		return c.defaultColumnWidth
	}
	colWidth, exists := c.columnsWidths[col]
	if !exists {
		return c.defaultColumnWidth
	}
	return colWidth
}*/

func (c *Table) updateInnerWidgetsLayout() {
	c.headerWidget.SetPosition(c.scrollX, c.scrollY)
	c.headerWidget.SetSize(c.innerWidth, c.headerHeight())

	for _, inWidget := range c.innerWidgets {
		widgetPosInPixelsX := c.columnOffset(inWidget.posCellX)
		widgetPosInPixelsY := c.rowOffset(inWidget.posCellY)
		widgetWidthInPixels := 0
		for i := 0; i < inWidget.widthInCells; i++ {
			widgetWidthInPixels += c.columnWidth(inWidget.posCellX + i)
		}
		widgetHeightInPixels := inWidget.heightInCells * c.rowHeight1
		inWidget.widget.SetPosition(widgetPosInPixelsX, widgetPosInPixelsY)
		inWidget.widget.SetSize(widgetWidthInPixels, widgetHeightInPixels)
	}
}

func (c *Table) RowCount() int {
	return c.rowCount
}

func (c *Table) SetRowCount(count int) {
	c.rowCount = count
	c.updateInnerSize()
	c.updateInnerWidgetsLayout()
}

func (c *Table) SetColumnCount(count int) {
	c.columnCount = count
	c.updateInnerSize()
	c.updateInnerWidgetsLayout()
}

func (c *Table) SetColumnWidth(col int, width int) {
	/*if col < 0 || col >= c.columnCount {
		return
	}
	colInfo, exists := c.cols[col]
	if !exists {
		colInfo = &tableColumn{name: "", width: c.defaultColumnWidth}
		c.cols[col] = colInfo
	}*/
	//colInfo.width = width

	c.columnsWidths[col] = width

	c.updateInnerSize()
	c.updateInnerWidgetsLayout()
}

func (c *Table) SetColumnCellName(col int, row int, name string) {
	if col < 0 || col >= c.columnCount {
		return
	}
	headerCell := c.headerCell(col, row)
	headerCell.name = name
	c.updateInnerSize()
}

func (c *Table) SetColumnName(col int, name string) {
	if col < 0 || col >= c.columnCount {
		return
	}
	headerCell := c.headerCell(col, 0)
	headerCell.name = name
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

	c.ScrollToCell(c.currentCellY, c.currentCellX)

	UpdateMainForm()
	if c.onSelectionChanged != nil {
		c.onSelectionChanged(c.currentCellX, c.currentCellY)
	}
}

func (c *Table) SetHeaderCellSpan(col int, row int, spanX int, spanY int) {
	if col < 0 || col >= c.columnCount || row < 0 || row >= c.headerRowsCount {
		return
	}
	headerCell := c.headerCell(col, row)
	headerCell.spanX = spanX
	headerCell.spanY = spanY
	c.updateInnerSize()
	c.updateInnerWidgetsLayout()
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
	leftTopY := c.rowOffset(row)
	c.ScrollEnsureVisible(leftTopX, leftTopY-c.rowHeight1)

	rightBottomX := leftTopX + c.columnWidth(col)
	rightBottomY := leftTopY + c.rowHeight1
	c.ScrollEnsureVisible(rightBottomX, rightBottomY)
	UpdateMainForm()
}

func (c *Table) onMouseDown(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
	headerColumnBorder := c.headerColumnBorderByPosition(x, y)
	if headerColumnBorder >= 0 {
		fmt.Println("Header column border clicked:", headerColumnBorder)
		c.columnResizingIndex = headerColumnBorder
		return true
	}

	headerColumn := c.headerColumnByPosition(x, y)
	if headerColumn >= 0 {
		fmt.Println("Header column clicked:", headerColumn)
		return true
	}

	col, row := c.cellByPosition(x, y)
	if row >= 0 && col >= 0 {
		c.SetCurrentCell(col, row)
		//fmt.Println("Cell clicked:", col, row, " at ", x, y)
	}
	return true
}

func (c *Table) onMouseUp(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
	c.columnResizingIndex = -1
	return true
}

func (c *Table) onMouseDblClick(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
	col, row := c.cellByPosition(x, y)
	fmt.Println("Cell double clicked:", col, row, " at ", x, y)
	if row >= 0 && col >= 0 {
		c.SetCurrentCell(col, row)
		if c.editTriggerDoubleClick {
			c.EditCurrentCell("")
		}
		return true
	}
	return false
}

func (c *Table) onFocused() {
}

func (c *Table) ProcessKeyDown(key nuikey.Key, mods nuikey.KeyModifiers) bool {
	processed := false

	if c.onKeyDown != nil {
		processed = c.onKeyDown(key, mods)
	}

	if processed {
		return processed
	}

	if key == nuikey.KeyArrowLeft {
		if c.currentCellX > 0 {
			c.SetCurrentCell(c.currentCellX-1, c.currentCellY)
			UpdateMainForm()
			processed = true
		}
	}

	if key == nuikey.KeyArrowRight {
		if c.currentCellX < c.columnCount-1 {
			c.SetCurrentCell(c.currentCellX+1, c.currentCellY)
			UpdateMainForm()
			processed = true
		}
	}

	if key == nuikey.KeyArrowUp {
		if c.currentCellY > 0 {
			selectRowIndex := c.currentCellY - 1
			if selectRowIndex >= c.rowCount {
				selectRowIndex = c.rowCount - 1
			}
			c.SetCurrentCell(c.currentCellX, selectRowIndex)
			UpdateMainForm()
			processed = true
		}
	}

	if key == nuikey.KeyArrowDown {
		if c.currentCellY < c.rowCount-1 {
			selectRowIndex := c.currentCellY + 1
			if selectRowIndex >= c.rowCount {
				selectRowIndex = c.rowCount - 1
			}
			c.SetCurrentCell(c.currentCellX, selectRowIndex)
			UpdateMainForm()
			processed = true
		}
	}

	if key == nuikey.KeyHome {
		c.SetCurrentCell(c.currentCellX, 0)
		UpdateMainForm()
		processed = true
	}

	if key == nuikey.KeyEnd {
		c.SetCurrentCell(c.currentCellX, c.rowCount-1)
		UpdateMainForm()
		processed = true
	}

	if key == nuikey.KeyEnter {
		if c.editTriggerEnter {
			c.EditCurrentCell("")
			UpdateMainForm()
			processed = true
		}
	}

	if key == nuikey.KeyF2 {
		if c.editTriggerF2 {
			c.EditCurrentCell("")
			UpdateMainForm()
			processed = true
		}
	}

	if key == nuikey.KeyPageUp {
		pageSizeInRows := c.Height() / c.rowHeight1
		targetRow := c.currentCellY - pageSizeInRows
		if targetRow < 0 {
			targetRow = 0
		}
		if targetRow != c.currentCellY {
			c.SetCurrentCell(c.currentCellX, targetRow)
			UpdateMainForm()
		}
		processed = true
	}

	if key == nuikey.KeyPageDown {
		pageSizeInRows := c.Height() / c.rowHeight1
		targetRow := c.currentCellY + pageSizeInRows
		if targetRow >= c.rowCount {
			targetRow = c.rowCount - 1
		}
		if targetRow != c.currentCellY {
			c.SetCurrentCell(c.currentCellX, targetRow)
			UpdateMainForm()
		}
		processed = true
	}

	return processed
}

func (c *Table) onKeyUp(key nuikey.Key, mods nuikey.KeyModifiers) bool {
	return true
}

func (c *Table) onMouseMoveHeader(x int, y int, mods nuikey.KeyModifiers) nuimouse.MouseCursor {
	if c.columnResizingIndex >= 0 {
		if c.columnResizingIndex < 0 || c.columnResizingIndex >= c.columnCount {
			return nuimouse.MouseCursorResizeHor
		}
		/*colInfo, exists := c.cols[c.columnResizingIndex]
		if !exists {
			colInfo = &tableColumn{name: "", width: c.defaultColumnWidth}
			c.cols[c.columnResizingIndex] = colInfo
		}*/

		colWidth := c.columnWidth(c.columnResizingIndex)

		newWidth := x - c.columnOffset(c.columnResizingIndex)
		if newWidth < 50 {
			newWidth = 50
		}
		if newWidth != colWidth {
			c.columnsWidths[c.columnResizingIndex] = newWidth
			c.updateInnerSize()
			c.updateInnerWidgetsLayout()
			c.SetMouseCursor(nuimouse.MouseCursorResizeHor)
			UpdateMainForm()
		}
		return nuimouse.MouseCursorResizeHor
	}

	headerColumnBorder := c.headerColumnBorderByPosition(x, y)
	if headerColumnBorder >= 0 {
		return nuimouse.MouseCursorResizeHor
	}
	return nuimouse.MouseCursorArrow
}

func (c *Table) onMouseMove(x int, y int, mods nuikey.KeyModifiers) bool {
	c.SetMouseCursor(nuimouse.MouseCursorArrow)
	return true
}

func (c *Table) SetShowSelection(show bool) {
	c.showSelection = show
}

func (c *Table) SetSelectingRow(selecting bool) {
	c.selectingRow = selecting
}

func (c *Table) SetSelectingCell(selecting bool) {
	c.selectingCell = selecting
}

func (c *Table) draw(cnv *Canvas) {
	yOffset := 0

	yOffset += c.headerHeight()

	visibleRow1, visibleRow2 := c.visibleRows()
	//fmt.Println("Visible rows:", visibleRow1, visibleRow2)

	yOffset += visibleRow1 * c.rowHeight1

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

					backColor := c.BackgroundColor()
					if c.showSelection {
						if rowIsSelected {
							backColor = c.BackgroundColorAccent1()
						}
						if cellIsSelected {
							backColor = c.BackgroundColorAccent2()
						}
					}
					cnv.FillRect(x, y, columnWidth, c.rowHeight1, backColor)

					cellText := ""
					if cellObj != nil {
						cellText = cellObj.text
					}

					cnv.SetHAlign(HAlignLeft)
					cnv.SetVAlign(VAlignCenter)
					cnv.SetFontFamily(c.FontFamily())
					cnv.SetFontSize(c.FontSize())
					cnv.SetColor(c.Color())
					cnv.DrawText(x+c.cellPadding, y+c.cellPadding, columnWidth-c.cellPadding*2, c.rowHeight1-c.cellPadding*2, cellText)
				}
			}
		}

		yOffset += c.rowHeight1
	}

	// Draw cell borders
	cnv.Save()
	cnv.SetDirectTranslateAndClip(cnv.state.translateX+c.scrollX, cnv.state.translateY+c.scrollY+c.headerHeight(), c.Width(), c.Height()-c.headerHeight())
	for rowIndex := visibleRow1; rowIndex < visibleRow2+1; rowIndex++ {
		x1 := 0
		y1 := rowIndex*c.rowHeight1 - c.scrollY
		x2 := c.innerWidth
		y2 := y1
		cnv.DrawLine(x1, y1, x2, y2, c.cellBorderWidth, c.cellBorderColor)
	}

	for colIndex := 0; colIndex < c.columnCount+1; colIndex++ {
		x1 := c.columnOffset(colIndex)
		y1 := visibleRow1 * c.rowHeight1
		x2 := x1
		y2 := visibleRow2 * c.rowHeight1
		cnv.DrawLine(x1, y1, x2, y2, c.cellBorderWidth, c.cellBorderColor)
	}

	cnv.Restore()
}

func (c *Table) drawPost(cnv *Canvas) {
	// Draw header
	for headerRowIndex := 0; headerRowIndex < c.headerRowsCount; headerRowIndex++ {
		for colIndex := 0; colIndex < c.columnCount; colIndex++ {
			needToDisplay := true
			if c.headerCellShadowed(colIndex, headerRowIndex) {
				needToDisplay = false
			}

			if !needToDisplay {
				continue
			}

			headerCell := c.headerCell(colIndex, headerRowIndex)

			cellSpanX := headerCell.SpanX()
			cellSpanY := headerCell.SpanY()

			headerRowOffset := c.headerRowOffset(headerRowIndex)
			//headerRowHeight := c.headerRowHeight(headerRowIndex)

			//cellWidth := c.columnWidth(colIndex)
			cellWidth := 0
			for i := 0; i < cellSpanX; i++ {
				cellWidth += c.columnWidth(colIndex + i)
			}

			cellHeight := 0
			for i := 0; i < cellSpanY; i++ {
				cellHeight += c.headerRowHeight(headerRowIndex + i)
			}
			//cellHeight := headerRowHeight

			x := c.columnOffset(colIndex)
			y := headerRowOffset + c.scrollY

			cnv.FillRect(x, y, cellWidth, cellHeight, c.BackgroundColor())
			cnv.SetHAlign(HAlignLeft)
			cnv.SetVAlign(VAlignCenter)
			cnv.SetColor(c.Color())
			cnv.SetFontFamily(c.FontFamily())
			cnv.SetFontSize(c.FontSize())
			cnv.DrawText(x+c.cellPadding, y+c.cellPadding, cellWidth-c.cellPadding*2, cellHeight-c.cellPadding*2, headerCell.name)
		}
	}

	/*
		for colIndex := 0; colIndex < c.columnCount; colIndex++ {
			colObj, exists := c.cols[colIndex]
			if exists {
				x := c.columnOffset(colIndex)
				cnv.FillRect(x, c.scrollY, colObj.width, c.headerHeight(), color.RGBA{R: 70, G: 80, B: 90, A: 255})
				cnv.DrawTextMultiline(x+c.cellPadding, c.scrollY+c.cellPadding, colObj.width-c.cellPadding*2, c.headerHeight()-c.cellPadding*2, HAlignLeft, VAlignCenter, colObj.name, color.RGBA{R: 200, G: 200, B: 200, A: 255}, "robotomono", 16, false)
			}
		}*/

	for colIndex := 0; colIndex < c.columnCount+1; colIndex++ {
		x1 := c.columnOffset(colIndex)
		y1 := c.scrollY
		x2 := x1
		y2 := c.headerHeight() + c.scrollY
		cnv.DrawLine(x1, y1, x2, y2, c.cellBorderWidth, c.cellBorderColor)
	}

	// Draw table border
	cnv.SetColor(c.BackgroundColorAccent2())
	cnv.DrawRect(c.scrollX, c.scrollY, c.Width(), c.Height())
}

func (c *Table) visibleRows() (min int, max int) {
	min = c.scrollY / c.rowHeight1
	max = min + (c.Height()-c.headerHeight())/c.rowHeight1
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

	colWidth, exists := c.columnsWidths[col]
	if !exists {
		return c.defaultColumnWidth
	}

	return colWidth
}

func (c *Table) columnOffset(col int) int {
	result := 0
	for i := 0; i < col; i++ {
		colWidth := c.defaultColumnWidth
		if colWidthValue, exists := c.columnsWidths[i]; exists {
			colWidth = colWidthValue
		}
		result += colWidth
	}
	return result
}

func (c *Table) updateInnerSize() {
	width := 0
	for i := 0; i < c.columnCount; i++ {
		colWidth := c.defaultColumnWidth
		if colWidthValue, exists := c.columnsWidths[i]; exists {
			colWidth = colWidthValue
		}
		width += colWidth
	}
	c.SetInnerSize(width, c.headerHeight()+c.rowCount*c.rowHeight1)
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
	row = (y - c.headerHeight()) / c.rowHeight1
	if row < 0 || row >= c.rowCount {
		return -1, -1
	}
	return col, row
}

func (c *Table) headerRowOffset(headerRowIndex int) int {
	result := 0
	for i := 0; i < headerRowIndex; i++ {
		rowHeight := c.headerRowHeight(i)
		result += rowHeight
	}
	return result
}

func (c *Table) headerRowHeight(headerRowIndex int) int {
	result := c.rowHeight1
	if headerRowHeight, exists := c.headerRowHeights[headerRowIndex]; exists {
		result = headerRowHeight
	}
	return result
}

func (c *Table) headerHeight() int {
	result := 0
	for i := 0; i < c.headerRowsCount; i++ {
		rowHeight := c.headerRowHeight(i)
		result += rowHeight
	}
	return result
}

func (c *Table) EditCurrentCell(enteredText string) {
	if c.editorTextBox != nil {
		c.RemoveWidget(c.editorTextBox)
		c.editorTextBox = nil
	}

	c.editorTextBox = NewTextBox()

	if len(enteredText) == 0 {
		enteredText = c.GetCellText(c.currentCellX, c.currentCellY)
	}

	c.editorTextBox.SetText(enteredText)
	c.editorTextBox.MoveCursorToEnd()
	c.editorTextBox.SelectAllText()
	c.editorTextBox.SetOnTextBoxKeyDown(func(key nuikey.Key, mods nuikey.KeyModifiers) bool {
		if key == nuikey.KeyEnter {
			if c.editorTextBox != nil {
				c.SetCellText(c.currentCellX, c.currentCellY, c.editorTextBox.Text())
				c.RemoveWidget(c.editorTextBox)
				c.editorTextBox = nil
				UpdateMainForm()
				c.Focus()
			}
			return true
		}
		if key == nuikey.KeyEsc {
			if c.editorTextBox != nil {
				c.RemoveWidget(c.editorTextBox)
				c.editorTextBox = nil
				UpdateMainForm()
				c.Focus()
			}
			return true
		}
		return false
	})
	c.editorTextBox.SetOnFocusLost(func() {
		if c.editorTextBox != nil {
			c.RemoveWidget(c.editorTextBox)
			c.editorTextBox = nil
			UpdateMainForm()
			c.Focus()
		}
	})
	c.AddWidgetOnTable(c.editorTextBox, c.currentCellX, c.currentCellY, 1, 1)
	c.editorTextBox.Focus()
}
