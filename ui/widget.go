package ui

import (
	"fmt"
	"image/color"
	"math"
	"strings"

	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
)

type Widget struct {
	id       string
	name     string
	typeName string
	userData interface{}

	// position
	x int
	y int

	// size
	w int
	h int

	// anchors
	anchorLeft   bool
	anchorTop    bool
	anchorRight  bool
	anchorBottom bool

	// inner widgets
	absolutePositioning bool
	widgets             []Widgeter
	cellPadding         int // Padding between cells in the grid
	panelPadding        int // Padding around the panel

	gridX int // Grid position
	gridY int // Grid position

	xExpandable bool // If the widget can expand in X direction
	yExpandable bool // If the widget can expand in Y direction

	minWidth  int // Minimum width
	maxWidth  int // Maximum width
	minHeight int // Minimum height
	maxHeight int // Maximum height

	allowScrollX   bool
	allowScrollY   bool
	hideScrollbarX bool
	hideScrollbarY bool
	scrollX        int
	scrollY        int
	innerWidth     int
	innerHeight    int

	scrollBarXColor           color.RGBA
	scrollBarXSize            int
	scrollingX                bool
	scrollingXInitial         int
	scrollingXInitialMousePos int

	scrollBarYColor           color.RGBA
	scrollBarYSize            int
	scrollingY                bool
	scrollingYInitial         int
	scrollingYInitialMousePos int

	props map[string]interface{}

	timers []*timer

	visible bool

	// temp
	lastMouseX       int // After scrolling
	lastMouseY       int // After scrolling
	lastMouseAbsPosX int // Last mouse position relative to the widget
	lastMouseAbsPosY int // Last mouse position relative to the widget

	mouseCursor nuimouse.MouseCursor

	backgroundColor color.RGBA

	// callbacks
	onCustomPaint func(cnv *Canvas)
	onMouseDown   func(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers)
	onMouseUp     func(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers)
	onMouseMove   func(x int, y int, mods nuikey.KeyModifiers)
	onMouseLeave  func()
	onMouseEnter  func()
	onKeyDown     func(key nuikey.Key, mods nuikey.KeyModifiers)
	onKeyUp       func(key nuikey.Key, mods nuikey.KeyModifiers)
	onChar        func(char rune, mods nuikey.KeyModifiers)
	onMouseWheel  func(deltaX, deltaY int)
	onClick       func(button nuimouse.MouseButton, x int, y int)
}

/*func NewWidget() *Widget {
	var c Widget
	c.InitWidget()
	return &c
}*/

type ContainerGridColumnInfo struct {
	minWidth   int
	maxWidth   int
	expandable bool
	width      int
	collapsed  bool
}

type ContainerGridRowInfo struct {
	minHeight  int
	maxHeight  int
	expandable bool
	height     int
	collapsed  bool
}

const MaxUint = ^uint(0)
const MinUint = 0

const MaxInt = int(^uint(0) >> 1)
const MinInt = -MaxInt - 1

const MAX_WIDTH = 100000
const MAX_HEIGHT = 100000

func (c *Widget) InitWidget() {
	c.id = NewId()
	c.typeName = "Widget"
	c.name = "Widget-" + c.id
	c.props = make(map[string]any)
	c.timers = make([]*timer, 0)
	c.x = 0
	c.y = 0
	c.w = 300
	c.h = 180
	c.minWidth = 0
	c.minHeight = 0
	c.maxWidth = MAX_WIDTH
	c.maxHeight = MAX_HEIGHT
	c.visible = true
	c.panelPadding = 2
	c.cellPadding = 2
	c.scrollBarXSize = 10
	c.scrollBarYSize = 10
	c.scrollBarXColor = color.RGBA{R: 150, G: 150, B: 150, A: 100}
	c.scrollBarYColor = color.RGBA{R: 150, G: 150, B: 150, A: 100}
	c.innerWidth = 0
	c.anchorLeft = true
	c.anchorTop = true
	c.anchorRight = false
	c.anchorBottom = false
	c.widgets = make([]Widgeter, 0)
	c.mouseCursor = nuimouse.MouseCursorArrow
	c.backgroundColor = color.RGBA{R: 0, G: 0, B: 0, A: 0} // transparent by default
}

func (c *Widget) Id() string {
	return c.id
}

func (c *Widget) SetName(name string) {
	c.name = name
}

func (c *Widget) SetVisible(visible bool) {
	if c.visible != visible {
		c.visible = visible
		c.updateLayout(c.w, c.h, c.w, c.h)
		UpdateMainForm()
	}
}

func (c *Widget) SetTypeName(typeName string) {
	c.typeName = typeName
}

func (c *Widget) TypeName() string {
	return c.typeName
}

func (c *Widget) IsVisible() bool {
	return c.visible
}

func (c *Widget) GridX() int {
	return c.gridX
}

func (c *Widget) GridY() int {
	return c.gridY
}

func (c *Widget) SetGridPosition(x, y int) {
	c.gridX = x
	c.gridY = y
}

