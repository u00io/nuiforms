package ex09textbox

import (
	"encoding/base64"

	"github.com/u00io/nuiforms/ui"
)

func Run(form *ui.Form) {
	OnButtonClick := func() {
		tb := form.Panel().FindWidgetByName("tb2").(*ui.TextBox)
		lbl := form.Panel().FindWidgetByName("lblBase64").(*ui.Label)
		text := tb.Text()
		encoded := base64.StdEncoding.EncodeToString([]byte(text))
		lbl.SetText(encoded)
	}

	fns := map[string]func(){
		"OnButtonClick": OnButtonClick,
	}

	_ = fns

	form.Panel().SetLayout(
		`
<column>
	<row>
		<label text="Enter Text:"/>
		<textbox id="tb2" text="" emptyText="Type here..."/>
		<button text="Convert To Base64" onclick="OnButtonClick"/>
	</row>
	<label id="lblBase64" text="Base64 Encoded Text Will Appear Here"/>
	<vspacer />
</column>
		`, fns, nil)

}
