package ui

import (
	"image"
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

	lastFreeMemoryTime time.Time
}

var MainForm *Form

type Widgeter interface {
	X() int
	Y() int
	Width() int
	Height() int

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

	AddWidget(widget any)
	AddWidgetOnGrid(widget any, gridX, gridY int)
	RemoveWidget(widget any)

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
	return &c
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
	c.wnd = nui.CreateWindow(c.title, c.width, c.height, true)

	MainForm = c

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

	c.wnd.Show()
	c.wnd.EventLoop()
}

func (c *Form) Update() {
	if c.wnd != nil {
		c.wnd.Update()
	}
}

func (c *Form) processPaint(rgba *image.RGBA) {
	cnv := NewCanvas(rgba)
	cnv.SetClip(0, 0, c.width, c.height)
	GetWidgeter(c.topWidget).processPaint(cnv)
}

func (c *Form) processResize(width, height int) {
	GetWidgeter(c.topWidget).SetSize(width, height)
	c.width = width
	c.height = height
	UpdateMainForm()
}

func (c *Form) processMouseDown(button nuimouse.MouseButton, x int, y int) {
	if button == nuimouse.MouseButtonLeft {
		c.mouseLeftButtonPressed = true
	}
	widgetAtCoords := GetWidgeter(c.topWidget).findWidgetAt(x, y)
	if c.mouseLeftButtonPressed {
		c.mouseLeftButtonPressedWidget = widgetAtCoords
	}
	if widgetAtCoords != nil {
		widgetAtCoords.Focus()
	}
	GetWidgeter(c.topWidget).processMouseDown(button, x, y, c.lastKeyboardModifiers)
	c.Update()
}

func (c *Form) processMouseUp(button nuimouse.MouseButton, x int, y int) {
	if button == nuimouse.MouseButtonLeft {
		c.mouseLeftButtonPressed = false
		c.mouseLeftButtonPressedWidget = nil
	}
	GetWidgeter(c.topWidget).processMouseUp(button, x, y, c.lastKeyboardModifiers)
	c.Update()
}

func (c *Form) processMouseMove(x int, y int) {
	/*if c.mouseLeftButtonPressedWidget != nil {
		c.mouseLeftButtonPressedWidget.processMouseMove(x, y, c.lastKeyboardModifiers)
		c.Update()
		return
	}*/

	GetWidgeter(c.topWidget).processMouseMove(x, y, c.lastKeyboardModifiers)
	c.lastMouseX = x
	c.lastMouseY = y
	hoverWidget := GetWidgeter(c.topWidget).findWidgetAt(x, y)
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
	GetWidgeter(c.topWidget).processMouseLeave()

	if c.hoverWidget != nil {
		c.hoverWidget.processMouseLeave()
		c.hoverWidget = nil
	}

	c.Update()
}

func (c *Form) processMouseEnter() {
	GetWidgeter(c.topWidget).processMouseEnter()
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
	GetWidgeter(c.topWidget).processKeyDown(keyCode, mods)
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
	GetWidgeter(c.topWidget).processKeyUp(keyCode, mods)
	c.Update()
}

func (c *Form) processMouseDblClick(button nuimouse.MouseButton, x int, y int) {
	if c.focusedWidget != nil {
		c.focusedWidget.processMouseDblClick(button, x, y, c.lastKeyboardModifiers)
		c.Update()
		return
	}
	GetWidgeter(c.topWidget).processMouseDblClick(button, x, y, c.lastKeyboardModifiers)
	c.Update()
}

func (c *Form) processChar(char rune) {
	if c.focusedWidget != nil {
		c.focusedWidget.processChar(char, c.lastKeyboardModifiers)
		c.Update()
		return
	}
	GetWidgeter(c.topWidget).processChar(char, c.lastKeyboardModifiers)
}

func (c *Form) processMouseWheel(deltaX int, deltaY int) {
	if c.lastKeyboardModifiers.Shift {
		deltaX, deltaY = deltaY, deltaX // Swap for horizontal scrolling
	}
	GetWidgeter(c.topWidget).processMouseWheel(deltaX, deltaY)
	c.Update()
}

func (c *Form) processTimer() {

	if time.Since(c.lastFreeMemoryTime) > 5*time.Second {
		c.freeMemory()
		c.lastFreeMemoryTime = time.Now()
	}

	GetWidgeter(c.topWidget).processTimer()
}

func (c *Form) freeMemory() {
	runtime.GC()
	debug.FreeOSMemory()
}
