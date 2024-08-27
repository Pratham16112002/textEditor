package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type gui struct {
	win   fyne.Window
	title binding.String
}

func makeBanner() fyne.CanvasObject {
	toolbar := widget.NewToolbar(widget.NewToolbarAction(theme.HomeIcon(), func() {
		fmt.Print("Home button is clicked")
	}))

	logo := canvas.NewImageFromResource(resourceLogoSvg)
	logo.FillMode = canvas.ImageFillContain
	return container.NewStack(toolbar, container.NewPadded(logo))
}

func (g *gui) makeGUI() fyne.CanvasObject {
	top := makeBanner()
	left := widget.NewLabel("Left")
	right := widget.NewLabel("Right")

	directory := widget.NewLabelWithData(g.title)
	content := container.NewStack(canvas.NewRectangle(color.Gray{Y: 0xee}), directory)
	seperators := [3]fyne.CanvasObject{
		widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(),
	}
	obj := []fyne.CanvasObject{
		content,
		top,
		left,
		right,
		seperators[0],
		seperators[1],
		seperators[2],
	}
	return container.New(CustomLayout(top, left, right, content, seperators), obj...)
}

func (g *gui) openProjectDialog() {
	dialog.ShowFolderOpen(func(dir fyne.ListableURI, e error) {
		if e != nil {
			dialog.ShowError(e, g.win)
			return
		}
		if dir == nil {
			return
		}
		g.openProject(dir)
	}, g.win)
}

func (g *gui) openProject(dir fyne.ListableURI) {
	g.title.Set(dir.Name())
	return
}
