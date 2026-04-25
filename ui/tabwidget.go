package ui

import (
	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
)

type TabWidget struct {
	Widget
	pages        []tabWidgetPage
	panelTop     *tabWidgetHeader
	panelContent *Panel
	currentPage  int

	headerHeight       int
	headerItemMinWidth int
}

type tabWidgetPage struct {
	name   string
	widget Widgeter
}

const (
	tabWidgetBorderAlpha = 48
)

func NewTabWidget() *TabWidget {
	var c TabWidget
	c.InitWidget()
	c.SetPanelPadding(0)
	c.SetCellPadding(0)

	c.headerHeight = 30
	c.headerItemMinWidth = 100

	c.SetTypeName("TabWidget")

	c.SetXExpandable(true)
	c.SetYExpandable(true)

	c.pages = make([]tabWidgetPage, 0)

	c.panelTop = newTabWidgetHeader(c.headerItemMinWidth, c.headerHeight)
	c.panelTop.SetMinHeight(c.headerHeight)
	c.panelTop.SetMaxHeight(c.headerHeight)
	c.panelTop.onTabChanged = c.onTabChanged
	c.AddWidgetOnGrid(c.panelTop, 0, 0)

	c.panelContent = NewPanel()
	c.AddWidgetOnGrid(c.panelContent, 1, 0)

	c.SetOnPostPaint(c.drawPost)

	return &c
}

func (c *TabWidget) onTabChanged(index int) {
	c.currentPage = index
	c.rebuildInterface()
}

func (c *TabWidget) SetLayoutXml(n *uiNode, eventProcessor interface{}, widgets map[string]Widgeter) {
	for _, child := range n.Nodes {
		if child.XMLName.Local == "tab" {
			childNode := child
			childNode.XMLName.Local = "panel"
			name := child.GetAttrValueByName("text", "")
			wContent := NewPanel()
			wContent.buildNode(&childNode, wContent, 0, 0, eventProcessor, widgets)
			c.AddPage(name, wContent)
		}
	}
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
	items := make([]string, len(c.pages))
	for i, page := range c.pages {
		items[i] = page.name
	}
	c.panelTop.SetItems(items)
	//c.panelTop.AddWidgetOnGrid(NewHSpacer(), 0, len(c.pages))

	// Build content
	c.panelContent.RemoveAllWidgets()
	c.panelContent.SetPanelPadding(0)
	//c.panelContent.SetAutoFillBackground(true)
	//c.panelContent.SetBackgroundColor(ColorFromHex("#222222"))
	if c.currentPage >= 0 && c.currentPage < len(c.pages) {
		c.panelContent.AddWidgetOnGrid(c.pages[c.currentPage].widget, 0, 0)
	}
	UpdateMainForm()
}

func (c *TabWidget) drawPost(cnv *Canvas) {
	// Draw border
	cnv.SetColor(c.ForegroundColor())

	borderColor := ThemeForegroundColor("")
	borderColor.A = tabWidgetBorderAlpha

	// Border left
	cnv.DrawLine(0, c.headerHeight, 0, c.Height(), 1, borderColor)
	// Border bottom
	cnv.DrawLine(0, c.Height()-1, c.Width(), c.Height()-1, 1, borderColor)
	// Border right
	cnv.DrawLine(c.Width()-1, c.headerHeight, c.Width()-1, c.Height(), 1, borderColor)
}

type tabWidgetHeader struct {
	Widget
	items        []string
	minWidth     int
	height       int
	currentIndex int

	itemsWidths []int

	onTabChanged       func(index int)
	onTabChangedCalled bool
}

func newTabWidgetHeader(minWidth int, height int) *tabWidgetHeader {
	var c tabWidgetHeader
	c.InitWidget()
	c.minWidth = minWidth
	c.height = height
	c.SetTypeName("TabWidgetHeader")
	c.SetOnPaint(c.draw)

	c.SetOnMouseDown(c.onMouseDown)
	c.SetOnMouseUp(c.onMouseUp)
	c.SetOnMouseMove(c.onMouseMove)
	return &c
}

func (c *tabWidgetHeader) setCurrentPage(index int) {
	c.currentIndex = index
	if c.onTabChanged != nil {
		if !c.onTabChangedCalled {
			c.onTabChangedCalled = true
			c.onTabChanged(index)
			c.onTabChangedCalled = false
		}
	}
	UpdateMainForm()
}

func (c *tabWidgetHeader) SetItems(items []string) {
	c.items = items
	UpdateMainForm()
}

func (c *tabWidgetHeader) itemByCoords(x int, y int) int {
	if y < 0 || y > c.height {
		return -1
	}

	xOffset := 0
	for i, width := range c.itemsWidths {
		if x >= xOffset && x < xOffset+width {
			return i
		}
		xOffset += width
	}

	return -1
}

func (c *tabWidgetHeader) onMouseDown(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
	index := c.itemByCoords(x, y)
	if index >= 0 && index < len(c.items) {
		c.setCurrentPage(index)
		return true
	}
	return false
}

func (c *tabWidgetHeader) onMouseUp(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
	return false
}

func (c *tabWidgetHeader) onMouseMove(x int, y int, mods nuikey.KeyModifiers) bool {
	// hover - mouse cursor pointer
	index := c.itemByCoords(x, y)
	if index >= 0 && index < len(c.items) {
		c.SetMouseCursor(nuimouse.MouseCursorPointer)
	} else {
		c.SetMouseCursor(nuimouse.MouseCursorNotDefined)
	}
	UpdateMainForm()
	return true
}

func (c *tabWidgetHeader) draw(cnv *Canvas) {
	cnv.SetColor(c.ForegroundColor())
	cnv.SetFontFamily(c.FontFamily())
	cnv.SetFontSize(c.FontSize())

	c.itemsWidths = make([]int, len(c.items))
	xOffset := 0
	for i, item := range c.items {
		width := c.minWidth

		// Measure text width
		textWidth, textHeight, err := MeasureText(c.FontFamily(), c.FontSize(), item)
		_ = textHeight
		if err == nil {
			textPadding := 20
			width = textWidth + textPadding
			width = max(width, c.minWidth)
		}

		c.itemsWidths[i] = width

		x := xOffset
		cnv.SetHAlign(HAlignCenter)
		cnv.SetVAlign(VAlignCenter)
		cnv.DrawText(x, 0, width, c.height, item)

		yOffset := 2
		if i == c.currentIndex {
			yOffset = 0
		}

		borderColor := ThemeForegroundColor("")
		borderColor.A = tabWidgetBorderAlpha

		// Left Border
		cnv.DrawLine(x, yOffset, x, c.height, 1, borderColor)
		// Top Border
		cnv.DrawLine(x, yOffset, x+width, yOffset, 1, borderColor)
		// Right Border
		cnv.DrawLine(x+width, yOffset, x+width, c.height, 1, borderColor)

		if i != c.currentIndex {
			// Bottom Border
			cnv.DrawLine(x, c.height-1, x+width, c.height-1, 1, borderColor)
		}

		xOffset += width
	}

	// Fill remaining space with bottom border
	if xOffset < c.Width() {
		borderColor := ThemeForegroundColor("")
		borderColor.A = tabWidgetBorderAlpha
		cnv.DrawLine(xOffset, c.height-1, c.Width(), c.height-1, 1, borderColor)
	}
}
