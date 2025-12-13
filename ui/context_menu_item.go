package ui

import (
	"time"

	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
)

type ContextMenuItem struct {
	Widget
	text                 string
	OnClick              func()
	parentMenu           *ContextMenu
	needToClosePopupMenu func()

	timerEnabled bool

	timerLastElapsedDTMSec int64

	innerMenu *ContextMenu
}

func NewContextMenuItem() *ContextMenuItem {
	var item ContextMenuItem
	item.InitWidget()
	item.SetAbsolutePositioning(true)
	item.SetMouseCursor(nuimouse.MouseCursorPointer)

	item.SetOnPaint(item.Draw)

	item.SetOnMouseDown(item.mouseDownHandler)
	item.SetOnMouseMove(item.MouseMove)
	item.SetOnMouseEnter(item.MouseEnter)
	item.SetOnMouseLeave(item.MouseLeave)

	item.AddTimer(200, item.timerShowInnerMenuHandler)
	return &item
}

func (c *ContextMenuItem) SetText(text string) {
	c.text = text
	UpdateMainForm()
}

func (c *ContextMenuItem) ControlType() string {
	return "PopupMenuItem"
}

func (c *ContextMenuItem) Draw(ctx *Canvas) {
	backColor := c.BackgroundColorAccent1()
	if c.IsHovered() {
		backColor = c.BackgroundColorAccent2()
	}
	ctx.FillRect(0, 0, c.InnerWidth(), c.InnerHeight(), backColor)
	ctx.SetHAlign(HAlignLeft)
	ctx.SetVAlign(VAlignCenter)
	ctx.SetColor(c.ForegroundColor())
	ctx.SetFontFamily(c.FontFamily())
	ctx.SetFontSize(c.FontSize())
	ctx.DrawText(0, 0, c.Width(), c.Height(), c.text)

	if c.innerMenu != nil {
		rectSize := c.Height()
		x := c.Width() - c.Height()
		y := 0
		ctx.SetHAlign(HAlignLeft)
		ctx.SetVAlign(VAlignCenter)
		ctx.SetColor(c.ForegroundColor())
		ctx.SetFontFamily(c.FontFamily())
		ctx.SetFontSize(c.FontSize())
		ctx.DrawText(x, y, rectSize, rectSize, "\u00BB")
	}
}

func (c *ContextMenuItem) mouseDownHandler(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
	c.timerEnabled = false

	if c.innerMenu != nil {
		x, y := c.RectClientAreaOnWindow()
		w := c.Width()
		c.innerMenu.showMenu(x+w, y, c.parentMenu)
		return true
	}

	if c.needToClosePopupMenu != nil {
		c.needToClosePopupMenu()
	}

	if c.OnClick != nil {
		c.OnClick()
	}
	return true
}

func (c *ContextMenuItem) MouseEnter() {
	MainForm.Panel().CloseAfterPopupWidget(c.parentMenu)
	if c.innerMenu != nil {
		c.timerEnabled = true
		c.timerLastElapsedDTMSec = time.Now().UnixNano() / 1000000
		return
	}
}

func (c *ContextMenuItem) MouseLeave() {
	c.timerEnabled = false
}

func (c *ContextMenuItem) MouseMove(x int, y int, mods nuikey.KeyModifiers) bool {
	UpdateMainForm()
	return true
}

func (c *ContextMenuItem) SetInnerMenu(menu *ContextMenu) {
	c.innerMenu = menu
}

func (c *ContextMenuItem) timerShowInnerMenuHandler() {
	if c.timerEnabled && time.Now().UnixNano()/1000000-c.timerLastElapsedDTMSec > 200 {
		c.timerEnabled = false
		x, y := c.parentMenu.RectClientAreaOnWindow()
		y += c.Y()
		w := c.Width()
		c.innerMenu.showMenu(x+w, y, c.parentMenu)
	}
}
