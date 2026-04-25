package ex13table

import (
	_ "embed"
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
		<table>
			<columns>
				<column text="ID" width="100" />
				<column text="Name" width="200" />
			</columns>
			<rows>
				<row>
					<cell>1</cell>
					<cell>John Doe</cell>
				</row>
				<row>
					<cell>2</cell>
					<cell>Jane Smith</cell>
				</row>
			</rows>
		</table>
	`, &c, nil)
	return &c
}