func (c *Widget) MinWidth() int {
	result := 0

	if !c.allowScrollX {
		_, _, _, allCellPadding := c.makeColumnsInfo(c.Width())
		columnsInfo, _, _, _ := c.makeColumnsInfo(c.Width() - (c.panelPadding + allCellPadding + c.panelPadding))
		for _, columnInfo := range columnsInfo {
			result += columnInfo.minWidth
		}
		result = result + c.panelPadding + allCellPadding + c.panelPadding
	}

	if c.minWidth > result {
		return c.minWidth
	}
	return result
}

func (c *Widget) MinHeight() int {
	result := 0

	if !c.allowScrollY {
		_, _, _, allCellPadding := c.makeRowsInfo(c.Height())
		rowsInfo, _, _, _ := c.makeRowsInfo(c.Height() - (c.panelPadding + allCellPadding + c.panelPadding))
		for _, rowInfo := range rowsInfo {
			result += rowInfo.minHeight
		}
		result += c.panelPadding + allCellPadding + c.panelPadding
	}

	if c.minHeight > result {
		return c.minHeight
	}
	return result
}

func (c *Widget) MaxWidth() int {
	return c.maxWidth
}

func (c *Widget) MaxHeight() int {
	return c.maxHeight
}

func (c *Widget) AddTimer(intervalMs int, callback func()) {
	t := &timer{
		intervalMs: intervalMs,
		callback:   callback,
	}
	c.timers = append(c.timers, t)
}

func (c *Widget) AddWidget(w Widgeter) {
	if _, exists := allwidgets[w.Id()]; exists {
		return
	}
	c.widgets = append(c.widgets, w)
	allwidgets[w.Id()] = w
	c.updateLayout(0, 0, 0, 0)
}

func (c *Widget) SetPanelPadding(padding int) {
	c.panelPadding = padding
}

func (c *Widget) AddWidgetOnGrid(w Widgeter, gridX int, gridY int) {
	if _, exists := allwidgets[w.Id()]; exists {
		return
	}
	w.SetGridPosition(gridX, gridY)
	c.widgets = append(c.widgets, w)
	allwidgets[w.Id()] = w
	c.updateLayout(0, 0, 0, 0)
	MainForm.Panel().updateLayout(0, 0, 0, 0) // Global Update Layout
}

func (c *Widget) RemoveWidget(w Widgeter) {
	delete(allwidgets, w.Id())
	for i, widget := range c.widgets {
		widgeter := widget
		if widgeter.Id() == w.Id() {
			c.widgets = append(c.widgets[:i], c.widgets[i+1:]...)
			return
		}
	}
	c.updateLayout(0, 0, 0, 0)
}

func (c *Widget) RemoveAllWidgets() {
	for _, w := range c.widgets {
		delete(allwidgets, w.Id())
	}
	c.widgets = make([]Widgeter, 0)
	c.updateLayout(0, 0, 0, 0)
	UpdateMainForm()
}

func (c *Widget) NextGridX() int {
	if len(c.widgets) == 0 {
		return 0
	}

	maxX := 0
	for _, w := range c.widgets {
		if w.GridX() >= maxX {
			maxX = w.GridX() + 1
		}
	}
	return maxX
}

func (c *Widget) NextGridY() int {
	if len(c.widgets) == 0 {
		return 0
	}

	maxY := 0
	for _, w := range c.widgets {
		if w.GridY() >= maxY {
			maxY = w.GridY() + 1
		}
	}
	return maxY
}

func (c *Widget) SetBackgroundColor(col color.RGBA) {
	c.backgroundColor = col
}

func (c *Widget) SetAllowScroll(allowX bool, allowY bool) {
	c.allowScrollX = allowX
	c.allowScrollY = allowY
}

func (c *Widget) SetMouseCursor(cursor nuimouse.MouseCursor) {
	c.mouseCursor = cursor
}

func (c *Widget) MouseCursor() nuimouse.MouseCursor {
	return c.mouseCursor
}

func (c *Widget) X() int {
	return c.x
}

func (c *Widget) Y() int {
	return c.y
}

func (c *Widget) Width() int {
	return c.w
}

func (c *Widget) Height() int {
	return c.h
}

func (c *Widget) InnerWidth() int {
	if c.innerWidth == 0 {
		return c.w
	}
	return c.innerWidth
}

func (c *Widget) InnerHeight() int {
	if c.innerHeight == 0 {
		return c.h
	}
	return c.innerHeight
}

func (c *Widget) SetInnerSize(width, height int) {
	c.innerWidth = width
	c.innerHeight = height
}

func (c *Widget) SetProp(key string, value interface{}) {
	c.props[key] = value
}

func (c *Widget) GetProp(key string) interface{} {
	if value, ok := c.props[key]; ok {
		return value
	}
	return nil
}

func (c *Widget) GetPropString(key string, defaultValue string) string {
	if value, ok := c.props[key]; ok {
		if strValue, ok := value.(string); ok {
			return strValue
		}
	}
	return defaultValue
}

func (c *Widget) GetPropInt(key string, defaultValue int) int {
	if value, ok := c.props[key]; ok {
		if intValue, ok := value.(int); ok {
			return intValue
		}
	}
	return defaultValue
}

func (c *Widget) GetPropBool(key string, defaultValue bool) bool {
	if value, ok := c.props[key]; ok {
		if boolValue, ok := value.(bool); ok {
			return boolValue
		}
	}
	return defaultValue
}

