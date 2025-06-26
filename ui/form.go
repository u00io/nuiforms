package ui

import (
	"fmt"
	"image"
	"image/color"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/u00io/nui/nui"
	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
)

type Form struct {
	wnd    nui.Window
	title  string
	width  int
	height int

	lastMouseX      int
	lastMouseY      int
	lastMouseCursor nuimouse.MouseCursor

	lastKeyboardModifiers nuikey.KeyModifiers

	mouseLeftButtonPressed       bool
	mouseLeftButtonPressedWidget Widgeter

	topWidget     *Panel
	hoverWidget   Widgeter
	focusedWidget Widgeter

	needUpdate         bool
	lastFreeMemoryTime time.Time

	lastUpdateTime time.Time
}

var MainForm *Form
var mainFormExecuted bool

var allwidgets map[string]Widgeter
var nextId int64

func init() {
	allwidgets = make(map[string]Widgeter)
}

type Widgeter interface {
	Id() string
	TypeName() string
	Name() string
	X() int
	Y() int
	Width() int
	Height() int
	InnerWidth() int
	InnerHeight() int

	SetName(name string)
	SetPosition(x, y int)
	SetSize(width, height int)
	SetAnchors(left, top, right, bottom bool)

	getWidgetAt(x, y int) Widgeter
	findWidgetAt(x, y int) Widgeter
	Focus()

	processPaint(cnv *Canvas)
	processMouseDown(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers)
	processMouseUp(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers)
	processMouseMove(x int, y int, mods nuikey.KeyModifiers)
	processMouseLeave()
	processMouseEnter()
	processKeyDown(keyCode nuikey.Key, mods nuikey.KeyModifiers)
	processKeyUp(keyCode nuikey.Key, mods nuikey.KeyModifiers)
	processMouseDblClick(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers)
	processChar(char rune, mods nuikey.KeyModifiers)
	processMouseWheel(deltaX int, deltaY int) bool
	processTimer()

	SetMouseCursor(cursor nuimouse.MouseCursor)
	MouseCursor() nuimouse.MouseCursor

	Anchors() (left, top, right, bottom bool)

	AddWidget(widget Widgeter)
	AddWidgetOnGrid(widget Widgeter, gridX, gridY int)
	RemoveWidget(widget Widgeter)

	IsVisible() bool

	GridX() int
	GridY() int
	SetGridPosition(x, y int)

	XExpandable() bool
	YExpandable() bool

	MinWidth() int
	MinHeight() int

	MaxWidth() int
	MaxHeight() int

	SetAbsolutePositioning(absolute bool)
	SetXExpandable(expandable bool)
	SetYExpandable(expandable bool)
}

func UpdateMainForm() {
	if MainForm != nil {
		MainForm.Update()
	}
}

func WidgetById(id string) Widgeter {
	if widget, exists := allwidgets[id]; exists {
		return widget
	}
	return nil
}

func NewId() string {
	id := fmt.Sprint(nextId)
	for len(id) < 3 {
		id = "0" + id
	}
	nextId++
	return id
}

func NewForm() *Form {
	var c Form
	c.title = "Form"
	c.width = 800
	c.height = 600
	topWidget := NewPanel()
	topWidget.SetName("FormTopWidget")
	topWidget.SetPosition(0, 0)
	topWidget.SetSize(c.width, c.height)
	topWidget.SetAnchors(true, true, true, true)
	c.topWidget = topWidget
	allwidgets[topWidget.Id()] = topWidget
	if MainForm != nil {
		panic("MainForm already exists, cannot create a new one")
	}
	MainForm = &c
	return &c
}

func (c *Form) Close() {
	if c.wnd != nil {
		c.wnd.Close()
		c.wnd = nil
	}
}

func (c *Form) SetTitle(title string) {
	c.title = title
	if c.wnd != nil {
		c.wnd.SetTitle(title)
	}
}

func (c *Form) SetSize(width, height int) {
	c.width = width
	c.height = height
	if c.wnd != nil {
		c.wnd.Resize(width, height)
	}
}

func (c *Form) Panel() *Panel {
	return c.topWidget
}

