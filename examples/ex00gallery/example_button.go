package ex00gallery

import (
	"fmt"

	"github.com/u00io/nuiforms/ui"
)

type ExamplePageButton struct {
	ui.Widget

	lblStatus  *ui.Label
	clickCount int
}

func NewExamplePageButton() *ExamplePageButton {
	var c ExamplePageButton
	c.InitWidget()

	row := 0

	c.lblStatus = c.AddLabel(row, 0, "Click a button below")
	c.lblStatus.SetUnderline(true)
	row = addSectionGap(&c.Widget, row+1)

	row = addSectionHeader(&c.Widget, row, "Basic button")
	c.AddButton(row, 0, "Click me", func() {
		c.clickCount++
		c.setStatus(fmt.Sprintf("Basic button clicked %d time(s)", c.clickCount))
	})
	row = addSectionGap(&c.Widget, row+1)

	row = addSectionHeader(&c.Widget, row, "Disabled button")
	btnDisabled := c.AddButton(row, 0, "Can't click me", func() {
		c.setStatus("This should never show up")
	})
	btnDisabled.SetEnabled(false)
	row = addSectionGap(&c.Widget, row+1)

	row = addSectionHeader(&c.Widget, row, "Button sized to its text")
	c.AddButton(row, 0, "A button with a much longer label", func() {
		c.setStatus("Long button clicked")
	})
	row = addSectionGap(&c.Widget, row+1)

	row = addSectionHeader(&c.Widget, row, "Multiple buttons, shared status")
	c.AddButton(row, 0, "Yes", func() { c.setStatus("You chose: Yes") })
	c.AddButton(row, 1, "No", func() { c.setStatus("You chose: No") })
	c.AddButton(row, 2, "Cancel", func() { c.setStatus("You chose: Cancel") })
	row = addSectionGap(&c.Widget, row+1)

	c.AddVSpacer(row, 0)

	return &c
}

func (c *ExamplePageButton) setStatus(text string) {
	c.lblStatus.SetText(text)
}
