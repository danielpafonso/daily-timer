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

type App struct {
	gui       *gocui.Gui
	timer     Timer
	users     TextUsers
	helpPopup TextPopup
}

func NewAppUI(config Configurations) *App {
	return &App{}
}

func (app *App) Start() error {
	var err error
	app.gui, err = gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		return err
	}
	defer app.gui.Close()

	maxX, maxY := app.gui.Size()
	// Create views
	//  timer view(s)
	app.timer = Timer{
		midX: maxX/2 - 2,
		topY: 0,
	}
	//  user list
	app.users = TextUsers{
		name: "users",
		x0:   0,
		y0:   8,
		x1:   maxX - 1,
		y1:   maxY - 2,
	}
	//  help popup
	app.helpPopup = TextPopup{
		name:    "help",
		x0:      maxX / 2,
		y0:      maxY / 2,
		x1:      maxX/2 + 20,
		y1:      maxY/2 + 3,
		visible: false,
		text:    "Hello Stella",
	}

	// Set Update Manager
	app.gui.SetManager(
		&app.helpPopup,
		&app.users,
		&app.timer,
	)

	//  Static view: help footer
	if view, err := app.gui.SetView("footer", -1, maxY-2, maxX, maxY, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		helpLine := "<q> exit    <h> help menu"
		view.SetWritePos(maxX/2-len(helpLine)/2, 0)
		view.WriteString(helpLine)
	}

	// Set keybindings
	if err := app.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}); err != nil {
		return err
	}
	if err := app.gui.SetKeybinding("", 'q', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}); err != nil {
		return err
	}

	if err := app.gui.SetKeybinding("", 'h', gocui.ModNone, app.helpPopup.ToogleVisible); err != nil {
		return err
	}
	if err := app.gui.SetKeybinding("", 'c', gocui.ModNone, app.helpPopup.Color); err != nil {
		return err
	}
	if err := app.gui.SetKeybinding("", 'i', gocui.ModNone, app.timer.Increment); err != nil {
		return err
	}

	// enter UI mainloop
	if err := app.gui.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		return err
	}
	return nil
}
