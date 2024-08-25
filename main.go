package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(NewCustomTheme())
	w := a.NewWindow("IDE")
	w.Resize(fyne.NewSize(1024, 768))
	w.SetContent(makeGUI())
	w.ShowAndRun()
}
