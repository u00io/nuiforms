package ui

import (
	"image"

	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
)

type ButtonImage struct {
	Widget

	pressed       bool
	img           image.Image
	onButtonClick func(btn *ButtonImage)
}

func NewButtonImage(img image.Image) *ButtonImage {
	var c ButtonImage
	c.InitWidget()
	c.SetTypeName("Button")
	c.SetMinSize(100, 30)
	c.SetMaxSize(10000, 30)
	c.SetMouseCursor(nuimouse.MouseCursorPointer)
	c.SetImage(img)
	c.SetCanBeFocused(true)

	c.SetOnPaint(c.draw)
	c.SetOnMouseDown(c.buttonProcessMouseDown)
	c.SetOnMouseUp(c.buttonProcessMouseUp)

	return &c
}

func (c *ButtonImage) Image() image.Image {
	return c.img
}

func (c *ButtonImage) SetImage(img image.Image) {
	c.img = img
	UpdateMainForm()
}

func (c *ButtonImage) SetOnButtonClick(fn func(btn *ButtonImage)) {
	c.onButtonClick = fn
}

func (c *ButtonImage) Press() {
	if c.onButtonClick != nil {
		c.onButtonClick(c)
	}
}

func (c *ButtonImage) draw(cnv *Canvas) {
	backColor := c.BackgroundColor()
	if c.IsHovered() {
		backColor = c.BackgroundColorWithAddElevation(-1)
	}
	if c.pressed {
		backColor = c.BackgroundColorWithAddElevation(2)
	}
	cnv.FillRect(0, 0, c.Width(), c.Height(), backColor)

	cnv.SetHAlign(HAlignCenter)
	cnv.SetVAlign(VAlignCenter)
	cnv.SetColor(c.ForegroundColor())
	cnv.SetFontFamily(c.FontFamily())
	cnv.SetFontSize(c.FontSize())
	cnv.DrawImage(0, 0, c.img)

	cnv.SetColor(c.BackgroundColorWithAddElevation(2))
	cnv.DrawRect(0, 0, c.Width(), c.Height())
}

func (c *ButtonImage) buttonProcessMouseDown(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
	c.pressed = true
	return true
}

func (c *ButtonImage) buttonProcessMouseUp(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
	c.pressed = false

	if x < 0 || x >= c.Width() || y < 0 || y >= c.Height() {
		// MouseUp outside the button area, ignore
		return false
	}

	hoverWidgeter := MainForm.hoverWidget
	var localWidgeter Widgeter = c
	if hoverWidgeter == localWidgeter {
		if c.onButtonClick != nil {
			c.onButtonClick(c)
		}
	}

	return true
}
