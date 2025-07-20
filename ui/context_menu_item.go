package ui

import (
	"image/color"
	"time"

	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
)

type ContextMenuItem struct {
	Widget
	text    string
	OnClick func()
	//Image          image.Image
	//ImageResource  []byte
	//KeyCombination string
	parentMenu           *ContextMenu
	needToClosePopupMenu func()

	//timerShowInnerMenu  *FormTimer
	//AdjustColorForImage bool

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

func (c *ContextMenuItem) OnInit() {
	//c.timerShowInnerMenu = c.ownWindow.NewTimer(200, c.timerShowInnerMenuHandler)
	//c.AdjustColorForImage = true
}

func (c *ContextMenuItem) SetText(text string) {
	c.text = text
	UpdateMainForm()
}

func (c *ContextMenuItem) ControlType() string {
	return "PopupMenuItem"
}

func (c *ContextMenuItem) Draw(ctx *Canvas) {
	backColor := GetThemeColor("background", DefaultBackground)
	if c.IsHovered() {
		backColor = color.RGBA{R: 20, G: 20, B: 50, A: 255}
	}
	ctx.FillRect(0, 0, c.InnerWidth(), c.InnerHeight(), backColor)
	ctx.DrawTextMultiline(0, 0, c.Width(), c.Height(), HAlignLeft, VAlignCenter, c.text, GetThemeColor("foreground", DefaultForeground), "robotomono", 16, false)

	//xOffset := 0
	/*if c.Image != nil || c.ImageResource != nil {
		imageSource := c.Image
		if c.ImageResource != nil {
			imageSource = uiresources.ResImgCol(c.ImageResource, DefaultForeground)
		}

		img := resize.Resize(24, 24, imageSource, resize.Bicubic)
		if c.AdjustColorForImage {
			img = canvas.AdjustImageForColor(img, img.Bounds().Max.X, img.Bounds().Max.Y, c.foregroundColor.Color())
		}

		height := img.Bounds().Max.Y
		yOffset := (c.Height() - height) / 2

		ctx.DrawImage(3, yOffset, img.Bounds().Max.X, height, img)
		xOffset += 32
	}*/

	/*ctx.SetColor(c.foregroundColor.Color())
	ctx.SetFontSize(c.fontSize.Float64())
	textWidth := c.InnerWidth()
	if c.innerMenu != nil {
		textWidth -= c.InnerHeight()
	}
	ctx.DrawText(xOffset+5, 0, textWidth, c.InnerHeight(), c.text)*/

	/*if c.innerMenu != nil {
		imgArrow := uiresources.ResImgCol(uiresources.R_icons_material4_png_av_play_arrow_materialicons_48dp_1x_baseline_play_arrow_black_48dp_png, c.ForeColor())
		ctx.DrawImage(c.InnerWidth()-c.InnerHeight(), 0, imgArrow.Bounds().Max.X, imgArrow.Bounds().Max.Y, resize.Resize(uint(c.InnerHeight()), uint(c.InnerHeight()), imgArrow, resize.Bicubic))
	}*/
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
