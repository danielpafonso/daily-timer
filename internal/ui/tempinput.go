package ui

import (
	"errors"

	"github.com/awesome-gocui/gocui"
)

type TextInput struct {
	Name    string
	x0, y0  int
	x1, y1  int
	Visible bool
	view    *gocui.View
}

// Layout creates/updates users widget
func (ti *TextInput) Layout(g *gocui.Gui) error {
	if v, err := g.SetView(ti.Name, ti.x0, ti.y0, ti.x1, ti.y1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = "Temp User:"
		v.Editable = true
		v.Visible = ti.Visible
		ti.view = v
	} else {
		v.Visible = ti.Visible
	}
	return nil
}

// Close helper function that extracts input and hides widget
func (ti *TextInput) Close() string {
	newUser := ti.view.Buffer()
	ti.view.Clear()
	ti.Visible = false
	return newUser
}
