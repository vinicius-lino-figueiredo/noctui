package tui

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	tc "github.com/gdamore/tcell/v2"
	tv "github.com/rivo/tview"
	"github.com/vinicius-lino-figueiredo/insomnium"
)

// NewURLEditor creates an element used to display and edit a request URL and
// HTTP method.
func NewURLEditor() *URLEditor {
	ue := &URLEditor{
		Flex: tv.NewFlex(),
	}
	ue.CreateElements()
	return ue
}

// URLEditor is a tview.Primitive that holds a URL editor (a *tview.Pages with a
// *tview.TextView to display a colored URL and a *tview.InputField to edit the
// URL) and an HTTP method selector (a *tview.DropDown).
type URLEditor struct {
	*tv.Flex
	methods                []string
	method                 *tv.DropDown
	pages                  *tv.Pages
	urlView                *MaskedTextView
	urlInput               *tv.InputField
	inputFieldStyle        tc.Style
	EnvironmentVarStyle    tc.Style
	TemplateVarStyle       tc.Style
	bg                     tc.Color
	fg                     tc.Color
	tagFunc                func(string, tc.Style, tc.Style) string
	getDropDownMethodStyle func(string) tc.Style
	focus                  func(tv.Primitive)
	ResetFocus             func()
}

// CreateElements loads instances for the tview elements.
func (ue *URLEditor) CreateElements() {
	if ue.Flex == nil {
		ue.Flex = tv.NewFlex()
	}
	ue.Flex.Clear()
	ue.CreateMethod()
	ue.CreatePages()

	urlPaddingFlex := tv.NewFlex()
	urlPaddingFlex.
		AddItem(ue.pages, 0, 1, false).
		SetBorderPadding(0, 0, 1, 1)
	ue.Flex.AddItem(ue.method, 5, 0, false).
		AddItem(urlPaddingFlex, 0, 1, false).
		SetBorderPadding(0, 0, 1, 0)
}

// CreateMethod loads the method *tview.DropDown element.
func (ue *URLEditor) CreateMethod() {
	if ue.method == nil {
		ue.method = tv.NewDropDown()
	}
	ue.method.SetFocusFunc(ue.methodFocus).
		SetInputCapture(ue.methodCapture)
}

// CreatePages loads the URL *tview.TextView/*tview.InputField page elements.
func (ue *URLEditor) CreatePages() {
	if ue.pages == nil {
		ue.pages = tv.NewPages()
	}
	ue.CreateURLView()
	ue.CreateURLInput()
	for _, page := range ue.pages.GetPageNames(false) {
		ue.pages.RemovePage(page)
	}
	ue.pages.AddPage("view", ue.urlView, true, true).
		AddPage("input", ue.urlInput, true, false)
}

// CreateURLView loads the *MaskedTextView element.
func (ue *URLEditor) CreateURLView() {
	if ue.pages == nil {
		ue.pages = tv.NewPages()
	}
	ue.urlView = NewMaskedTextView()
	ue.urlView.SetMask(ue.viewMask).
		SetDynamicColors(true)
}

// CreateURLInput loads the *tview.InputField element.
func (ue *URLEditor) CreateURLInput() {
	if ue.urlInput == nil {
		ue.urlInput = tv.NewInputField()
	}
	ue.urlInput.SetChangedFunc(ue.inputChanged).
		SetDoneFunc(ue.inputDone).
		SetFocusFunc(ue.inputFocus).
		SetBlurFunc(ue.inputBlur)
}

// inputFocus is called when the *tview.InputField receives focus and switches
// the URL page to it, hiding the *MaskedTextView.
func (ue *URLEditor) inputFocus() {
	ue.pages.SwitchToPage("input")
}

// inputChanged is called when the value of the *tview.InputField changes and
// applies the changes in the *MaskedTextView too.
func (ue *URLEditor) inputChanged(text string) {
	ue.urlView.SetText(text)
}

// inputDone is called when the user finishes writing in the URL
// *tview.InputField widget and releases the app's focus.
func (ue *URLEditor) inputDone(_ tc.Key) {
	if ue.ResetFocus != nil {
		ue.ResetFocus()
	}
}

// inputBlur is called when the *tview.InputField element loses focus and sets
// the *MaskedTextView URL element visible.
func (ue *URLEditor) inputBlur() {
	ue.pages.SwitchToPage("view")
}

// methodFocus is called when the method *tview.DropDown receives focus and
// opens the method DropDown list.
func (ue *URLEditor) methodFocus() {
	ue.openMethodDropDown()
}

// methodCapture is called on keyboard input while the method *tview.DropDown
// widget is on focus. it uses 'j' and 'k' as motions and avoids other inputs
func (ue *URLEditor) methodCapture(event *tc.EventKey) *tc.EventKey {
	r := event.Rune()
	k := event.Key()
	switch {
	case k == tc.KeyEnter:
		return event
	case k == tc.KeyDown:
		return event
	case r == 'j':
		return tc.NewEventKey(tc.KeyDown, 0, event.Modifiers())
	case k == tc.KeyUp:
		return event
	case r == 'k':
		return tc.NewEventKey(tc.KeyUp, 0, event.Modifiers())
	case r == 'q' || event.Key() == tc.KeyEsc:
		ue.ResetFocus()
		return nil
	default:
		return nil
	}
}

// GetURLElement returns the URL text input element.
func (ue *URLEditor) GetURLElement() tv.Primitive {
	return ue.urlInput
}

// GetMethodElement returns the method dropdown element.
func (ue *URLEditor) GetMethodElement() tv.Primitive {
	return ue.method
}

