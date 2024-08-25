package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

const sideWidth = 220

type customLayout struct {
	top, left, right, content fyne.CanvasObject
	seperators                [3]fyne.CanvasObject
}

func CustomLayout(top, left, right, content fyne.CanvasObject, seperators [3]fyne.CanvasObject) fyne.Layout {
	return &customLayout{top: top, left: left, right: right, content: content, seperators: seperators}
}

func (l *customLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	topHeight := l.top.MinSize().Height
	l.top.Resize(fyne.NewSize(size.Width, topHeight))

	l.left.Move(fyne.NewPos(0, topHeight))
	l.left.Resize(fyne.NewSize(sideWidth, size.Height-topHeight))

	l.right.Move(fyne.NewPos(size.Width-sideWidth, topHeight))
	l.right.Resize(fyne.NewSize(sideWidth, size.Height-topHeight))

	l.content.Move(fyne.NewPos(sideWidth, topHeight))
	l.content.Resize(fyne.NewSize(size.Width-2*sideWidth, size.Height-topHeight))

	seperatorThickness := theme.SeparatorThicknessSize()

	l.seperators[0].Move(fyne.NewPos(0, topHeight))
	l.seperators[0].Resize(fyne.NewSize(size.Width, seperatorThickness))

	l.seperators[1].Move(fyne.NewPos(sideWidth, topHeight))
	l.seperators[1].Resize(fyne.NewSize(seperatorThickness, size.Height-topHeight))

	l.seperators[2].Move(fyne.NewPos(size.Width-sideWidth, topHeight))
	l.seperators[2].Resize(fyne.NewSize(seperatorThickness, size.Height-topHeight))

}

func (l *customLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	borders := fyne.NewSize(sideWidth*2, l.top.MinSize().Height)

	return borders.AddWidthHeight(100, 100)
}
