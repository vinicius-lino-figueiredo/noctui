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
	// fgColor is the default font color for the TUI.
	fgColor = tc.ColorWhite

	// bgColor is the default background color for the TUI.
	bgColor = tc.ColorBlack

	// CmdFieldStyle is the style used for the command input field.
	CmdFieldStyle = tc.StyleDefault.Foreground(fgColor).Background(bgColor)

	// BtnStyle is the style used for buttons.
	BtnStyle = tc.StyleDefault.Foreground(fgColor).Background(bgColor)

	// ActiveBtnStyle is the style used for active buttons.
	ActiveBtnStyle = BtnStyle

	// CmdFieldErrorStyle is the style for command field errors.
	CmdFieldErrorStyle = CmdFieldStyle.Foreground(tc.ColorRed)

	// environmentVarStyle is the style for env variable tags in the URL view.
	environmentVarStyle = tc.StyleDefault.Background(tc.ColorMediumPurple).Foreground(tc.ColorBlack)

	// templateVarStyle is the style for template variable tags in the URL view.
	templateVarStyle = tc.StyleDefault.Background(tc.ColorSkyblue).Foreground(tc.ColorBlack)

	// urlInputFieldStyle is the style for the URL input field.
	urlInputFieldStyle = tc.StyleDefault.Background(tc.ColorBlack).Foreground(tc.ColorWhite)

	// DefaultMethodStyle is the default style for HTTP methods.
	DefaultMethodStyle = tc.StyleDefault.Background(tc.ColorBlack)

	// DefaultMethodColor is the default font color for HTTP methods.
	DefaultMethodColor = tc.ColorWhite

	// MethodGetColor is the color for the GET method.
	MethodGetColor = tc.ColorPurple

	// MethodPostColor is the color for the POST method.
	MethodPostColor = tc.ColorGreen

	// MethodPutColor is the color for the PUT method.
	MethodPutColor = tc.ColorOrange

	// MethodPatchColor is the color for the PATCH method.
	MethodPatchColor = tc.ColorYellow

	// MethodDeleteColor is the color for the DELETE method.
	MethodDeleteColor = tc.ColorRed

	// MethodOptionsColor is the color for the OPTIONS method.
	MethodOptionsColor = tc.ColorBlue

	// MethodHeadColor is the color for the HEAD method.
	MethodHeadColor = tc.ColorBlue

	// methods is the list of supported HTTP methods.
	methods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

	// unselectedMehodStyle is the style for an unselected HTTP method.
	unselectedMehodStyle = tc.StyleDefault.Background(tc.ColorWhite).Foreground(tc.ColorBlack)

	// selectedMehodStyle is the style for a selected HTTP method.
	selectedMehodStyle = tc.StyleDefault.Background(tc.ColorBlack).Foreground(tc.ColorWhite)

	// tagStart defines the starting delimiter for a tag.
	tagStart = "▐"

	// tagEnd defines the ending delimiter for a tag.
	tagEnd = "▌"
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
	// The workspace viewer page
	wrk *WrkPage
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
		s.NeutralizeFocus()
		return nil
	}
	return event
}

// CreateMainPane generates the app main widgets (main frame, input field and
// the place where the different pages wil be).
func (s *Screen) CreateMainPane() (*tv.Flex, *tv.Pages, *CmdLine) {
	s.CreateProjsScreen()
	s.CreateWrkScreen()
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
		"wrk":   s.wrk,
	}
}

// CreateProjsScreen generates a workspace selecting screenfunc
func (s *Screen) CreateProjsScreen() *ProjPage {
	s.proj = NewProjPage(s.app, s.inso).
		SetWrkFunc(s.OpenWorkspace).
		SetNeutralizeFocusFunc(s.NeutralizeFocus).
		SetBackgroundColor(bgColor).
		SetForegroundColor(fgColor)
	return s.proj
}

// GetDropDownMethodStyle receives a HTTP method name and returns a style that
// should be used to decorate that method name in the screen.
func (s *Screen) GetDropDownMethodStyle(method string) tc.Style {
	methodColor := s.GetMethodColor(method)
	return DefaultMethodStyle.Foreground(methodColor)
}

// GetMethodColor receives a HTTP method name and returns the color that should
// be used to decorate that method.
func (s *Screen) GetMethodColor(method string) tc.Color {
	switch method {
	case "GET":
		return MethodGetColor
	case "POST":
		return MethodPostColor
	case "PUT":
		return MethodPutColor
	case "PATCH":
		return MethodPatchColor
	case "DELETE":
		return MethodDeleteColor
	case "OPTIONS":
		return MethodOptionsColor
	case "HEAD":
		return MethodHeadColor
	default:
		return DefaultMethodColor
	}
}

