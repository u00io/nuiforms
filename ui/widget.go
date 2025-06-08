package ui

import (
	"crypto/rand"
	"encoding/hex"
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

	props map[string]interface{}

	// temp
	lastMouseX int
	lastMouseY int

	mouseCursor nuimouse.MouseCursor

	// callbacks
	onCustomPaint func(cnv *Canvas)
	onMouseDown   func(button nuimouse.MouseButton, x int, y int)
	onMouseUp     func(button nuimouse.MouseButton, x int, y int)
	onMouseMove   func(x int, y int)
	onMouseLeave  func()
	onMouseEnter  func()
	onKeyDown     func(key nuikey.Key)
	onKeyUp       func(key nuikey.Key)
	onMouseWheel  func(deltaX, deltaY int)
	onClick       func(button nuimouse.MouseButton, x int, y int)
}

func NewWidget() *Widget {
	var c Widget
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)
	c.name = "Widget-" + strings.ToUpper(hex.EncodeToString(randomBytes))
	c.props = make(map[string]interface{})
	return &c
}

func (c *Widget) SetName(name string) {
	c.name = name
}

func (c *Widget) AddWidget(w *Widget) {
	c.widgets = append(c.widgets, w)
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

func (c *Widget) SetOnMouseDown(f func(button nuimouse.MouseButton, x int, y int)) {
	c.onMouseDown = f
}

func (c *Widget) SetOnMouseUp(f func(button nuimouse.MouseButton, x int, y int)) {
	c.onMouseUp = f
}

func (c *Widget) SetOnMouseMove(f func(x int, y int)) {
	c.onMouseMove = f
}

func (c *Widget) SetOnMouseLeave(f func()) {
	c.onMouseLeave = f
}

func (c *Widget) SetOnMouseEnter(f func()) {
	c.onMouseEnter = f
}

func (c *Widget) SetOnKeyDown(f func(key nuikey.Key)) {
	c.onKeyDown = f
}

func (c *Widget) SetOnKeyUp(f func(key nuikey.Key)) {
	c.onKeyUp = f
}

func (c *Widget) SetOnMouseWheel(f func(deltaX, deltaY int)) {
	c.onMouseWheel = f
}

func (c *Widget) getWidgetAt(x, y int) *Widget {
	for _, w := range c.widgets {
		if x >= w.x && x < w.x+w.w && y >= w.y && y < w.y+w.h {
			return w
		}
	}
	return nil
}

func (c *Widget) findWidgetAt(x, y int) *Widget {
	innerWidget := c.getWidgetAt(x, y)
	if innerWidget != nil {
		return innerWidget.findWidgetAt(x-innerWidget.x, y-innerWidget.y)
	}
	return c
}

func (c *Widget) processPaint(cnv *Canvas) {
	if c.onCustomPaint != nil {
		c.onCustomPaint(cnv)
	}

	for _, w := range c.widgets {
		cnv.Save()
		cnv.Translate(w.x, w.y)
		cnv.SetClip(w.x, w.y, w.w, w.h)
		w.processPaint(cnv)
		cnv.Restore()
	}

	/*if c.IsHovered() {
		cnv.SetColor(color.RGBA{255, 255, 255, 255})
		for x := 0; x < c.w; x++ {
			cnv.SetPixel(x, 0)
			cnv.SetPixel(x, 1)
		}
	}*/

	/*if c.IsFocused() {
		cnv.SetColor(color.RGBA{0, 255, 0, 255})
		for x := 0; x < c.w; x++ {
			cnv.SetPixel(x, 0)
			cnv.SetPixel(x, c.h-1)
		}
		for y := 0; y < c.h; y++ {
			cnv.SetPixel(0, y)
			cnv.SetPixel(c.w-1, y)
		}
	}*/
}

func (c *Widget) processMouseDown(button nuimouse.MouseButton, x int, y int) {
	if c.onMouseDown != nil {
		c.onMouseDown(button, x, y)
	}

	for _, w := range c.widgets {
		if x >= w.x && x < w.x+w.w && y >= w.y && y < w.y+w.h {
			w.processMouseDown(button, x-w.x, y-w.y)
		}
	}

	//c.focused = true
}

func (c *Widget) processMouseUp(button nuimouse.MouseButton, x int, y int) {
	if c.onMouseUp != nil {
		c.onMouseUp(button, x, y)
	}

	for _, w := range c.widgets {
		w.processMouseUp(button, x-w.x, y-w.y)
	}
}

func (c *Widget) processMouseMove(x int, y int) {
	c.lastMouseX = x
	c.lastMouseY = y
	if c.onMouseMove != nil {
		c.onMouseMove(x, y)
	}

	for _, w := range c.widgets {
		//widgetInArea := false
		if x >= w.x && x < w.x+w.w && y >= w.y && y < w.y+w.h {
			w.processMouseMove(x-w.x, y-w.y)
			//widgetInArea = true
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

func (c *Widget) processKeyDown(key nuikey.Key) {
	if c.onKeyDown != nil {
		c.onKeyDown(key)
	}

	for _, w := range c.widgets {
		w.processKeyDown(key)
	}
}

func (c *Widget) processKeyUp(key nuikey.Key) {
	if c.onKeyUp != nil {
		c.onKeyUp(key)
	}

	for _, w := range c.widgets {
		w.processKeyUp(key)
	}
}

func (c *Widget) processMouseDblClick(button nuimouse.MouseButton, x int, y int) {
	if c.onMouseUp != nil {
		c.onMouseUp(button, x, y)
	}

	for _, w := range c.widgets {
		if x >= w.x && x < w.x+w.w && y >= w.y && y < w.y+w.h {
			w.processMouseDblClick(button, x-w.x, y-w.y)
		}
	}
}

func (c *Widget) processChar(char rune) {
	for _, w := range c.widgets {
		w.processChar(char)
	}
}

func (c *Widget) processMouseWheel(deltaX, deltaY int) {
	if c.onMouseWheel != nil {
		c.onMouseWheel(deltaX, deltaY)
	}

	for _, w := range c.widgets {
		w.processMouseWheel(deltaX, deltaY)
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
