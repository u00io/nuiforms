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
	c.form.Update()
}

func (c *ContextMenuItem) ControlType() string {
	return "PopupMenuItem"
}

// contextMenuItemPadding keeps the item's text and submenu arrow off its
// edges - the item itself still spans the full menu width so the hover
// highlight reaches border to border, as in standard menu styling.
const contextMenuItemPadding = 10

func (c *ContextMenuItem) Draw(ctx *Canvas) {
	backColor := c.BackgroundColor()
	if c.IsHovered() {
		backColor = c.BackgroundColorWithAddElevation(2)
	}
	ctx.FillRect(0, 0, c.InnerWidth(), c.InnerHeight(), backColor)

	textAreaWidth := c.Width() - contextMenuItemPadding*2
	if c.innerMenu != nil {
		textAreaWidth -= c.Height() + contextMenuItemPadding
	}
	displayText := truncateTextToWidth(c.FontFamily(), c.FontSize(), c.text, textAreaWidth)

	ctx.SetHAlign(HAlignLeft)
	ctx.SetVAlign(VAlignCenter)
	ctx.SetColor(c.ForegroundColor())
	ctx.SetFontFamily(c.FontFamily())
	ctx.SetFontSize(c.FontSize())
	ctx.DrawText(contextMenuItemPadding, 0, c.Width()-contextMenuItemPadding*2, c.Height(), displayText)

	if c.innerMenu != nil {
		rectSize := c.Height()
		x := c.Width() - rectSize - contextMenuItemPadding
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
	c.form.Panel().CloseAfterPopupWidget(c.parentMenu)
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
	c.form.Update()
	return true
}

func (c *ContextMenuItem) SetInnerMenu(menu *ContextMenu) {
	c.innerMenu = menu
}

// attachToForm also attaches innerMenu, which AddItemWithSubmenu stores
// directly on the item rather than adding it as a child widget, so the
// generic Widget.attachToForm cascade would otherwise never reach it -
// leaving its form nil and crashing (nil c.form.Panel()) the first time
// the submenu is opened.
func (c *ContextMenuItem) attachToForm(form *Form) {
	c.Widget.attachToForm(form)
	if c.innerMenu != nil {
		c.innerMenu.attachToForm(form)
	}
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
