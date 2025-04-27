package tui

import (
	"main/internal/component"

	tc "github.com/gdamore/tcell/v2"
	tv "github.com/rivo/tview"
	"github.com/vinicius-lino-figueiredo/insomnium"
)

// NewResponseHeaders returns a new instance of a response header viewer.
func NewResponseHeaders() *ResponseHeaders {
	rh := &ResponseHeaders{}
	rh.Load()
	return rh
}

// ResponseHeaders is a widget what holds the headers of a http response.
type ResponseHeaders struct {
	*component.Matrix
	inputStyle tc.Style
	focus      func(tv.Primitive)
	ResetFocus func()
}

// Load initializes and sets up the widget.
func (rh *ResponseHeaders) Load() {
	if rh.Matrix == nil {
		rh.Matrix = component.NewMatrix()
	}
	rh.SetWidth(1).
		SetHeight(5).
		SetInputCapture(rh.MatrixCapture)
}

// NewHeader returns a new widget to represent a response header.
func (rh *ResponseHeaders) NewHeader(key string, value string) *Header {
	return NewHeader(key, value).
		SetBackgroundColor(rh.GetBackgroundColor()).
		SetForegroundColor(rh.GetForegroundColor()).
		SetInputStyle(rh.inputStyle)
}

// MatrixCapture is called on input event and moves the selection.
func (rh *ResponseHeaders) MatrixCapture(event *tc.EventKey) *tc.EventKey {
	switch {
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

// Focus implements tview.Primitive.
func (rh *ResponseHeaders) Focus(fn func(tv.Primitive)) {
	rh.focus = fn
	rh.Matrix.Focus(fn)
}

// SetResetFocusFunc sets a function that is called when the widget wants to
// reset focus.
func (rh *ResponseHeaders) SetResetFocusFunc(fn func()) *ResponseHeaders {
	rh.ResetFocus = fn
	return rh
}

// SetResponse sets the current response.
func (rh *ResponseHeaders) SetResponse(req *insomnium.Response) *ResponseHeaders {
	rh.Matrix.Clear()
	if req == nil {
		rh.Matrix.Refresh()
		return rh
	}
	for _, header := range req.Headers {
		b := rh.NewHeader(header.Name, header.Value)
		rh.Matrix.AddItem(b)
	}
	rh.Matrix.Refresh()
	return rh
}
