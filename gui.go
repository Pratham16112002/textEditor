package main

import (
	"errors"
	"fmt"
	"image/color"
	"io"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Pratham16112002/textEditor.git/internal/dialogs"
	"github.com/Pratham16112002/textEditor.git/internal/editors"
)

type gui struct {
	win      fyne.Window
	title    binding.String
	fileTree binding.URITree
	content  *container.DocTabs
	openTabs map[fyne.URI]*container.TabItem
}

func (g *gui) makeBanner() fyne.CanvasObject {
	title := canvas.NewText("IDE", theme.ForegroundColor())
	title.TextSize = 14
	title.TextStyle = fyne.TextStyle{Bold: true}
	g.title.AddListener(binding.NewDataListener(func() {
		name, _ := g.title.Get()
		if name == "" {
			name = "IDE"
		}
		title.Text = name
		title.Refresh()
	}))
	home := widget.NewButtonWithIcon("Home", theme.HomeIcon(), func() {
	})
	new_file := widget.NewButtonWithIcon("New", theme.FileIcon(), func() {
		new_file_dialog := dialog.NewFileSave(func(f fyne.URIWriteCloser, err error) {
			new_file_editor := editors.NewEditor("")
			new_file_uri := f.URI()
			item := container.NewStack(new_file_editor.TextGrid, new_file_editor.Entry)
			new_file_item := container.NewTabItemWithIcon(
				new_file_uri.Name(),
				theme.FileIcon(),
				item,
			)
			new_file_editor.Show()
			g.content.Append(new_file_item)
			g.openTabs[new_file_uri] = new_file_item
			g.content.Select(new_file_item)
		}, g.win)
		new_file_dialog.Show()
	})
	left := container.NewHBox(home, new_file, title)

	logo := canvas.NewImageFromResource(resourceLogoSvg)
	logo.FillMode = canvas.ImageFillContain
	return container.NewStack(left, container.NewPadded(logo))
}

func (g *gui) makeGUI() fyne.CanvasObject {
	top := g.makeBanner()
	g.fileTree = binding.NewURITree()
	files := widget.NewTreeWithData(g.fileTree, func(isBranch bool) fyne.CanvasObject {
		return widget.NewLabel("filename.jpg")
	}, func(data binding.DataItem, isBranch bool, obj fyne.CanvasObject) {
		l := obj.(*widget.Label)
		u, _ := data.(binding.URI).Get()

		name := u.Name()
		l.SetText(name)
	})
	left := widget.NewAccordion(
		widget.NewAccordionItem("Files", files),
		widget.NewAccordionItem("Screens", widget.NewLabel("TODO Screens")),
	)
	left.Open(0)
	left.MultiOpen = true
	files.OnSelected = func(id widget.TreeNodeID) {
		u, err := g.fileTree.GetValue(id)
		if err != nil {
			dialog.ShowError(err, g.win)
		}
		g.openFile(u)
	}
	project_title, _ := g.title.Get()
	welcome := widget.NewRichTextFromMarkdown(`
    # Welcome
    Open a file from file tree
    `)
	preview := container.NewBorder(nil, nil, nil, nil, welcome)
	content := container.NewStack(canvas.NewRectangle(color.Gray{Y: 0xee}), preview)
	g.content = container.NewDocTabs(
		container.NewTabItemWithIcon(project_title, theme.FileIcon(), content),
	)
	g.content.CloseIntercept = func(i *container.TabItem) {
		var itemURI fyne.URI
		for uri, item := range g.openTabs {
			if i == item {
				itemURI = uri
			}
		}
		if itemURI != nil {
			delete(g.openTabs, itemURI)
		}
		g.content.Remove(i)
	}
	seperators := [2]fyne.CanvasObject{
		widget.NewSeparator(), widget.NewSeparator(),
	}
	obj := []fyne.CanvasObject{
		g.content,
		top,
		left,
		seperators[0],
		seperators[1],
	}
	g.win.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		fmt.Print(key.Physical)
	})
	return container.New(CustomLayout(top, left, nil, g.content, seperators), obj...)
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

func (g *gui) ShowCreate(w fyne.Window) {
	var wizard *dialogs.Wizard
	intro := widget.NewLabel("Create a new project or open an existing")
	create_button := widget.NewButton("Create", func() {
		wizard.Push("Project Details", g.projectDetails(wizard))
	})
	create_button.Importance = widget.HighImportance
	open_button := widget.NewButton("Open", func() {
		wizard.Hide()
		g.openProjectDialog()
	})
	buttons := container.NewGridWithColumns(2, create_button, open_button)
	home := container.NewVBox(intro, buttons)

	wizard = dialogs.NewWizard(
		"Open Project",
		home,
	)
	wizard.Show(w)
	wizard.Resize(home.MinSize().AddWidthHeight(40, 100))
}

func (g *gui) openFile(uri fyne.URI) {
	listable, err := storage.CanList(uri)
	if listable || err != nil {
		if err != nil {
			dialog.ShowError(err, g.win)
			return
		}
		return
	}
	if g.openTabs == nil {
		g.openTabs = make(map[fyne.URI]*container.TabItem)
	}
	if item, ok := g.openTabs[uri]; ok {
		g.content.Select(item)
		g.content.Refresh()
		return
	}
	fileReader, _ := storage.Reader(uri)
	fileContent, _ := io.ReadAll(fileReader)
	edit := editors.NewEditor(string(fileContent))
	item := container.NewTabItemWithIcon(uri.Name(), theme.FileIcon(), tab_item)

	g.openTabs[uri] = item
	for _, tab := range g.content.Items {
		if tab.Text != uri.Name() {
			continue
		}
		for c_uri, child := range g.openTabs {
			if child != tab {
				continue
			}
			parent, _ := storage.Parent(c_uri)
			child.Text = parent.Name() + string([]rune{filepath.Separator}) + child.Text
		}
		parent, _ := storage.Parent(uri)
		item.Text = parent.Name() + string([]rune{filepath.Separator}) + item.Text
		break
	}
	g.content.Append(item)
	g.content.Select(item)
	g.content.Refresh()
}

func (g *gui) projectDetails(wizard *dialogs.Wizard) fyne.CanvasObject {
	homeDir, _ := os.UserHomeDir()
	homeURL := storage.NewFileURI(homeDir)
	chosen, _ := storage.ListerForURI(homeURL)
	name := widget.NewEntry()
	name.Validator = func(in string) error {
		if in == "" {
			return errors.New("Not a valid Project Name")
		}
		return nil
	}
	var dir *widget.Button
	dir = widget.NewButton(chosen.Name(), func() {
		d := dialog.NewFolderOpen(func(f fyne.ListableURI, err error) {
			if err != nil || f == nil {
				return
			}
			dir.SetText(f.Name())
			chosen = f
		}, g.win)
		d.SetLocation(chosen)
		d.Show()
		// TODO open diaglog
	})

	form := widget.NewForm(
		widget.NewFormItem("Name ", name),
		widget.NewFormItem("Parent directory", dir),
	)
	form.OnSubmit = func() {
		if name.Text == "" {
			return
		}
		project, err := createProject(name.Text, chosen)
		if err != nil {
			dialog.ShowError(err, g.win)
			return
		}
		wizard.Hide()
		g.openProject(project)
	}
	return form
}
