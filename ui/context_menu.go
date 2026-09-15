package ui

const ContextMenuItemHeight = 32

type ContextMenu struct {
	Widget
	menuWidth  int
	menuHeight int
	items      []*ContextMenuItem
	CloseEvent func()
	parentMenu *ContextMenu
}

const contextMenuBorderAlpha = 80

// Adaptive width bounds: the menu shrinks to fit short item text and grows
// for long item text, but never past these limits.
const contextMenuMinWidth = 140
const contextMenuMaxWidth = 420

func NewContextMenu(parent Widgeter) *ContextMenu {
	var c ContextMenu
	c.InitWidget()
	c.SetAbsolutePositioning(true)
	c.SetTypeName("ContextMenu")
	c.SetName("PopupMenuPanel")
	c.SetElevation(3)
	c.SetAutoFillBackground(true)
	c.SetOnPostPaint(c.drawBorder)
	return &c
}

// drawBorder gives the popup a subtle outline so it reads as a distinct
// surface instead of blending into whatever is behind it.
func (c *ContextMenu) drawBorder(cnv *Canvas) {
	borderColor := ThemeForegroundColor("")
	borderColor.A = contextMenuBorderAlpha
	cnv.SetColor(borderColor)
	cnv.DrawRect(0, 0, c.Width(), c.Height())
}

func (c *ContextMenu) ShowMenu(x int, y int) {
	c.SetPosition(x, y)
	c.rebuildVisualElements()
	c.form.Panel().AppendPopupWidget(c)
	c.form.Update()
}

func (c *ContextMenu) showMenu(x int, y int, parentMenu *ContextMenu) {
	c.CloseAfterPopupWidget(parentMenu)
	c.parentMenu = parentMenu
	c.SetPosition(x, y)
	c.rebuildVisualElements()
	//c.Window().AppendPopup(c)
	c.form.Panel().AppendPopupWidget(c)
}

func (c *ContextMenu) ClosePopup() {
	if c.CloseEvent != nil {
		c.CloseEvent()
	}
}

func (c *ContextMenu) AddItem(text string, onClick func()) *ContextMenuItem {
	item := NewContextMenuItem()
	item.parentWidgetId = c.Id()
	item.SetText(text)
	item.OnClick = onClick

	c.items = append(c.items, item)
	c.AddWidget(0, 0, item)
	return item
}

func (c *ContextMenu) AddItemWithSubmenu(text string, innerMenu *ContextMenu) *ContextMenuItem {
	item := NewContextMenuItem()
	item.parentWidgetId = c.Id()
	item.SetText(text)
	item.innerMenu = innerMenu
	c.items = append(c.items, item)
	c.AddWidget(0, 0, item)
	return item
}

func (c *ContextMenu) RemoveAllItems() {
	c.RemoveAllWidgets()
	c.rebuildVisualElements()
	c.form.Update()
}

func (c *ContextMenu) OnInit() {
	c.rebuildVisualElements()
}

func (c *ContextMenu) needToClose() {
	c.form.Panel().CloseTopPopup()
	if c.parentMenu != nil {
		c.parentMenu.needToClose()
	}
}

func (c *ContextMenu) rebuildVisualElements() {
	menuWidth := c.contentWidth()

	yOffset := 0
	for _, item := range c.items {
		item.needToClosePopupMenu = c.needToClose
		item.parentMenu = c
		item.SetPosition(0, yOffset)
		item.SetSize(menuWidth, ContextMenuItemHeight)
		yOffset += ContextMenuItemHeight
	}
	c.SetSize(menuWidth, yOffset)
	c.menuWidth = menuWidth
	c.menuHeight = yOffset
}

// contentWidth measures the widest item text (reserving room for the
// submenu arrow on items that have one) and clamps the result between
// contextMenuMinWidth and contextMenuMaxWidth.
func (c *ContextMenu) contentWidth() int {
	width := contextMenuMinWidth
	for _, item := range c.items {
		textWidth, _, err := MeasureText(item.FontFamily(), item.FontSize(), item.text)
		if err != nil {
			continue
		}
		itemWidth := contextMenuItemPadding*2 + textWidth
		if item.innerMenu != nil {
			itemWidth += ContextMenuItemHeight + contextMenuItemPadding
		}
		if itemWidth > width {
			width = itemWidth
		}
	}
	if width > contextMenuMaxWidth {
		width = contextMenuMaxWidth
	}
	return width
}
