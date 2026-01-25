package ex03textbox

import (
	"github.com/u00io/nuiforms/ui"
)

type MainWidget struct {
	ui.Widget
}

func NewMainWidget() *MainWidget {
	var c MainWidget
	c.InitWidget()
	c.SetLayout(`
		<column>
			<textbox id="txtLine" />
			<textbox id="txtMultiline" multiline="true" />
		</column>
	`, &c, nil)
	return &c
}

func Run(form *ui.Form) {
	form.SetTitle("Hello, World!")
	form.SetMainWidget(NewMainWidget())
}
