package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
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
	openTabs map[fyne.URI]*tabItem
	curDir   fyne.ListableURI
	menu     *fyne.MainMenu
	eContent fyne.CanvasObject
}

type tabItem struct {
	editor  editors.Editor
	tabItem *container.TabItem
}

func (g *gui) makeMenu() {
	save := fyne.NewMenuItem("Save", func() {
		current_selected := g.content.Selected()
		for _, item := range g.openTabs {
			if item.tabItem == current_selected {
				err := item.editor.Save()
				if err != nil {
					dialog.ShowError(err, g.win)
				}
			}
		}
	})
	save.Shortcut = &desktop.CustomShortcut{
		KeyName:  fyne.KeyS,
		Modifier: fyne.KeyModifierShortcutDefault,
	}
	file := fyne.NewMenu("File",
		fyne.NewMenuItem("Open project", g.openProjectDialog),
		fyne.NewMenuItem("New File", func() {
			new_file_dialog := dialog.NewFileSave(func(f fyne.URIWriteCloser, err error) {
				if err != nil || f == nil {
					return
				}
				new_file_uri := f.URI()

				parent_URI, err := storage.Parent(new_file_uri)
				if err != nil {
					dialog.ShowError(err, g.win)
					return
				}
				fmt.Println(g.curDir.Name())
				fmt.Println(parent_URI.Name())
				err = g.fileTree.Append(parent_URI.String(), new_file_uri.String(), new_file_uri)
				if g.curDir.Name() == parent_URI.Name() {
					err = g.fileTree.Append(
						binding.DataTreeRootID,
						new_file_uri.String(),
						new_file_uri,
					)
				} else {
					err = g.fileTree.Append(parent_URI.String(), new_file_uri.String(), new_file_uri)
				}
				if err != nil {
					dialog.ShowError(err, g.win)
					return
				}
				g.win.Content().Refresh()
				g.openFile(new_file_uri)
			}, g.win)
			new_file_dialog.SetLocation(g.curDir)
			new_file_dialog.Show()
		}),
		save,
	)
	if g.menu == nil {
		g.menu = fyne.NewMainMenu()
	}
	g.menu.Items = append(g.menu.Items, file)
}

func updateTime(clock *widget.Label) {
	formattedTime := time.Now().Format("03:04:05")
	canvasTime := canvas.NewText(formattedTime, theme.ForegroundColor())
	canvasTime.TextSize = 14
	clock.SetText(canvasTime.Text)
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
	timeWidget := widget.NewLabel("")
	go func() {
		for range time.Tick(time.Second) {
			updateTime(timeWidget)
		}
	}()
	left := container.NewHBox(home, title)
	wrapper := container.NewHBox(left, layout.NewSpacer(), timeWidget)
	logo := canvas.NewImageFromResource(resourceLogoSvg)
	logo.FillMode = canvas.ImageFillContain
	stack := container.NewStack(wrapper, container.NewPadded(logo))
	return stack
}

func (g *gui) makeGUI() fyne.CanvasObject {
	top := g.makeBanner()
	g.fileTree = binding.NewURITree()
	files := widget.NewTreeWithData(g.fileTree, func(isBranch bool) fyne.CanvasObject {
		temp := binding.NewString()
		temp.Set("filename.jpg")
		fileName, _ := temp.Get()
		return widget.NewLabel(fileName)
	}, func(data binding.DataItem, isBranch bool, obj fyne.CanvasObject) {
		l := obj.(*widget.Label)
		bindingURL := data.(binding.URI)
		u, err := bindingURL.Get()
		if err != nil {
			dialog.ShowError(err, g.win)
			return
		}
		l.SetText(u.Name())
	})
	left := widget.NewAccordion(
		widget.NewAccordionItem("Files", container.NewScroll(files)),
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
	welcome := widget.NewRichTextFromMarkdown(`# Welcome
    Open a file from file tree
`)
	g.content = container.NewDocTabs(
		container.NewTabItem("Welcome", welcome),
	)
	g.content.CloseIntercept = func(i *container.TabItem) {
		var itemURI fyne.URI
		for uri, item := range g.openTabs {
			if i == item.tabItem {
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
	g.win.Canvas().SetOnTypedKey(func(s *fyne.KeyEvent) {
	})
	g.eContent = container.New(CustomLayout(top, left, nil, g.content, seperators), obj...)

	return g.eContent
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
	edit, err := editors.ForURI(uri)
	if err != nil {
		dialog.ShowError(err, g.win)
		return
	}

	if g.openTabs == nil {
		g.openTabs = make(map[fyne.URI]*tabItem)
	}

	if g.openTabs[uri] != nil {
		g.content.Select(g.openTabs[uri].tabItem)
		return
	}

	item := container.NewTabItemWithIcon(
		uri.Name(),
		theme.FileIcon(),
		edit.Content(),
	)
	dirty := edit.Edited()
	dirty.AddListener(binding.NewDataListener(func() {
		isEdited, _ := dirty.Get()
		if isEdited {
			item.Text = uri.Name() + "*"
		} else {
			item.Text = uri.Name()
		}
		g.content.Refresh()
	}))
	g.openTabs[uri] = &tabItem{
		editor:  edit,
		tabItem: item,
	}

	for _, tab := range g.content.Items {
		if tab.Text != uri.Name() {
			continue
		}
		for c_uri, child := range g.openTabs {
			if child.tabItem != tab {
				continue
			}
			child_parent, _ := storage.Parent(c_uri)
			child.tabItem.Text = child_parent.Name() + string(
				[]rune{filepath.Separator},
			) + child.tabItem.Text
		}
		parent, _ := storage.Parent(uri)
		tab.Text = parent.Name() + string([]rune{filepath.Separator}) + tab.Text
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
		g.curDir = chosen
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
