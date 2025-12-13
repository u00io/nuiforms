package ui

import (
	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
)

type RadioButton struct {
	Widget
	checked        bool
	text           string
	onStateChanged func(btn *RadioButton, checked bool)
}

func NewRadioButton(text string) *RadioButton {
	var c RadioButton
	c.InitWidget()
	c.SetTypeName("RadioButton")
	c.SetMinSize(150, 30)
	c.SetMaxSize(10000, 30)
	c.SetMouseCursor(nuimouse.MouseCursorPointer)
	c.SetText("RadioButton")
	c.SetCanBeFocused(true)

	c.SetOnPaint(c.draw)
	c.SetOnMouseDown(c.buttonProcessMouseDown)
	c.SetOnMouseUp(c.buttonProcessMouseUp)

	c.SetText(text)

	return &c
}

func (c *RadioButton) Text() string {
	return c.text
}

func (c *RadioButton) SetText(text string) {
	c.text = text
	UpdateMainForm()
}

func (c *RadioButton) SetOnStateChanged(fn func(btn *RadioButton, checked bool)) {
	c.onStateChanged = fn
}

func (c *RadioButton) SetChecked(checked bool) {
	if c.checked == checked {
		return
	}

	if checked {
		// Uncheck all other radio buttons in the same group
		parentWidget := c.ParentWidget()
		if parentWidget != nil {
			for _, child := range parentWidget.Widgets() {
				if radioButton, ok := child.(*RadioButton); ok && radioButton.Id() != c.Id() {
					radioButton.SetChecked(false)
				}
			}
		}
	}

	c.checked = checked
	if c.onStateChanged != nil {
		c.onStateChanged(c, checked)
	}
}

func (c *RadioButton) Checked() bool {
	return c.checked
}

func (c *RadioButton) draw(cnv *Canvas) {
	backColor := c.BackgroundColor()
	cnv.FillRect(0, 0, c.Width(), c.Height(), backColor)

	boxAndTextSpace := 0
	cnv.SetHAlign(HAlignLeft)
	cnv.SetVAlign(VAlignCenter)
	cnv.SetColor(c.ForegroundColor())
	cnv.SetFontFamily(c.FontFamily())
	cnv.SetFontSize(c.FontSize())
	cnv.DrawText(30+boxAndTextSpace, 0, c.Width()-30-boxAndTextSpace, c.Height(), c.text)

	padding := 5

	cnv.SetColor(c.BackgroundColorWithAddElevation(-1))
	cnv.DrawRect(padding, padding, 30-padding*2, 30-padding*2)
	if c.checked {
		cnv.FillRect(padding*2, padding*2, 30-padding*4, 30-padding*4, c.ForegroundColor())
	}
}

func (c *RadioButton) buttonProcessMouseDown(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
	return true
}

func (c *RadioButton) buttonProcessMouseUp(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
	if x < 0 || x >= c.Width() || y < 0 || y >= c.Height() {
		// MouseUp outside the button area, ignore
		return false
	}

	hoverWidgeter := MainForm.hoverWidget
	var localWidgeter Widgeter = c
	if hoverWidgeter == localWidgeter {
		if !c.checked {
			c.SetChecked(true)
		}
	}

	return true
}