func (c *Form) Exec() {
	if mainFormExecuted {
		panic("MainForm already executed, cannot execute again")
	}
	mainFormExecuted = true

	c.wnd = nui.CreateWindow(c.title, c.width, c.height, true)

	c.wnd.OnPaint(c.processPaint)
	c.wnd.OnResize(c.processResize)
	c.wnd.OnMouseButtonDown(c.processMouseDown)
	c.wnd.OnMouseButtonUp(c.processMouseUp)
	c.wnd.OnMouseButtonDblClick(c.processMouseDblClick)
	c.wnd.OnMouseMove(c.processMouseMove)
	c.wnd.OnMouseWheel(c.processMouseWheel)
	c.wnd.OnMouseLeave(c.processMouseLeave)
	c.wnd.OnMouseEnter(c.processMouseEnter)
	c.wnd.OnKeyDown(c.processKeyDown)
	c.wnd.OnKeyUp(c.processKeyUp)
	c.wnd.OnChar(c.processChar)
	c.wnd.OnTimer(c.processTimer)
	c.wnd.OnMove(c.processWindowMove)

	c.wnd.Show()
	c.wnd.EventLoop()
}

func (c *Form) realUpdate() {
	if c.wnd != nil && c.needUpdate {
		c.wnd.Update()
		c.needUpdate = false
		c.lastUpdateTime = time.Now()
	}
}

func (c *Form) Update() {
	c.needUpdate = true
	if time.Since(c.lastUpdateTime) > 50*time.Millisecond {
		c.realUpdate()
	}
}

func (c *Form) forceUpdate() {
	c.needUpdate = true
	c.realUpdate()
}

func (c *Form) processPaint(rgba *image.RGBA) {
	cnv := NewCanvas(rgba)
	cnv.SetDirectTranslateAndClip(0, 0, c.width, c.height)
	c.topWidget.processPaint(cnv)
	if c.hoverWidget != nil {
		c.DrawWidgetDebugInfo(c.hoverWidget, cnv)
	}
}

func (c *Form) DrawWidgetDebugInfo(w Widgeter, cnv *Canvas) {
	if w == nil {
		return
	}
	lines := make([]string, 0)
	posX := c.lastMouseX + 16
	posY := c.lastMouseY + 16
	col := color.RGBA{R: 0, G: 200, B: 200, A: 255}
	lines = append(lines, fmt.Sprintf("ID: %s", w.Id()))
	lines = append(lines, fmt.Sprintf("Name: %s", w.Name()))
	lines = append(lines, fmt.Sprintf("Type: %s", w.TypeName()))
	lines = append(lines, fmt.Sprintf("Position: (%d, %d)", w.X(), w.Y()))
	lines = append(lines, fmt.Sprintf("Size: %dx%d", w.Width(), w.Height()))
	lines = append(lines, fmt.Sprintf("Inner Size: %dx%d", w.InnerWidth(), w.InnerHeight()))
	lines = append(lines, fmt.Sprintf("Grid Position: (%d, %d)", w.GridX(), w.GridY()))
	lines = append(lines, fmt.Sprintf("Expandable: %t %t", w.XExpandable(), w.YExpandable()))
	lines = append(lines, fmt.Sprintf("Min Size: %dx%d", w.MinWidth(), w.MinHeight()))
	lines = append(lines, fmt.Sprintf("Max Size: %dx%d", w.MaxWidth(), w.MaxHeight()))

	for _, line := range lines {
		cnv.FillRect(posX, posY, 200, 20, color.RGBA{R: 0, G: 0, B: 0, A: 150})
		cnv.DrawText(posX, posY, line, "roboto", 16, col, false)
		posY += 20
	}

	/*for y := 0; y < c.height; y += 10 {
		cnv.DrawLine(0, y, c.width, y, 1, color.RGBA{R: 0, G: 100, B: 0, A: 150})
		if y%50 == 0 {
			cnv.DrawText(0, y, fmt.Sprintf("%d", y), "roboto", 12, color.RGBA{R: 0, G: 200, B: 0, A: 255}, false)
		}
	}*/

}

func (c *Form) processResize(width, height int) {
	c.topWidget.SetSize(width, height)
	c.width = width
	c.height = height
	c.forceUpdate()
}

func (c *Form) processMouseDown(button nuimouse.MouseButton, x int, y int) {
	if button == nuimouse.MouseButtonLeft {
		c.mouseLeftButtonPressed = true
	}
	widgetAtCoords := c.topWidget.findWidgetAt(x, y)
	if c.mouseLeftButtonPressed {
		c.mouseLeftButtonPressedWidget = widgetAtCoords
	}
	if widgetAtCoords != nil {
		widgetAtCoords.Focus()
	}
	c.topWidget.processMouseDown(button, x, y, c.lastKeyboardModifiers)
	c.Update()
}

