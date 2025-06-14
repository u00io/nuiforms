package ui

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"image/color"
	"strings"

	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
)

type Widget struct {
	name     string
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
	widgets []*Widget

	allowScrollX bool
	allowScrollY bool
	scrollX      int
	scrollY      int
	innerWidth   int
	innerHeight  int

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

func NewWidget() *Widget {
	var c Widget
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)
	c.name = "Widget-" + strings.ToUpper(hex.EncodeToString(randomBytes))
	c.props = make(map[string]any)
	c.timers = make([]*timer, 0)
	c.x = 0
	c.y = 0
	c.w = 300
	c.h = 200
	c.scrollBarXSize = 10
	c.scrollBarYSize = 10
	c.scrollBarXColor = color.RGBA{R: 150, G: 150, B: 150, A: 100}
	c.scrollBarYColor = color.RGBA{R: 150, G: 150, B: 150, A: 100}
	c.innerWidth = 0
	c.anchorLeft = true
	c.anchorTop = true
	c.anchorRight = false
	c.anchorBottom = false
	c.widgets = make([]*Widget, 0)
	c.mouseCursor = nuimouse.MouseCursorArrow
	c.backgroundColor = color.RGBA{R: 0, G: 0, B: 0, A: 0} // transparent by default
	return &c
}

func (c *Widget) SetName(name string) {
	c.name = name
}

func (c *Widget) AddTimer(intervalMs int, callback func()) {
	t := &timer{
		intervalMs: intervalMs,
		callback:   callback,
	}
	c.timers = append(c.timers, t)
}

func (c *Widget) AddWidget(w *Widget) {
	c.widgets = append(c.widgets, w)
}

func (c *Widget) RemoveWidget(w *Widget) {
	for i, widget := range c.widgets {
		if widget == w {
			c.widgets = append(c.widgets[:i], c.widgets[i+1:]...)
			return
		}
	}
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

func (c *Widget) X() int {
	return c.x
}

func (c *Widget) Y() int {
	return c.y
}

func (c *Widget) W() int {
	return c.w
}

func (c *Widget) H() int {
	return c.h
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
	MainForm.focusedWidget = c
	MainForm.Update()
}

func (c *Widget) IsFocused() bool {
	return MainForm.focusedWidget == c
}

func (c *Widget) IsHovered() bool {
	return MainForm.hoverWidget == c
}

func (c *Widget) Name() string {
	return c.name
}

func (c *Widget) SetPosition(x, y int) {
	c.x = x
	c.y = y
}

func (c *Widget) SetSize(w, h int) {
	c.updateLayout(c.w, c.h, w, h)
	c.w = w
	c.h = h
	c.checkScrolls() // Update scroll position if needed
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

func (c *Widget) getWidgetAt(x, y int) *Widget {
	x += c.scrollX
	y += c.scrollY

	for _, w := range c.widgets {
		if x >= w.x && x < w.x+w.w && y >= w.y && y < w.y+w.h {
			return w
		}
	}
	return nil
}

func (c *Widget) findWidgetAt(x, y int) *Widget {
	// if it is the bar area, return self
	if c.allowScrollX && c.innerWidth > c.w && y >= c.h-c.scrollBarXSize {
		return c
	}
	if c.allowScrollY && c.innerHeight > c.h && x >= c.w-c.scrollBarYSize {
		return c
	}

	innerWidget := c.getWidgetAt(x, y)
	if innerWidget != nil {
		return innerWidget.findWidgetAt(x-innerWidget.x, y-innerWidget.y)
	}
	return c
}

func (c *Widget) processPaint(cnv *Canvas) {
	// Draw the background color if set
	if c.backgroundColor.A > 0 {
		cnv.SetColor(c.backgroundColor)
		cnv.FillRect(0, 0, c.w, c.h, c.backgroundColor)
	}

	// Draw using the custom paint function if set
	cnv.Save()
	cnv.Translate(-c.scrollX, -c.scrollY)

	if c.onCustomPaint != nil {
		c.onCustomPaint(cnv)
	}

	// Draw all child widgets
	for _, w := range c.widgets {
		cnv.Save()
		cnv.Translate(w.x, w.y)
		//cnv.SetClip(0, 0, 1000, 1000)
		w.processPaint(cnv)
		cnv.Restore()
	}

	cnv.Restore()

	// Draw ScrollBarX
	if c.allowScrollX && c.innerWidth > c.w {
		scrollBarWidth := c.w * c.w / c.innerWidth
		scrollBarX := c.scrollX * (c.w - scrollBarWidth) / (c.innerWidth - c.w)

		barColor := c.scrollBarXColor
		if c.lastMouseAbsPosY >= c.h-c.scrollBarXSize && c.lastMouseAbsPosY < c.h {
			barColor = color.RGBA{R: barColor.R, G: barColor.G, B: barColor.B, A: 200} // Darker color when hovered
		}

		cnv.FillRect(scrollBarX, c.h-c.scrollBarXSize, scrollBarWidth, c.scrollBarXSize, barColor)
	}

	// Draw ScrollBarY
	if c.allowScrollY && c.innerHeight > c.h {
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
		if x >= w.x && x < w.x+w.w && y >= w.y && y < w.y+w.h {
			w.processMouseDown(button, x-w.x, y-w.y, mods)
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
		w.processMouseUp(button, x-w.x, y-w.y, mods)
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
		if x >= w.x && x < w.x+w.w && y >= w.y && y < w.y+w.h {
			w.processMouseMove(x-w.x, y-w.y, mods)
		}
	}
}

func (c *Widget) processMouseLeave() {
	fmt.Println("Widget", c.name, "mouse leave")
	if c.onMouseLeave != nil {
		c.onMouseLeave()
	}
	MainForm.Update()
}

func (c *Widget) processMouseEnter() {
	fmt.Println("Widget", c.name, "mouse enter")
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
		if x >= w.x && x < w.x+w.w && y >= w.y && y < w.y+w.h {
			w.processMouseDblClick(button, x-w.x, y-w.y, mods)
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

	if c.allowScrollX || c.allowScrollY {
		if c.allowScrollX {
			c.scrollX -= deltaX * 30
		}
		if c.allowScrollY {
			c.scrollY -= deltaY * 30
		}
		c.checkScrolls()
		return true
	}

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

func (c *Widget) updateLayout(oldWidth, oldHeight, newWidth, newHeight int) {
	for _, w := range c.widgets {
		deltaWidth := newWidth - oldWidth
		deltaHeight := newHeight - oldHeight

		newX := w.X()
		newY := w.Y()
		newW := w.W()
		newH := w.H()

		if w.anchorLeft && w.anchorRight {
			newW += deltaWidth
		}
		if !w.anchorLeft && w.anchorRight {
			newX += deltaWidth
		}

		if w.anchorTop && w.anchorBottom {
			newH += deltaHeight
		}

		if !w.anchorTop && w.anchorBottom {
			newY += deltaHeight
		}

		w.SetSize(newW, newH)
		w.SetPosition(newX, newY)
	}
}