func (c *Widget) Focus() {
	MainForm.focusedWidget = WidgetById(c.Id())
	MainForm.Update()
}

func (c *Widget) IsFocused() bool {
	return MainForm.focusedWidget == WidgetById(c.Id())
}

func (c *Widget) IsHovered() bool {
	return MainForm.hoverWidget == WidgetById(c.Id())
}

func (c *Widget) Name() string {
	return c.name
}

func (c *Widget) SetPosition(x, y int) {
	c.x = x
	c.y = y
}

func (c *Widget) SetSize(w, h int) {
	oldW := c.w
	oldH := c.h
	c.w = w
	c.h = h
	c.updateLayout(oldW, oldH, w, h)
	c.checkScrolls()
}

func (c *Widget) SetMinSize(minWidth, minHeight int) {
	c.minWidth = minWidth
	c.minHeight = minHeight
	c.checkScrolls()
}

func (c *Widget) SetMaxSize(maxWidth, maxHeight int) {
	c.maxWidth = maxWidth
	c.maxHeight = maxHeight
	c.checkScrolls()
}

func (c *Widget) SetAnchors(left, top, right, bottom bool) {
	c.anchorLeft = left
	c.anchorTop = top
	c.anchorRight = right
	c.anchorBottom = bottom
}

func (c *Widget) SetOnClick(f func(button nuimouse.MouseButton, x int, y int)) {
	c.onClick = f
}

func (c *Widget) SetOnPaint(f func(cnv *Canvas)) {
	c.onCustomPaint = f
}

func (c *Widget) SetOnMouseDown(f func(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers)) {
	c.onMouseDown = f
}

func (c *Widget) SetOnMouseUp(f func(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers)) {
	c.onMouseUp = f
}

func (c *Widget) SetOnMouseMove(f func(x int, y int, mods nuikey.KeyModifiers)) {
	c.onMouseMove = f
}

func (c *Widget) SetOnMouseLeave(f func()) {
	c.onMouseLeave = f
}

func (c *Widget) SetOnMouseEnter(f func()) {
	c.onMouseEnter = f
}

func (c *Widget) SetOnKeyDown(f func(key nuikey.Key, mods nuikey.KeyModifiers)) {
	c.onKeyDown = f
}

func (c *Widget) SetOnKeyUp(f func(key nuikey.Key, mods nuikey.KeyModifiers)) {
	c.onKeyUp = f
}

func (c *Widget) SetOnChar(f func(char rune, mods nuikey.KeyModifiers)) {
	c.onChar = f
}

func (c *Widget) SetOnMouseWheel(f func(deltaX, deltaY int)) {
	c.onMouseWheel = f
}

func (c *Widget) ScrollEnsureVisible(x1, y1 int) {

	if y1 < c.scrollY {
		c.scrollY = y1
	}
	if y1 > c.scrollY+c.Height() {
		c.scrollY = y1 - c.Height()
	}

	if x1 < c.scrollX {
		c.scrollX = x1
	}
	if x1 > c.scrollX+c.Width() {
		c.scrollX = x1 - c.Width()
	}
}

func (c *Widget) getWidgetAt(x, y int) Widgeter {
	for _, w := range c.widgets {
		innerWidth := w.Width()
		innerHeight := w.Height()
		//innerWidth := w.InnerWidth()
		//innerHeight := w.InnerHeight()
		if x >= w.X() && x < w.X()+innerWidth && y >= w.Y() && y < w.Y()+innerHeight {
			//fmt.Println("Widget found at", w.Name(), "at position", w.X(), w.Y(), "with size", innerWidth, innerHeight)
			return w
		}
	}
	return nil
}

func (c *Widget) findWidgetAt(x, y int) Widgeter {

	// if it is the bar area, return self
	if c.allowScrollX && c.innerWidth > c.w && y >= c.h-c.scrollBarXSize {
		return c
	}
	if c.allowScrollY && c.innerHeight > c.h && x >= c.w-c.scrollBarYSize {
		return c
	}

	x += c.scrollX
	y += c.scrollY

	innerWidget := c.getWidgetAt(x, y)
	if innerWidget != nil {
		return innerWidget.findWidgetAt(x-innerWidget.X(), y-innerWidget.Y())
	}
	return WidgetById(c.Id())
	//return c
}

