package tui

import (
	tc "github.com/gdamore/tcell/v2"
	tv "github.com/rivo/tview"
	"github.com/vinicius-lino-figueiredo/insomnium"
)

// NewResponseView creates a new ResponseView instance and initializes it.
func NewResponseView() *ResponseView {
	rv := &ResponseView{}
	rv.Load()
	return rv
}

// ResponseView represents a view to display HTTP response details.
type ResponseView struct {
	*tv.Flex
	response   *insomnium.Response
	Header     *ResponseBar
	Tabs       *ResponseTabs
	resetFocus func()
}

// Load initializes the layout and components of the response view.
func (rv *ResponseView) Load() {
	if rv.Header == nil {
		rv.LoadHeader()
	}
	if rv.Tabs == nil {
		rv.LoadTabs()
	}
	if rv.Flex == nil {
		rv.Flex = tv.NewFlex()
	}
	rv.Flex.Clear()
	rv.Flex.SetDirection(tv.FlexRow)
	rv.Flex.AddItem(rv.Header, 3, 0, false).
		AddItem(rv.Tabs, 0, 2, false)
}

// LoadHeader initializes the response header component.
func (rv *ResponseView) LoadHeader() {
	rv.Header = NewResponseBar()
	rv.Header.SetBorder(true)
	rv.Header.SetResponseFunc(rv.ResponseFunc)
}

// LoadTabs initlalizes the response tabs component.
func (rv *ResponseView) LoadTabs() {
	rv.Tabs = NewResponseTabs()
	rv.Tabs.OpenTab("body")
}

// ResponseFunc processes the response when a new one is set.
func (rv *ResponseView) ResponseFunc(res *insomnium.Response) {
	rv.SetResponse(res)
	if rv.resetFocus != nil {
		rv.resetFocus()
	}
}

// SetResetFocusFunc sets the function to reset the focus.
func (rv *ResponseView) SetResetFocusFunc(fn func()) *ResponseView {
	rv.resetFocus = fn
	rv.Tabs.SetResetFocusFunc(fn)
	return rv
}

// SetTagFunc sets a custom function for rendering tags in the header.
func (rv *ResponseView) SetTagFunc(fn TagFunc) *ResponseView {
	rv.Header.SetTagFunc(fn)
	return rv
}

// SetTagStyle sets the style for tags in the response header.
func (rv *ResponseView) SetTagStyle(style tc.Style) *ResponseView {
	rv.Header.SetTagStyle(style)
	return rv
}

// SetResponseOptionSelectedStyle sets the style for selected response options.
func (rv *ResponseView) SetResponseOptionSelectedStyle(style tc.Style) *ResponseView {
	rv.Header.SetResponseOptionSelectedStyle(style)
	return rv
}

// SetSelectedTagStyle sets the style for the selected tag in the header.
func (rv *ResponseView) SetSelectedTagStyle(style tc.Style) *ResponseView {
	rv.Header.SetSelectedTagStyle(style)
	return rv
}

// SetResponse sets the response data for the view.
func (rv *ResponseView) SetResponse(res *insomnium.Response) *ResponseView {
	rv.response = res
	rv.Header.SetResponse(res)
	rv.Tabs.SetResponse(res)
	return rv
}

// SetResponses sets the available responses for a given request.
func (rv *ResponseView) SetResponses(req *insomnium.Request, r []*insomnium.Response) *ResponseView {
	rv.Header.SetResponses(req, r)
	return rv
}

// SetResponseStyleFn sets the function for customizing response styles.
func (rv *ResponseView) SetResponseStyleFn(fn func(*insomnium.Response) tc.Style) *ResponseView {
	rv.Header.SetResponseStyleFn(fn)
	return rv
}

// SetBackgroundColor sets the background color for the response view.
func (rv *ResponseView) SetBackgroundColor(bg tc.Color) *ResponseView {
	rv.Flex.SetBackgroundColor(bg)
	rv.Header.SetBackgroundColor(bg)
	rv.Tabs.SetBackgroundColor(bg)
	return rv
}

// SetForegroundColor sets the foreground color for the response view.
func (rv *ResponseView) SetForegroundColor(fg tc.Color) *ResponseView {
	rv.Flex.SetBorderColor(fg)
	rv.Header.SetForegroundColor(fg)
	return rv
}

// SetResponseElapsedTimeTagStyle sets the style for the elapsed time tag.
func (rv *ResponseView) SetResponseElapsedTimeTagStyle(style tc.Style) *ResponseView {
	rv.Header.SetResponseElapsedTimeTagStyle(style)
	return rv
}

// SetResponseBytesReadTagStyle sets the style for the bytes read tag.
func (rv *ResponseView) SetResponseBytesReadTagStyle(style tc.Style) *ResponseView {
	rv.Header.SetResponseBytesReadTagStyle(style)
	return rv
}

// SetResponseLastCallTagStyle sets the style for the last call tag.
func (rv *ResponseView) SetResponseLastCallTagStyle(style tc.Style) *ResponseView {
	rv.Header.SetResponseLastCallTagStyle(style)
	return rv
}

// SetResponseEmptyBarRune sets the rune for the empty progress bar.
func (rv *ResponseView) SetResponseEmptyBarRune(r rune) *ResponseView {
	rv.Header.SetResponseEmptyBarRune(r)
	return rv
}

// SetResponseFilledBarRune sets the rune for the filled progress bar.
func (rv *ResponseView) SetResponseFilledBarRune(r rune) *ResponseView {
	rv.Header.SetResponseFilledBarRune(r)
	return rv
}

// SetResponseEmptyBarStyle sets the style for the empty progress bar.
func (rv *ResponseView) SetResponseEmptyBarStyle(style tc.Style) *ResponseView {
	rv.Header.SetResponseEmptyBarStyle(style)
	return rv
}

// SetResponseFilledBarStyle sets the style for the filled progress bar.
func (rv *ResponseView) SetResponseFilledBarStyle(style tc.Style) *ResponseView {
	rv.Header.SetResponseFilledBarStyle(style)
	return rv
}

// SetStyleTextFunc sets the function that decorates a text with a given style.
func (rv *ResponseView) SetStyleTextFunc(fn func(string, tc.Style) string) *ResponseView {
	rv.Tabs.SetStyleTextFunc(fn)
	return rv
}

// SetGetResponseBodyFunc sets the function that returns a response body.
func (rv *ResponseView) SetGetResponseBodyFunc(fn func(*insomnium.Response) (string, error)) *ResponseView {
	rv.Tabs.SetGetResponseBodyFunc(fn)
	return rv
}

// SetResponseBodyTheme sets the
func (rv *ResponseView) SetResponseBodyTheme(theme string) *ResponseView {
	rv.Tabs.SetResponseBodyTheme(theme)
	return rv
}

// SetSelectedResponseTabColor sets the color of the selected tab in the
// response view.
func (rv *ResponseView) SetSelectedResponseTabColor(c tc.Color) *ResponseView {
	rv.Tabs.SetSelectedTabColor(c)
	return rv
}

// SetUnselectedResponseTabColor sets the color of the unselected tab in the
// response view.
func (rv *ResponseView) SetUnselectedResponseTabColor(c tc.Color) *ResponseView {
	rv.Tabs.SetUnselectedTabColor(c)
	return rv
}
