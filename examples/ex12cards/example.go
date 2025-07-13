package ex12cards

import "github.com/u00io/nuiforms/ui"

type Cards struct {
	ui.Widget

	panelCards *ui.Panel
}

func NewCards() *Cards {
	var c Cards
	c.InitWidget()
	c.SetTypeName("Cards")
	c.SetMaxHeight(600)

	c.AddWidgetOnGrid(ui.NewLabel("Cards Example"), 0, 0)

	c.panelCards = ui.NewPanel()
	c.panelCards.SetName("CardsPanel")
	c.panelCards.SetXExpandable(true)
	c.panelCards.SetYExpandable(true)
	c.panelCards.SetAllowScroll(true, true)
	c.AddWidgetOnGrid(c.panelCards, 0, 1)

	c.AddWidgetOnGrid(ui.NewLabel("Bottom Panel"), 0, 2)

	c.loadCards()

	return &c
}

func (c *Cards) loadCards() {
	cardNames := []string{
		"Category 1",
		"Category 2",
		"Category 3",
	}

	c.panelCards.RemoveAllWidgets()
	for _, name := range cardNames {
		widget := NewCategoryWidget(name)
		c.panelCards.AddWidgetOnGrid(widget, 0, c.panelCards.NextGridY())
	}

}