func (c *Widget) processPaint(cnv *Canvas) {
	// Draw the background color if set
	if c.backgroundColor.A > 0 {
		cnv.SetColor(c.backgroundColor)
		cnv.FillRect(0, 0, c.w, c.h, c.backgroundColor)
	}

	// Draw using the custom paint function if set
	cnv.Save()

	cnv.state.translateX -= c.scrollX
	cnv.state.translateY -= c.scrollY

	if c.onCustomPaint != nil {
		c.onCustomPaint(cnv)
	}

	// Draw all child widgets
	for _, w := range c.widgets {

		cnv.Save()
		cnv.TranslateAndClip(w.X(), w.Y(), w.Width(), w.Height())
		w.processPaint(cnv)
		cnv.Restore()
	}

	cnv.Restore()

	// Draw ScrollBarX
	if !c.hideScrollbarX && c.allowScrollX && c.innerWidth > c.w {
		scrollBarWidth := c.w * c.w / c.innerWidth
		scrollBarX := c.scrollX * (c.w - scrollBarWidth) / (c.innerWidth - c.w)

		barColor := c.scrollBarXColor
		if c.lastMouseAbsPosY >= c.h-c.scrollBarXSize && c.lastMouseAbsPosY < c.h {
			barColor = color.RGBA{R: barColor.R, G: barColor.G, B: barColor.B, A: 200} // Darker color when hovered
		}

		cnv.FillRect(scrollBarX, c.h-c.scrollBarXSize, scrollBarWidth, c.scrollBarXSize, barColor)
	}

	// Draw ScrollBarY
	if !c.hideScrollbarY && c.allowScrollY && c.innerHeight > c.h {
		scrollBarHeight := c.h * c.h / c.innerHeight
		scrollBarY := c.scrollY * (c.h - scrollBarHeight) / (c.innerHeight - c.h)

		barColor := c.scrollBarYColor
		if c.lastMouseAbsPosX >= c.w-c.scrollBarYSize && c.lastMouseAbsPosX < c.w {
			barColor = color.RGBA{R: barColor.R, G: barColor.G, B: barColor.B, A: 200} // Darker color when hovered
		}

		cnv.FillRect(c.w-c.scrollBarYSize, scrollBarY, c.scrollBarYSize, scrollBarHeight, barColor)
	}

}

func (c *Widget) processMouseDown(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) {
	// Determine if the click is within the horizontal scroll bar area
	if c.allowScrollX && c.innerWidth > c.w && y >= c.h-c.scrollBarXSize {
		isLeftBar := x < c.w*c.scrollX/c.innerWidth
		if isLeftBar {
			// Clicked in the left part of the scroll bar
			pageSize := c.w * c.w / c.innerWidth
			c.scrollX -= pageSize // Scroll left
			if c.scrollX < 0 {
				c.scrollX = 0
			}
			c.checkScrolls()
			return
		}

		isRightBar := x >= c.w*(c.scrollX+c.w)/c.innerWidth
		if isRightBar {
			// Clicked in the right part of the scroll bar
			pageSize := c.w * c.w / c.innerWidth
			c.scrollX += pageSize // Scroll right
			if c.scrollX > c.innerWidth-c.w {
				c.scrollX = c.innerWidth - c.w
			}
			c.checkScrolls()
			return
		}

		// Clicked in the scroll bar
		scrollBarWidth := c.w * c.w / c.innerWidth
		scrollBarX := c.scrollX * (c.w - scrollBarWidth) / (c.innerWidth - c.w)
		if x >= scrollBarX && x < scrollBarX+scrollBarWidth {
			c.scrollingX = true
			c.scrollingXInitial = c.scrollX
			c.scrollingXInitialMousePos = x
			fmt.Println("Started scrollingX", c.scrollingXInitial, c.scrollingXInitialMousePos)
			return
		}
	}

	// Determine if the click is within the vertical scroll bar area
	if c.allowScrollY && c.innerHeight > c.h && x >= c.w-c.scrollBarYSize {
		isUpperBar := y < c.h*c.scrollY/c.innerHeight
		if isUpperBar {
			// Clicked in the upper part of the scroll bar
			pageSize := c.h * c.h / c.innerHeight
			c.scrollY -= pageSize // Scroll up
			if c.scrollY < 0 {
				c.scrollY = 0
			}
			c.checkScrolls()
			return
		}

		isLowerBar := y >= c.h*(c.scrollY+c.h)/c.innerHeight
		if isLowerBar {
			// Clicked in the lower part of the scroll bar
			pageSize := c.h * c.h / c.innerHeight
			c.scrollY += pageSize // Scroll down
			if c.scrollY > c.innerHeight-c.h {
				c.scrollY = c.innerHeight - c.h
			}
			c.checkScrolls()
			return
		}

		// Clicked in the scroll bar
		scrollBarHeight := c.h * c.h / c.innerHeight
		scrollBarY := c.scrollY * (c.h - scrollBarHeight) / (c.innerHeight - c.h)
		if y >= scrollBarY && y < scrollBarY+scrollBarHeight {
			c.scrollingY = true
			c.scrollingYInitial = c.scrollY
			c.scrollingYInitialMousePos = y
			fmt.Println("Started scrollingY", c.scrollingYInitial, c.scrollingYInitialMousePos)
			return
		}
	}

	x += c.scrollX
	y += c.scrollY
	if c.onMouseDown != nil {
		c.onMouseDown(button, x, y, mods)
	}

	for _, w := range c.widgets {
		if x >= w.X() && x < w.X()+w.Width() && y >= w.Y() && y < w.Y()+w.Height() {
			w.processMouseDown(button, x-w.X(), y-w.Y(), mods)
		}
	}

	//c.focused = true
}

func (c *Widget) processMouseUp(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) {
	// If scrolling is active, stop it
	if c.scrollingX {
		c.scrollingX = false
		return
	}

	// If scrolling is active, stop it
	if c.scrollingY {
		c.scrollingY = false
		return
	}

	x += c.scrollX
	y += c.scrollY

	if c.onMouseUp != nil {
		c.onMouseUp(button, x, y, mods)
	}

	for _, w := range c.widgets {
		w.processMouseUp(button, x-w.X(), y-w.Y(), mods)
	}
}

