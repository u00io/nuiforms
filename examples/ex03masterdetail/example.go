package ex03masterdetail

import (
	"github.com/u00io/nuiforms/ui"
)

func Run(form *ui.Form) {
	form.SetTitle("Hello, World!")
	form.SetSize(800, 600)
	masterWidget := NewMasterWidget()
	form.Panel().AddWidget(masterWidget)
}
