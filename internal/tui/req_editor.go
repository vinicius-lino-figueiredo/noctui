package tui

import (
	"main/internal/component"

	tc "github.com/gdamore/tcell/v2"
	tv "github.com/rivo/tview"
	"github.com/vinicius-lino-figueiredo/insomnium"
)

// NewReqEditor creates a new instance of the request editor widget.
func NewReqEditor() *ReqEditor {
	re := &ReqEditor{
		Flex: tv.NewFlex(),
	}
	re.CreateContent()
	return re
}

// ReqEditor is a tview.Primitive used to view and edit a request. It contains a
// URL editor element.
type ReqEditor struct {
	*tv.Flex
	URLEditor          *URLEditor
	tabs               *ReqTabs
	bg                 tc.Color
	fg                 tc.Color
	MatrixInputCapture func(*component.Matrix) InputFn
}

// CreateContent initializes the layout and adds child elements. It creates the
// UrlEditor if it is nil.
func (re *ReqEditor) CreateContent() {
	re.Flex.Clear()
	re.Flex.SetDirection(tv.FlexRow).
		SetBackgroundColor(re.bg)

	if re.URLEditor == nil {
		re.CreateURLEditor()
	}
	if re.tabs == nil {
		re.CreateTabs()
	}

	re.Flex.AddItem(re.URLEditor, 3, 0, false).
		AddItem(re.tabs, 0, 1, false)
}

// CreateURLEditor initializes the URL editor and sets its border.
func (re *ReqEditor) CreateURLEditor() {
	re.URLEditor = NewURLEditor()
	re.URLEditor.SetBorder(true)
}

// CreateTabs loads the tabs widget.
func (re *ReqEditor) CreateTabs() {
	re.tabs = NewReqTabs()
}

// GetURLElement returns the tview.Primitive used for URL editing.
func (re *ReqEditor) GetURLElement() tv.Primitive {
	return re.URLEditor.GetURLElement()
}

// GetMethodElement returns the tview.Primitive used for method editing.
func (re *ReqEditor) GetMethodElement() tv.Primitive {
	return re.URLEditor.GetMethodElement()
}

// SetRequest updates the request data displayed in the editor.
func (re *ReqEditor) SetRequest(req *insomnium.Request) *ReqEditor {
	re.SetReqTitle(req)
	re.URLEditor.SetRequest(req)
	re.tabs.SetRequest(req)
	return re
}

// SetReqTitle updates the editor title with the request name. If req is nil,
// the title is cleared.
func (re *ReqEditor) SetReqTitle(req *insomnium.Request) *ReqEditor {
	if req == nil {
		re.Flex.SetTitle("")
	} else {
		re.Flex.SetTitle(req.Name)
	}
	return re
}

// SetEnvironmentVarStyle sets the style for environment variables.
func (re *ReqEditor) SetEnvironmentVarStyle(style tc.Style) *ReqEditor {
	re.URLEditor.SetEnvironmentVarStyle(style)
	return re
}

// SetTemplateVarStyle sets the style for template variables.
func (re *ReqEditor) SetTemplateVarStyle(style tc.Style) *ReqEditor {
	re.URLEditor.SetTemplateVarStyle(style)
	return re
}

// SetURLInputFieldStyle sets the base style for the URL input widget.
func (re *ReqEditor) SetURLInputFieldStyle(style tc.Style) *ReqEditor {
	re.URLEditor.SetURLInputFieldStyle(style)
	return re
}

// SetMethodStyleFunc sets a function that returns a style for an HTTP method.
func (re *ReqEditor) SetMethodStyleFunc(fn func(string) tc.Style) *ReqEditor {
	re.URLEditor.SetMethodStyleFunc(fn)
	return re
}

// SetBackgroundColor sets the background color for the editor.
func (re *ReqEditor) SetBackgroundColor(bg tc.Color) *ReqEditor {
	re.bg = bg
	re.Flex.SetBackgroundColor(bg)
	re.URLEditor.SetBackgroundColor(bg)
	re.tabs.SetBackgroundColor(bg)
	return re
}

// SetForegroundColor sets the foreground color for the editor.
func (re *ReqEditor) SetForegroundColor(fg tc.Color) *ReqEditor {
	re.bg = fg
	re.Flex.SetBorderColor(fg)
	re.Flex.SetTitleColor(fg)
	re.URLEditor.SetForegroundColor(fg)
	re.tabs.SetForegroundColor(fg)
	return re
}

// SetSelectedReqTabColor sets the color of the selected tab in the request
// editor tabs.
func (re *ReqEditor) SetSelectedReqTabColor(c tc.Color) *ReqEditor {
	re.tabs.SetSelectedTabColor(c)
	return re
}

// SetUnselectedReqTabColor sets the color of the unselected tabs in the request
// editor tabs.
func (re *ReqEditor) SetUnselectedReqTabColor(c tc.Color) *ReqEditor {
	re.tabs.SetUnselectedTabColor(c)
	return re
}

// SetResetFocusFunc sets the function to reset the app's focus.
func (re *ReqEditor) SetResetFocusFunc(fn func()) *ReqEditor {
	re.URLEditor.SetResetFocusFunc(fn)
	re.tabs.SetResetFocusFunc(fn)
	return re
}

// SetMethods sets the HTTP method options in the dropdown.
func (re *ReqEditor) SetMethods(methods []string) *ReqEditor {
	re.URLEditor.SetMethods(methods)
	return re
}

// SetMethodStyles sets the styles for selected and unselected methods.
func (re *ReqEditor) SetMethodStyles(unselected, selected tc.Style) *ReqEditor {
	re.URLEditor.SetMethodStyles(unselected, selected)
	return re
}

// SetTagFunc sets a function that formats text with styles into a tag.
func (re *ReqEditor) SetTagFunc(fn TagFunc) *ReqEditor {
	re.URLEditor.SetTagFunc(fn)
	return re
}

// SetHeaderInputStyle sets the header input field style.
func (re *ReqEditor) SetHeaderInputStyle(style tc.Style) *ReqEditor {
	re.tabs.SetHeaderInputStyle(style)
	return re
}

// SetStyleTextFunc sets the function that styles a given text.
func (re *ReqEditor) SetStyleTextFunc(fn func(string, tc.Style) string) *ReqEditor {
	re.tabs.SetStyleTextFunc(fn)
	return re
}

// SetRequestBodyInputStyle sets the request body input style.
func (re *ReqEditor) SetRequestBodyInputStyle(style tc.Style) *ReqEditor {
	re.tabs.SetRequestBodyInputStyle(style)
	return re
}

// SetRequestBodyTheme sets the request body highlight theme.
func (re *ReqEditor) SetRequestBodyTheme(theme string) *ReqEditor {
	re.tabs.SetRequestBodyTheme(theme)
	return re
}

// TagFunc defines a function type that formats a text string using two styles
// and returns a formatted tag.
type TagFunc func(string, tc.Style, tc.Style) string