func (c *Widget) processMouseMove(x int, y int, mods nuikey.KeyModifiers) {
	if c.scrollingX {
		if c.allowScrollX && c.innerWidth > c.w {
			k := float64(c.innerWidth) / float64(c.w)
			newScrollFloat64 := float64(c.scrollingXInitial) + float64(x-c.scrollingXInitialMousePos)*k
			c.scrollX = int(newScrollFloat64)
			c.checkScrolls()
			return
		}
		return
	}

	if c.scrollingY {
		if c.allowScrollY && c.innerHeight > c.h {
			k := float64(c.innerHeight) / float64(c.h)
			newScrollFloat64 := float64(c.scrollingYInitial) + float64(y-c.scrollingYInitialMousePos)*k
			c.scrollY = int(newScrollFloat64)
			c.checkScrolls()
			return
		}
		return
	}

	c.lastMouseAbsPosX = x
	c.lastMouseAbsPosY = y

	x += c.scrollX
	y += c.scrollY

	c.lastMouseX = x
	c.lastMouseY = y

	if c.onMouseMove != nil {
		c.onMouseMove(x, y, mods)
	}

	for _, w := range c.widgets {
		// Temporary process in the all widgets - perrormance issue
		inWidget := true
		//inWidget := x >= w.X() && x < w.X()+w.Width() && y >= w.Y() && y < w.Y()+w.Height()
		if inWidget {
			w.processMouseMove(x-w.X(), y-w.Y(), mods)
		}
	}
}

func (c *Widget) processMouseLeave() {
	if c.onMouseLeave != nil {
		c.onMouseLeave()
	}
	MainForm.Update()
}

func (c *Widget) processMouseEnter() {
	if c.onMouseEnter != nil {
		c.onMouseEnter()
	}
	MainForm.Update()
}

func (c *Widget) processKeyDown(key nuikey.Key, mods nuikey.KeyModifiers) {
	if c.onKeyDown != nil {
		c.onKeyDown(key, mods)
	}

	for _, w := range c.widgets {
		w.processKeyDown(key, mods)
	}
}

func (c *Widget) processKeyUp(key nuikey.Key, mods nuikey.KeyModifiers) {
	if c.onKeyUp != nil {
		c.onKeyUp(key, mods)
	}

	for _, w := range c.widgets {
		w.processKeyUp(key, mods)
	}
}

func (c *Widget) processMouseDblClick(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) {
	x += c.scrollX
	y += c.scrollY

	if c.onMouseUp != nil {
		c.onMouseUp(button, x, y, mods)
	}

	for _, w := range c.widgets {
		if x >= w.X() && x < w.X()+w.Width() && y >= w.Y() && y < w.Y()+w.Height() {
			w.processMouseDblClick(button, x-w.X(), y-w.Y(), mods)
		}
	}
}

func (c *Widget) processChar(char rune, mods nuikey.KeyModifiers) {
	if c.onChar != nil {
		c.onChar(char, mods)
	}

	for _, w := range c.widgets {
		w.processChar(char, mods)
	}
}

func (c *Widget) processMouseWheel(deltaX, deltaY int) bool {
	hoverWidget := c.getWidgetAt(c.lastMouseX, c.lastMouseY)
	if hoverWidget != nil {
		processed := hoverWidget.processMouseWheel(deltaX, deltaY)
		if processed {
			return true
		}
	}

	if c.onMouseWheel != nil {
		c.onMouseWheel(deltaX, deltaY)
		return true
	}

	if deltaY != 0 && c.allowScrollY && c.InnerHeight() > c.h {
		c.scrollY -= deltaY * 30 // Adjust the scroll speed as needed
		c.checkScrolls()
		return true
	}

	if deltaX != 0 && c.allowScrollX && c.InnerWidth() > c.w {
		c.scrollX -= deltaX * 30 // Adjust the scroll speed as needed
		c.checkScrolls()
		return true
	}

	/*if (c.allowScrollX || c.allowScrollY) && (c.innerWidth > c.w || c.innerHeight > c.h) {
		if c.allowScrollX {
			c.scrollX -= deltaX * 30
		}
		if c.allowScrollY {
			c.scrollY -= deltaY * 30
			fmt.Println("WidgetName:", c.name, "ScrollY:", c.scrollY)
		}
		c.checkScrolls()
		return true
	}*/

	return false
}

func (c *Widget) processTimer() {
	for _, t := range c.timers {
		t.tick()
	}

	for _, w := range c.widgets {
		w.processTimer()
	}
}

func (c *Widget) checkScrolls() {
	if c.allowScrollX {
		if c.scrollX > c.innerWidth-c.w {
			c.scrollX = c.innerWidth - c.w
		}
		if c.scrollX < 0 {
			c.scrollX = 0
		}
		if c.innerWidth < c.w {
			c.scrollX = 0
		}
	}

	if c.allowScrollY {
		if c.scrollY > c.innerHeight-c.h {
			c.scrollY = c.innerHeight - c.h
		}
		if c.scrollY < 0 {
			c.scrollY = 0
		}
		if c.innerHeight < c.h {
			c.scrollY = 0
		}
	}
}

