package editors

import (
	"errors"
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type Editor interface {
	Content() fyne.CanvasObject
	Edited() binding.Bool
	Save() error
}

var mimes = map[string]func(fyne.URI) (Editor, error){
	"text/plain": makeTxt,
}

var extensions = map[string]func(fyne.URI) (Editor, error){
	".go":  makeGo,
	".txt": makeTxt,
	".md":  makeMd,
	".png": makeImg,
}

func makeMd(u fyne.URI) (Editor, error) {
	code, err := makeTxt(u)
	if code == nil || err != nil {
		return nil, err
	}
	txt := code.(*SimpleEditor).content.(*codeEntry)
	txt.TextStyle = fyne.TextStyle{Monospace: true}
	txt.Refresh()

	preview := widget.NewRichTextFromMarkdown(txt.Text)
	dirty := txt.OnChanged
	txt.OnChanged = func(s string) {
		preview.ParseMarkdown(s)
		dirty(s)
	}
	code.(*SimpleEditor).content = container.NewHSplit(txt, container.NewScroll(preview))
	return code, err
}

func makeImg(u fyne.URI) (Editor, error) {
	img := canvas.NewImageFromURI(u)
	img.FillMode = canvas.ImageFillContain
	return &SimpleEditor{content: img}, nil
}

func ForURI(u fyne.URI) (Editor, error) {
	name := strings.ToLower(u.Name())
	var matched func(fyne.URI) (Editor, error)
	for ext, edit := range extensions {
		pos := strings.LastIndex(name, ext)
		if pos == -1 || pos != len(name)-len(ext) {
			continue
		}
		matched = edit
		break
	}
	if matched == nil {
		edit, ok := mimes[u.MimeType()]
		if !ok {

			warning := fmt.Sprintf("Unable to open edtor for %s mime : %s", u.Name(), u.MimeType())
			return nil, errors.New(warning)
		}
		return edit(u)
	}
	return matched(u)
}

func CreateEditor() *widget.Entry {
	entry := widget.NewMultiLineEntry()
	entry.TextStyle.Monospace = true
	entry.MultiLine = true
	entry.Wrapping = fyne.TextWrapWord
	return entry
}

func makeGo(u fyne.URI) (Editor, error) {
	code, err := makeTxt(u)
	if code != nil {
		code.(*SimpleEditor).content.(*codeEntry).TextStyle = fyne.TextStyle{Monospace: true}
	}
	return code, err
}

type SimpleEditor struct {
	content fyne.CanvasObject
	edited  binding.Bool
	save    func() error
}

func (s *SimpleEditor) Content() fyne.CanvasObject {
	return s.content
}

func (s *SimpleEditor) Edited() binding.Bool {
	if s.edited == nil {
		s.edited = binding.NewBool()
	}
	return s.edited
}

func (s *SimpleEditor) Save() error {
	if s.save == nil {
		return nil
	}
	err := s.save()
	if err == nil {
		s.edited.Set(false)
	}
	return err
}
