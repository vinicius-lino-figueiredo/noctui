package tui

import (
	"fmt"
	"main/internal/component"
	"time"
	"unicode/utf8"

	"github.com/dustin/go-humanize"
	tc "github.com/gdamore/tcell/v2"
	tv "github.com/rivo/tview"
	"github.com/vinicius-lino-figueiredo/insomnium"
)

// NewResponseBar creates and initializes a new ResponseBar instance.
func NewResponseBar() *ResponseBar {
	rb := &ResponseBar{}
	rb.Load()
	return rb
}

// ResponseBar is a horizontal UI component that displays summary info about
// HTTP responses, including status, elapsed time, size, and request metadata.
// It also allows switching between multiple responses.
type ResponseBar struct {
	*tv.Flex
	responses                   []*insomnium.Response
	response                    *insomnium.Response
	StatusCode                  *component.Tag
	Elapsed                     *component.Tag
	BytesRead                   *component.Tag
	box                         *tv.Box
	LastCall                    *component.Tag
	responseDropDown            *component.DropDown
	responseOptionStyle         tc.Style
	responseOptionSelectedStyle tc.Style
	ResponseStyleFn             func(*insomnium.Response) tc.Style
	tagFunc                     func(string, tc.Style, tc.Style) string
	NoResponseColor             tc.Color
	CreateTag                   TagFunc
	bg                          tc.Color
	fg                          tc.Color
	tagStyle                    tc.Style
	selectedTagStyle            tc.Style
	responseFn                  func(*insomnium.Response)
}

// Load initializes and lays out the components in the ResponseBar.
func (rb *ResponseBar) Load() {
	if rb.StatusCode == nil {
		rb.LoadStatusCode()
	}
	if rb.Elapsed == nil {
		rb.LoadElapsed()
	}
	if rb.BytesRead == nil {
		rb.LoadSize()
	}
	if rb.box == nil {
		rb.box = tv.NewBox()
	}
	if rb.LastCall == nil {
		rb.LastCall = component.NewTag()
	}
	if rb.responseDropDown == nil {
		rb.responseDropDown = component.NewDropDown(rb.LastCall)
		rb.responseDropDown.SetSelectedFunc(rb.DropDownFunc)
	}
	if rb.Flex == nil {
		rb.Flex = tv.NewFlex()
	}
	rb.Flex.Clear()
	rb.Flex.AddItem(rb.StatusCode, 0, 0, false)
	rb.Flex.AddItem(rb.Elapsed, 0, 0, false)
	rb.Flex.AddItem(rb.BytesRead, 0, 0, false)
	rb.Flex.AddItem(rb.box, 0, 1, false)
	rb.Flex.AddItem(rb.responseDropDown, 0, 0, false)
}

// LoadStatusCode initializes the status code tag component.
func (rb *ResponseBar) LoadStatusCode() {
	rb.StatusCode = component.NewTag()
}

// LoadElapsed initializes the elapsed time tag component.
func (rb *ResponseBar) LoadElapsed() {
	rb.Elapsed = component.NewTag()
}

// LoadSize initializes the bytes read tag component.
func (rb *ResponseBar) LoadSize() {
	rb.BytesRead = component.NewTag()
}

// DropDownFunc is triggered when a response is selected from the dropdown. It
// calls the configured response handler with the selected response.
func (rb *ResponseBar) DropDownFunc(_ int, o *component.DropDownOption) {
	if o == nil {
		return
	}

	if rb.responseFn == nil {
		return
	}

	res, ok := o.GetRef().(*insomnium.Response)
	if !ok {
		return
	}

	rb.responseFn(res)
}

// SetTagFunc sets the function used to render styled tags.
func (rb *ResponseBar) SetTagFunc(fn TagFunc) *ResponseBar {
	rb.tagFunc = fn
	return rb
}

// SetResponseFunc sets the callback function to handle response selection.
func (rb *ResponseBar) SetResponseFunc(fn func(*insomnium.Response)) *ResponseBar {
	rb.responseFn = fn
	return rb
}

// SetTagStyle sets the style used for response option tags.
func (rb *ResponseBar) SetTagStyle(style tc.Style) *ResponseBar {
	rb.tagStyle = style
	return rb
}

// SetResponseOptionSelectedStyle sets the style for selected dropdown items.
func (rb *ResponseBar) SetResponseOptionSelectedStyle(style tc.Style) *ResponseBar {
	rb.responseOptionSelectedStyle = style
	return rb
}

// SetSelectedTagStyle sets the style for selected tags in dropdowns.
func (rb *ResponseBar) SetSelectedTagStyle(style tc.Style) *ResponseBar {
	rb.selectedTagStyle = style
	return rb
}

// SetResponses populates the dropdown with response options for the request.
func (rb *ResponseBar) SetResponses(req *insomnium.Request, r []*insomnium.Response) *ResponseBar {
	rb.responses = r
	rb.responseDropDown.Clear()
	for _, res := range r {
		txt, sel := rb.CreateResponseOption(req, res)
		do := component.NewDropDownOption(txt, sel, res)
		rb.responseDropDown.AddOption(do)
	}
	return rb
}

