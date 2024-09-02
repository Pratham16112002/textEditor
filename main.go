package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/storage"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(NewCustomTheme())
	w := a.NewWindow("IDE")
	w.Resize(fyne.NewSize(1024, 768))
	win := &gui{win: w, title: binding.NewString()}
	w.SetContent(win.makeGUI())
	w.SetMainMenu(win.makeMenu())
	win.title.AddListener(binding.NewDataListener(func() {
		name, _ := win.title.Get()
		w.SetTitle("IDE : " + name)
	}))
	flag.Usage = func() {
		fmt.Println("Usage : IDE project [directory]")
	}
	flag.Parse()
	if len(flag.Args()) > 0 {
		dirPath := flag.Args()[0]
		dirPath, err := filepath.Abs(dirPath)
		if err != nil {
			fmt.Println("Error resolving the path : ", dirPath)
			return
		}
		dirURI := storage.NewFileURI(dirPath)
		dir, err := storage.ListerForURI(dirURI)
		if err != nil {
			fmt.Println("Error opening project", err)
			return
		}
		win.openProject(dir)
	} else {
		win.ShowCreate(w)
	}
	w.ShowAndRun()
}

func (g *gui) makeMenu() *fyne.MainMenu {
	file := fyne.NewMenu("File", fyne.NewMenuItem("Open project", func() {
		g.openProjectDialog()
	}))

	return fyne.NewMainMenu(file)
}
