package main

import (
	"bytes"
	"fmt"
	"strings"
	"sync"

	tc "github.com/gdamore/tcell/v2"
	tv "github.com/rivo/tview"
	"github.com/vinicius-lino-figueiredo/insomnium"
)

// EscapeRune precedes a rune that will be escaped
const EscapeRune = '\\'

// Define styles and colors used in the TUI interface.
var (
	FontColor = tc.ColorWhite

	// Style for command field
	CmdFieldStyle = tc.StyleDefault.Foreground(FontColor)

	// Style for buttons
	BtnStyle = tc.StyleDefault.Foreground(FontColor)

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
	mainPane := tv.NewFlex().SetDirection(tv.FlexRow)
	pages := s.CreatePages()
	cmdField := s.CreateCmdField()

	mainPane.AddItem(pages, 0, 1, false).
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

// CreateProjsScreen generates a workspace selecting screen.
func (s *Screen) CreateProjsScreen() *tv.Flex {
	projGrid := tv.NewGrid()
	projGrid.SetBorder(true).
		SetTitle("Projects").
		SetInputCapture(s.IgnoreNonMotion)
	wrkGrid := tv.NewGrid()
	wrkGrid.SetBorder(true).
		SetTitle("Workspaces").
		SetInputCapture(s.IgnoreNonMotion)

	projMatrix := NewMatrix(projGrid)
	wrkMatrix := NewMatrix(wrkGrid)

	s.PopulateProjects(projMatrix, wrkMatrix)

	s.PopulateWorkspaces("", wrkMatrix)
	workspacePane := tv.NewFlex().
		AddItem(projGrid, 0, 1, false).
		AddItem(wrkGrid, 0, 5, false)
	return workspacePane
}

// PopulateProjects loads the project grid with all the projects found in the
// *insomnium.Insomnium global instance.
func (s *Screen) PopulateProjects(projMatrix, wrkMatrix *Matrix) {
	once := &sync.Once{}
	for n, proj := range inso.Projects {
		projectButton := s.CreateProjBtn(proj, wrkMatrix)
		item := projMatrix.Set(projectButton, 0, n)
		projectButton.SetStyle(BtnStyle).
			SetActivatedStyle(ActiveBtnStyle).
			SetSelectedFunc(s.SelectProjBtn(proj.ID, wrkMatrix)).
			SetBorder(true).
			SetFocusFunc(s.FocusProjBtn(proj.ID, wrkMatrix)).
			SetInputCapture(s.ProjInputFn(item, wrkMatrix))
		once.Do(func() { s.focus = projectButton })
	}
	for _, i := range projMatrix.All() {
		projMatrix.grid.AddItem(i.t, i.y, i.x, 1, 1, 3, 3, false)
	}
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
			s.app.SetFocus(fb.t)
		}
	}
}

// ProjInputFn is called for every input when a project button is focused. It
// detects motion (arrow keys and hjkl), selection (enter) and quit (Escape key
// and q) commands.
func (s *Screen) ProjInputFn(projItem *MatrixItem, wrkMatrix *Matrix) InputFn {
	return func(key *tc.EventKey) *tc.EventKey {
		fn := s.MotionInput(projItem)
		fn(key)
		switch {
		case key.Key() == tc.KeyEsc || key.Rune() == 'q':
			s.PopulateWorkspaces("", wrkMatrix)
		case key.Key() == tc.KeyEnter:
			return key
		default:
		}
		return nil
	}
}

// MotionInput detects arrow keys and vim motions (hjkl) and changes the
// selected item of a *Matrix.
func (s *Screen) MotionInput(sb *MatrixItem) InputFn {
	return func(key *tc.EventKey) *tc.EventKey {
		var p tv.Primitive
		switch {
		case key.Key() == tc.KeyDown || key.Rune() == 'j':
			p = sb.Down().t
		case key.Key() == tc.KeyUp || key.Rune() == 'k':
			p = sb.Up().t
		case key.Key() == tc.KeyRight || key.Rune() == 'l':
			p = sb.Right().t
		case key.Key() == tc.KeyLeft || key.Rune() == 'h':
			p = sb.Left().t
		case key.Key() == tc.KeyEsc || key.Rune() == 'q':
			p = s.cmd
		default:
		}
		s.app.SetFocus(p)
		return key
	}
}

