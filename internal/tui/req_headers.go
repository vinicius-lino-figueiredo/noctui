package tui

import (
	tc "github.com/gdamore/tcell/v2"
	tv "github.com/rivo/tview"
	"github.com/vinicius-lino-figueiredo/insomnium"
)

// NewReqHeaders creates and initializes a new ReqHeaders widget, which allows
// editing HTTP headers.
func NewReqHeaders() *ReqHeaders {
	rh := &ReqHeaders{}
	rh.Load()
	return rh
}

// ReqHeaders is a tview widget for managing and editing HTTP headers. It embeds
// a Matrix to organize the header fields in a grid layout.
type ReqHeaders struct {
	*Matrix
	inputStyle      tc.Style
	focus           func(tv.Primitive)
	NeutralizeFocus func()
}

// Load initializes the ReqHeaders layout and sets the default behavior for
// input capture and dimensions.
func (rh *ReqHeaders) Load() {
	if rh.Matrix == nil {
		rh.Matrix = NewMatrix()
	}
	rh.SetWidth(1).
		SetHeight(5).
		SetInputCapture(rh.MatrixCapture)
}

// NewHeader creates a new Header widget with the current styling and
// focus/blur behavior.
func (rh *ReqHeaders) NewHeader(key string, value string) *Header {
	return NewHeader(key, value).
		SetBackgroundColor(rh.bg).
		SetForegroundColor(rh.fg).
		SetFocusFunc(rh.headerFocus).
		SetBlurFunc(rh.headerBlur).
		SetInputStyle(rh.inputStyle)
}

// headerFocus disables matrix-level input capture to allow editing the header.
func (rh *ReqHeaders) headerFocus() {
	rh.SetInputCapture(nil)
}

// headerBlur restores the matrix-level input capture after header editing.
func (rh *ReqHeaders) headerBlur() {
	rh.SetInputCapture(rh.MatrixCapture)
}

// MatrixCapture handles keyboard navigation and actions within the matrix.
// - Enter opens the current header for editing.
// - Esc/q neutralizes focus.
// - Arrow keys or hjkl navigate.
func (rh *ReqHeaders) MatrixCapture(event *tc.EventKey) *tc.EventKey {
	switch {
	case event.Key() == tc.KeyEnter:
		h := rh.GetCurrentPrimitive().(*Header)
		rh.OpenHeader(h)
		return nil
	case event.Key() == tc.KeyEsc || event.Rune() == 'q':
		rh.NeutralizeFocus()
		return nil
	case event.Key() == tc.KeyLeft || event.Rune() == 'h':
		rh.Left()
	case event.Key() == tc.KeyRight || event.Rune() == 'l':
		rh.Right()
	case event.Key() == tc.KeyUp || event.Rune() == 'k':
		rh.Up()
	case event.Key() == tc.KeyDown || event.Rune() == 'j':
		rh.Down()
	default:
		return event
	}
	p := rh.GetCurrentPrimitive()
	if p != nil {
		rh.focus(p)
	}
	return nil
}

// OpenHeader sets focus on the key input field of the given Header.
func (rh *ReqHeaders) OpenHeader(h *Header) {
	rh.focus(h.KeyInput)
}

// Focus sets the focus function for the ReqHeaders and its Matrix.
func (rh *ReqHeaders) Focus(fn func(tv.Primitive)) {
	rh.focus = fn
	rh.Matrix.Focus(fn)
}

// SetNeutralizeFocusFunc defines a function to be called to reset the focus
// state.
func (rh *ReqHeaders) SetNeutralizeFocusFunc(fn func()) *ReqHeaders {
	rh.NeutralizeFocus = fn
	return rh
}

// SetRequest replaces the current headers with those from the given request.
func (rh *ReqHeaders) SetRequest(req *insomnium.Request) *ReqHeaders {
	rh.Matrix.itms = rh.Matrix.itms[:0]
	if req == nil {
		rh.Matrix.Refresh()
		return rh
	}
	for _, header := range req.Headers {
		b := rh.NewHeader(header.Name, header.Value)
		rh.Matrix.itms = append(rh.Matrix.itms, b)
	}
	rh.Matrix.Refresh()
	return rh
}
