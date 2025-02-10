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
	maxX, maxY := g.Size()
	// Set mininimal width equal to timer lenght
	inputX0 := maxX / 3
	inputX1 := 2 * maxX / 3
	if inputX1-inputX0 < 42 {
		inputX0 = maxX/2 - 21
		inputX1 = maxX/2 + 21
	}

	if v, err := g.SetView(ti.Name, inputX0, maxY/2-1, inputX1, maxY/2+1, 0); err != nil {
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
