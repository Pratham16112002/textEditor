package editors

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type codeEntry struct {
	widget.Entry
	win    fyne.Window
	save   func() error
	rename func() error
}

func newCodeEntry(w fyne.Window) *codeEntry {
	c := &codeEntry{win: w}
	c.ExtendBaseWidget(c)
	c.MultiLine = true
	return c
}

func (c *codeEntry) TypedShortcut(s fyne.Shortcut) {
	if sh, ok := s.(*desktop.CustomShortcut); ok {
		if sh.KeyName == fyne.KeyS && sh.Modifier == fyne.KeyModifierShortcutDefault {
			c.save()
			return
		}
	}
	c.Entry.TypedShortcut(s)
}

func makeTxt(u fyne.URI) (Editor, error) {
	var code *codeEntry
	save := func() error {
		return saveTxt(u, code.Text)
	}
	code = newCodeEntry(fyne.CurrentApp().Driver().AllWindows()[0])
	r, err := storage.Reader(u)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	data, err := io.ReadAll(r)
	code.SetText(string(data))
	bindingURI := binding.NewURI()
	bindingURI.Set(u)
	edit := &SimpleEditor{content: code, save: save, uri: bindingURI}
	code.OnChanged = func(_ string) {
		edit.Edited().Set(true)
	}
	code.save = edit.Save
	rename := func() error {
		return renameFile(u, code.win, edit)
	}
	code.rename = rename

	return edit, err
}

func saveTxt(u fyne.URI, s string) error {
	w, err := storage.Writer(u)
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = io.WriteString(w, s)
	return err
}

func validate(s string) error {
	if len(strings.TrimSpace(s)) == 0 || s == "" {
		return errors.New("Enter a valid file name")
	}
	if strings.Contains(s, ".") {
		return errors.New("File cannot contain a dot.")
	}
	return nil
}

func renameFile(u fyne.URI, win fyne.Window, editor *SimpleEditor) error {
	newFileNameEntry := widget.NewEntry()
	newFileNameEntry.TextStyle.Bold = true
	newFileNameEntry.Validator = validate
	renameItem := widget.NewFormItem("New Name", container.NewPadded(newFileNameEntry))
	renameItem.Widget.MinSize().AddWidthHeight(500, 100)
	var formItems []*widget.FormItem
	formItems = append(formItems, renameItem)
	dialog.ShowForm("Rename File", "Confirm", "Cancel", formItems, func(ok bool) {
		if ok {
			var err error
			if edited, _ := editor.edited.Get(); edited {
				dialog.ShowError(errors.New("First save the file"), win)
				return
			}
			oldFileName := strings.Split(u.Path(), "/")
			newFileParentPath := strings.Join(
				oldFileName[:len(oldFileName)-1], "/")
			newFilePath := newFileParentPath + "/" + newFileNameEntry.Text + u.Extension()

			newFilePathURI := storage.NewFileURI(newFilePath)
			fmt.Println(newFilePathURI)
			editor.uri.Set(newFilePathURI)
			err = os.Rename(u.Path(), newFilePath)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			return
		}
	}, win)
	return nil
}