func (c *Widget) Anchors() (left, top, right, bottom bool) {
	return c.anchorLeft, c.anchorTop, c.anchorRight, c.anchorBottom
}

func (c *Widget) SetAbsolutePositioning(absolute bool) {
	c.absolutePositioning = absolute
}

func (c *Widget) SetXExpandable(expandable bool) {
	c.xExpandable = expandable
}

func (c *Widget) SetYExpandable(expandable bool) {
	c.yExpandable = expandable
}

func (c *Widget) SetMinWidth(minWidth int) {
	c.minWidth = minWidth
}

func (c *Widget) SetMinHeight(minHeight int) {
	c.minHeight = minHeight
}

func (c *Widget) SetMaxWidth(maxWidth int) {
	c.maxWidth = maxWidth
}

func (c *Widget) SetMaxHeight(maxHeight int) {
	c.maxHeight = maxHeight
}

func (c *Widget) updateLayout(oldWidth, oldHeight, newWidth, newHeight int) {
	if c.absolutePositioning {
		for _, w := range c.widgets {
			deltaWidth := newWidth - oldWidth
			deltaHeight := newHeight - oldHeight

			newX := w.X()
			newY := w.Y()
			newW := w.Width()
			newH := w.Height()

			anchorLeft, anchorTop, anchorRight, anchorBottom := w.Anchors()

			if anchorLeft && anchorRight {
				newW += deltaWidth
			}
			if !anchorLeft && anchorRight {
				newX += deltaWidth
			}

			if anchorTop && anchorBottom {
				newH += deltaHeight
			}

			if !anchorTop && anchorBottom {
				newY += deltaHeight
			}

			w.SetSize(newW, newH)
			w.SetPosition(newX, newY)
		}
	} else {
		fullWidth := c.w
		fullHeight := c.h

		_, minX, maxX, allCellPaddingX := c.makeColumnsInfo(fullWidth)
		columnsInfo, _, _, _ := c.makeColumnsInfo(fullWidth - (c.panelPadding + allCellPaddingX + c.panelPadding))

		_, minY, maxY, allCellPaddingY := c.makeRowsInfo(fullHeight)
		rowsInfo, _, _, _ := c.makeRowsInfo(fullHeight - (c.panelPadding + allCellPaddingY + c.panelPadding))

		if strings.Contains(c.name, "Top") {
			fmt.Println("RowsInfo:")
			for yy := minY; yy <= maxY; yy++ {
				if rowInfo, ok := rowsInfo[yy]; ok {
					fmt.Printf("Row %d: minHeight=%d, maxHeight=%d, expandable=%t, height=%d, collapsed=%t\n",
						yy, rowInfo.minHeight, rowInfo.maxHeight, rowInfo.expandable, rowInfo.height, rowInfo.collapsed)
				}
			}
		}

		xOffset := c.panelPadding //+ c.LeftBorderWidth()
		for x := minX; x <= maxX; x++ {
			if colInfo, ok := columnsInfo[x]; ok {
				yOffset := c.panelPadding // + c.TopBorderWidth()
				for y := minY; y <= maxY; y++ {
					if rowInfo, ok := rowsInfo[y]; ok {
						w := c.getWidgetInGridCell(x, y)
						if w != nil {

							cX := xOffset
							cY := yOffset

							wWidth := colInfo.width
							if wWidth > w.MaxWidth() {
								wWidth = w.MaxWidth()
							}
							wHeight := rowInfo.height
							if wHeight > w.MaxHeight() {
								wHeight = w.MaxHeight()
							}

							// Place widget in the center of the cell
							//cX += (colInfo.width - wWidth) / 2
							//cY += (rowInfo.height - wHeight) / 2

							w.SetPosition(cX, cY)

							if w.IsVisible() {
								w.SetSize(wWidth, wHeight)
							} else {
								w.SetSize(0, 0)
							}
						}

						yOffset += rowInfo.height
						if rowInfo.height > 0 && y < maxY {
							yOffset += c.cellPadding
						}
					}
				}

				xOffset += colInfo.width
				if colInfo.width > 0 && x < maxX {
					xOffset += c.cellPadding
				}
			}
		}

		for _, w := range c.widgets {
			if !w.IsVisible() {
				w.SetSize(0, 0)
			}
		}

		if len(c.widgets) > 0 {
			// Set InnerSize
			innerWidth := 0
			innerHeight := 0

			for _, w := range c.widgets {
				if w.IsVisible() {
					if w.X()+w.Width() > innerWidth {
						innerWidth = w.X() + w.Width()
					}
					if w.Y()+w.Height() > innerHeight {
						innerHeight = w.Y() + w.Height()
					}
				}
			}

			if innerWidth < c.w {
				innerWidth = c.w
			}
			if innerHeight < c.h {
				innerHeight = c.h
			}
			c.innerWidth = innerWidth
			c.innerHeight = innerHeight
			c.checkScrolls()
		}

	}

	// fmt.Println("Widget", c.name, "layout updated:", "Width:", c.w, "Height:", c.h, "InnerWidth:", c.innerWidth, "InnerHeight:", c.innerHeight)
}

