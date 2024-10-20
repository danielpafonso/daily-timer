package internal

import (
	"errors"
	"time"

	"github.com/awesome-gocui/gocui"
)

var (
	// Colors
	colorStop    = gocui.ColorGreen
	colorRun     = gocui.ColorWhite
	colorWarning = gocui.ColorYellow
	colorOver    = gocui.ColorRed
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
	// color limits
	warning int
	limit   int
	// Timer
	value    int
	nextTick time.Time
	running  bool
}

func (tm *Timer) setColor(color gocui.Attribute) {
	tm.minute10.view.FgColor = color
	tm.minute1.view.FgColor = color
	tm.dots.FgColor = color
	tm.second10.view.FgColor = color
	tm.second1.view.FgColor = color
}

func (tm *Timer) displayTimer() {
	minutes := tm.value / 60
	seconds := tm.value % 60

	tm.minute10.value = minutes / 10
	tm.minute1.value = minutes % 10

	tm.second10.value = seconds / 10
	tm.second1.value = seconds % 10
}

func (tm *Timer) Increment(g *gocui.Gui, v *gocui.View) error {
	tm.minute1.value += 1
	if tm.minute1.value == 10 {
		tm.minute1.value = 0
	}
	return nil
}

func (tm *Timer) Toogle(g *gocui.Gui, v *gocui.View) error {
	if tm.running {
		tm.running = false
		tm.setColor(colorStop)
	} else {
		tm.nextTick = time.Now().Add(time.Second)
		tm.running = true
		if tm.value >= tm.warning {
			if tm.value >= tm.limit {
				tm.setColor(colorOver)
			} else {
				tm.setColor(colorWarning)
			}
		} else {
			tm.setColor(colorRun)
		}
	}
	return nil
}

func (tm *Timer) internalTicket(updateCh chan<- func(g *gocui.Gui) error) {
	for {
		if tm.running && time.Now().After(tm.nextTick) {
			tm.nextTick = tm.nextTick.Add(time.Second)
			tm.value += 1
			// not the most efficiency
			if tm.value >= tm.warning {
				if tm.value >= tm.limit {
					tm.setColor(colorOver)
				} else {
					tm.setColor(colorWarning)
				}
			}
			tm.displayTimer()

			// signal for gui update
			updateCh <- tm.Layout
		}
		time.Sleep(250 * time.Millisecond)
	}
}

func (tm *Timer) Layout(g *gocui.Gui) error {
	// minute 10s
	diff := 19
	if view, err := g.SetView("minute-10", tm.midX-diff, tm.topY, tm.midX+9-diff, tm.topY+7, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		view.FgColor = colorStop
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
		view.FgColor = colorStop
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
		view.FgColor = colorStop
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
		view.FgColor = colorStop
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
		view.FgColor = colorStop
		view.WriteString(Digits[0])
		view.Frame = false
		tm.second1.view = view
	} else {
		tm.second1.view.Clear()
		tm.second1.view.WriteString(Digits[tm.second1.value])
	}
	return nil
}
