package main

import (
	"bytes"
	"fmt"
	"strings"

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
	CmdFieldErrorStyle = tc.StyleDefault.Foreground(tc.ColorRed)
)

// Screen represents the TUI screen, holding the main elements of the interface.
type Screen struct {
	app    *tv.Application // The tview application instance
	cmd    *tv.InputField  // Input field for commands
	layout *tv.Flex        // Layout container for organizing the screen
	pages  *tv.Pages       // Pages for switching between different views
	focus  tv.Primitive    // The currently focused primitive in the TUI
}

// NewScreen creates a new *Scree instance with default layout.
func NewScreen(app *tv.Application) *Screen {
	s := &Screen{
		app: app,
	}
	s.layout, s.pages, s.cmd = s.CreateMainPane()
	return s
}

// CreateMainPane generates the app main widgets (main frame, input field and
// the place where the different pages wil be).
func (s *Screen) CreateMainPane() (*tv.Flex, *tv.Pages, *tv.InputField) {
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
		"projs": s.CreateProjsScreen(),
	}
}

// CreateProjsScreen generates a workspace selecting screenfunc
func (s *Screen) CreateProjsScreen() *tv.Flex {
	projMatrix := NewMatrix()
	wrkMatrix := NewMatrix()
	projMatrix.
		SetWidth(1).
		SetHeight(6).
		SetBorder(true).
		SetFocusFunc(s.MatrixOnFocus(projMatrix)).
		SetInputCapture(s.ProjMatrixInputCapture(projMatrix, wrkMatrix))
	wrkMatrix.
		SetWidth(5).
		SetHeight(5).
		SetBorder(true).
		SetFocusFunc(s.MatrixOnFocus(wrkMatrix)).
		SetInputCapture(s.WrkMatrixInputCapture(wrkMatrix))

	s.focus = projMatrix
	s.PopulateProjects(projMatrix, wrkMatrix)
	s.PopulateWorkspaces("", wrkMatrix)

	workspacePane := tv.NewFlex().
		AddItem(projMatrix, 0, 1, false).
		AddItem(wrkMatrix, 0, 5, false)
	return workspacePane
}

// MatrixOnFocus returns a func to be used as a callback for a Matrix focus.
func (s *Screen) MatrixOnFocus(m *Matrix) func() {
	return func() {
		if len(m.itms) != 0 {
			s.app.SetFocus(m.itms[0])
		}
	}
}

// ProjMatrixInputCapture returns a input capture func that will wait for a
// Escape or 'q' key press and will reload all the workspaces. It will ignore
// any other input and pass it to the default matrix input capture function.
func (s *Screen) ProjMatrixInputCapture(projMatrix, wrkMatrix *Matrix) InputFn {
	return func(event *tc.EventKey) *tc.EventKey {
		if event.Key() == tc.KeyEsc || event.Rune() == 'q' {
			s.PopulateWorkspaces("", wrkMatrix)
			s.app.SetFocus(s.cmd)
			return nil
		}
		fn := s.MatrixInputCapture(projMatrix)
		return fn(event)
	}
}

// WrkMatrixInputCapture returns a input capture func that will wait for a
// Escape or 'q' key press and will reload all the workspaces. It will ignore
// any other input and pass it to the default matrix input capture function.
func (s *Screen) WrkMatrixInputCapture(wrkMatrix *Matrix) InputFn {
	return func(event *tc.EventKey) *tc.EventKey {
		if event.Key() == tc.KeyEsc || event.Rune() == 'q' {
			s.PopulateWorkspaces("", wrkMatrix)
			s.app.SetFocus(s.cmd)
			return nil
		}
		fn := s.MatrixInputCapture(wrkMatrix)
		return fn(event)
	}
}

