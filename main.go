package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("PPU")

	tabs := container.NewAppTabs(
		container.NewTabItem("Categories", widget.NewLabel("Category selection goes here")),
		container.NewTabItem("Entries", widget.NewLabel("Entry selection goes here")),
		container.NewTabItem("Entry", widget.NewLabel("Current Entry bits go here")),
	)

	w.SetContent(tabs)
	w.ShowAndRun()
}
