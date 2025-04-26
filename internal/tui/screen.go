// Package tui provides the terminal user interface components for the Noctui
// application. It includes reusable primitives such as editable text views,
// tabbed layouts, tree views, command-line input handlers, and
// project/workspace selectors.
//
// The package leverages tview and tcell to build a dynamic and
// keyboard-friendly environment for interacting with Insomnium API definitions.
// Each component is designed to be modular, themeable, and fully integrated
// with the application's workflow.
//
// This package is the core of the user experience, enabling navigation,
// editing, and inspection of API projects and HTTP requests directly in the
// terminal.
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

	styleDefault = tc.StyleDefault.Foreground(fgColor).Background(bgColor)

	// CmdFieldStyle is the style used for the command input field.
	CmdFieldStyle = styleDefault

	// BtnStyle is the style used for buttons.
	BtnStyle = styleDefault

	// ActiveBtnStyle is the style used for active buttons.
	ActiveBtnStyle = BtnStyle

	// CmdFieldErrorStyle is the style for command field errors.
	CmdFieldErrorStyle = CmdFieldStyle.Foreground(tc.ColorRed)

	// environmentVarStyle is the style for env variable tags in the URL view.
	environmentVarStyle = styleDefault.Background(tc.ColorMediumPurple).Foreground(tc.ColorBlack)

	// templateVarStyle is the style for template variable tags in the URL view.
	templateVarStyle = styleDefault.Background(tc.ColorSkyblue).Foreground(tc.ColorBlack)

	// urlInputFieldStyle is the style for the URL input field.
	urlInputFieldStyle = styleDefault

	// DefaultMethodStyle is the default style for HTTP methods.
	DefaultMethodStyle = styleDefault

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
	unselectedMehodStyle = styleDefault.Background(tc.ColorWhite).Foreground(tc.ColorBlack)

	// selectedMehodStyle is the style for a selected HTTP method.
	selectedMehodStyle = styleDefault.Background(tc.ColorBlack).Foreground(tc.ColorWhite)

	// tagStart defines the starting delimiter for a tag.
	tagStart = "▐"

	// tagEnd defines the ending delimiter for a tag.
	tagEnd = "▌"

	// requestBodyPlaceholderStyle is the style used for placeholder text in
	// the request body text area.
	requestBodyPlaceholderStyle = CmdFieldStyle.Foreground(tc.ColorDarkGray)

	// requestTabSelectedColor is the foreground color of the selected
	// request tab.
	requestTabSelectedColor = fgColor

	// requestTabUnselectedColor is the foreground color of unselected
	// request tabs.
	requestTabUnselectedColor = tc.ColorDarkSlateGray

	// headerInputStyle defines the style for the header input field.
	headerInputStyle = styleDefault.Background(fgColor).Foreground(bgColor)

	// requestBodyInputStyle defines the style for the body input field.
	requestBodyInputStyle = styleDefault.Background(bgColor).Foreground(fgColor)

	// reqBodyTheme is the theme used for the request body editor.
	reqBodyTheme = "catppuccin-macchiato"

	// statusCodeFG is the foreground color for status code tags.
	statusCodeFG = tc.ColorWhite

	// statusNoResponseColor is the color used when there's no response.
	statusNoResponseColor = tc.ColorBlack

	// status0Color is the color used for status code 0.
	status0Color = tc.ColorDarkRed

	// status1xxColor is the color used for 1xx HTTP status codes.
	status1xxColor = tc.ColorRed

	// status2xxColor is the color used for 2xx HTTP status codes.
	status2xxColor = tc.ColorDarkGreen

	// status3xxColor is the color used for 3xx HTTP status codes.
	status3xxColor = tc.ColorRed

	// status4xxColor is the color used for 4xx HTTP status codes.
	status4xxColor = tc.ColorOrange

	// status5xxColor is the color used for 5xx HTTP status codes.
	status5xxColor = tc.ColorDarkRed

	// statusOthersColor is the color used for unknown HTTP status codes.
	statusOthersColor = tc.ColorDarkRed

	// infoTagStyle is the default style used in informational tags.
	infoTagStyle = styleDefault.Background(tc.ColorDarkSlateGray).Foreground(tc.ColorWhite)

	// responseElapsedTimeTagStyle is the style used for elapsed time tags.
	responseElapsedTimeTagStyle = infoTagStyle

	// responseBytesReadTagStyle is the style used for response size tags.
	responseBytesReadTagStyle = infoTagStyle

	// responseLastCallTagStyle is the style used for the last call tag.
	responseLastCallTagStyle = infoTagStyle

	// responseDropDownTagStyle is the tag style for unselected tags in the
	// response dropdown.
	responseDropDownTagStyle = infoTagStyle

	// responseDrowDownSelectedTagStyle is the tag style for selected tags
	// in the response dropdown.
	responseDrowDownSelectedTagStyle = infoTagStyle.Background(tc.ColorDarkGray)

	// selectedResponseOptionStyle is the style for selected dropdown
	// options.
	selectedResponseOptionStyle = infoTagStyle.Background(tc.ColorLightGray)

	// responseEmptyBarRune is the rune used to draw the scrollbar's empty
	// space.
	responseEmptyBarRune = '█'

	// responseFilledBarRune is the rune used to draw the scrollbar's filled
	// space.
	responseFilledBarRune = '█'

	// responseEmptyBarStyle is the style for the empty scrollbar segment.
	responseEmptyBarStyle = styleDefault.Foreground(tc.ColorGray)

	// responseFilledBarStyle is the style for the filled scrollbar segment.
	responseFilledBarStyle = styleDefault.Foreground(tc.ColorBlue)
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
		s.ResetFocus()
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
		SetResetFocusFunc(s.ResetFocus).
		SetBackgroundColor(bgColor).
		SetForegroundColor(fgColor).
		SetMatrixInputCapture(s.MatrixInputCapture)
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

