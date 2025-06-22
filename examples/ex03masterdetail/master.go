package ex03masterdetail

import (
	"fmt"
	"image/color"

	"github.com/u00io/nuiforms/ui"
)

type MasterWidget struct {
	widget ui.Widget

	currentTabWidget any

	panelLeft  *ui.Panel
	panelRight *ui.Panel
}

func (c *MasterWidget) Widgeter() any {
	return &c.widget
}

func NewMasterWidget() *MasterWidget {
	var c MasterWidget
	c.widget.SetBackgroundColor(color.RGBA{R: 40, G: 40, B: 60, A: 255})
	c.widget.InitWidget()

	c.panelLeft = ui.NewPanel()
	c.widget.AddWidgetOnGrid(c.panelLeft, 0, 0)
	c.panelLeft.SetMinWidth(300)
	c.panelLeft.SetMaxWidth(300)
	c.panelLeft.SetMaxHeight(300)

	c.panelRight = ui.NewPanel()
	c.panelRight.SetName("panelRight")
	c.widget.AddWidgetOnGrid(c.panelRight, 1, 0)

	{
		btnOpenTab1 := ui.NewButton()
		btnOpenTab1.SetText("Tab 1")
		btnOpenTab1.SetOnButtonClick(c.onBtn1Click)
		btnOpenTab1.SetMinSize(300, 50)
		btnOpenTab1.SetMaxSize(ui.MaxInt, 50)
		c.panelLeft.AddWidgetOnGrid(btnOpenTab1, 0, 0)
	}

	for i := 0; i < 20; i++ {
		btnOpenTab2 := ui.NewButton()
		btnOpenTab2.SetName("btnOpenTab2_" + fmt.Sprint(i))
		btnOpenTab2.SetText("Tab " + fmt.Sprint(i))
		btnOpenTab2.SetOnButtonClick(func(btn *ui.Button) {
			c.closeCurrentTab()
			w := NewTab2Widget("PAGE " + btn.Text())
			c.currentTabWidget = w
			c.panelRight.AddWidgetOnGrid(w, 0, 0)
		})
		btnOpenTab2.SetMinSize(300, 50)
		btnOpenTab2.SetMaxSize(ui.MaxInt, 50)
		c.panelLeft.AddWidgetOnGrid(btnOpenTab2, 0, c.panelLeft.NextGridY())
	}

	c.onBtn1Click(nil)

	return &c
}

func (c *MasterWidget) closeCurrentTab() {
	if c.currentTabWidget != nil {
		c.panelRight.RemoveWidget(c.currentTabWidget)
		c.currentTabWidget = nil
	}
}

func (c *MasterWidget) onBtn1Click(btn *ui.Button) {
	c.closeCurrentTab()
	w := NewTab1Widget()
	c.currentTabWidget = w
	c.panelRight.AddWidgetOnGrid(w, 0, 0)
}
