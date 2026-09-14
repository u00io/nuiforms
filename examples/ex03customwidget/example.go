package ex03customwidget

import "github.com/u00io/nuiforms/ui"

type Content struct {
	ui.Widget

	panelLeft  *ui.Panel
	panelRight *ui.Panel
}

func NewContent() *Content {
	var c Content
	c.InitWidget()

	panel := ui.NewPanel()
	c.panelLeft = panel.AddPanel(0, 0)
	c.panelLeft.AddButton(0, 0, "Set Mode 0", func() { c.SetMode("0") })
	c.panelLeft.AddButton(1, 0, "Set Mode 1", func() { c.SetMode("1") })
	c.panelLeft.AddVSpacer(10, 0)
	c.panelLeft.SetBackgroundColor(ui.ColorFromHex("#555555"))
	c.panelLeft.SetAutoFillBackground(true)

	c.panelRight = panel.AddPanel(0, 1)
	c.panelRight.AddLabel(0, 0, "12313123")
	c.panelRight.SetBackgroundColor(ui.ColorFromHex("#777777"))
	c.panelRight.SetAutoFillBackground(true)

	menu := ui.NewContextMenu(c.panelRight)
	menu.AddItem("123", nil)
	menu.AddItem("456", nil)
	c.panelRight.SetContextMenu(menu)

	c.AddWidget(panel, 0, 0)

	return &c
}

func (c *Content) SetMode(mode string) {
	c.panelRight.RemoveAllWidgets()
	var w ui.Widgeter
	switch mode {
	case "0":
		w = ui.NewCheckbox("123")
	case "1":
		b := ui.NewButton("POPUP")
		b.SetOnClick(func() {
			p := ui.NewPanel()
			p.SetAbsolutePositioning(true)
			p.SetSize(100, 90)
			p.SetPosition(10, 30)
			p.AddLabel(0, 0, "88888888888")
			c.panelRight.AppendPopupWidget(p)
			p.SetBackgroundColor(ui.ColorFromHex("#888800"))
			c.SetAutoFillBackground(true)
		})
		w = b
	}
	if w != nil {
		c.panelRight.AddWidget(w, 0, 0)
		c.panelRight.AddVSpacer(1, 0)
	}
}

func NewExample() *ui.Form {
	form := ui.NewForm()
	form.Panel().AddWidget(NewContent(), 0, 0)
	return form
}