// MatrixInputCapture will return an input capture function and will accept
// vim motions to control the selected element in the grid.
func (s *Screen) MatrixInputCapture(m *Matrix) InputFn {
	return func(event *tc.EventKey) *tc.EventKey {
		switch {
		case event.Key() == tc.KeyEsc || event.Rune() == 'q':
			s.ResetFocus()
		case event.Key() == tc.KeyLeft || event.Rune() == 'h':
			m.Left()
		case event.Key() == tc.KeyRight || event.Rune() == 'l':
			m.Right()
		case event.Key() == tc.KeyUp || event.Rune() == 'k':
			m.Up()
		case event.Key() == tc.KeyDown || event.Rune() == 'j':
			m.Down()
		default:
			return event
		}
		p := m.GetCurrentPrimitive()
		if p != nil {
			s.app.SetFocus(p)
		}
		return nil
	}
}

// CreateWrkScreen creates the page for managing a workspace.
func (s *Screen) CreateWrkScreen() *WrkPage {
	s.wrk = NewWrkPage(s.inso).
		SetResetFocusFunc(s.ResetFocus).
		SetEnvironmentVarStyle(environmentVarStyle).
		SetTemplateVarStyle(templateVarStyle).
		SetURLInputFieldStyle(urlInputFieldStyle).
		SetMethods(methods).
		SetMethodStyles(unselectedMehodStyle, selectedMehodStyle).
		SetMethodStyleFunc(s.GetDropDownMethodStyle).
		SetTagFunc(s.TextToTag).
		SetBackgroundColor(bgColor).
		SetForegroundColor(fgColor).
		SetSelectedReqTabColor(requestTabSelectedColor).
		SetUnselectedReqTabColor(requestTabUnselectedColor).
		SetHeaderInputStyle(headerInputStyle).
		SetStyleTextFunc(s.StyleText).
		SetRequestBodyInputStyle(requestBodyInputStyle).
		SetRequestBodyTheme(reqBodyTheme).
		SetResponseStyleFn(s.ResponseStyleFn).
		SetResponseElapsedTimeTagStyle(responseElapsedTimeTagStyle).
		SetResponseBytesReadTagStyle(responseBytesReadTagStyle).
		SetResponseLastCallTagStyle(responseLastCallTagStyle).
		SetTagStyle(responseDropDownTagStyle).
		SetSelectedTagStyle(responseDrowDownSelectedTagStyle).
		SetResponseOptionSelectedStyle(selectedResponseOptionStyle).
		SetResponseEmptyBarRune(responseEmptyBarRune).
		SetResponseEmptyBarStyle(responseEmptyBarStyle).
		SetResponseFilledBarRune(responseFilledBarRune).
		SetResponseFilledBarStyle(responseFilledBarStyle)

	return s.wrk
}

