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
	mouseLeftButtonPressedWidget *Widget

	topWidget     *Widget
	hoverWidget   *Widget
	focusedWidget *Widget

	timers             []*timer
	lastFreeMemoryTime time.Time
}

var MainForm *Form

func NewForm() *Form {
	var c Form
	c.timers = make([]*timer, 0)
	c.title = "Form"
	c.width = 800
	c.height = 600
	c.topWidget = NewWidget()
	c.topWidget.SetName("FormTopWidget")
	c.topWidget.SetPosition(0, 0)
	c.topWidget.SetSize(c.width, c.height)
	c.topWidget.SetAnchors(true, true, true, true)
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

func (c *Form) Panel() *Widget {
	return c.topWidget
}

func (c *Form) AddTimer(intervalMs int, callback func()) {
	t := &timer{
		intervalMs: intervalMs,
		callback:   callback,
	}
	c.timers = append(c.timers, t)
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
	c.topWidget.processPaint(cnv)
}

func (c *Form) processResize(width, height int) {
	c.topWidget.SetSize(width, height)
	c.width = width
	c.height = height
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
	c.topWidget.processMouseDown(button, x, y)
	c.Update()
}

func (c *Form) processMouseUp(button nuimouse.MouseButton, x int, y int) {
	if button == nuimouse.MouseButtonLeft {
		c.mouseLeftButtonPressed = false
		c.mouseLeftButtonPressedWidget = nil
	}
	c.topWidget.processMouseUp(button, x, y)
}

func (c *Form) processMouseMove(x int, y int) {
	if c.mouseLeftButtonPressedWidget != nil {
		c.mouseLeftButtonPressedWidget.processMouseMove(x, y)
		c.Update()
		return
	}

	c.topWidget.processMouseMove(x, y)
	c.lastMouseX = x
	c.lastMouseY = y
	hoverWidget := c.topWidget.findWidgetAt(x, y)
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
		if c.hoverWidget.mouseCursor != nuimouse.MouseCursorNotDefined {
			newCursor = c.hoverWidget.mouseCursor
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
		c.focusedWidget.processKeyDown(keyCode)
		return
	}
	c.topWidget.processKeyDown(keyCode)
}

func (c *Form) processKeyUp(keyCode nuikey.Key, mods nuikey.KeyModifiers) {
	if c.lastKeyboardModifiers != mods {
		c.lastKeyboardModifiers = mods
	}
	if c.focusedWidget != nil {
		c.focusedWidget.processKeyUp(keyCode)
		return
	}
	c.topWidget.processKeyUp(keyCode)
}

func (c *Form) processMouseDblClick(button nuimouse.MouseButton, x int, y int) {
	if c.focusedWidget != nil {
		c.focusedWidget.processMouseDblClick(button, x, y)
		return
	}
	c.topWidget.processMouseDblClick(button, x, y)
}

func (c *Form) processChar(char rune) {
	if c.focusedWidget != nil {
		c.focusedWidget.processChar(char)
		return
	}
	c.topWidget.processChar(char)
}

func (c *Form) processMouseWheel(deltaX int, deltaY int) {
	if c.lastKeyboardModifiers.Shift {
		deltaX, deltaY = deltaY, deltaX // Swap for horizontal scrolling
	}
	c.topWidget.processMouseWheel(deltaX, deltaY)
	c.Update()
}

func (c *Form) processTimer() {
	for _, t := range c.timers {
		t.tick()
	}

	if time.Since(c.lastFreeMemoryTime) > 5*time.Second {
		c.freeMemory()
		c.lastFreeMemoryTime = time.Now()
	}
}

func (c *Form) freeMemory() {
	runtime.GC()
	debug.FreeOSMemory()
}
