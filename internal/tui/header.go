package tui

import (
	tc "github.com/gdamore/tcell/v2"
	tv "github.com/rivo/tview"
)

// NewHeader creates and initializes a new Header component with the given key
// and value. It loads the internal layout and sets the initial input values.
func NewHeader(key string, value string) *Header {
	h := &Header{}
	h.Load()
	h.SetKey(key).
		SetValue(value)
	return h
}

// Header represents a UI component for a key-value input pair, built using Flex
// layouts and styled boxes.
type Header struct {
	*tv.Flex
	Inner      *tv.Flex
	KeyInput   *tv.InputField
	ValueInput *tv.InputField
	box        *tv.Box
	bg         tc.Color
	fg         tc.Color
	focus      func(tv.Primitive)
}

// Load initializes all components of the Header if they are not already set.
// It builds the full UI layout, including borders, colors, and input capturing.
func (h *Header) Load() {
	if h.KeyInput == nil {
		h.LoadKey()
	}
	if h.ValueInput == nil {
		h.LoadValue()
	}
	if h.Flex == nil {
		h.Flex = tv.NewFlex()
	}
	if h.box == nil {
		h.box = tv.NewBox().SetBackgroundColor(h.bg)
	}
	if h.Inner == nil {
		h.Inner = tv.NewFlex()
	}
	h.Flex.SetDirection(tv.FlexRow).
		AddItem(h.box, 0, 1, false).
		AddItem(h.Inner, 1, 0, false).
		AddItem(h.box, 0, 1, false).
		SetBorder(true).
		SetBackgroundColor(h.bg).
		SetInputCapture(h.InputCapture)
	h.Inner.SetDirection(tv.FlexColumn).
		AddItem(h.box, 1, 0, false).
		AddItem(h.KeyInput, 0, 1, false).
		AddItem(h.box, 1, 0, false).
		AddItem(h.ValueInput, 0, 1, false).
		AddItem(h.box, 1, 0, false).
		SetBackgroundColor(h.bg).
		SetBorderColor(h.fg)
}

// LoadKey initializes the KeyInput field with a label and a DoneFunc.
func (h *Header) LoadKey() {
	h.KeyInput = tv.NewInputField().
		SetLabel("Key: ").
		SetDoneFunc(h.InputDoneFunc)
}

// LoadValue initializes the ValueInput field with a label and a DoneFunc.
func (h *Header) LoadValue() {
	h.ValueInput = tv.NewInputField().
		SetLabel("Value: ").
		SetDoneFunc(h.InputDoneFunc)
}

// InputDoneFunc handles the event when input is completed (e.g. Enter is
// pressed). It returns focus to the outer Flex container.
func (h *Header) InputDoneFunc(_ tc.Key) {
	h.focus(h.Flex)
}

// Focus is called when the app focuses on the header. It stores the focus
// function and delegates the focus to the embeded flex.
func (h *Header) Focus(fn func(tv.Primitive)) {
	h.focus = fn
	h.Flex.Focus(fn)
}

// InputCapture intercepts keyboard events. It cycles focus between inputs when
// Tab or Backtab is pressed.
func (h *Header) InputCapture(event *tc.EventKey) *tc.EventKey {
	switch event.Key() {
	case tc.KeyTab, tc.KeyBacktab:
		h.CycleFocus()
		return nil
	default:
		return event
	}
}

// CycleFocus switches focus between KeyInput and ValueInput fields.
func (h *Header) CycleFocus() {
	if h.KeyInput.HasFocus() {
		h.focus(h.ValueInput)
	} else if h.ValueInput.HasFocus() {
		h.focus(h.KeyInput)
	} else {
	}
}

// SetKey sets the text of the KeyInput field.
func (h *Header) SetKey(key string) *Header {
	h.KeyInput.SetText(key)
	return h
}

// SetValue sets the text of the ValueInput field.
func (h *Header) SetValue(value string) *Header {
	h.ValueInput.SetText(value)
	return h
}

// SetBackgroundColor sets the background color for the entire Header, including
// labels and spacing boxes.
func (h *Header) SetBackgroundColor(bg tc.Color) *Header {
	h.bg = bg
	h.Flex.SetBackgroundColor(bg)
	h.KeyInput.SetLabelStyle(tc.StyleDefault.Background(bg).Foreground(h.fg))
	h.ValueInput.SetLabelStyle(tc.StyleDefault.Background(bg).Foreground(h.fg))
	h.box.SetBackgroundColor(bg)
	return h
}

// SetForegroundColor sets the foreground (border/label) color for the Header.
func (h *Header) SetForegroundColor(fg tc.Color) *Header {
	baseStyle := tc.StyleDefault.Background(h.bg)
	h.fg = fg
	h.Flex.SetBorderColor(fg)
	h.KeyInput.SetLabelStyle(baseStyle.Foreground(fg))
	h.ValueInput.SetLabelStyle(baseStyle.Foreground(fg))
	return h
}

// SetFocusFunc sets a function to be called when either input field gets focus.
func (h *Header) SetFocusFunc(fn func()) *Header {
	h.KeyInput.SetFocusFunc(fn)
	h.ValueInput.SetFocusFunc(fn)
	return h
}

// SetBlurFunc sets a function to be called when either input field loses focus.
func (h *Header) SetBlurFunc(fn func()) *Header {
	h.KeyInput.SetBlurFunc(fn)
	h.ValueInput.SetBlurFunc(fn)
	return h
}

// SetInputStyle sets the visual style of the input fields.
func (h *Header) SetInputStyle(style tc.Style) *Header {
	h.KeyInput.SetFieldStyle(style)
	h.ValueInput.SetFieldStyle(style)
	return h
}