// MatrixInputCapture will return an input capture function and will accept
// vim motions to control the selected element in the grid.
func (s *Screen) MatrixInputCapture(m *Matrix) InputFn {
	return func(event *tc.EventKey) *tc.EventKey {
		switch {
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

// PopulateProjects loads the project grid with all the projects found in the
// *insomnium.Insomnium global instance.
func (s *Screen) PopulateProjects(projMatrix, wrkMatrix *Matrix) {
	for _, proj := range inso.Projects {
		projectButton := s.CreateProjBtn(proj, wrkMatrix)
		projectButton.SetStyle(BtnStyle).
			SetActivatedStyle(ActiveBtnStyle).
			SetSelectedFunc(s.SelectProjBtn(proj.ID, wrkMatrix)).
			SetBorder(true).
			SetFocusFunc(s.FocusProjBtn(proj.ID, wrkMatrix))
		projMatrix.itms = append(projMatrix.itms, projectButton)
	}

	projMatrix.Refresh()
}

// PopulateWorkspaces uses an id to populate the workspace matrix with
// workspaces whose ParentID equals the given id. When "" is passed,
// the matrix is populated with all workspaces.
func (s *Screen) PopulateWorkspaces(id string, m *Matrix) {
	m.itms = []tv.Primitive{}
	for _, wrk := range inso.Workspaces {
		if id != "" && id != wrk.ParentID {
			continue
		}
		btn := tv.NewButton(wrk.Name)
		btn.SetStyle(BtnStyle).
			SetActivatedStyle(ActiveBtnStyle).
			SetSelectedFunc(s.WrkBtnSelFunc).
			SetBorder(true)
		m.itms = append(m.itms, btn)
	}

	m.Refresh()
}

// CreateProjBtn generates a button that selects the given project.
func (s *Screen) CreateProjBtn(proj insomnium.Project, wrkMatrix *Matrix) *tv.Button {
	btn := tv.NewButton(proj.Name)
	btn.SetStyle(BtnStyle).
		SetActivatedStyle(ActiveBtnStyle).
		SetSelectedFunc(s.SelectProjBtn(proj.ID, wrkMatrix)).
		SetBorder(true).
		SetFocusFunc(s.FocusProjBtn(proj.ID, wrkMatrix))
	return btn
}

// IgnoreNonMotion is a event function that checks if the input is a motion
// input (arrow keys, hjlk, q for quit, enter to select) and ignores and ignores
// is when it is not.
func (s *Screen) IgnoreNonMotion(event *tc.EventKey) *tc.EventKey {
	switch event.Key() {
	case tc.KeyEnter, tc.KeyEsc, tc.KeyUp, tc.KeyDown, tc.KeyLeft, tc.KeyRight:
		return event
	case tc.KeyRune:
		if strings.ContainsRune("hjlkq", event.Rune()) {
			return event
		}
	}
	return nil
}

// FocsProjBtn returns a func that is called when the app sets focus to the
// given button. The function returned repopulates the workspaces grid with
// workspaces whose ParentID is equal to the ID of the focused button.
func (s *Screen) FocusProjBtn(id string, m *Matrix) func() {
	return func() {
		s.PopulateWorkspaces(id, m)
	}
}

// SelectProjBtn is called when a project button is selected and it sets focus
// on the fist workspace in the grid. If there is no workspace in that grid, it
// sets focus to the command input field.
func (s *Screen) SelectProjBtn(id string, m *Matrix) func() {
	return func() {
		fb, ok := m.Get(0, 0)
		if !ok {
			s.app.SetFocus(s.cmd)
		} else {
			s.app.SetFocus(fb)
		}
	}
}

// WrkBtnSelFunc is called when a workspace is selected.
func (s *Screen) WrkBtnSelFunc() {
	s.app.SetFocus(s.cmd)
}

// Run switches the page and starts the app.
func (s *Screen) Run(page string) error {
	s.pages.SwitchToPage(page)
	return s.app.SetRoot(s.layout, true).SetFocus(s.cmd).Run()
}

// CreateCmdField creates an input field that accepts commands.
func (s *Screen) CreateCmdField() *tv.InputField {
	return tv.NewInputField().
		SetFieldStyle(CmdFieldStyle).
		SetAcceptanceFunc(s.CmdAcceptanceFunc).
		SetDoneFunc(s.CmdDoneFunc)
}

// CmdAcceptanceFunc ignores input text to the cmd input field when the command
// does not start with a colon.
func (s *Screen) CmdAcceptanceFunc(textToCheck string, lastChar rune) bool {
	return strings.HasPrefix(textToCheck, ":")
}

// CmdDoneFunc removes the starting colon from the input field, clears the
// content and executes the command.
func (s *Screen) CmdDoneFunc(key tc.Key) {
	if key != tc.KeyEnter {
		return
	}
	t, _ := strings.CutPrefix(s.cmd.GetText(), ":")
	s.cmd.SetText("")
	if t != "" {
		s.Execute(s.ReadCommand(t))
	}
}

// Execute finds and runs the function for the given command name or outputs an
// error in the input field.
func (s *Screen) Execute(cmd string, args []string) {
	commands := s.GetCommands()
	fn, ok := commands[cmd]
	if !ok {
		s.ErrorMessage(fmt.Sprintf("Command %q not found", cmd))
		return
	}
	fn(args)
}

// GetCommands returns all the existing commands.
func (s *Screen) GetCommands() map[string]func(args []string) {
	return map[string]func(args []string){
		"startinsert": s.StartInsert,
		"q":           s.Quit,
	}
}

// StartInsert focuses on a given widget.
func (s *Screen) StartInsert(args []string) {
	if s.focus != nil {
		s.app.SetFocus(s.focus)
	}
}

// Quit closes the app.
func (s *Screen) Quit(args []string) {
	s.app.Stop()
}

// ReadCommand takes the first word as the command and passes the rest of the
// text to be processed as a set of arguments
func (s *Screen) ReadCommand(input string) (cmd string, args []string) {
	end := strings.Index(input, " ")
	if end < 0 {
		return input, nil
	}
	return input[:end], s.ReadArgs(input[end:])
}

// ReadArgs splits the given text as a list of arguments, respecting quotes.
func (s *Screen) ReadArgs(input string) (res []string) {
	buf := bytes.NewBuffer(nil)
	var inWord, inQuotes, inBreak bool
	var quoteType rune
	for _, c := range input {
		switch {
		case inWord && strings.ContainsRune(" \t\n", c):
			inWord = false
			buf.Reset()
		case inBreak:
			buf.WriteRune(c)
		case inQuotes && c == quoteType:
			inQuotes = false
			res = append(res, buf.String())
			buf.Reset()
		case inQuotes && c == EscapeRune:
			inBreak = true
		case inQuotes && c == quoteType:
			inQuotes = false
			res = append(res, buf.String())
			buf.Reset()
		case !inWord && !inQuotes && strings.ContainsRune("\"", c):
			inQuotes = true
		default:
			buf.WriteRune(c)
		}
	}
	if buf.Len() > 0 {
		res = append(res, buf.String())
	}
	return
}

// ErrorMessage sets the input field text as the given message and changes the
// input fg color to red. It also locks the input until the user presses Enter.
func (s *Screen) ErrorMessage(msg string) {
	s.cmd.SetFieldStyle(CmdFieldErrorStyle)
	s.cmd.SetText(msg)
	s.cmd.SetInputCapture(func(key *tc.EventKey) *tc.EventKey {
		if key.Key() == tc.KeyEnter {
			s.cmd.SetFieldStyle(CmdFieldStyle)
			s.cmd.SetText("")
			s.cmd.SetInputCapture(nil)
		}
		return nil
	})
}

// NewMatrix returns a new instance of a matrix that holds a *tview.Grid
func NewMatrix() *Matrix {
	grid := tv.NewGrid()
	buttonUp := tv.NewButton("△")
	buttonUp.SetStyle(BtnStyle)
	buttonDown := tv.NewButton("▽")
	buttonDown.SetStyle(BtnStyle)
	flex := tv.NewFlex().
		SetDirection(tv.FlexRow).
		AddItem(buttonUp, 0, 0, false).
		AddItem(grid, 0, 1, false).
		AddItem(buttonDown, 0, 0, false)
	return &Matrix{
		Flex:       flex,
		Grid:       grid,
		ButtonUp:   buttonUp,
		ButtonDown: buttonDown,
	}
}

// Matrix holds a *tview.Grid and is used to help setting a motion between the
// elements without loosing track of the individual objects.
type Matrix struct {
	*tv.Flex
	Grid         *tv.Grid
	ButtonUp     *tv.Button
	ButtonDown   *tv.Button
	currX, currY int
	skip         int
	w, h         int
	itms         []tv.Primitive
}

// SetWidth sets the number of columns in the matrix grid.
func (m *Matrix) SetWidth(w int) *Matrix {
	m.w = w
	return m
}

// SetHeight sets the number of rows in the matrix grid.
func (m *Matrix) SetHeight(h int) *Matrix {
	m.h = h
	return m
}

// Regresh reloads and reorders the grid elements.
func (m *Matrix) Refresh() {
	m.Grid.Clear()
	for n := m.skip * m.w; n < (m.h*m.w)+m.skip*m.w; n++ {
		var itm tv.Primitive
		if n < len(m.itms) {
			itm = m.itms[n]
		} else {
			itm = tv.NewBox()
		}
		x := n % m.w
		y := n / m.w
		m.Grid.AddItem(itm, y-m.skip, x, 1, 1, 1, 1, false)
	}
	if m.skip > 0 {
		m.Flex.ResizeItem(m.ButtonUp, 1, 0)
	} else {
		m.Flex.ResizeItem(m.ButtonUp, 0, 0)
	}
	a := ((len(m.itms) + m.w - 1) / m.w) - m.skip - m.h
	if a > 0 {
		m.Flex.ResizeItem(m.ButtonDown, 1, 0)
	} else {
		m.Flex.ResizeItem(m.ButtonDown, 0, 0)
	}
}

// RefreshY checks if the selected cell is within the visible range. If it is
// not, it adjusts the range so the element can be displayed and then refreshes
// the grid.
func (m *Matrix) RefreshY() {
	if m.currY > m.skip+m.h-1 {
		m.skip++
		m.Refresh()
	} else if m.currY < m.skip {
		m.skip--
		m.Refresh()
	}
}

// Left moves the selected x position one cell to the left. If the current row
// is the last one, it means the current x position might be beyond the last
// element in the list, because the last row is the only one that may not be
// fully populated. In this case, it adjusts the x position to refer to the last
// element in the list.
func (m *Matrix) Left() {
	if m.currX > 0 {
		m.currX--
	}
	lastRow := len(m.itms) / m.w
	lastRowLimit := (len(m.itms) % m.w) - 1
	if m.currY == lastRow {
		m.currX = min(m.currX, lastRowLimit)
	}
}

// Right moves the selected x position one cell to the right. If the new cell is
// beyond the current row limit, it does nothing.
func (m *Matrix) Right() {
	lastRow := len(m.itms) / m.w
	currLineLimit := m.w
	if m.currY == lastRow {
		currLineLimit = (len(m.itms) % m.w) - 1
	}
	if m.currX < currLineLimit {
		m.currX++
	}
}

// Up moves the selected y position one cell up, unless it would result in a
// negative y position. It also calls the function that updates the range, which
// adjusts it if necessary and refreshes the grid.
func (m *Matrix) Up() {
	if m.currY > 0 {
		m.currY--
	}
	m.RefreshY()
}

// Down moves the selected y position one cell down, unless it would exceed the
// last row. It also calls the function that updates the range, which adjusts it
// if necessary and refreshes the grid.
func (m *Matrix) Down() {
	lastRow := len(m.itms) / m.w
	if m.currY < lastRow {
		m.currY++
	}
	m.RefreshY()
}

// Get returns the item for the given location. If there is no item, it returns
// nil and the returned boolean is set to false.
func (m *Matrix) Get(x, y int) (tv.Primitive, bool) {
	if len(m.itms) == 0 {
		return nil, false
	}
	i := (y * m.w) + x
	if i >= len(m.itms) {
		if i <= m.w*((len(m.itms)+m.w-1)/m.w) {
			return m.itms[len(m.itms)-1], true
		} else {
			return nil, false
		}
	}
	return m.itms[i], true
}

// GetCurrentPrimitive returns the currently selected cell. If the position
// refers to a blank spot in the grid, it adjusts it accordingly.
func (m *Matrix) GetCurrentPrimitive() tv.Primitive {
	totalRows := (len(m.itms) + m.w - 1) / m.w
	m.currX = max(m.currX, 0)
	m.currX = min(m.currX, m.w-1)
	m.currY = max(m.currY, 0)
	m.currY = min(m.currY, totalRows-1)
	curr, _ := m.Get(m.currX, m.currY)
	return curr
}

// InputFn defines a function type used by tview to capture input, taking a
// pointer to a tcell.EventKey and returning a pointer to a tcell.EventKey.
type InputFn = func(*tc.EventKey) *tc.EventKey
