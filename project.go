package main

import (
	"fmt"
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/storage"
)

func createProject(name string, parent fyne.ListableURI) (fyne.ListableURI, error) {
	dir, err := storage.Child(parent, name)
	if err != nil {
		return nil, err
	}

	err = storage.CreateListable(dir)
	if err != nil {
		return nil, err
	}
	mod, err := storage.Child(dir, "go.mod")
	if err != nil {
		return nil, err
	}
	w, err := storage.Writer(mod)
	if err != nil {
		return nil, err
	}

	defer w.Close()
	_, err = io.WriteString(w, fmt.Sprintf(`module %s 

go 1.22.6`, name))
	list, _ := storage.ListerForURI(dir)
	return list, err
}

func (g *gui) openProject(dir fyne.ListableURI) {
	g.title.Set(dir.Name())
	g.curDir = dir
	// Reseting the filetree before loading a new project.
	g.fileTree.Set(map[string][]string{}, map[string]fyne.URI{})
	addFilesToTree(dir, g.fileTree, binding.DataTreeRootID)
}

func addFilesToTree(dir fyne.ListableURI, tree binding.URITree, root string) {
	items, _ := dir.List()
	for _, uri := range items {
		nodeID := uri.String()
		tree.Append(root, nodeID, uri)
		isDir, err := storage.CanList(uri)
		if err != nil {
			fmt.Print(err)
			return
		}
		if isDir {
			child, _ := storage.ListerForURI(uri)
			addFilesToTree(child, tree, nodeID)
		}
	}
}
