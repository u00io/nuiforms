package ui

import (
	"image/color"
)

type TabWidget struct {
	Widget
	pages        []tabWidgetPage
	panelTop     *Panel
	panelContent *Panel
	currentPage  int
}

type tabWidgetPage struct {
	name   string
	widget Widgeter
}

func NewTabWidget() *TabWidget {
	var c TabWidget
	c.InitWidget()
	c.SetTypeName("TabWidget")

	c.SetXExpandable(true)
	c.SetYExpandable(true)

	c.pages = make([]tabWidgetPage, 0)

	c.panelTop = NewPanel()
	c.AddWidgetOnGrid(c.panelTop, 0, 0)

	c.panelContent = NewPanel()
	c.AddWidgetOnGrid(c.panelContent, 0, 1)

	return &c
}

func (c *TabWidget) AddPage(name string, widgeter Widgeter) {
	if c.pages == nil {
		c.pages = make([]tabWidgetPage, 0)
	}

	var p tabWidgetPage
	p.name = name
	p.widget = widgeter
	c.pages = append(c.pages, p)

	c.rebuildInterface()
}

func (c *TabWidget) rebuildInterface() {
	// Build buttons
	c.panelTop.RemoveAllWidgets()
	for i, page := range c.pages {
		btn := NewTabWidgetButton(page.name)
		btn.SetOnButtonClick(func(btn *tabWidgetButton) {
			c.currentPage = i
			c.rebuildInterface()
		})
		if i == c.currentPage {
			btn.SetBackgroundColor(color.RGBA{R: 0, G: 70, B: 200, A: 255}) // Highlight color
		} else {
			btn.SetBackgroundColor(color.RGBA{R: 0, G: 0, B: 0, A: 255}) // Default color
		}
		c.panelTop.AddWidgetOnGrid(btn, i, 0)
	}

	c.panelTop.AddWidgetOnGrid(NewHSpacer(), len(c.pages), 0)

	// Build content
	c.panelContent.RemoveAllWidgets()
	if c.currentPage >= 0 && c.currentPage < len(c.pages) {
		c.panelContent.AddWidgetOnGrid(c.pages[c.currentPage].widget, 0, 0)
	}
	UpdateMainForm()
}