// CreateResponseOption builds a pair of styled text options for display and
// selection in the dropdown list.
func (rb *ResponseBar) CreateResponseOption(req *insomnium.Request, res *insomnium.Response) (string, string) {
	statusCodeStyle := rb.ResponseStyleFn(res)

	statusCode := fmt.Sprint(res.StatusCode)
	if res.StatusCode == 0 && res.StatusMessage != "" {
		statusCode = res.StatusMessage
	}
	elapsed := time.Duration(int64(res.ElapsedTime)).String()
	size := humanize.Bytes(uint64(max(0, res.BytesContent)))
	URL := fmt.Sprintf("%s %s", req.Method, req.URL)

	txt := rb.tagFunc(statusCode, rb.responseOptionStyle, statusCodeStyle)
	txt += rb.tagFunc(URL, rb.responseOptionStyle, rb.tagStyle)
	txt += rb.tagFunc(elapsed, rb.responseOptionStyle, rb.tagStyle)
	txt += rb.tagFunc(size, rb.responseOptionStyle, rb.tagStyle)

	sel := rb.tagFunc(statusCode, rb.responseOptionSelectedStyle, statusCodeStyle)
	sel += rb.tagFunc(res.URL, rb.responseOptionSelectedStyle, rb.selectedTagStyle)
	sel += rb.tagFunc(elapsed, rb.responseOptionSelectedStyle, rb.selectedTagStyle)
	sel += rb.tagFunc(size, rb.responseOptionSelectedStyle, rb.selectedTagStyle)

	return txt, sel
}

// SetResponse updates the ResponseBar UI elements with the given response.
func (rb *ResponseBar) SetResponse(res *insomnium.Response) *ResponseBar {
	rb.response = res
	statusCode := ""
	elapsed := ""
	size := ""
	lastCall := ""
	if res != nil {
		statusCode = fmt.Sprint(res.StatusCode)
		elapsed = time.Duration(int64(res.ElapsedTime)).String()
		size = humanize.Bytes(uint64(max(0, res.BytesContent)))
		lastCall = humanize.RelTime(time.UnixMilli(int64(res.Created)), time.Now(), "ago", "from now")
		if res.StatusMessage != "" && res.StatusCode == 0 {
			statusCode = res.StatusMessage
		}
		rb.Flex.ResizeItem(rb.StatusCode, len(statusCode)+2, 0)
		rb.Flex.ResizeItem(rb.Elapsed, utf8.RuneCountInString(elapsed)+2, 0)
		rb.Flex.ResizeItem(rb.BytesRead, len(size)+2, 0)
		rb.Flex.ResizeItem(rb.responseDropDown, len(lastCall)+2, 0)
	} else {
		rb.Flex.ResizeItem(rb.StatusCode, 0, 0)
		rb.Flex.ResizeItem(rb.Elapsed, 0, 0)
		rb.Flex.ResizeItem(rb.BytesRead, 0, 0)
		rb.Flex.ResizeItem(rb.responseDropDown, 0, 0)
	}
	rb.StatusCode.SetText(statusCode)
	rb.Elapsed.SetText(elapsed)
	rb.BytesRead.SetText(size)
	rb.LastCall.SetText(lastCall)
	rb.StatusCode.SetTagStyle(rb.ResponseStyleFn(res))
	return rb
}

// SetResponseStyleFn sets the function to compute styles for response tags.
func (rb *ResponseBar) SetResponseStyleFn(fn func(*insomnium.Response) tc.Style) *ResponseBar {
	rb.ResponseStyleFn = fn
	return rb
}

// SetBackgroundColor sets the background color for all elements.
func (rb *ResponseBar) SetBackgroundColor(bg tc.Color) *ResponseBar {
	rb.Flex.SetBackgroundColor(bg)
	rb.StatusCode.SetBackgroundColor(bg)
	rb.Elapsed.SetBackgroundColor(bg)
	rb.BytesRead.SetBackgroundColor(bg)
	rb.box.SetBackgroundColor(bg)
	rb.LastCall.SetBackgroundColor(bg)
	return rb
}

// SetForegroundColor sets the border color for the ResponseBar.
func (rb *ResponseBar) SetForegroundColor(fg tc.Color) *ResponseBar {
	rb.Flex.SetBorderColor(fg)
	return rb
}

// SetResponseElapsedTimeTagStyle sets the style for the elapsed time tag.
func (rb *ResponseBar) SetResponseElapsedTimeTagStyle(style tc.Style) *ResponseBar {
	rb.Elapsed.SetTagStyle(style)
	return rb
}

// SetResponseBytesReadTagStyle sets the style for the bytes read tag.
func (rb *ResponseBar) SetResponseBytesReadTagStyle(style tc.Style) *ResponseBar {
	rb.BytesRead.SetTagStyle(style)
	return rb
}

// SetResponseLastCallTagStyle sets the style for the last call tag.
func (rb *ResponseBar) SetResponseLastCallTagStyle(style tc.Style) *ResponseBar {
	rb.LastCall.SetTagStyle(style)
	return rb
}

// SetResponseEmptyBarRune sets the rune used for the empty scrollbar bar.
func (rb *ResponseBar) SetResponseEmptyBarRune(r rune) *ResponseBar {
	rb.responseDropDown.SetEmptyBarRune(r)
	return rb
}

// SetResponseFilledBarRune sets the rune used for the filled scrollbar bar.
func (rb *ResponseBar) SetResponseFilledBarRune(r rune) *ResponseBar {
	rb.responseDropDown.SetFilledBarRune(r)
	return rb
}

// SetResponseEmptyBarStyle sets the style for the empty scrollbar segment.
func (rb *ResponseBar) SetResponseEmptyBarStyle(style tc.Style) *ResponseBar {
	rb.responseDropDown.SetEmptyBarStyle(style)
	return rb
}

// SetResponseFilledBarStyle sets the style for the filled scrollbar segment.
func (rb *ResponseBar) SetResponseFilledBarStyle(style tc.Style) *ResponseBar {
	rb.responseDropDown.SetFilledBarStyle(style)
	return rb
}
