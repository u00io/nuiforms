package ex02label

import (
	"github.com/u00io/nuiforms/ui"
)

type MainForm struct {
	ui.Widget
}

func NewMainForm() *MainForm {
	var c MainForm
	c.InitWidget()
	c.SetLayout(`
		<column>
			<label id="lblMain" text="This is a label example."/>
			<vspacer />
			<row>
				<button text="Set Empty Text" onclick="BtnSetEmptyText"/>
				<button text="Set Normal Text" onclick="BtnSetNormalText"/>
				<button text="Set Big Text" onclick="BtnSetBigText"/>
			</row>
			<row>
				<button text="Set Expandable" onclick="BtnSetExpandable"/>			
				<button text="Set Non-Expandable" onclick="BtnSetNonExpandable"/>
				<hspacer />
			</row>
		</column>
	`, &c, nil)
	return &c
}

func (c *MainForm) BtnSetEmptyText() {
	lbl := c.FindWidgetByName("lblMain").(*ui.Label)
	if lbl != nil {
		lbl.SetText("")
	}
}

func (c *MainForm) BtnSetNormalText() {
	lbl := c.FindWidgetByName("lblMain").(*ui.Label)
	if lbl != nil {
		lbl.SetText("Data")
	}
}

func (c *MainForm) BtnSetBigText() {
	lbl := c.FindWidgetByName("lblMain").(*ui.Label)
	if lbl != nil {
		lbl.SetText("This is a very big text example to show how the label widget can handle larger amounts of text. It should properly wrap and display all the content without any issues.")
	}
}

func (c *MainForm) BtnSetExpandable() {
	lbl := c.FindWidgetByName("lblMain").(*ui.Label)
	if lbl != nil {
		lbl.SetXExpandable(true)
	}
}

func (c *MainForm) BtnSetNonExpandable() {
	lbl := c.FindWidgetByName("lblMain").(*ui.Label)
	if lbl != nil {
		lbl.SetXExpandable(false)
	}
}

func Run(form *ui.Form) {
	form.SetTitle("Expample 02 - Label")
	form.Panel().AddWidgetOnGrid(NewMainForm(), 0, 0)
}
