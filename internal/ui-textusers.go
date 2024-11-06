package internal

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/awesome-gocui/gocui"
)

type TextUsers struct {
	name      string
	view      *gocui.View
	users     []Stats
	padding   int
	current   int
	showStats bool
	x0, y0    int
	x1, y1    int
}

func (tu *TextUsers) RandomizeOrder() {
	rand.Shuffle(len(tu.users), func(i, j int) { tu.users[i], tu.users[j] = tu.users[j], tu.users[i] })
}

func (tu *TextUsers) UserLine(idx int) string {
	if idx < 0 || idx >= len(tu.users) {
		return ""
	}
	// Cursor/active timer
	prefix := "  "
	if tu.current == idx {
		prefix = " >"
	}
	// User name
	current := "--:--"
	if tu.users[idx].Active {
		current = fmt.Sprintf("%02d:%02d", tu.users[idx].Current/60, tu.users[idx].Current%60)
	}
	user := fmt.Sprintf("%-*s  %s", tu.padding, tu.users[idx].Name, current)
	// stats
	stats := ""
	if tu.showStats {
		maxString := fmt.Sprintf("%02d:%02d", tu.users[idx].Max/60, tu.users[idx].Max%60)
		avgString := fmt.Sprintf("%02d:%02d", tu.users[idx].Average/60, tu.users[idx].Average%60)
		stats = fmt.Sprintf("max: %s, avg: %s", maxString, avgString)

	}
	return fmt.Sprintf("%s %s    %s [%v]", prefix, user, stats, tu.users[idx].Active)
}

func (tu *TextUsers) calculatePadding() {
	size := 0
	for _, person := range tu.users {
		size = max(size, len(person.Name))
	}
	tu.padding = size
}

func (tu *TextUsers) ChangeUser(delta int, timer int) int {
	old := tu.current
	tu.current += delta
	// no loop change
	if tu.current < 0 {
		tu.current = 0
	}
	if tu.current == len(tu.users) {
		tu.current = len(tu.users) - 1
	}

	// update old
	tu.users[old].Current = timer
	oldLine := tu.UserLine(old)
	_ = oldLine

	// update new
	newLine := tu.UserLine(tu.current)
	_ = newLine

	return tu.users[tu.current].Current
}

func (tu *TextUsers) ToggleStats(g *gocui.Gui, v *gocui.View) error {
	tu.showStats = !tu.showStats

	return nil
}

func (tu *TextUsers) ToogleActive(g *gocui.Gui, v *gocui.View) error {
	tu.users[tu.current].Active = !tu.users[tu.current].Active

	return nil
}

// Layout creates/updates users widget
func (tu *TextUsers) Layout(g *gocui.Gui) error {
	if view, err := g.SetView("users", tu.x0, tu.y0, tu.x1, tu.y1, 0); err != nil {
		// Create view
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		view.Frame = false
		view.Wrap = false
		// view.WriteString("one\ntwo\nthree")
		lines := make([]string, len(tu.users))
		for idx := range tu.users {
			lines[idx] = tu.UserLine(idx)
		}
		view.WriteString(strings.Join(lines, "\n"))
		tu.view = view
	} else {
		// update view
		lines := make([]string, len(tu.users))
		for idx := range tu.users {
			lines[idx] = tu.UserLine(idx)
		}

		tu.view.Clear()
		tu.view.WriteString(strings.Join(lines, "\n"))
	}
	return nil
}
