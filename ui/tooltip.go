package ui

import "time"

const (
	tooltipDelay   = 600 * time.Millisecond
	tooltipPadding = 6
)

// tooltipState tracks the tooltip of the widget under the mouse cursor.
// The tooltip appears after the mouse rests on a widget for tooltipDelay,
// and is hidden when the hovered widget changes, the mouse leaves the
// window or a mouse button is pressed (until the hover changes again).
type tooltipState struct {
	hoverStart time.Time
	visible    bool
	suppressed bool
	x, y       int
}

// SetTooltip sets the text shown when the mouse rests on the widget.
// An empty string disables the tooltip.
func (c *Widget) SetTooltip(text string) {
	c.SetProp("tooltip", text)
}

func (c *Widget) Tooltip() string {
	return c.GetPropString("tooltip", "")
}

func (c *Form) tooltipText() string {
	if c.hoverWidget == nil {
		return ""
	}
	return c.hoverWidget.GetPropString("tooltip", "")
}

// tooltipHoverChanged restarts the delay when the mouse moves to another widget.
func (c *Form) tooltipHoverChanged() {
	c.tooltip.hoverStart = time.Now()
	c.tooltip.suppressed = false
	c.tooltipHide()
}

// tooltipSuppress hides the tooltip until the mouse moves to another widget.
func (c *Form) tooltipSuppress() {
	c.tooltip.suppressed = true
	c.tooltipHide()
}

func (c *Form) tooltipHide() {
	if c.tooltip.visible {
		c.tooltip.visible = false
		c.Update()
	}
}

func (c *Form) tooltipProcessTimer() {
	if c.tooltip.visible || c.tooltip.suppressed || c.mouseLeftButtonPressed {
		return
	}
	if c.tooltipText() == "" || time.Since(c.tooltip.hoverStart) < tooltipDelay {
		return
	}
	c.tooltip.visible = true
	c.tooltip.x = c.lastMouseX
	c.tooltip.y = c.lastMouseY
	c.Update()
}

func (c *Form) tooltipPaint(cnv *Canvas) {
	if !c.tooltip.visible {
		return
	}
	text := c.tooltipText()
	if text == "" {
		return
	}

	fontFamily := ThemeFontFamily()
	fontSize := ThemeFontSize()
	textWidth, textHeight, err := MeasureText(fontFamily, fontSize, text)
	if err != nil {
		return
	}
	w := textWidth + tooltipPadding*4
	h := textHeight + tooltipPadding*2

	// Below-right of the cursor, kept inside the window
	x := c.tooltip.x + 12
	y := c.tooltip.y + 20
	if x+w > c.width {
		x = c.width - w
	}
	if y+h > c.height {
		y = c.tooltip.y - h - 4
	}
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	cnv.SetColor(ThemeBackgroundColor(8, ""))
	cnv.FillRoundedRect(x, y, w, h, 4)
	cnv.SetColor(ThemeBackgroundColor(14, ""))
	cnv.DrawRoundedRect(x, y, w, h, 4)

	cnv.SetHAlign(HAlignCenter)
	cnv.SetVAlign(VAlignCenter)
	cnv.SetColor(ThemeForegroundColor(""))
	cnv.SetFontFamily(fontFamily)
	cnv.SetFontSize(fontSize)
	cnv.DrawText(x, y, w, h, text)
}
