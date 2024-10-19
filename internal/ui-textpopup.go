package internal

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

func (tp *TextPopup) Layout(g *gocui.Gui) error {
	if v, err := g.SetView(tp.name, tp.x0, tp.y0, tp.x1, tp.y1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.FrameColor = gocui.ColorRed
		v.WriteString(tp.text)
		v.Visible = tp.visible
	} else {
		v.Visible = tp.visible
	}
	return nil
}

func (tp *TextPopup) ToogleVisible(g *gocui.Gui, v *gocui.View) error {
	tp.visible = !tp.visible
	return nil
}

func (tp *TextPopup) Color(g *gocui.Gui, v *gocui.View) error {
	view, _ := g.View(tp.name)
	if view.FgColor == gocui.ColorCyan {
		view.FgColor = gocui.ColorYellow
	} else {
		view.FgColor = gocui.ColorCyan
	}
	return nil
}
