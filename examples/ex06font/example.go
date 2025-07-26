package ex06font

import (
	"bytes"
	_ "embed"
	"image"

	"github.com/u00io/nuiforms/ui"
)

// Embed file: test_image.png
//
//go:embed test_image.png
var testImage []byte

func Run(form *ui.Form) {
	img, _, _ := image.Decode(bytes.NewReader(testImage))

	panel := ui.NewPanel()
	imgBox := ui.NewImageBox()
	imgBox.SetImage(img)
	panel.AddWidgetOnGrid(imgBox, 0, 0)
	form.Panel().AddWidgetOnGrid(panel, 0, 0)
}
