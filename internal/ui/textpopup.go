package ui

import (
	"errors"

	"github.com/awesome-gocui/gocui"
)

type TextPopup struct {
	name    string
	x0, y0  int
	x1, y1  int
	visible bool
	text    string
}

// Layout creates/updates help popup widget
func (tp *TextPopup) Layout(g *gocui.Gui) error {
	if v, err := g.SetView(tp.name, tp.x0, tp.y0, tp.x1, tp.y1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		g.SetViewOnTop(tp.name)
		v.WriteString(tp.text)
		v.Visible = tp.visible
	} else {
		v.Visible = tp.visible
	}
	return nil
}

// ToggleVisible shows/hides the help widget
func (tp *TextPopup) ToggleVisible(g *gocui.Gui, v *gocui.View) error {
	tp.visible = !tp.visible
	return nil
}
