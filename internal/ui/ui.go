package ui

import (
	"errors"
	"fmt"

	"github.com/awesome-gocui/gocui"

	"daily-timer/internal"
)

type App struct {
	gui         *gocui.Gui
	timer       Timer
	warning     TextWarning
	users       TextUsers
	helpPopup   TextPopup
	hiddenPopup TextPopup
	configPopup TextPopup
	inputTemp   TextInput
}

// NewAppUI initiates new UI
func NewAppUI(config internal.Configurations, stats *[]internal.Stats) *App {
	minimalX := 50
	minimalY := 21

	newApp := App{
		timer: Timer{
			warning:   config.Warning,
			limit:     config.Time,
			running:   false,
			stopwatch: config.Stopwatch,
			minimalX:  minimalX,
			minimalY:  minimalY,
		},
		warning: TextWarning{
			Name:     "warningView",
			minimalX: minimalX,
			minimalY: minimalY,
		},
		users: TextUsers{
			Name:      "users",
			showStats: config.Status.Display,
			users:     stats,
			minimalX:  minimalX,
			minimalY:  minimalY,
		},
		configPopup: TextPopup{
			name:    "configHelp",
			x0:      -10,
			y0:      -4,
			x1:      10,
			y1:      5,
			visible: false,
			text: fmt.Sprintf(`      Configs

 time: %d
 warning: %d
 random: %v
 stopwatch: %v
 add temp: %v
 last dailies: %d`, config.Time, config.Warning, config.Random, config.Stopwatch, config.AddTemp, config.Status.LastDailies),
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

// NextUser selects next user
func (app *App) NextUser(g *gocui.Gui, v *gocui.View) error {
	currentTimer := app.timer.value
	newTimer := app.users.ChangeUser(1, currentTimer, app.timer.running)
	if newTimer >= 0 {
		app.timer.value = newTimer
		app.timer.ResetTimer()
		app.gui.Update(app.timer.Layout)
	}
	return nil
}

// PrevUser selects previous user
func (app *App) PrevUser(g *gocui.Gui, v *gocui.View) error {
	currentTimer := app.timer.value
	newTimer := app.users.ChangeUser(-1, currentTimer, app.timer.running)
	if newTimer >= 0 {
		app.timer.value = newTimer
		app.timer.ResetTimer()
		app.gui.Update(app.timer.Layout)
	}
	return nil
}

// ToggleOnActive toggle timer only on active users
func (app *App) ToggleOnActive(g *gocui.Gui, v *gocui.View) error {
	if (*app.users.users)[app.users.current].Active {
		app.timer.Toggle()
	}
	return nil
}

// OpenTempUser opens text input and sets it as active/current view
func (app *App) OpenTempUser(g *gocui.Gui, v *gocui.View) error {
	app.inputTemp.Visible = true
	g.SetViewOnTop(app.inputTemp.Name)
	g.SetCurrentView(app.inputTemp.Name)
	return nil
}

// AddTempUser closes text input and adds a temp user if input isn't empty
func (app *App) AddTempUser(g *gocui.Gui, v *gocui.View) error {
	newUser := app.inputTemp.Close()
	g.SetCurrentView(app.users.Name)
	if newUser != "" {
		app.users.AddTempUser(newUser)
	}
	return nil
}

// CloseTempUser closes text input
func (app *App) CloseTempUser(g *gocui.Gui, v *gocui.View) error {
	app.inputTemp.Close()
	g.SetCurrentView(app.users.Name)
	return nil
}

// Start starts the application
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
	app.users.x1 = -1
	app.users.y1 = -2

	//  help popup
	app.helpPopup = TextPopup{
		name:    "help",
		x0:      -17,
		y0:      -7,
		x1:      17,
		y1:      8,
		visible: false,
		text: fmt.Sprintf(`          Key  Mapping

 <h> Show/Hide this menu

 <Space> Toggle timer on/off

 <down>/<j> Next user
 <up>/<k>   Previous user

 <s> Show/Hide statistics
 <a> Toggle user active/inactive
 <i> Add temporary user

         version: %s`, version),
	}

	//  hidden help popup
	app.hiddenPopup = TextPopup{
		name:    "hidden",
		x0:      -17,
		y0:      -6,
		x1:      17,
		y1:      6,
		visible: false,
		text: fmt.Sprintf(`          Key  Mapping

 <alt+h> Show/Hide this menu

 <alt+r> Randomize users

 <alt+f> Toggle Flash users

 <alt+c> Show configuration

         version: %s`, version),
	}

	// Set minimal width equal to timer lenght
	inputX0 := maxX / 3
	inputX1 := 2 * maxX / 3
	if inputX1-inputX0 < 42 {
		inputX0 = maxX/2 - 21
		inputX1 = maxX/2 + 21
	}
	// temp user input
	app.inputTemp = TextInput{
		Name:    "tempuser",
		x0:      inputX0,
		y0:      maxY/2 - 1,
		x1:      inputX1,
		y1:      maxY/2 + 1,
		Visible: false,
	}

	// Set Update Manager, order is required
	app.gui.SetManager(
		&app.warning,
		&app.users,
		&app.helpPopup,
		&app.hiddenPopup,
		&app.configPopup,
		&app.timer,
		&app.inputTemp,
	)

	// Set keybindings
	//  exit application
	if err := app.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		// Clear stats so nothing is writting
		*app.users.users = []internal.Stats{}
		return gocui.ErrQuit
	}); err != nil {
		return err
	}
	if err := app.gui.SetKeybinding(app.users.Name, 'q', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		// update current user if timer is running
		if app.timer.running {
			(*app.users.users)[app.users.current].Current = app.timer.value
		}
		return gocui.ErrQuit
	}); err != nil {
		return err
	}

	//  toggle help popup
	if err := app.gui.SetKeybinding(app.users.Name, 'h', gocui.ModNone, app.helpPopup.ToggleVisible); err != nil {
		return err
	}
	// toggle hidden help Popup
	if err := app.gui.SetKeybinding(app.users.Name, 'h', gocui.ModAlt, app.hiddenPopup.ToggleVisible); err != nil {
		return err
	}
	// toggle config Popup
	if err := app.gui.SetKeybinding(app.users.Name, 'c', gocui.ModAlt, app.configPopup.ToggleVisible); err != nil {
		return err
	}

	// Start/stop timer
	if err := app.gui.SetKeybinding(app.users.Name, gocui.KeySpace, gocui.ModNone, app.ToggleOnActive); err != nil {
		return err
	}
	if err := app.gui.SetKeybinding(app.users.Name, gocui.KeyEnter, gocui.ModNone, app.ToggleOnActive); err != nil {
		return err
	}

	//  User list controls:
	//  next
	if err := app.gui.SetKeybinding(app.users.Name, gocui.KeyArrowDown, gocui.ModNone, app.NextUser); err != nil {
		return err
	}
	if err := app.gui.SetKeybinding(app.users.Name, 'j', gocui.ModNone, app.NextUser); err != nil {
		return err
	}
	//  previous
	if err := app.gui.SetKeybinding(app.users.Name, gocui.KeyArrowUp, gocui.ModNone, app.PrevUser); err != nil {
		return err
	}
	if err := app.gui.SetKeybinding(app.users.Name, 'k', gocui.ModNone, app.PrevUser); err != nil {
		return err
	}
	//  show/hide user statistic
	if err := app.gui.SetKeybinding(app.users.Name, 's', gocui.ModNone, app.users.ToggleStats); err != nil {
		return err
	}
	//  randomize user list
	if err := app.gui.SetKeybinding(app.users.Name, 'r', gocui.ModAlt, func(*gocui.Gui, *gocui.View) error {
		app.users.RandomizeOrder()
		return nil
	}); err != nil {
		return err
	}
	//  toggle active/inactive users
	if err := app.gui.SetKeybinding(app.users.Name, 'a', gocui.ModNone, app.users.ToggleActive); err != nil {
		return err
	}
	//  opens window to insert a temp user
	if err := app.gui.SetKeybinding(app.users.Name, 'i', gocui.ModNone, app.OpenTempUser); err != nil {
		return err
	}
	//  toggle user list flash
	if err := app.gui.SetKeybinding(app.users.Name, 'f', gocui.ModAlt, app.users.ToogleFlash); err != nil {
		return nil
	}

	// User input keybindings
	if err := app.gui.SetKeybinding(app.inputTemp.Name, gocui.KeyEnter, gocui.ModNone, app.AddTempUser); err != nil {
		return err
	}
	if err := app.gui.SetKeybinding(app.inputTemp.Name, gocui.KeyEsc, gocui.ModNone, app.CloseTempUser); err != nil {
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