// SetRequest sets the current selected request.
func (ue *URLEditor) SetRequest(req *insomnium.Request) {
	if req == nil {
		ue.urlInput.SetText("")
		ue.SetMethod("")
		return
	}
	ue.urlInput.SetText(req.URL)
	ue.SetMethod(req.Method)
}

// SetURLInputFieldStyle sets the base style for the URL input widget.
func (ue *URLEditor) SetURLInputFieldStyle(style tc.Style) *URLEditor {
	ue.inputFieldStyle = style
	ue.urlInput.SetFieldStyle(style)
	return ue
}

// SetMethod sets the current method at the *tview.DropDown element.
func (ue *URLEditor) SetMethod(method string) {
	ue.updateMethodStyle(method)
	// unsetting the selection func because it would unfocus the dropdown
	ue.method.SetSelectedFunc(nil)
	ue.method.SetCurrentOption(slices.Index(ue.methods, method))
	ue.method.SetSelectedFunc(ue.methodSelected)
}

// updateMethodStyle resizes the element in the flex so the content fits in and
// reapplies the style with the corresponding method style.
func (ue *URLEditor) updateMethodStyle(method string) {
	ue.Flex.ResizeItem(ue.method, len(method), 0)
	if ue.getDropDownMethodStyle != nil {
		style := ue.getDropDownMethodStyle(method)
		fg, bg, _ := style.Decompose()
		ue.method.SetFieldTextColor(fg).
			SetFieldBackgroundColor(bg).
			SetLabelColor(tc.ColorDarkRed)
	}
}

// SetEnvironmentVarStyle sets the style used to highlight an environment
// variable in the URL text view.
func (ue *URLEditor) SetEnvironmentVarStyle(style tc.Style) *URLEditor {
	ue.EnvironmentVarStyle = style
	return ue
}

// SetTemplateVarStyle sets the style used to highlight an template variable.
func (ue *URLEditor) SetTemplateVarStyle(style tc.Style) *URLEditor {
	ue.TemplateVarStyle = style
	return ue
}

// SetForegroundColor sets the foreground color for the url editor.
func (ue *URLEditor) SetForegroundColor(fg tc.Color) *URLEditor {
	ue.fg = fg
	ue.Flex.SetBorderColor(fg)
	return ue
}

// viewMask is a text mask that displays environment variables and template
// variables as highlighted tags.
func (ue *URLEditor) viewMask(text string) string {
	envVars := regexp.MustCompile(`{{\s?_\.\w+\s?}}`)
	templVars := regexp.MustCompile(`{%\s?\w+\s'\w+'.*\s?%}`)
	text = envVars.ReplaceAllStringFunc(text, func(s string) string {
		s = strings.Trim(s, "{ }")
		s = s[2:]
		s = ue.tagFunc(s, ue.inputFieldStyle, ue.EnvironmentVarStyle)
		return s
	})
	text = templVars.ReplaceAllStringFunc(text, func(s string) string {
		s = strings.Trim(s, "{ }%")
		words := strings.Split(s, " ")
		s = fmt.Sprintf("%s", words[0])
		s = ue.tagFunc(s, ue.inputFieldStyle, ue.TemplateVarStyle)
		return s
	})
	return text
}

// SetMethodStyleFunc sets a function that is called to receive an HTTP method
// name and return a style for that method.
func (ue *URLEditor) SetMethodStyleFunc(fn func(string) tc.Style) *URLEditor {
	ue.getDropDownMethodStyle = fn
	_, method := ue.method.GetCurrentOption()
	ue.updateMethodStyle(method)
	return ue
}

// SetFocusFunc sets a function that is called to set focus on an element.
func (ue *URLEditor) SetFocusFunc(fn func(tv.Primitive)) *URLEditor {
	ue.focus = fn
	return ue
}

// SetResetFocusFunc sets the function that is called to reset the app
// focus.
func (ue *URLEditor) SetResetFocusFunc(fn func()) *URLEditor {
	ue.ResetFocus = fn
	return ue
}

// SetMethods sets the HTTP methods options in the method dropdown element.
func (ue *URLEditor) SetMethods(methods []string) *URLEditor {
	ue.methods = methods
	ue.method.SetOptions(methods, nil)
	return ue
}

// SetMethodStyles sets the style for the selected and the unselected methods in
// the *tview.DropDown list.
func (ue *URLEditor) SetMethodStyles(unselected, selected tc.Style) *URLEditor {
	ue.method.SetListStyles(unselected, selected)
	return ue
}

// SetTagFunc sets a function that formats text with styles into a tag.
func (ue *URLEditor) SetTagFunc(fn TagFunc) *URLEditor {
	ue.tagFunc = fn
	return ue
}

// methodSelected is called when a method is selected in the *tview.DropDown and
// it updates the method field style and resets the app's focus.
func (ue *URLEditor) methodSelected(method string, _ int) {
	ue.updateMethodStyle(method)
	if ue.ResetFocus != nil {
		ue.ResetFocus()
	}
}

// openMethodDropDown currently simulates an Enter key press to open the HTTP
// method dropdown list. At the moment, this is necessary because the dropdown's
// openList function is unexported and can only be triggered through the
// MouseHandler or InputHandler functions.
func (ue *URLEditor) openMethodDropDown() {
	if ue.method.IsOpen() {
		return
	}
	ih := ue.method.InputHandler()
	ih(tc.NewEventKey(tc.KeyEnter, rune(tc.KeyEnter), tc.ModNone), ue.focus)
}
