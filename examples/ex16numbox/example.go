package ex16numbox

import (
	"fmt"

	"github.com/u00io/nuiforms/ui"
)

type MainWidget struct {
	ui.Widget

	num    *ui.NumBox
	lblV   *ui.Label
	lblTx  *ui.Label
	lblEvt *ui.Label
}

func NewMainWidget() *MainWidget {
	var c MainWidget
	c.InitWidget()

	title := ui.NewLabel("NumBox demo (float64)")
	title.SetUnderline(true)

	lbl1 := ui.NewLabel("Value:")
	lbl1.SetXExpandable(false)

	c.lblV = ui.NewLabel("")
	c.lblV.SetXExpandable(true)

	lbl2 := ui.NewLabel("Text:")
	lbl2.SetXExpandable(false)

	c.lblTx = ui.NewLabel("")
	c.lblTx.SetXExpandable(true)

	lbl3 := ui.NewLabel("OnChanged event:")
	lbl3.SetXExpandable(false)

	c.lblEvt = ui.NewLabel("")
	c.lblEvt.SetXExpandable(true)

	c.num = ui.NewNumBox()
	c.num.SetDecimals(3)
	c.num.SetMin(-10)
	c.num.SetMax(10)
	c.num.SetStep(0) // default from decimals
	c.num.SetValue(1.234)

	c.num.SetOnChanged(func() {
		// NumBox pushes EventNumBoxValueChanged into CurrentEvent().
		ev := ui.CurrentEvent()
		if ev != nil {
			if e, ok := ev.Parameter.(*ui.EventNumBoxValueChanged); ok {
				c.lblEvt.SetText(fmt.Sprintf("value=%.10g", e.Value))
			}
		}
		c.refresh()
	})

	btnSetPi := ui.NewButton("Set π")
	btnSetPi.SetOnClick(func() {
		c.num.SetValue(3.141592653589793)
	})

	btnSetMin := ui.NewButton("Set Min")
	btnSetMin.SetOnClick(func() {
		c.num.SetValue(c.num.Min())
	})

	btnSetMax := ui.NewButton("Set Max")
	btnSetMax.SetOnClick(func() {
		c.num.SetValue(c.num.Max())
	})

	btnGet := ui.NewButton("GetValue() -> label")
	btnGet.SetOnClick(func() {
		c.refresh()
	})

	// Layout (grid)
	row := 0
	c.AddWidgetOnGrid(title, row, 0)
	row++

	c.AddWidgetOnGrid(c.num, row, 0)
	row++

	// Buttons row
	btnRow := ui.NewPanel()
	btnRow.SetAbsolutePositioning(false)
	btnRow.SetXExpandable(true)
	btnRow.SetYExpandable(false)
	btnRow.AddWidgetOnGrid(btnSetPi, 0, 0)
	btnRow.AddWidgetOnGrid(btnSetMin, 0, 1)
	btnRow.AddWidgetOnGrid(btnSetMax, 0, 2)
	btnRow.AddWidgetOnGrid(btnGet, 0, 3)
	c.AddWidgetOnGrid(btnRow, row, 0)
	row++

	// Value/Text info rows
	info := ui.NewPanel()
	info.SetXExpandable(true)
	info.SetYExpandable(false)
	info.AddWidgetOnGrid(lbl1, 0, 0)
	info.AddWidgetOnGrid(c.lblV, 0, 1)
	info.AddWidgetOnGrid(lbl2, 1, 0)
	info.AddWidgetOnGrid(c.lblTx, 1, 1)
	info.AddWidgetOnGrid(lbl3, 2, 0)
	info.AddWidgetOnGrid(c.lblEvt, 2, 1)
	c.AddWidgetOnGrid(info, row, 0)
	row++

	c.AddWidgetOnGrid(ui.NewVSpacer(), row, 0)

	c.refresh()
	return &c
}

func (c *MainWidget) refresh() {
	c.lblV.SetText(fmt.Sprintf("%.*f", c.num.Decimals(), c.num.Value()))
	c.lblTx.SetText(c.num.Text())
}

func Run(form *ui.Form) {
	form.SetTitle("Example 16 - NumBox")
	form.SetMainWidget(NewMainWidget())
}

