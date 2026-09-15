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

/*
Properties:
- halign: string - Horizontal alignment of the image within the button ["left", "center", "right"].
- valign: string - Vertical alignment of the image within the button ["top", "center", "bottom"].
*/

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
	c.form.Update()
}

func (c *ButtonImage) SetOnButtonClick(fn func(btn *ButtonImage)) {
	c.onButtonClick = fn
}

// ImageHAlign returns the horizontal alignment of the image within the button.
func (c *ButtonImage) ImageHAlign() HAlign {
	return c.GetHAlign("halign", HAlignCenter)
}

// SetImageHAlign sets the horizontal alignment of the image within the button.
func (c *ButtonImage) SetImageHAlign(align HAlign) {
	c.SetProp("halign", align.String())
}

// ImageVAlign returns the vertical alignment of the image within the button.
func (c *ButtonImage) ImageVAlign() VAlign {
	return c.GetVAlign("valign", VAlignCenter)
}

// SetImageVAlign sets the vertical alignment of the image within the button.
func (c *ButtonImage) SetImageVAlign(align VAlign) {
	c.SetProp("valign", align.String())
}

// imagePosition computes the top-left drawing position for the image given
// the current alignment props and the image's own size vs. the button's.
func (c *ButtonImage) imagePosition() (int, int) {
	b := c.img.Bounds()
	imgWidth := b.Dx()
	imgHeight := b.Dy()

	x := 0
	switch c.ImageHAlign() {
	case HAlignCenter:
		x = (c.Width() - imgWidth) / 2
	case HAlignRight:
		x = c.Width() - imgWidth
	}

	y := 0
	switch c.ImageVAlign() {
	case VAlignCenter:
		y = (c.Height() - imgHeight) / 2
	case VAlignBottom:
		y = c.Height() - imgHeight
	}

	return x, y
}

// ProcessPropChange repaints when an alignment prop (or any other prop) changes.
func (c *ButtonImage) ProcessPropChange(key string, value interface{}) {
	c.form.Update()
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

	if c.img != nil {
		x, y := c.imagePosition()
		cnv.DrawImage(x, y, c.img)
	}

	cnv.SetColor(c.BackgroundColorWithAddElevation(2))
	cnv.DrawRect(0, 0, c.Width(), c.Height())
}

func (c *ButtonImage) buttonProcessMouseDown(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
	c.pressed = true
	return true
}

func (c *ButtonImage) buttonProcessMouseUp(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
	wasPressed := c.pressed
	c.pressed = false

	if x < 0 || x >= c.Width() || y < 0 || y >= c.Height() {
		// Released outside the button - cancel the click.
		return false
	}

	// Compare via bounds + wasPressed rather than "c.form.hoverWidget == c":
	// for a ButtonImage nested inside a custom composite widget, the
	// interface value stored as hoverWidget can carry the promoted *Widget
	// type instead of *ButtonImage, so a direct comparison against c never
	// matches even though the mouse is genuinely over this button.
	if wasPressed && c.onButtonClick != nil {
		c.onButtonClick(c)
	}

	return true
}
