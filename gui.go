package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func makeBanner() fyne.CanvasObject {
	toolbar := widget.NewToolbar(widget.NewToolbarAction(theme.HomeIcon(), func() {
		fmt.Print("Home button is clicked")
	}))

	logo := canvas.NewImageFromResource(resourceLogoSvg)
	logo.FillMode = canvas.ImageFillContain
	return container.NewStack(toolbar, container.NewPadded(logo))
}

func makeGUI() fyne.CanvasObject {
	top := makeBanner()
	left := widget.NewLabel("Left")
	right := widget.NewLabel("Right")
	content := canvas.NewRectangle(color.Gray{Y: 0xee})
	seperators := [3]fyne.CanvasObject{
		widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(),
	}
	obj := []fyne.CanvasObject{content, top, left, right, seperators[0], seperators[1], seperators[2]}
	return container.New(CustomLayout(top, left, right, content, seperators), obj...)
}
