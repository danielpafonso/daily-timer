package internal

import (
	"errors"

	"github.com/awesome-gocui/gocui"
)

type TimerDigit struct {
	name  string
	view  *gocui.View
	value int
}

type Timer struct {
	// views
	minute10 TimerDigit
	minute1  TimerDigit
	second10 TimerDigit
	second1  TimerDigit
	dots     *gocui.View
	// middle coord
	midX, topY int
	// Timer
	value  int
	toogle chan 
}

func (tm *Timer) Increment(g *gocui.Gui, v *gocui.View) error {
	tm.minute1.value += 1
	if tm.minute1.value == 10 {
		tm.minute1.value = 0
	}
	return nil
}

func (tm *Timer) Layout(g *gocui.Gui) error {
	// minute 10s
	diff := 19
	if view, err := g.SetView("minute-10", tm.midX-diff, tm.topY, tm.midX+9-diff, tm.topY+7, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		view.WriteString(Digits[0])
		view.Frame = false
		tm.minute10.view = view
	} else {
		tm.minute10.view.Clear()
		tm.minute10.view.WriteString(Digits[tm.minute10.value])
	}
	// minute 1s
	diff = 9
	if view, err := g.SetView("minute-1", tm.midX-diff, tm.topY, tm.midX+9-diff, tm.topY+7, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		view.WriteString(Digits[0])
		view.Frame = false
		tm.minute1.view = view
	} else {
		tm.minute1.view.Clear()
		tm.minute1.view.WriteString(Digits[tm.minute1.value])
	}

	// dots view
	if view, err := g.SetView("dots", tm.midX, tm.topY, tm.midX+5, tm.topY+7, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		view.WriteString(Dots)
		view.Frame = false
		tm.dots = view
	}
	// second 10s
	diff = 5
	if view, err := g.SetView("second-10", tm.midX+diff, tm.topY, tm.midX+9+diff, tm.topY+7, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		view.WriteString(Digits[0])
		view.Frame = false
		tm.second10.view = view
	} else {
		tm.second10.view.Clear()
		tm.second10.view.WriteString(Digits[tm.second10.value])
	}
	// second 1s
	diff = 15
	if view, err := g.SetView("second-1", tm.midX+diff, tm.topY, tm.midX+9+diff, tm.topY+7, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		view.WriteString(Digits[0])
		view.Frame = false
		tm.second1.view = view
	} else {
		tm.second1.view.Clear()
		tm.second1.view.WriteString(Digits[tm.second1.value])
	}
	return nil
}