// ResponseStyleFn returns the style to display a response, based on its status
// code.
func (s *Screen) ResponseStyleFn(res *insomnium.Response) tc.Style {
	color := s.GetStatusCodeColor(res)
	return styleDefault.Foreground(statusCodeFG).Background(color)
}

// GetStatusCodeColor maps a response status code to a background color.
func (s *Screen) GetStatusCodeColor(res *insomnium.Response) tc.Color {
	switch {
	case res == nil:
		return statusNoResponseColor
	case res.StatusCode == 0:
		return status0Color
	case res.StatusCode < 200:
		return status1xxColor
	case res.StatusCode < 300:
		return status2xxColor
	case res.StatusCode < 400:
		return status3xxColor
	case res.StatusCode < 500:
		return status4xxColor
	case res.StatusCode < 600:
		return status5xxColor
	default:
		return statusOthersColor
	}
}

// ResetFocus sets the focus back to the command input field.
func (s *Screen) ResetFocus() {
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
	return styleDefault.Background(bg).Foreground(fg)
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
		"q":       s.Quit,
		"projs":   s.ProjCommand,
		"wrk":     s.WrkCommand,
		"req":     s.ReqCommand,
		"url":     s.URLCommand,
		"method":  s.MethodCommand,
		"body":    s.MethodBody,
		"headers": s.MethodHeaders,
		"res":     s.MethodRes,
	}
}

// Quit closes the app.
func (s *Screen) Quit(_ []string) {
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

// URLCommand handles the rmethod command, focussing in the url input field
// element. The command is ignored when the current page is not the wrk screen.
func (s *Screen) URLCommand(args []string) {
	pageName, _ := s.pages.GetFrontPage()
	if len(args) == 0 && pageName == "wrk" && s.wrk.GetRequest() != nil {
		focus := s.wrk.GetURLElement()
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

// MethodBody moves the focus to the request body editor element. The command is
// ignored if the current page is not the wrk screen or if no request is
// selected.
func (s *Screen) MethodBody(args []string) {
	pageName, _ := s.pages.GetFrontPage()
	if len(args) == 1 && pageName == "wrk" && s.wrk.GetRequest() != nil {
		s.wrk.MiddlePanel.tabs.OpenTab("Body")
		switch args[0] {
		case "edit":
			s.wrk.MiddlePanel.tabs.BodyTab.SwitchToPage("edit")
			focus := s.wrk.MiddlePanel.tabs.BodyTab.editor
			s.app.SetFocus(focus)
		case "view":
			s.wrk.MiddlePanel.tabs.BodyTab.SwitchToPage("view")
			focus := s.wrk.MiddlePanel.tabs.BodyTab.viewer
			s.app.SetFocus(focus)
		default:
			s.app.Stop()
			fmt.Printf("%q", args[0])
		}
	}
}

// MethodHeaders moves the focus to the request headers editor element. The
// command is ignored if the current page is not the wrk screen or if no request
// is selected.
func (s *Screen) MethodHeaders(args []string) {
	pageName, _ := s.pages.GetFrontPage()
	if len(args) == 0 && pageName == "wrk" && s.wrk.GetRequest() != nil {
		s.wrk.MiddlePanel.tabs.OpenTab("Headers")
		focus := s.wrk.MiddlePanel.tabs.HeadersTab
		s.app.SetFocus(focus)
	}
}

// MethodRes focuses the response dropdown for the current request.
func (s *Screen) MethodRes(_ []string) {
	pageName, _ := s.pages.GetFrontPage()
	req := s.wrk.GetRequest()
	if pageName != "wrk" || len(s.wrk.getResponses(req)) == 0 {
		return
	}
	s.wrk.RightPanel.Header.responseDropDown.SetOpen(true)
	s.app.SetFocus(s.wrk.RightPanel.Header.responseDropDown)
}

// InputFn defines a function type used by tview to capture input, taking a
// pointer to a tcell.EventKey and returning a pointer to a tcell.EventKey.
type InputFn = func(*tc.EventKey) *tc.EventKey
