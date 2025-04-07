package tui

import (
	"bytes"
	"strings"

	tc "github.com/gdamore/tcell/v2"
	tv "github.com/rivo/tview"
)

const (
	CmdLineInstanceNormal = iota
	CmdLineInstanceError
)

// NewCmdLine creates a *CmdLine, a tview.Primitive used for user input. It is a
// *tview.Flex that holds focus until ":" is pressed, then transfers focus to
// the *tview.InputField within the main Flex. When the input field is cleared,
// focus returns to the Flex. On Enter, the widget reads the input, extracts the
// command and arguments, and passes them to the defined execFunc.
func NewCmdLine(app *tv.Application) *CmdLine {
	c := &CmdLine{
		Flex:     tv.NewFlex(),
		Input:    tv.NewInputField(),
		app:      app,
		Instance: CmdLineInstanceNormal,
	}

	c.Flex.AddItem(c.Input, 0, 1, false).
		SetInputCapture(c.flexCapture)
	c.Input.SetChangedFunc(c.inputChanged).
		SetDoneFunc(c.inputDone)

	return c
}

// CmdLine is a tview.Primitive that handles command-line style user input. It
// consists of a *tview.Flex containing a *tview.InputField. The widget holds
// focus until ":" is pressed, activating the input field. When cleared, focus
// returns to the Flex. On Enter, it parses the input into a command and
// arguments, passing them to the defined execFunc.
type CmdLine struct {
	*tv.Flex
	Input             *tv.InputField
	app               *tv.Application
	execFunc          func(string, []string)
	FieldStyle        tc.Style
	ErrorMessageStyle tc.Style
	Instance          int
}

// SetExecFunc sets the function to be executed when the user submits a command.
// The function receives the parsed command and its arguments. Returns the
// CmdLine itself for chaining.
func (cl *CmdLine) SetExecFunc(fn func(cmd string, args []string)) *CmdLine {
	cl.execFunc = fn
	return cl
}

// flexCapture handles key events for the CmdLine. If ":" is pressed and the
// input field is not focused, it clears the input and moves focus to the input
// field.
func (cl *CmdLine) flexCapture(event *tc.EventKey) *tc.EventKey {
	if event.Rune() == ':' && cl.app.GetFocus() != cl.Input {
		cl.SetText("")
		cl.app.SetFocus(cl.Input)
		cl.Instance = CmdLineInstanceNormal
		cl.RefreshStyle()
	}
	return event
}

// inputChanged is triggered when the input field's text changes. If the input
// field is focused and becomes empty, focus returns to the Flex.
func (cl *CmdLine) inputChanged(text string) {
	if cl.app.GetFocus() == cl.Input && text == "" {
		cl.app.SetFocus(cl.Flex)
	}
}

// inputDone is triggered when the user submits input. If Enter is pressed, it
// moves focus back to the Flex, clears the input field, and processes the
// entered command. Any other key is ignored.
func (cl *CmdLine) inputDone(key tc.Key) {
	if key != tc.KeyEnter {
		return
	}
	cl.app.SetFocus(cl.Flex)
	t := cl.GetText()
	cl.SetText("")
	cl.parseCommand(t)
}

// parseCommand processes the given command string. It trims the leading ":",
// extracts the command and arguments using ReadCommand, and executes the
// command if it's valid and an execFunc is set.
func (cl *CmdLine) parseCommand(command string) {
	command = strings.TrimLeft(command, ":")
	cmd, args := cl.ReadCommand(command)
	if cmd != "" && cl.execFunc != nil {
		cl.execFunc(cmd, args)
	}
}

// ReadCommand takes the first word as the command and passes the rest of the
// text to be processed as a set of arguments.
func (cl *CmdLine) ReadCommand(input string) (cmd string, args []string) {
	end := strings.Index(input, " ")
	if end < 0 {
		return input, nil
	}
	return input[:end], cl.ReadArgs(input[end:])
}

// ReadArgs parses the input string into a slice of arguments, handling spaces,
// quotes, and escape sequences. It supports both single and double quotes and
// preserves quoted text. The result is returned as a slice of strings, each
// representing an argument.
func (cl *CmdLine) ReadArgs(input string) (res []string) {
	buf := bytes.NewBuffer(nil)
	var inWord, inQuotes, inBreak bool
	var quoteType rune
	for _, c := range input {
		switch {
		case !inWord && !inQuotes && strings.ContainsRune(" \t\n", c):
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

// RefreshStyle updates the input field's style based on the current instance
// mode. It applies the standard style for normal mode and an error style when
// in error mode.
func (cl *CmdLine) RefreshStyle() {
	switch cl.Instance {
	case CmdLineInstanceNormal:
		cl.Input.SetFieldStyle(cl.FieldStyle)
	case CmdLineInstanceError:
		cl.Input.SetFieldStyle(cl.ErrorMessageStyle)
	default:
	}
}

// DisplayErrorMessage sets the input field's text to the given error message,
// changes the instance to error mode, and updates the field's style to reflect
// the error state.
func (cl *CmdLine) DisplayErrorMessage(message string) {
	cl.SetText(message)
	cl.Instance = CmdLineInstanceError
	cl.RefreshStyle()
}

// SetFieldStyle sets the style for the input field and updates the display to
// reflect the new style. Returns the CmdLine itself for method chaining.
func (cl *CmdLine) SetFieldStyle(style tc.Style) *CmdLine {
	cl.FieldStyle = style
	cl.RefreshStyle()
	return cl
}

// SetErrorMessageStyle sets the style for error messages in the input field
// and updates the display to reflect the new style. Returns the CmdLine itself
// for method chaining.
func (cl *CmdLine) SetErrorMessageStyle(style tc.Style) *CmdLine {
	cl.ErrorMessageStyle = style
	cl.RefreshStyle()
	return cl
}

// SetText sets the input field's text to the given value and returns the
// CmdLine itself for method chaining.
func (cl *CmdLine) SetText(text string) *CmdLine {
	cl.Input.SetText(text)
	return cl
}

// GetText retrieves the current text from the input field.
func (cl *CmdLine) GetText() string {
	return cl.Input.GetText()
}
