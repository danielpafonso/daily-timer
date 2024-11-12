package internal

import (
	"errors"
	"fmt"

	"github.com/awesome-gocui/gocui"
)

type App struct {
	gui       *gocui.Gui
	timer     Timer
	users     TextUsers
	helpPopup TextPopup
}

func NewAppUI(config Configurations, stats *[]Stats) *App {
	newApp := App{
		timer: Timer{
			warning:   config.Warning,
			limit:     config.Time,
			running:   false,
			stopwatch: config.Stopwatch,
		},
		users: TextUsers{
			showStats: config.Status.Display,
			users:     *stats,
		},
	}
	// calculate user padding
	newApp.users.calculatePadding()
	// randomize order if desired
	if config.Random {
		newApp.users.RandomizeOrder()
	}

	return &newApp
}

func (app *App) NextUser(g *gocui.Gui, v *gocui.View) error {
	currentTimer := app.timer.value
	newTimer := app.users.ChangeUser(1, currentTimer, app.timer.running)
	app.timer.value = newTimer
	app.timer.ResetTimer()
	app.gui.Update(app.timer.Layout)

	return nil
}

func (app *App) PrevUser(g *gocui.Gui, v *gocui.View) error {
	currentTimer := app.timer.value
	newTimer := app.users.ChangeUser(-1, currentTimer, app.timer.running)
	app.timer.value = newTimer
	app.timer.ResetTimer()
	app.gui.Update(app.timer.Layout)

	return nil
}

func (app *App) ToogleOnActive(g *gocui.Gui, v *gocui.View) error {
	if app.users.users[app.users.current].Active {
		app.timer.Toogle()
	}
	return nil
}

func (app *App) Start(version string) error {
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
	app.timer.midX = maxX/2 - 2
	app.timer.topY = 0
	// run display to update if in stopwatch
	if !app.timer.stopwatch {
		app.timer.displayTimer()
	}

	//  user list
	app.users.x0 = 0
	app.users.y0 = 8
	app.users.x1 = maxX - 1
	app.users.y1 = maxY - 2

	//  help popup
	app.helpPopup = TextPopup{
		name:    "help",
		x0:      maxX/2 - 17,
		y0:      maxY/2 - 7,
		x1:      maxX/2 + 17,
		y1:      maxY/2 + 7,
		visible: false,
		text: fmt.Sprintf(`          Key  Mapping

 <h> Show/Hide this menu

 <Space> Toggle timer on/off

 <down>/<j> Next user
 <up>/<k>   Previous user

 <s> Show/Hide statistics
 <a> Toggle user active/inactive

         version: %s`, version),
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
		// update current user if timer is running
		if app.timer.running {
			app.users.users[app.users.current].Current = app.timer.value
		}
		return gocui.ErrQuit
	}); err != nil {
		return err
	}

	//  toogle help popup
	if err := app.gui.SetKeybinding("", 'h', gocui.ModNone, app.helpPopup.ToogleVisible); err != nil {
		return err
	}

	// Start/stop timer
	if err := app.gui.SetKeybinding("", gocui.KeySpace, gocui.ModNone, app.ToogleOnActive); err != nil {
		return err
	}
	if err := app.gui.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, app.ToogleOnActive); err != nil {
		return err
	}

	//  User list controls:
	//  next
	if err := app.gui.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, app.NextUser); err != nil {
		return err
	}
	if err := app.gui.SetKeybinding("", 'j', gocui.ModNone, app.NextUser); err != nil {
		return err
	}
	//  previous
	if err := app.gui.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, app.PrevUser); err != nil {
		return err
	}
	if err := app.gui.SetKeybinding("", 'k', gocui.ModNone, app.PrevUser); err != nil {
		return err
	}
	//  show/hide user statistic
	if err := app.gui.SetKeybinding("", 's', gocui.ModNone, app.users.ToggleStats); err != nil {
		return err
	}
	// randomize user list
	if err := app.gui.SetKeybinding("", 'r', gocui.ModAlt, func(*gocui.Gui, *gocui.View) error {
		app.users.RandomizeOrder()
		return nil
	}); err != nil {
		return err
	}
	// toogle active/inactive users
	if err := app.gui.SetKeybinding("", 'a', gocui.ModNone, app.users.ToogleActive); err != nil {
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
