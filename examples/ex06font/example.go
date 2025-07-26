package ex06font

import (
	"fmt"

	"github.com/u00io/nuiforms/ui"
)

func Run(form *ui.Form) {
	panel := ui.NewPanel()

	cb1 := ui.NewCheckbox("Checkbox 1")
	cb1.SetOnStateChanged(func(btn *ui.Checkbox, checked bool) {
		fmt.Println("Checkbox 1 checked:", checked)
	})

	rb1 := ui.NewRadioButton("Option 1")
	rb1.SetOnStateChanged(func(btn *ui.RadioButton, checked bool) {
		fmt.Println("RadioButton 1 checked:", checked)
		cb1.SetChecked(checked)
	})
	panel.AddWidgetOnGrid(rb1, 0, 0)

	rb2 := ui.NewRadioButton("Option 2")
	rb2.SetOnStateChanged(func(btn *ui.RadioButton, checked bool) {
		fmt.Println("RadioButton 2 checked:", checked)
	})
	panel.AddWidgetOnGrid(rb2, 0, 1)

	rb1.SetChecked(true) // Set the first radio button as checked by default

	panel.AddWidgetOnGrid(cb1, 0, 2)

	form.Panel().AddWidgetOnGrid(panel, 0, 0)
}
