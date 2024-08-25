//go:generate fyne bundle -o bundled.go assets

package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type customTheme struct {
	fyne.Theme
}

func NewCustomTheme() fyne.Theme {
	return &customTheme{Theme: theme.DefaultTheme()}
}

func (t *customTheme) Color(name fyne.ThemeColorName, varient fyne.ThemeVariant) color.Color {
	return t.Theme.Color(name, theme.VariantLight)
}

func (t *customTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameText {
		return 12.0
	}
	return t.Theme.Size(name)
}
