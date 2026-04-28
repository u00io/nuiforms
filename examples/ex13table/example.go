package ex13table

import (
	_ "embed"
	"fmt"
	"image"
	"image/png"
	"strings"

	"github.com/u00io/nuiforms/ui"
)

func Run(form *ui.Form) {
	form.Panel().RemoveAllWidgets()
	mainWidget := NewTableWidget()
	form.Panel().AddWidgetOnGrid(mainWidget, 0, 0)
}

//go:embed testImage.png
var testImage []byte

func TestImage() image.Image {
	im, _ := png.Decode(strings.NewReader(string(testImage)))
	return im
}

type TableWidget struct {
	ui.Widget
}

func NewTableWidget() *TableWidget {
	var c TableWidget
	c.InitWidget()
	c.SetTypeName("TableWidget")

	c.SetLayout(`
		<table id="lvItems">
			<columns>
				<column text="ID" width="100" />
				<column text="Name" width="200" />
			</columns>
			<rows>
				<row>
					<cell text="1"></cell>
					<cell text="John Doe"></cell>
				</row>
			</rows>
		</table>
	`, &c, nil)

	lvItems := c.FindWidgetByName("lvItems").(*ui.Table)
	lvItems.SetRowCount(100)
	for i := 0; i < 100; i++ {
		lvItems.SetCellText2(i, 0, fmt.Sprint(i+1))
		lvItems.SetCellText2(i, 1, "Name "+fmt.Sprint(i+1))
	}

	return &c
}
