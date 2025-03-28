package tui

import (
	"fmt"

	tc "github.com/gdamore/tcell/v2"
	tv "github.com/rivo/tview"
	"github.com/vinicius-lino-figueiredo/insomnium"
)

// EscapeRune precedes a rune that will be escaped
const EscapeRune = '\\'

// Define styles and colors used in the TUI interface.
var (
	// Default font color for the tui
	FontColor = tc.ColorWhite

	// Default font color for the tui
	bgColor = tc.ColorBlack

	// Style for command field
	CmdFieldStyle = tc.StyleDefault.Foreground(FontColor).Background(bgColor)

	// Style for buttons
	BtnStyle = tc.StyleDefault.Foreground(FontColor).Background(bgColor)

	// Style for active buttons
	ActiveBtnStyle = BtnStyle

	// Style for command field errors
	CmdFieldErrorStyle = CmdFieldStyle.Foreground(tc.ColorRed)
)

// Screen represents the TUI screen, holding the main elements of the interface.
type Screen struct {
	// The tview application instance
	app *tv.Application
	// The Insomnium instance for handling storage
	inso *insomnium.Insomnium
	// Input field for commands
	cmd *CmdLine
	// Layout container for organizing the screen
	layout *tv.Flex
	// Pages for switching between different views
	pages *tv.Pages
	// The currently focused primitive in the TUI
	focus tv.Primitive
	// The project/workspace selection page
	proj *ProjPage
}

// NewScreen creates a new *Scree instance with default layout.
func NewScreen(app *tv.Application, inso *insomnium.Insomnium) *Screen {
	s := &Screen{
		app:  app,
		inso: inso,
	}
	app = app.SetInputCapture(s.appInputCapture)
	s.layout, s.pages, s.cmd = s.CreateMainPane()
	return s
}

// appInputCapture handles the Ctrl+C key press by preventing it from closing
// the application.
func (s *Screen) appInputCapture(event *tc.EventKey) *tc.EventKey {
	if event.Key() == tc.KeyCtrlC {
		return nil
	}
	return event
}

// CreateMainPane generates the app main widgets (main frame, input field and
// the place where the different pages wil be).
func (s *Screen) CreateMainPane() (*tv.Flex, *tv.Pages, *CmdLine) {
	s.CreateProjsScreen()
	pages := s.CreatePages()
	cmdField := s.CreateCmdField()

	mainPane := tv.NewFlex().
		SetDirection(tv.FlexRow).
		AddItem(pages, 0, 1, false).
		AddItem(cmdField, 1, 0, false)

	return mainPane, pages, cmdField
}

// CreatePages creates a new instance of the widget that holds the different
// pages of the app.
func (s *Screen) CreatePages() *tv.Pages {
	pages := tv.NewPages()
	for n, v := range s.GetPages() {
		pages.AddPage(n, v, true, false)
	}
	return pages
}

// GetPages returns a map with the id and the instance of the widget to be put
// in the page holder.
func (s *Screen) GetPages() map[string]tv.Primitive {
	return map[string]tv.Primitive{
		"projs": s.proj,
	}
}

// CreateProjsScreen generates a workspace selecting screenfunc
func (s *Screen) CreateProjsScreen() *ProjPage {
	s.proj = NewProjPage(s.app, s.inso).
		SetWrkFunc(s.OpenWorkspace).
		SetNeutralizeFocusFunc(s.NeutralizeFocus)
	return s.proj
}

// NeutralizeFocus sets the focus back to the command input field.
func (s *Screen) NeutralizeFocus() {
	s.app.SetFocus(s.cmd)
}

// OpenWorkspace is called to open the specified workspace by its name.
func (s *Screen) OpenWorkspace(wrk *insomnium.Workspace) {
	// TODO: implement
}

// Run switches the page and starts the app.
func (s *Screen) Run(page string) error {
	s.pages.SwitchToPage(page)
	return s.app.SetRoot(s.layout, true).SetFocus(s.cmd).Run()
}

// CreateCmdField creates an input field that accepts commands.
func (s *Screen) CreateCmdField() *CmdLine {
	return NewCmdLine(s.app).
		SetFieldStyle(CmdFieldStyle).
		SetErrorMessageStyle(CmdFieldErrorStyle).
		SetExecFunc(s.Execute)
}

// Execute finds and runs the function for the given command name or outputs an
// error in the input field.
func (s *Screen) Execute(cmd string, args []string) {
	commands := s.GetCommands()
	fn, ok := commands[cmd]
	if !ok {
		errorMessage := fmt.Sprintf("Command %q not found", cmd)
		s.cmd.DisplayErrorMessage(errorMessage)
		return
	}
	fn(args)
}

// GetCommands returns all the existing commands.
func (s *Screen) GetCommands() map[string]func(args []string) {
	return map[string]func(args []string){
		"q":     s.Quit,
		"projs": s.ProjCommand,
		"wrk":   s.WrkCommand,
	}
}

// Quit closes the app.
func (s *Screen) Quit(args []string) {
	s.app.Stop()
}

// ProjCommand handles the project command, switching to the "projs" page if no
// arguments are provided and setting focus to the current project element.
func (s *Screen) ProjCommand(args []string) {
	if len(args) == 0 {
		s.pages.SwitchToPage("projs")
		focus := s.proj.GetProjElement()
		s.app.SetFocus(focus)
	}
}

// WrkCommand handles the workspace command, switching to the "projs" page if no
// arguments are provided and setting focus to the current workspace element.
func (s *Screen) WrkCommand(args []string) {
	if len(args) == 0 {
		s.pages.SwitchToPage("projs")
		focus := s.proj.GetWrkElement()
		s.app.SetFocus(focus)
	}
}

// InputFn defines a function type used by tview to capture input, taking a
// pointer to a tcell.EventKey and returning a pointer to a tcell.EventKey.
type InputFn = func(*tc.EventKey) *tc.EventKey
