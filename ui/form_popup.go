package ui

import (
	"image"

	"github.com/u00io/nui/nui"
	"github.com/u00io/nui/nuimouse"
)

// nativePopupWidget is implemented by popup widgets that are shown in their
// own native window, so they aren't clipped by the form's bounds.
type nativePopupWidget interface {
	nativePopup() bool
}

// popupPlacer is implemented by popup widgets that know a better position
// than being pushed back onto the screen when they don't fit on it.
type popupPlacer interface {
	// popupFlipped returns, in client coordinates, the X to use when the
	// popup doesn't fit to the right and the Y to use when it doesn't fit below
	popupFlipped() (x, y int)
}

// popupHost is the native window of an open popup widget. The widget keeps
// its position in the form's client coordinates, which may now be outside
// the client area. Mouse events of the window are translated to those
// coordinates and processed by the form like its own, so the popup widgets
// work the same way whether they have a window or not.
type popupHost struct {
	widget        Widgeter
	wnd           nui.PopupWindow
	width, height int
}

func wantsPopupWindow(w Widgeter) bool {
	n, ok := w.(nativePopupWidget)
	return ok && n.nativePopup()
}

// syncPopupWindows shows windows for newly opened popup widgets and hides
// the windows of closed ones. Called whenever the popup list changes.
func (c *Form) syncPopupWindows() {
	if c.wnd == nil {
		return
	}

	open := c.popupHosts[:0]
	for _, h := range c.popupHosts {
		if c.isPopupWidgetOpen(h.widget) {
			open = append(open, h)
			continue
		}
		h.wnd.Hide()
		c.freePopupWindows = append(c.freePopupWindows, h.wnd)
		if c.popupUnderMouse == h {
			c.popupUnderMouse = nil
		}
	}
	c.popupHosts = open

	for _, w := range c.topWidget.PopupWidgets {
		if !wantsPopupWindow(w) || c.popupHostOf(w) != nil {
			continue
		}
		wnd := c.takePopupWindow()
		if wnd == nil {
			continue // no popups on this platform: drawn inside the form
		}
		h := &popupHost{widget: w, wnd: wnd}
		c.bindPopupWindow(h)
		c.popupHosts = append(c.popupHosts, h)
		c.placePopupWindow(h)
	}
}

func (c *Form) isPopupWidgetOpen(w Widgeter) bool {
	for _, popupWidget := range c.topWidget.PopupWidgets {
		if popupWidget == w {
			return true
		}
	}
	return false
}

func (c *Form) popupHostOf(w Widgeter) *popupHost {
	for _, h := range c.popupHosts {
		if h.widget == w {
			return h
		}
	}
	return nil
}

func (c *Form) isPopupHostOpen(h *popupHost) bool {
	for _, openHost := range c.popupHosts {
		if openHost == h {
			return true
		}
	}
	return false
}

// takePopupWindow reuses a hidden popup window or creates a new one.
func (c *Form) takePopupWindow() nui.PopupWindow {
	if n := len(c.freePopupWindows); n > 0 {
		wnd := c.freePopupWindows[n-1]
		c.freePopupWindows = c.freePopupWindows[:n-1]
		return wnd
	}
	if c.popupWindowsUnavailable {
		return nil
	}
	wnd := nui.CreatePopupWindow(c.wnd, true)
	if wnd == nil {
		c.popupWindowsUnavailable = true
	}
	return wnd
}

func (c *Form) bindPopupWindow(h *popupHost) {
	h.wnd.OnPaint(func(rgba *image.RGBA) {
		cnv := NewCanvas(rgba)
		cnv.SetDirectTranslateAndClip(0, 0, h.widget.Width(), h.widget.Height())
		h.widget.ProcessPaint(cnv)
	})
	h.wnd.OnMouseMove(func(x, y int) {
		if !c.isPopupHostOpen(h) {
			return
		}
		c.popupUnderMouse = h
		c.processMouseMove(x+h.widget.X(), y+h.widget.Y())
	})
	h.wnd.OnMouseButtonDown(func(btn nuimouse.MouseButton, x, y int) {
		if !c.isPopupHostOpen(h) {
			return
		}
		c.popupUnderMouse = h
		c.processMouseDown(btn, x+h.widget.X(), y+h.widget.Y())
		if !c.isPopupHostOpen(h) {
			// The click closed the popup (e.g. chose a menu item): its
			// window is hidden, so the button release may never come
			c.mouseLeftButtonPressed = false
			c.mouseLeftButtonPressedWidget = nil
		}
	})
	h.wnd.OnMouseButtonUp(func(btn nuimouse.MouseButton, x, y int) {
		if !c.isPopupHostOpen(h) {
			return
		}
		c.processMouseUp(btn, x+h.widget.X(), y+h.widget.Y())
	})
	h.wnd.OnMouseWheel(func(deltaX, deltaY int) {
		if !c.isPopupHostOpen(h) {
			return
		}
		c.processMouseWheel(deltaX, deltaY)
	})
	h.wnd.OnMouseLeave(func() {
		if c.popupUnderMouse == h {
			c.popupUnderMouse = nil
			c.processMouseLeave()
		}
	})
}

// placePopupWindow shows the window over the widget's position, moved onto
// the screen if needed. The widget is moved along, so the form's hit
// testing matches what is on the screen.
func (c *Form) placePopupWindow(h *popupHost) {
	w := h.widget
	width, height := w.Width(), w.Height()
	originX, originY := c.wnd.ClientToScreen(0, 0)
	x, y := originX+w.X(), originY+w.Y()
	areaX, areaY, areaW, areaH := c.wnd.ScreenWorkArea(x, y)

	if placer, ok := w.(popupPlacer); ok {
		flippedX, flippedY := placer.popupFlipped()
		if x+width > areaX+areaW {
			x = originX + flippedX
		}
		if y+height > areaY+areaH {
			y = originY + flippedY
		}
	}
	x = max(areaX, min(x, areaX+areaW-width))
	y = max(areaY, min(y, areaY+areaH-height))

	w.SetPosition(x-originX, y-originY)
	h.width, h.height = width, height
	h.wnd.ShowAt(x, y, width, height)
}

// updatePopupWindows repaints the popup windows along with the form and
// follows the size changes of their widgets.
func (c *Form) updatePopupWindows() {
	for _, h := range c.popupHosts {
		if h.widget.Width() != h.width || h.widget.Height() != h.height {
			c.placePopupWindow(h)
		}
		h.wnd.Update()
	}
}

// closePopups closes all popup widgets, e.g. when the form loses activation:
// their windows would otherwise stay above other applications.
func (c *Form) closePopups() {
	if len(c.topWidget.PopupWidgets) > 0 {
		c.topWidget.CloseAllPopup()
	}
}

// destroyPopupWindows runs before the form's window goes away.
func (c *Form) destroyPopupWindows() {
	for _, h := range c.popupHosts {
		h.wnd.Close()
	}
	for _, wnd := range c.freePopupWindows {
		wnd.Close()
	}
	c.popupHosts = nil
	c.freePopupWindows = nil
	c.popupUnderMouse = nil
}

func (c *Form) processDeactivate() {
	c.tooltipHide()
	c.closePopups()
}
