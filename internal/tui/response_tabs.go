package tui

import (
	"main/internal/component"

	tc "github.com/gdamore/tcell/v2"
	"github.com/vinicius-lino-figueiredo/insomnium"
)

// NewResponseTabs creates and initializes a new ResponseTabs instance.
func NewResponseTabs() *ResponseTabs {
	rt := &ResponseTabs{}
	rt.Load()
	return rt
}

// ResponseTabs represents a tab view with response info.
type ResponseTabs struct {
	*component.Tabs
	Body       *ResponseBody
	Headers    *ResponseHeaders
	styleText  func(string, tc.Style) string
	response   *insomnium.Response
	resetFocus func()
}

// Load initializes the tabs and loads body and headers if needed.
func (rt *ResponseTabs) Load() {
	if rt.Body == nil {
		rt.LoadBody()
	}
	if rt.Headers == nil {
		rt.LoadHeaders()
	}
	rt.Tabs = component.NewTabs()
	rt.Tabs.AddTab("Body", rt.Body)
	rt.Tabs.AddTab("Headers", rt.Headers)
	rt.OpenTab("Headers")
}

// LoadBody creates and configures the response body component.
func (rt *ResponseTabs) LoadBody() {
	rt.Body = NewResponseBody()
	rt.Body.SetBorder(true)
}

// LoadHeaders creates the response headers component.
func (rt *ResponseTabs) LoadHeaders() {
	rt.Headers = NewResponseHeaders()
}

// SetResponse updates both body and headers with a new response.
func (rt *ResponseTabs) SetResponse(res *insomnium.Response) *ResponseTabs {
	rt.Body.SetResponse(res)
	rt.Headers.SetResponse(res)
	return rt
}

// SetStyleTextFunc sets the function to style the text in the body.
func (rt *ResponseTabs) SetStyleTextFunc(fn func(string, tc.Style) string) *ResponseTabs {
	rt.Body.SetStyleTextFunc(fn)
	return rt
}

// SetGetResponseBodyFunc sets the function to retrieve the response body.
func (rt *ResponseTabs) SetGetResponseBodyFunc(fn func(*insomnium.Response) (string, error)) *ResponseTabs {
	rt.Body.SetGetResponseBodyFunc(fn)
	return rt
}

// SetResponseBodyTheme sets the syntax highlighting theme for the body.
func (rt *ResponseTabs) SetResponseBodyTheme(theme string) *ResponseTabs {
	rt.Body.SetResponseBodyTheme(theme)
	return rt
}

// SetBackgroundColor sets the background color for both body and headers.
func (rt *ResponseTabs) SetBackgroundColor(bg tc.Color) *ResponseTabs {
	rt.Body.SetBackgroundColor(bg)
	rt.Headers.SetBackgroundColor(bg)
	return rt
}

// SetResetFocusFunc sets the function called to reset focus from tabs.
func (rt *ResponseTabs) SetResetFocusFunc(fn func()) *ResponseTabs {
	rt.resetFocus = fn
	rt.Headers.SetResetFocusFunc(fn)
	rt.Body.SetResetFocusFunc(fn)
	return rt
}
