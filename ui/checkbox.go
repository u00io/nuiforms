package ui

import (
	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
)

type Checkbox struct {
	Widget
	//checked        bool
	//text           string
	//onStateChanged func(btn *Checkbox, checked bool)
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
	return c.GetPropString("text", "")
}

func (c *Checkbox) SetText(text string) {
	c.SetProp("text", text)
	c.form.Update()
}

func (c *Checkbox) SetOnStateChanged(fn func()) {
	c.SetPropFunction("onstatechanged", fn)
}

type EventCheckboxStateChanged struct {
	Checkbox *Checkbox
	Checked  bool
}

func (c *Checkbox) SetChecked(checked bool) {
	if c.GetPropBool("checked", false) == checked {
		return
	}

	c.SetProp("checked", checked)
	f := c.GetPropFunction("onstatechanged")
	if f != nil {
		var ev EventCheckboxStateChanged
		ev.Checkbox = c
		ev.Checked = checked
		PushEvent(&ev)
		f()
		PopEvent()
	}
}

func (c *Checkbox) Checked() bool {
	return c.GetPropBool("checked", false)
}

func (c *Checkbox) draw(cnv *Canvas) {
	//cnv.FillRect(0, 0, c.Width(), c.Height())

	boxAndTextSpace := 0
	cnv.SetHAlign(HAlignLeft)
	cnv.SetVAlign(VAlignCenter)
	cnv.SetColor(c.ForegroundColor())
	cnv.SetFontFamily(c.FontFamily())
	cnv.SetFontSize(c.FontSize())
	cnv.DrawText(30+boxAndTextSpace, 0, c.Width()-30-boxAndTextSpace, c.Height(), c.Text())

	padding := 5
	boxSize := 30 - padding*2

	cnv.SetColor(c.BackgroundColorWithAddElevation(-2))
	cnv.FillRoundedRect(padding, padding, boxSize, boxSize, 3)

	if c.Checked() {
		tickColor := c.ForegroundColor()
		tickWidth := 2
		// A short stroke down to the tick's low point, then a longer stroke
		// up to its top-right end - the usual checkmark shape, sized as
		// fractions of boxSize so it stays right if padding/boxSize change.
		x1, y1 := padding+boxSize*3/20, padding+boxSize*11/20
		x2, y2 := padding+boxSize*8/20, padding+boxSize*16/20
		x3, y3 := padding+boxSize*17/20, padding+boxSize*4/20
		cnv.DrawLine(x1, y1, x2, y2, tickWidth, tickColor)
		cnv.DrawLine(x2, y2, x3, y3, tickWidth, tickColor)
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

	// Compare via bounds rather than "c.form.hoverWidget == c": for a
	// Checkbox nested inside a custom composite widget, the interface value
	// stored as hoverWidget can carry the promoted *Widget type instead of
	// *Checkbox, so a direct comparison against c never matches even though
	// the mouse is genuinely over this checkbox. The bounds check above
	// already confirms that, so toggle unconditionally here.
	c.SetChecked(!c.Checked())

	return true
}
