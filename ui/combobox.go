package ui

import (
	"image/color"

	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
)

type ComboBox struct {
	Widget
	items         []*ComboBoxItem
	selectedIndex int
}

type ComboBoxItem struct {
	text string
	data interface{}
}

// DefaultComboBoxMinWidth keeps a bare ComboBox (no width ever set by the
// caller) from collapsing to an unusably thin trigger.
const DefaultComboBoxMinWidth = 150

func NewComboBox() *ComboBox {
	var c ComboBox
	c.InitWidget()
	// Height is fixed (min == max): without a max, a ComboBox alone in a
	// row - with nothing else in the page actually YExpandable - absorbs
	// all leftover vertical space via the layout's "grow every row"
	// fallback, ballooning into a huge empty box instead of a compact
	// dropdown trigger.
	c.SetMinSize(DefaultComboBoxMinWidth, 32)
	c.SetMaxSize(10000, 32)

	c.SetOnMouseDown(func(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
		if button == nuimouse.MouseButtonLeft {
			c.OpenPopup()
			return true
		}
		return false
	})

	c.SetTypeName("ComboBox")
	c.SetMouseCursor(nuimouse.MouseCursorPointer)
	c.SetOnPaint(c.draw)
	return &c
}

func (c *ComboBox) AddItem(text string, data interface{}) {
	var item ComboBoxItem
	item.text = text
	item.data = data
	c.items = append(c.items, &item)
}

func (c *ComboBox) SetSelectedIndex(index int) {
	if index < 0 || index >= len(c.items) {
		return
	}
	c.selectedIndex = index
	c.form.Update()
}

func (c *ComboBox) SelectedItemText() string {
	if c.selectedIndex < 0 || c.selectedIndex >= len(c.items) {
		return ""
	}
	return c.items[c.selectedIndex].text
}

func (c *ComboBox) SelectedItemData() interface{} {
	if c.selectedIndex < 0 || c.selectedIndex >= len(c.items) {
		return nil
	}
	return c.items[c.selectedIndex].data
}

func (c *ComboBox) OpenPopup() {
	popup := NewComboBoxPopup()
	// The popup is a fresh widget, never added via AddWidget, so nothing
	// else ever attaches it to a form - ShowPopup's c.form.Panel() would
	// nil-panic without this (the same class of bug as an unattached
	// ContextMenu submenu).
	popup.attachToForm(c.form)
	// The dropdown should never look narrower than the control it drops
	// from, even though it's free to grow wider to fit long item text.
	popup.triggerWidth = c.Width()
	// So the popup can highlight whichever item is currently selected.
	popup.selectedIndex = c.selectedIndex
	for _, item := range c.items {
		popup.AddItem(item.text, func(index int) {
			c.SetSelectedIndex(index)
			c.form.Update()
		})
	}
	x, y := c.RectClientAreaOnWindow()
	popup.ShowPopup(x, y+c.Height())
}

// Size and placement of the dropdown arrow drawn at the right of the control.
const (
	comboBoxArrowWidth   = 10
	comboBoxArrowHeight  = 5
	comboBoxArrowPadding = 10
)

func (c *ComboBox) draw(cnv *Canvas) {
	backColor := c.BackgroundColorWithAddElevation(-1)
	if c.IsHovered() {
		backColor = c.BackgroundColorWithAddElevation(2)
	}
	foreColor := c.ForegroundColor()

	itemText := c.SelectedItemText()

	cnv.FillRect(0, 0, c.Width(), c.Height(), backColor)
	cnv.SetHAlign(HAlignLeft)
	cnv.SetVAlign(VAlignCenter)
	cnv.SetColor(foreColor)
	cnv.SetFontFamily(c.FontFamily())
	cnv.SetFontSize(c.FontSize())
	textAreaWidth := c.Width() - comboBoxArrowPadding*2 - comboBoxArrowWidth
	cnv.DrawText(0, 0, textAreaWidth, c.Height(), itemText)

	c.drawArrow(cnv, foreColor)
}

// drawArrow paints a small downward-pointing triangle - the usual dropdown
// indicator - at the right edge of the control. Canvas has no filled-polygon
// primitive, so it's built from one DrawLine per row, narrowing symmetrically
// until the last row is a single point.
func (c *ComboBox) drawArrow(cnv *Canvas, arrowColor color.Color) {
	x := c.Width() - comboBoxArrowPadding - comboBoxArrowWidth
	y := (c.Height() - comboBoxArrowHeight) / 2
	halfWidth := comboBoxArrowWidth / 2

	for row := 0; row < comboBoxArrowHeight; row++ {
		inset := row * halfWidth / (comboBoxArrowHeight - 1)
		cnv.DrawLine(x+inset, y+row, x+comboBoxArrowWidth-inset, y+row, 1, arrowColor)
	}
}

