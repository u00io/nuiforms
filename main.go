package main

import (
	"time"

	"github.com/u00io/nuiforms/examples"
	"github.com/u00io/nuiforms/ui"
)

func thPrintAllWidgets() {
	time.Sleep(1000 * time.Millisecond)
	for {
		ui.PrintAllWidgets()
		time.Sleep(1000 * time.Millisecond)
	}
}

func main() {
	go thPrintAllWidgets()
	examples.RunSample()
}