// CreateWrkScreen creates the page for managing a workspace.
func (s *Screen) CreateWrkScreen() *WrkPage {
	s.wrk = NewWrkPage(s.inso).
		SetFocusFunc(s.Focus).
		SetNeutralizeFocusFunc(s.NeutralizeFocus).
		SetEnvironmentVarStyle(environmentVarStyle).
		SetTemplateVarStyle(templateVarStyle).
		SetUrlInputFieldStyle(urlInputFieldStyle).
		SetMethods(methods).
		SetMethodStyles(unselectedMehodStyle, selectedMehodStyle).
		SetMethodStyleFunc(s.GetDropDownMethodStyle).
		SetTagFunc(s.TextToTag).
		SetBackgroundColor(bgColor).
		SetForegroundColor(fgColor)
	return s.wrk
}

// Focus sets focus in the given element.
func (s *Screen) Focus(element tv.Primitive) {
	s.app.SetFocus(element)
}

// NeutralizeFocus sets the focus back to the command input field.
func (s *Screen) NeutralizeFocus() {
	s.app.SetFocus(s.cmd)
}

// OpenWorkspace is called to open the specified workspace by its name.
func (s *Screen) OpenWorkspace(wrk *insomnium.Workspace) {
	s.pages.SwitchToPage("wrk")
	s.wrk.SetWrk(wrk)
}

// TextToTag accepts a text, a base style and a tag style and uses those values
// to create a dynamically colored tag.
func (s *Screen) TextToTag(txt string, baseStyle, tagStyle tc.Style) string {
	halfBlockStyle := s.GetHalfBlockStyle(tagStyle, baseStyle)
	left := s.StyleText(tagStart, halfBlockStyle)
	middle := s.StyleText(txt, tagStyle)
	right := s.StyleText(tagEnd, halfBlockStyle)
	fixColor := s.StyleText("", baseStyle)
	return left + middle + right + fixColor
}

// GetHalfBlockStyle creates the style for the half blocks that complete
// the start and end of a tag. It uses a base style to match the widget's
// background and applies the tag style to set the foreground, making it
// look like a tag completion.
func (s *Screen) GetHalfBlockStyle(tagStyle, noStyle tc.Style) tc.Style {
	_, fg, _ := tagStyle.Decompose()
	_, bg, _ := noStyle.Decompose()
	return tc.StyleDefault.Background(bg).Foreground(fg)
}

// StyleText returns the styled version of the given text.
func (s *Screen) StyleText(text string, style tc.Style) string {
	fg, bg, _ := style.Decompose()
	text = fmt.Sprintf("[%s:%s]%s", fg, bg, text)
	return text
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
		"q":      s.Quit,
		"projs":  s.ProjCommand,
		"wrk":    s.WrkCommand,
		"req":    s.ReqCommand,
		"url":    s.UrlCommand,
		"method": s.MethodCommand,
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

// ReqCommand handles the request command, focussing in the request tree view
// element. The command is ignored when the current page is not the wrk screen.
func (s *Screen) ReqCommand(args []string) {
	pageName, _ := s.pages.GetFrontPage()
	if len(args) == 0 && pageName == "wrk" {
		focus := s.wrk.GetReqTreeElement()
		s.app.SetFocus(focus)
	}
}

// MethodCommand handles the rmethod command, focussing in the url input field
// element. The command is ignored when the current page is not the wrk screen.
func (s *Screen) UrlCommand(args []string) {
	pageName, _ := s.pages.GetFrontPage()
	if len(args) == 0 && pageName == "wrk" && s.wrk.GetRequest() != nil {
		focus := s.wrk.GetUrlElement()
		s.app.SetFocus(focus)
	}
}

// MethodCommand handles the rmethod command, focussing in the url input field
// element. The command is ignored when the current page is not the wrk screen.
func (s *Screen) MethodCommand(args []string) {
	pageName, _ := s.pages.GetFrontPage()
	if len(args) == 0 && pageName == "wrk" && s.wrk.GetRequest() != nil {
		focus := s.wrk.GetMethodElement()
		s.app.SetFocus(focus)
	}
}

// InputFn defines a function type used by tview to capture input, taking a
// pointer to a tcell.EventKey and returning a pointer to a tcell.EventKey.
type InputFn = func(*tc.EventKey) *tc.EventKey