func (c *Widget) makeColumnsInfo(fullWidth int) (map[int]*ContainerGridColumnInfo, int, int, int) {

	minX := MaxInt
	minY := MaxInt

	maxX := MinInt
	maxY := MinInt

	// Detect range of grid coordinates
	for _, w := range c.widgets {
		if w.GridX() < minX {
			minX = w.GridX()
		}
		if w.GridX() > maxX {
			maxX = w.GridX()
		}
		if w.GridY() < minY {
			minY = w.GridY()
		}
		if w.GridY() > maxY {
			maxY = w.GridY()
		}
	}

	columnsInfo := make(map[int]*ContainerGridColumnInfo)
	hasExpandableColumns := false

	// Fill columnsInfo
	for x := minX; x <= maxX; x++ {
		var colInfo ContainerGridColumnInfo
		colInfo.minWidth = MinInt
		colInfo.maxWidth = MaxInt
		colInfo.expandable = false
		found := false

		for y := minY; y <= maxY; y++ {
			w := c.getWidgetInGridCell(x, y)
			if w != nil {
				if w.XExpandable() {
					colInfo.expandable = true // Found expandable by X
					hasExpandableColumns = true
				}
				found = true
			}
		}

		if colInfo.expandable {
			colInfo.minWidth = MinInt
			colInfo.maxWidth = MinInt

			for y := minY; y <= maxY; y++ {
				w := c.getWidgetInGridCell(x, y)
				if w != nil {
					wMinWidth := w.MinWidth()
					if wMinWidth > colInfo.minWidth {
						colInfo.minWidth = wMinWidth
					}
					wMaxWidth := w.MaxWidth()
					if wMaxWidth > colInfo.maxWidth {
						colInfo.maxWidth = wMaxWidth
					}
				}
			}

		} else {
			colInfo.minWidth = MinInt
			colInfo.maxWidth = MinInt

			for y := minY; y <= maxY; y++ {
				w := c.getWidgetInGridCell(x, y)
				if w != nil {
					wMinWidth := w.MinWidth()
					if wMinWidth > colInfo.minWidth {
						colInfo.minWidth = w.MinWidth()
					}
					if wMinWidth > colInfo.maxWidth {
						colInfo.maxWidth = w.MaxWidth()
					}
					/*if w.MaxWidth() < colInfo.maxWidth {
						colInfo.maxWidth = w.MaxWidth()
					}*/
				}
			}
		}

		if found {
			columnsInfo[x] = &colInfo
		}
	}

	if hasExpandableColumns {
		hasNonExpandable := false
		for _, colInfo := range columnsInfo {
			if !colInfo.expandable {
				hasNonExpandable = true
				break
			}
		}
		if hasNonExpandable {
			for _, colInfo := range columnsInfo {
				if !colInfo.expandable {
					colInfo.width = colInfo.minWidth
					colInfo.collapsed = true
				}
			}
		}
	}

	width := fullWidth

	for {
		readyWidth := 0
		for _, colInfo := range columnsInfo {
			readyWidth += colInfo.width
		}
		deltaWidth := width - readyWidth
		countOfColumnCanChange := 0
		for _, colInfo := range columnsInfo {
			if deltaWidth > 0 {
				if colInfo.width < colInfo.maxWidth {
					if !colInfo.collapsed {
						countOfColumnCanChange++
					}
				}
			} else {
				if deltaWidth < 0 {
					if colInfo.width > colInfo.minWidth {
						if !colInfo.collapsed {
							countOfColumnCanChange++
						}
					}
				}
			}
		}

		if countOfColumnCanChange > 0 && deltaWidth != 0 {
			pixForOne := deltaWidth / countOfColumnCanChange
			if math.Abs(float64(pixForOne)) < 1 {
				break
			}
			for _, colInfo := range columnsInfo {
				if !colInfo.collapsed {
					colInfo.width += pixForOne
				}
			}
		} else {
			break
		}

		for _, colInfo := range columnsInfo {
			if colInfo.width > colInfo.maxWidth {
				colInfo.width = colInfo.maxWidth
			}
			if colInfo.width < colInfo.minWidth {
				colInfo.width = colInfo.minWidth
			}
		}
	}

	allCellPadding := 0
	for _, colInfo := range columnsInfo {
		if colInfo.width > 0 {
			allCellPadding++
		}
	}
	allCellPadding--
	allCellPadding *= c.cellPadding
	if allCellPadding < 0 {
		allCellPadding = 0
	}

	return columnsInfo, minX, maxX, allCellPadding

}