// PopulateWorkspaces loads in the grid all the workspaces whose ParentID are
// equal to the given id. If the id arg is "", the func populates the grid with
// no filter.
func (s *Screen) PopulateWorkspaces(id string, m *Matrix) {
	m.Clear()
	m.grid.Clear()
	var n int
	for _, wrk := range inso.Workspaces {
		if id != "" && wrk.ParentID != id {
			continue
		}
		btn := tv.NewButton(wrk.Name)
		lmi := m.Set(btn, n%5, n/5)
		btn.SetStyle(BtnStyle).
			SetActivatedStyle(ActiveBtnStyle).
			SetSelectedFunc(s.WrkBtnSelFunc).
			SetBorder(true).
			SetInputCapture(s.MotionInput(lmi))
		n++
	}
	for _, i := range m.All() {
		if i != nil {
			m.grid.AddItem(i.t, i.y, i.x, 1, 1, 1, 1, false)
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
func NewMatrix(grid *tv.Grid) *Matrix {
	return &Matrix{
		grid: grid,
	}
}

// Matrix holds a *tview.Grid and is used to help setting a motion between the
// elements without loosing track of the individual objects.
type Matrix struct {
	grid  *tv.Grid
	items [][]*MatrixItem
}

// Get returns the item for the given location. If there is no item, it returns
// nil and the returned boolean is set to false.
func (lm *Matrix) Get(x, y int) (*MatrixItem, bool) {
	if y >= len(lm.items) {
		return nil, false
	}
	if x >= len(lm.items[y]) {
		return nil, false
	}
	return lm.items[y][x], true
}

// Set is used to set a button as the item in the given location.
func (lm *Matrix) Set(t *tv.Button, x, y int) *MatrixItem {
	if cap(lm.items) <= y {
		n := make([][]*MatrixItem, y+1)
		copy(n, lm.items)
		lm.items = n
	}
	if cap(lm.items[y]) <= x {
		n := make([]*MatrixItem, x+1)
		copy(n, lm.items[y])
		lm.items[y] = n
	}
	lmi := &MatrixItem{parent: lm, x: x, y: y, t: t}
	lm.items[y][x] = lmi
	return lmi
}

// All returns all of the matrix items as a list.
func (lm *Matrix) All() (res []*MatrixItem) {
	for _, row := range lm.items {
		for _, item := range row {
			res = append(res, item)
		}
	}
	return
}

// Clear removes all of the matrix items
func (lm *Matrix) Clear() {
	lm.items = nil
}

// MatrixItem should be used as a sub element of Matrix and it stores a
// *tview.Button, wich is the type that is used in the app as a grid cell.
type MatrixItem struct {
	parent *Matrix
	x, y   int
	t      *tv.Button
}

// Right returns the element to the right, or the leftmost element if the
// current one is the rightmost.
func (lmi *MatrixItem) Right() *MatrixItem {
	if lmi.x == len(lmi.parent.items[lmi.y])-1 {
		return lmi.parent.items[lmi.y][0]
	}
	return lmi.parent.items[lmi.y][lmi.x+1]
}

// Left returns the element to the left, or the rightmost element if the current
// one is the leftmost.
func (lmi *MatrixItem) Left() *MatrixItem {
	if lmi.x == 0 {
		return lmi.parent.items[lmi.y][len(lmi.parent.items[lmi.y])-1]
	}
	return lmi.parent.items[lmi.y][lmi.x-1]
}

// Up returns the element above, or the lowermost element if the current one is
// the topmost.
func (lmi *MatrixItem) Up() *MatrixItem {
	if lmi.y == 0 {
		lastRow := lmi.parent.items[len(lmi.parent.items)-1]
		return lastRow[min(lmi.x, len(lastRow)-1)]
	}
	lastRow := lmi.parent.items[lmi.y-1]
	return lastRow[min(lmi.x, len(lastRow)-1)]
}

// Down returns the element below, or the uppermost element if the current one
// is the bottommost.
func (lmi *MatrixItem) Down() *MatrixItem {
	if lmi.y == len(lmi.parent.items)-1 {
		firstRow := lmi.parent.items[0]
		return firstRow[min(len(firstRow), lmi.x)]
	}
	firstRow := lmi.parent.items[lmi.y+1]
	return firstRow[min(len(firstRow)-1, lmi.x)]
}

// InputFn defines a function type used by tview to capture input, taking a
// pointer to a tcell.EventKey and returning a pointer to a tcell.EventKey.
type InputFn = func(*tc.EventKey) *tc.EventKey
