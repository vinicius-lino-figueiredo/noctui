package tui

import (
	tc "github.com/gdamore/tcell/v2"
	tv "github.com/rivo/tview"
	"github.com/vinicius-lino-figueiredo/insomnium"
)

// NewReqTabs creates a new instance of the tab used to manage the request
// editor tabs
func NewReqTabs() *ReqTabs {
	rt := &ReqTabs{}
	rt.Load()
	rt.Tabs.OpenTab("Body")
	return rt
}

// ReqTabs is a tview widget that organizes request editing tools into tabs. It
// embeds a *Tabs to switch between different request editors like body and
// headers.
type ReqTabs struct {
	*Tabs
	BodyTab          *ReqBody
	HeadersTab       *ReqHeaders
	bg               tc.Color
	fg               tc.Color
	headerInputStyle tc.Style
	focus            func(tv.Primitive)
	ResetFocus       func()
}

// Load initializes the ReqTabs and it's content.
func (rt *ReqTabs) Load() {
	if rt.BodyTab == nil {
		rt.LoadBodyTab()
	}
	if rt.HeadersTab == nil {
		rt.LoadHeadersTab()
	}
	if rt.Tabs == nil {
		rt.Tabs = NewTabs()
	}
	rt.AddTab("Body", rt.BodyTab)
	rt.AddTab("Headers", rt.HeadersTab)
}

// LoadBodyTab loads the request body editor widget.
func (rt *ReqTabs) LoadBodyTab() {
	rt.BodyTab = NewReqBody()
}

// LoadHeadersTab loads the request header editor widget.
func (rt *ReqTabs) LoadHeadersTab() {
	rt.HeadersTab = NewReqHeaders()
}

// Focus is called when the app focuses on the request tabs widget and it passes
// the focus to the embeded *Tabs field.
func (rt *ReqTabs) Focus(fn func(tv.Primitive)) {
	rt.focus = fn
	rt.Tabs.Focus(fn)
}

// SetRequest loads the provided request into all the tabs.
func (rt *ReqTabs) SetRequest(req *insomnium.Request) *ReqTabs {
	rt.HeadersTab.SetRequest(req)
	rt.BodyTab.SetRequest(req)
	return rt
}

// SetBackgroundColor sets the background color for all elements of the widget,
// including the Body and Headers editors.
func (rt *ReqTabs) SetBackgroundColor(bg tc.Color) *ReqTabs {
	rt.bg = bg
	rt.Tabs.SetBackgroundColor(bg)
	rt.BodyTab.SetBackgroundColor(bg)
	rt.HeadersTab.SetBackgroundColor(bg)
	for _, header := range rt.HeadersTab.GetAll() {
		h := header.(*Header)
		h.SetBackgroundColor(bg)
	}
	return rt
}

// SetForegroundColor sets the foreground color for the request tabs.
func (rt *ReqTabs) SetForegroundColor(fg tc.Color) *ReqTabs {
	rt.fg = fg
	rt.Tabs.SetForegroundColor(fg)
	rt.BodyTab.SetBorderColor(fg)
	rt.HeadersTab.SetForegroundColor(fg)
	for _, header := range rt.HeadersTab.itms {
		h := header.(*tv.Button)
		h.SetBorderColor(fg)
	}
	return rt
}

// SetResetFocusFunc sets a callback used to neutralize focus state. It
// applies this function to both the Body and Headers tabs.
func (rt *ReqTabs) SetResetFocusFunc(fn func()) *ReqTabs {
	rt.ResetFocus = fn
	rt.HeadersTab.SetResetFocusFunc(fn)
	rt.BodyTab.SetResetFocus(fn)
	return rt
}

// SetHeaderInputStyle defines the style to be applied to input fields inside
// the Headers tab.
func (rt *ReqTabs) SetHeaderInputStyle(style tc.Style) *ReqTabs {
	rt.headerInputStyle = style
	for _, header := range rt.HeadersTab.itms {
		h := header.(*Header)
		h.SetInputStyle(style)
	}
	return rt
}

// SetStyleTextFunc sets a function used to format styled text inside the Body
// editor.
func (rt *ReqTabs) SetStyleTextFunc(fn func(string, tc.Style) string) *ReqTabs {
	rt.BodyTab.SetStyleTextFunc(fn)
	return rt
}

// SetRequestBodyInputStyle sets the style of the Body tab's input editor.
func (rt *ReqTabs) SetRequestBodyInputStyle(style tc.Style) *ReqTabs {
	rt.BodyTab.SetStyle(style)
	return rt
}

// SetRequestBodyTheme applies a theme ("monokai", etc.) to the Body tab.
func (rt *ReqTabs) SetRequestBodyTheme(theme string) *ReqTabs {
	rt.BodyTab.SetTheme(theme)
	return rt
}
