package internal

import (
	"errors"
	"fmt"
	"time"

	"github.com/awesome-gocui/gocui"
)

type App struct {
	gui       *gocui.Gui
	timer     Timer
	users     TextUsers
	helpPopup TextPopup
	warning   int
	limit     int
	stats     []Stats
}

func NewAppUI(config Configurations, stats *[]Stats) *App {
	return &App{
		warning: config.Warning,
		limit:   config.Time,
		stats:   *stats,
	}
}

func (app *App) Start() error {
	var err error
	app.gui, err = gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		return err
	}
	defer app.gui.Close()
	// defer write stats

	maxX, maxY := app.gui.Size()
	// Create views
	//  timer view(s)
	app.timer = Timer{
		midX:    maxX/2 - 2,
		topY:    0,
		warning: app.warning,
		limit:   app.limit,
		// timer:   *time.NewTimer(time.Second),
		// running: false,
		nextTick: time.Now(),
		running:  false,
	}

	//  user list
	app.users = TextUsers{
		name:      "users",
		x0:        0,
		y0:        8,
		x1:        maxX - 1,
		y1:        maxY - 2,
		current:   0,
		showStats: true,
		users:     app.stats,
	}
	app.users.calculatePadding()

	//  help popup
	app.helpPopup = TextPopup{
		name:    "help",
		x0:      maxX / 2,
		y0:      maxY / 2,
		x1:      maxX/2 + 20,
		y1:      maxY/2 + 3,
		visible: false,
		text:    fmt.Sprintf("Hello Stella\nwarn: %d, limit: %d", app.warning, app.limit),
	}

	// Set Update Manager, order is required
	app.gui.SetManager(
		&app.users,
		&app.helpPopup,
		&app.timer,
	)

	if view, err := app.gui.SetView("footer", -1, maxY-2, maxX, maxY, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		helpLine := "<q> exit    <h> help menu"
		view.SetWritePos(maxX/2-len(helpLine)/2, 0)
		view.WriteString(helpLine)
	}

	// Set keybindings
	//  exit application
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

	//  toogle help popup
	if err := app.gui.SetKeybinding("", 'h', gocui.ModNone, app.helpPopup.ToogleVisible); err != nil {
		return err
	}

	// Start/stop timer
	if err := app.gui.SetKeybinding("", gocui.KeySpace, gocui.ModNone, app.timer.Toogle); err != nil {
		return err
	}
	if err := app.gui.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, app.timer.Toogle); err != nil {
		return err
	}

	//  User list controls:
	// //  next
	// if err := app.gui.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, app.users.NextUser); err != nil {
	// 	return err
	// }
	// if err := app.gui.SetKeybinding("", 'j', gocui.ModNone, app.users.NextUser); err != nil {
	// 	return err
	// }
	// //  previous
	// if err := app.gui.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, app.users.PrevUser); err != nil {
	// 	return err
	// }
	// if err := app.gui.SetKeybinding("", 'k', gocui.ModNone, app.users.PrevUser); err != nil {
	// 	return err
	// }
	// //  show/hide user statistic
	// if err := app.gui.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, app.users.ToggleStats); err != nil {
	// 	return err
	// }

	// debug/test
	if err := app.gui.SetKeybinding("", 'c', gocui.ModNone, app.helpPopup.Color); err != nil {
		return err
	}
	if err := app.gui.SetKeybinding("", 'i', gocui.ModNone, app.timer.Increment); err != nil {
		return err
	}

	// channl to trigger gui update
	updateChannel := make(chan func(g *gocui.Gui) error)

	go func() {
		// blocking channel read loop
		for {
			layoutFunc := <-updateChannel
			app.gui.Update(layoutFunc)
		}
	}()

	// start internal ticker
	go app.timer.internalTicket(updateChannel)

	// enter UI mainloop
	if err := app.gui.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		return err
	}
	return nil
}
