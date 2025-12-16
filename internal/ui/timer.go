package ui

import (
	"errors"
	"time"

	"github.com/awesome-gocui/gocui"

	"daily-timer/internal"
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
	value     int
	nextTick  time.Time
	running   bool
	stopwatch bool
	minimalX  int
	minimalY  int
}

func (tm *Timer) setColor(color gocui.Attribute) {
	tm.minute10.view.FgColor = color
	tm.minute1.view.FgColor = color
	tm.dots.FgColor = color
	tm.second10.view.FgColor = color
	tm.second1.view.FgColor = color
}

func (tm *Timer) displayTimer() {
	value := tm.value
	if !tm.stopwatch {
		value = tm.limit - value
		if value < 0 {
			value = value * -1
		}
	}
	minutes := value / 60
	seconds := value % 60

	tm.minute10.value = minutes / 10
	tm.minute1.value = minutes % 10

	tm.second10.value = seconds / 10
	tm.second1.value = seconds % 10
}

// ResetTimer writes current value to display, set display color equals to state and resets next ticker
func (tm *Timer) ResetTimer() {
	// set display numbers
	tm.displayTimer()
	// set colors
	if tm.running {
		if tm.value >= tm.warning {
			if tm.value >= tm.limit {
				tm.setColor(colorOver)
			} else {
				tm.setColor(colorWarning)
			}
		} else {
			tm.setColor(colorRun)
		}
	} else {
		tm.setColor(colorStop)
	}
	// reset next ticker
	tm.nextTick = time.Now().Add(time.Second)
}

// Toggle starts/stops the timer
func (tm *Timer) Toggle() {
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
}

// Reset the timer
func (tm *Timer) Reset() {
	tm.value = 0
	tm.displayTimer()
	if tm.running {
		tm.nextTick = time.Now().Add(time.Second)
	}
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
			} else {
				tm.setColor(colorRun)
			}
			tm.displayTimer()

			// signal for gui update
			updateCh <- tm.Layout
		}
		time.Sleep(250 * time.Millisecond)
	}
}

// Layout creates/updates timer widget
func (tm *Timer) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	tm.midX = maxX/2 - 2
	// minute 10s
	diff := 19
	if view, err := g.SetView("minute-10", tm.midX-diff, tm.topY, tm.midX+9-diff, tm.topY+7, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		view.FgColor = colorStop
		view.WriteString(internal.Digits[tm.minute10.value])
		view.Frame = false
		g.SetViewOnBottom("minute-10")
		tm.minute10.view = view
	} else {
		tm.minute10.view.Clear()
		tm.minute10.view.WriteString(internal.Digits[tm.minute10.value])
	}
	// minute 1s
	diff = 9
	if view, err := g.SetView("minute-1", tm.midX-diff, tm.topY, tm.midX+9-diff, tm.topY+7, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		view.FgColor = colorStop
		view.WriteString(internal.Digits[tm.minute1.value])
		view.Frame = false
		g.SetViewOnBottom("minute-1")
		tm.minute1.view = view
	} else {
		tm.minute1.view.Clear()
		tm.minute1.view.WriteString(internal.Digits[tm.minute1.value])
	}

	// dots view
	if view, err := g.SetView("dots", tm.midX, tm.topY, tm.midX+5, tm.topY+7, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		view.FgColor = colorStop
		view.WriteString(internal.Dots)
		view.Frame = false
		g.SetViewOnBottom("dots")
		tm.dots = view
	}
	// second 10s
	diff = 5
	if view, err := g.SetView("second-10", tm.midX+diff, tm.topY, tm.midX+9+diff, tm.topY+7, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		view.FgColor = colorStop
		view.WriteString(internal.Digits[tm.second10.value])
		view.Frame = false
		tm.second10.view = view
		g.SetViewOnBottom("second-10")
	} else {
		tm.second10.view.Clear()
		tm.second10.view.WriteString(internal.Digits[tm.second10.value])
	}
	// second 1s
	diff = 15
	if view, err := g.SetView("second-1", tm.midX+diff, tm.topY, tm.midX+9+diff, tm.topY+7, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		view.FgColor = colorStop
		view.WriteString(internal.Digits[tm.second1.value])
		view.Frame = false
		tm.second1.view = view
		g.SetViewOnBottom("second-1")
	} else {
		tm.second1.view.Clear()
		tm.second1.view.WriteString(internal.Digits[tm.second1.value])
	}

	minimalSize := maxX >= tm.minimalX && maxY >= tm.minimalY
	tm.minute10.view.Visible = minimalSize
	tm.minute1.view.Visible = minimalSize
	tm.second10.view.Visible = minimalSize
	tm.second1.view.Visible = minimalSize
	tm.dots.Visible = minimalSize
	return nil
}