// Adaptive popup width bounds: never wider than this, regardless of how
// long an item's text is.
const comboBoxPopupMaxWidth = 420
const comboBoxItemPadding = 10

type comboBoxPopup struct {
	Widget

	items []*comboBoxPopupItem
	// triggerWidth is the owning ComboBox's own width - the floor for the
	// popup's width, set by ComboBox.OpenPopup before ShowPopup runs.
	triggerWidth int
	// selectedIndex is the owning ComboBox's current selection, also set by
	// ComboBox.OpenPopup, so the matching item can be highlighted.
	selectedIndex int
}

func NewComboBoxPopup() *comboBoxPopup {
	var c comboBoxPopup
	c.InitWidget()
	c.SetTypeName("ComboBoxPopup")
	c.SetAbsolutePositioning(true)
	c.SetElevation(3)
	c.SetAutoFillBackground(true)
	c.SetOnPostPaint(c.drawBorder)
	return &c
}

// drawBorder matches ContextMenu's treatment - a subtle outline so the
// dropdown reads as a distinct surface instead of blending into whatever is
// behind it.
func (c *comboBoxPopup) drawBorder(cnv *Canvas) {
	borderColor := ThemeForegroundColor("")
	borderColor.A = contextMenuBorderAlpha
	cnv.SetColor(borderColor)
	cnv.DrawRect(0, 0, c.Width(), c.Height())
}

func (c *comboBoxPopup) ShowPopup(x int, y int) {
	c.SetPosition(x, y)
	c.rebuildVisualElements()
	c.form.Panel().AppendPopupWidget(c)
	c.form.Update()
}

func (c *comboBoxPopup) AddItem(text string, onClick func(index int)) {
	index := len(c.items)
	item := newComboBoxPopupItem(index, text)
	item.parentWidgetId = c.Id()
	item.OnClick = onClick
	item.selected = index == c.selectedIndex
	c.items = append(c.items, item)
	c.AddWidget(0, 0, item)
}

func (c *comboBoxPopup) rebuildVisualElements() {
	width := c.contentWidth()

	yOffset := 0
	for _, item := range c.items {
		item.SetPosition(0, yOffset)
		item.SetSize(width, ContextMenuItemHeight)
		yOffset += ContextMenuItemHeight
	}
	c.SetSize(width, yOffset)
}

// contentWidth is never narrower than triggerWidth (the owning ComboBox's
// own width), grows to fit the widest item text when that text wouldn't
// otherwise fit, and never exceeds comboBoxPopupMaxWidth.
func (c *comboBoxPopup) contentWidth() int {
	width := c.triggerWidth
	for _, item := range c.items {
		textWidth, _, err := MeasureText(item.FontFamily(), item.FontSize(), item.text)
		if err != nil {
			continue
		}
		itemWidth := comboBoxItemPadding*2 + textWidth
		if itemWidth > width {
			width = itemWidth
		}
	}
	if width > comboBoxPopupMaxWidth {
		width = comboBoxPopupMaxWidth
	}
	if width < c.triggerWidth {
		width = c.triggerWidth
	}
	return width
}

type comboBoxPopupItem struct {
	Widget
	index    int
	text     string
	selected bool
	OnClick  func(index int)
}

func newComboBoxPopupItem(index int, text string) *comboBoxPopupItem {
	var item comboBoxPopupItem
	item.InitWidget()
	item.SetTypeName("ComboBoxPopupItem")
	item.SetAbsolutePositioning(true)
	item.SetMouseCursor(nuimouse.MouseCursorPointer)
	item.SetOnPaint(item.Draw)
	item.SetOnMouseDown(item.mouseDownHandler)

	item.index = index
	item.text = text
	return &item
}

func (c *comboBoxPopupItem) Draw(ctx *Canvas) {
	backColor := c.BackgroundColorWithAddElevation(-1)
	if c.selected {
		backColor = c.BackgroundColorWithAddElevation(1)
	}
	if c.IsHovered() {
		backColor = c.BackgroundColorWithAddElevation(2)
	}
	ctx.FillRect(0, 0, c.InnerWidth(), c.InnerHeight(), backColor)
	ctx.SetHAlign(HAlignLeft)
	ctx.SetVAlign(VAlignCenter)
	ctx.SetColor(c.ForegroundColor())
	ctx.SetFontFamily(c.FontFamily())
	ctx.SetFontSize(c.FontSize())
	ctx.DrawText(comboBoxItemPadding, 0, c.Width()-comboBoxItemPadding*2, c.Height(), c.text)
}

func (c *comboBoxPopupItem) mouseDownHandler(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
	if c.OnClick != nil {
		c.OnClick(c.index)
	}
	c.form.Panel().CloseTopPopup()
	return true
}
