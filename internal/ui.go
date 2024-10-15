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

type App struct {
	gui       *gocui.Gui
	timer     *gocui.Gui
	users     *gocui.Gui
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
	//  user list
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
	app.gui.SetManager(&app.helpPopup)

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

	// enter UI mainloop
	if err := app.gui.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		return err
	}
	return nil
}
