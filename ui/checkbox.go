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
	UpdateMainForm()
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
		PushEvent(&Event{Parameter: ev})
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

	cnv.SetColor(c.BackgroundColorWithAddElevation(-2))
	cnv.FillRoundedRect(padding, padding, 30-padding*2, 30-padding*2, 3)

	//cnv.SetColor(c.BackgroundColorWithAddElevation(5))
	//cnv.DrawRect(padding, padding, 30-padding*2, 30-padding*2)
	if c.Checked() {
		cnv.SetColor(c.ForegroundColor())
		internalPadding := 9
		cnv.FillRoundedRect(internalPadding, internalPadding, 30-internalPadding*2, 30-internalPadding*2, 3)
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
		c.SetChecked(!c.Checked())
	}

	return true
}
