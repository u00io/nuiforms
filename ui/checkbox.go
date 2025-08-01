package ui

import (
	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
)

type Checkbox struct {
	Widget
	checked        bool
	text           string
	onStateChanged func(btn *Checkbox, checked bool)
}

func NewCheckbox(text string) *Checkbox {
	var c Checkbox
	c.InitWidget()
	c.SetTypeName("Checkbox")
	c.SetMinSize(150, 30)
	c.SetMaxSize(10000, 30)
	c.SetMouseCursor(nuimouse.MouseCursorPointer)
	c.SetText("Checkbox")
	c.SetCanBeFocused(true)

	c.SetOnPaint(c.draw)
	c.SetOnMouseDown(c.buttonProcessMouseDown)
	c.SetOnMouseUp(c.buttonProcessMouseUp)

	c.SetText(text)

	return &c
}

func (c *Checkbox) Text() string {
	return c.text
}

func (c *Checkbox) SetText(text string) {
	c.text = text
	UpdateMainForm()
}

func (c *Checkbox) SetOnStateChanged(fn func(btn *Checkbox, checked bool)) {
	c.onStateChanged = fn
}

func (c *Checkbox) SetChecked(checked bool) {
	if c.checked == checked {
		return
	}

	c.checked = checked
	if c.onStateChanged != nil {
		c.onStateChanged(c, checked)
	}
}

func (c *Checkbox) Checked() bool {
	return c.checked
}

func (c *Checkbox) draw(cnv *Canvas) {
	backColor := c.BackgroundColor()
	cnv.FillRect(0, 0, c.Width(), c.Height(), backColor)

	boxAndTextSpace := 0
	cnv.SetHAlign(HAlignLeft)
	cnv.SetVAlign(VAlignCenter)
	cnv.SetColor(c.Color())
	cnv.SetFontFamily(c.FontFamily())
	cnv.SetFontSize(c.FontSize())
	cnv.DrawText(30+boxAndTextSpace, 0, c.Width()-30-boxAndTextSpace, c.Height(), c.text)

	padding := 5

	cnv.SetColor(c.BackgroundColorAccent1())
	cnv.DrawRect(padding, padding, 30-padding*2, 30-padding*2)
	if c.checked {
		cnv.FillRect(padding*2, padding*2, 30-padding*4, 30-padding*4, c.Color())
	}
}

func (c *Checkbox) buttonProcessMouseDown(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
	return true
}

func (c *Checkbox) buttonProcessMouseUp(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
	if x < 0 || x >= c.Width() || y < 0 || y >= c.Height() {
		// MouseUp outside the button area, ignore
		return false
	}

	hoverWidgeter := MainForm.hoverWidget
	var localWidgeter Widgeter = c
	if hoverWidgeter == localWidgeter {
		c.SetChecked(!c.checked)
	}

	return true
}
