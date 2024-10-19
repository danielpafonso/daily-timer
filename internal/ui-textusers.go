package internal

import (
	"errors"

	"github.com/awesome-gocui/gocui"
)

type Users struct {
	active bool
	name   string
	timer  int
}

type TextUsers struct {
	name   string
	view   *gocui.View
	users  []Users
	x0, y0 int
	x1, y1 int
}

func (tu *TextUsers) Layout(g *gocui.Gui) error {
	if view, err := g.SetView(tu.name, tu.x0, tu.y0, tu.x1, tu.y1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		view.Frame = false
		view.Wrap = false
		view.WriteString("one\ntwo\nthree")
		tu.view = view
	}
	return nil
}
