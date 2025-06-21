package ex01base

import "github.com/u00io/nuiforms/ui"

func Run() {
	form := ui.NewForm()
	form.SetTitle("Example 01 - Base Form")
	form.Exec()
}