func (c *Widget) makeRowsInfo(fullHeight int) (map[int]*ContainerGridRowInfo, int, int, int) {

	// Определяем минимальный и максимальный индекс строк
	minX := MaxInt
	minY := MaxInt
	maxX := MinInt
	maxY := MinInt
	for _, w := range c.widgets {
		if w.GridX() < minX {
			minX = w.GridX()
		}
		if w.GridX() > maxX {
			maxX = w.GridX()
		}
		if w.GridY() < minY {
			minY = w.GridY()
		}
		if w.GridY() > maxY {
			maxY = w.GridY()
		}
	}

	// Подготовка
	rowsInfo := make(map[int]*ContainerGridRowInfo)
	hasExpandableRows := false

	// Главный цикл по строкам
	for y := minY; y <= maxY; y++ {
		var rowInfo ContainerGridRowInfo
		rowInfo.minHeight = MinInt // Минимальная высота строки пока 0
		rowInfo.maxHeight = MaxInt // Максимальная высота строки пока ... максимум
		rowInfo.expandable = false // Пока думаем, что строка не мажорная
		found := false             // Признак того, что вообще есть в строке контролы

		// If any widget in the row is expandable, set the expandable flag for the row
		for x := minX; x <= maxX; x++ {
			w := c.getWidgetInGridCell(x, y)
			if w != nil {
				if w.YExpandable() {
					rowInfo.expandable = true // Found expandable by Y
					hasExpandableRows = true
				}
				found = true
			}
		}

		if rowInfo.expandable {
			rowInfo.minHeight = MinInt
			rowInfo.maxHeight = MinInt

			for x := minX; x <= maxX; x++ {
				w := c.getWidgetInGridCell(x, y)
				if w != nil {
					wMinHeight := w.MinHeight()
					if wMinHeight > rowInfo.minHeight {
						rowInfo.minHeight = wMinHeight
					}
					wMaxHeight := w.MaxHeight()
					if wMaxHeight > rowInfo.maxHeight {
						rowInfo.maxHeight = wMaxHeight
					}
				}
			}

		} else {
			rowInfo.minHeight = MinInt
			rowInfo.maxHeight = MinInt

			for x := minX; x <= maxX; x++ {
				w := c.getWidgetInGridCell(x, y)
				if w != nil {
					wMinHeight := w.MinHeight()
					if wMinHeight > rowInfo.minHeight {
						rowInfo.minHeight = wMinHeight
					}
					if wMinHeight > rowInfo.maxHeight {
						rowInfo.maxHeight = w.MaxHeight()
					}
					/*if w.MaxWidth() < colInfo.maxWidth {
						colInfo.maxWidth = w.MaxWidth()
					}*/
				}
			}
		}

		if found {
			rowsInfo[y] = &rowInfo
		}
	}

	if hasExpandableRows {
		hasNonExpandable := false
		for _, rowInfo := range rowsInfo {
			if !rowInfo.expandable {
				hasNonExpandable = true
				break
			}
		}
		if hasNonExpandable {
			for _, rowsInfo := range rowsInfo {
				if !rowsInfo.expandable {
					rowsInfo.height = rowsInfo.minHeight
					rowsInfo.collapsed = true
				}
			}
		}
	}

	height := fullHeight

	for {
		readyHeight := 0
		for _, rowInfo := range rowsInfo {
			readyHeight += rowInfo.height
		}
		deltaHeight := height - readyHeight
		countOfRowCanChange := 0
		for _, rowInfo := range rowsInfo {
			if deltaHeight > 0 {
				if rowInfo.height < rowInfo.maxHeight {
					if !rowInfo.collapsed {
						countOfRowCanChange++
					}
				}
			} else {
				if deltaHeight < 0 {
					if rowInfo.height > rowInfo.minHeight {
						if !rowInfo.collapsed {
							countOfRowCanChange++
						}
					}
				}
			}
		}

		if countOfRowCanChange > 0 && deltaHeight != 0 {
			pixForOne := deltaHeight / countOfRowCanChange
			if math.Abs(float64(pixForOne)) < 1 {
				break
			}
			for _, rowInfo := range rowsInfo {
				if !rowInfo.collapsed {
					rowInfo.height += pixForOne
				}
			}
		} else {
			break
		}

		for _, rowInfo := range rowsInfo {
			if rowInfo.height > rowInfo.maxHeight {
				rowInfo.height = rowInfo.maxHeight
			}
			if rowInfo.height < rowInfo.minHeight {
				rowInfo.height = rowInfo.minHeight
			}
		}
	}

	allCellPadding := 0
	for _, rowInfo := range rowsInfo {
		if rowInfo.height > 0 {
			allCellPadding++
		}
	}
	allCellPadding--
	allCellPadding *= c.cellPadding
	if allCellPadding < 0 {
		allCellPadding = 0
	}

	return rowsInfo, minY, maxY, allCellPadding
}

func (c *Widget) getWidgetInGridCell(x, y int) Widgeter {
	for _, w := range c.widgets {
		if w.GridX() == x && w.GridY() == y {
			if w.IsVisible() {
				return w
			}
		}
	}
	return nil
}

func (c *Widget) XExpandable() bool {
	if len(c.widgets) == 0 {
		return c.xExpandable
	}

	colsInfo, _, _, _ := c.makeColumnsInfo(1000)
	for _, ci := range colsInfo {
		if ci.expandable {
			return true
		}
	}

	return false
}

func (c *Widget) YExpandable() bool {
	if len(c.widgets) == 0 {
		return c.yExpandable
	}

	rowsInfo, _, _, _ := c.makeRowsInfo(1000)
	for _, ri := range rowsInfo {
		if ri.expandable {
			return true
		}
	}

	return false
}
