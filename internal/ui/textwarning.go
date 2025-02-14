package ui

import (
	"errors"
	"fmt"

	"github.com/awesome-gocui/gocui"
)

type TextWarning struct {
	Name     string
	view     *gocui.View
	minimalX int
	minimalY int
}

// Layout creates/updates users widget
func (tw *TextWarning) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	minimalSize := maxX < tw.minimalX || maxY < tw.minimalY

	// warning view
	if view, err := g.SetView(tw.Name, 0, 0, tw.minimalX-1, tw.minimalY-2, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		// view.Frame = false
		view.FrameColor = gocui.ColorRed
		view.SetWritePos(tw.minimalX/2-11, tw.minimalY/2-2)
		view.WriteString("Mininal size required:")
		view.SetWritePos(tw.minimalX/2-3, tw.minimalY/2)
		view.WriteString(fmt.Sprintf("%d x %d", tw.minimalX, tw.minimalY))
		view.Visible = minimalSize
	} else {
		view.Visible = minimalSize
	}
	return nil
}