func (c *Form) processMouseUp(button nuimouse.MouseButton, x int, y int) {
	if button == nuimouse.MouseButtonLeft {
		c.mouseLeftButtonPressed = false
		c.mouseLeftButtonPressedWidget = nil
	}
	c.topWidget.processMouseUp(button, x, y, c.lastKeyboardModifiers)
	c.Update()
}

func (c *Form) processMouseMove(x int, y int) {

	// TODO:
	/*if c.mouseLeftButtonPressedWidget != nil {
		widgetX, widgetY := GetWidgeter(c.topWidget).absolutePositionOfWidget(0, 0, GetWidgeter(c.mouseLeftButtonPressedWidget))
		c.mouseLeftButtonPressedWidget.processMouseMove(x-widgetX, y-widgetY, c.lastKeyboardModifiers)
		c.Update()
		return
	}*/

	c.topWidget.processMouseMove(x, y, c.lastKeyboardModifiers)
	c.lastMouseX = x
	c.lastMouseY = y
	hoverWidget := c.topWidget.findWidgetAt(x, y)
	if hoverWidget == nil {
		hoverWidget = c.topWidget
	}

	if hoverWidget != c.hoverWidget {
		if c.hoverWidget != nil {
			c.hoverWidget.processMouseLeave()
		}
		c.hoverWidget = hoverWidget
		if c.hoverWidget != nil {
			c.hoverWidget.processMouseEnter()
		}
	}

	newCursor := nuimouse.MouseCursorArrow
	if c.hoverWidget != nil {
		if c.hoverWidget.MouseCursor() != nuimouse.MouseCursorNotDefined {
			newCursor = c.hoverWidget.MouseCursor()
		}
	}

	if c.lastMouseCursor != newCursor {
		c.wnd.SetMouseCursor(newCursor)
		c.lastMouseCursor = newCursor
	}

	c.Update()
}

func (c *Form) processMouseLeave() {
	c.topWidget.processMouseLeave()

	if c.hoverWidget != nil {
		c.hoverWidget.processMouseLeave()
		c.hoverWidget = nil
	}

	c.Update()
}

func (c *Form) processMouseEnter() {
	c.topWidget.processMouseEnter()
}

func (c *Form) processKeyDown(keyCode nuikey.Key, mods nuikey.KeyModifiers) {
	if c.lastKeyboardModifiers != mods {
		c.lastKeyboardModifiers = mods
	}

	if c.focusedWidget != nil {
		c.focusedWidget.processKeyDown(keyCode, mods)
		c.Update()
		return
	}
	c.topWidget.processKeyDown(keyCode, mods)
	c.Update()
}

func (c *Form) processKeyUp(keyCode nuikey.Key, mods nuikey.KeyModifiers) {
	if c.lastKeyboardModifiers != mods {
		c.lastKeyboardModifiers = mods
	}
	if c.focusedWidget != nil {
		c.focusedWidget.processKeyUp(keyCode, mods)
		c.Update()
		return
	}
	c.topWidget.processKeyUp(keyCode, mods)
	c.Update()
}

func (c *Form) processMouseDblClick(button nuimouse.MouseButton, x int, y int) {
	if c.focusedWidget != nil {
		c.focusedWidget.processMouseDblClick(button, x, y, c.lastKeyboardModifiers)
		c.Update()
		return
	}
	c.topWidget.processMouseDblClick(button, x, y, c.lastKeyboardModifiers)
	c.Update()
}

func (c *Form) processChar(char rune) {
	if c.focusedWidget != nil {
		c.focusedWidget.processChar(char, c.lastKeyboardModifiers)
		c.Update()
		return
	}
	c.topWidget.processChar(char, c.lastKeyboardModifiers)
}

func (c *Form) processMouseWheel(deltaX int, deltaY int) {
	if c.lastKeyboardModifiers.Shift {
		deltaX, deltaY = deltaY, deltaX // Swap for horizontal scrolling
	}
	c.topWidget.processMouseWheel(deltaX, deltaY)
	c.Update()
}

func (c *Form) processTimer() {

	if time.Since(c.lastFreeMemoryTime) > 5*time.Second {
		c.freeMemory()
		c.lastFreeMemoryTime = time.Now()
	}

	c.topWidget.processTimer()

	if c.needUpdate {
		c.realUpdate()
	}
}

func (c *Form) processWindowMove(x, y int) {
	c.forceUpdate()
}

func (c *Form) freeMemory() {
	runtime.GC()
	debug.FreeOSMemory()
}
