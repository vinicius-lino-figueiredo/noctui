package tui

import (
	"cmp"
	"main/internal/component"
	"slices"

	tc "github.com/gdamore/tcell/v2"
	tv "github.com/rivo/tview"
	"github.com/vinicius-lino-figueiredo/insomnium"
)

// NewWrkPage creates and returns a new instance of WrkPage, initializing its
// structure with the reference to the Insomnium application.
func NewWrkPage(inso *insomnium.Insomnium) *WrkPage {
	wp := &WrkPage{
		inso: inso,
	}
	wp.CreateFlex()
	return wp
}

// WrkPage represents the main workspace layout, organizing the interface into
// three panels (LeftPanel, MiddlePanel, RightPanel) within a Flex container.
type WrkPage struct {
	*tv.Flex
	inso               *insomnium.Insomnium
	WrkTree            *WrkTree
	LeftPanel          *tv.Flex
	MiddlePanel        *ReqEditor
	RightPanel         *ResponseView
	request            *insomnium.Request
	bg                 tc.Color
	fg                 tc.Color
	ResetFocus         func()
	MatrixInputCapture func(*component.Matrix) InputFn
}

// GetReqTreeElement returns the tree view element.
func (wp *WrkPage) GetReqTreeElement() tv.Primitive {
	return wp.WrkTree.GetReqTreeElement()
}

// GetURLElement returns the element that is focused to edit the request url.
func (wp *WrkPage) GetURLElement() tv.Primitive {
	return wp.MiddlePanel.GetURLElement()
}

// GetMethodElement returns the element that receives focus when the user tries
// to edit the curren request method
func (wp *WrkPage) GetMethodElement() tv.Primitive {
	return wp.MiddlePanel.GetMethodElement()
}

// CreateFlex creates the main *tview.Flex that stores the other elements.
func (wp *WrkPage) CreateFlex() {
	wp.CreateLeftPanel()
	wp.CreateMiddlePanel()
	wp.CreateRightPanel()
	wp.Flex = tv.NewFlex().
		AddItem(wp.LeftPanel, 0, 5, false).
		AddItem(wp.MiddlePanel, 0, 9, false).
		AddItem(wp.RightPanel, 0, 9, false)
}

// CreateLeftPanel creates the *tview.Flex that holds the workspace tree view.
func (wp *WrkPage) CreateLeftPanel() {
	wp.WrkTree = NewWrkTree(wp.inso).
		SetOpenRequestFunc(wp.SetRequest)
	wp.LeftPanel = tv.NewFlex()
	wp.LeftPanel.
		AddItem(wp.WrkTree, 0, 1, false).
		SetBorder(true)
}

// CreateMiddlePanel creates the *tview.Flex that holds the request editor.
func (wp *WrkPage) CreateMiddlePanel() {
	wp.MiddlePanel = NewReqEditor()
	wp.MiddlePanel.SetBorder(true)
}

// CreateRightPanel creates the *tview.Flex that holds the response viewer.
func (wp *WrkPage) CreateRightPanel() {
	wp.RightPanel = NewResponseView()
	wp.RightPanel.SetBorder(true)
}

// SetRequest sets the middle pannel request.
func (wp *WrkPage) SetRequest(req *insomnium.Request) {
	wp.request = req
	wp.MiddlePanel.SetRequest(req)
	r := wp.getResponses(req)
	res := wp.getOldestRes(r)
	wp.RightPanel.SetResponses(req, r)
	wp.RightPanel.SetResponse(res)
}

func (wp *WrkPage) getResponses(req *insomnium.Request) []*insomnium.Response {
	if req == nil {
		return nil
	}
	responses := make([]*insomnium.Response, 0, len(wp.inso.Responses))
	for _, res := range wp.inso.Responses {
		if res.ParentID == req.ID {
			responses = append(responses, &res)
		}
	}
	return responses
}

func (wp *WrkPage) getOldestRes(r []*insomnium.Response) *insomnium.Response {
	if len(r) == 0 {
		return nil
	}
	cmpFn := func(a, b *insomnium.Response) int {
		return cmp.Compare(a.Modified, b.Modified)
	}
	return slices.MaxFunc(r, cmpFn)
}

// GetRequest returns the current loaded request.
func (wp *WrkPage) GetRequest() *insomnium.Request {
	return wp.request
}

// SetWrk updates the page with a new workspace.
func (wp *WrkPage) SetWrk(wrk *insomnium.Workspace) {
	wp.WrkTree.SetWrk(wrk)
}

// SetResetFocusFunc sets the func that is called to reset the app's focus.
func (wp *WrkPage) SetResetFocusFunc(fn func()) *WrkPage {
	wp.ResetFocus = fn
	wp.WrkTree.SetResetFocusFunc(fn)
	wp.MiddlePanel.SetResetFocusFunc(fn)
	wp.RightPanel.SetResetFocusFunc(fn)
	return wp
}

// SetEnvironmentVarStyle sets the style for the env url tag.
func (wp *WrkPage) SetEnvironmentVarStyle(style tc.Style) *WrkPage {
	wp.MiddlePanel.SetEnvironmentVarStyle(style)
	return wp
}

// SetTemplateVarStyle sets the value for the url tags other than the envs.
func (wp *WrkPage) SetTemplateVarStyle(style tc.Style) *WrkPage {
	wp.MiddlePanel.SetTemplateVarStyle(style)
	return wp
}

// SetURLInputFieldStyle sets the url input field style.
func (wp *WrkPage) SetURLInputFieldStyle(style tc.Style) *WrkPage {
	wp.MiddlePanel.SetURLInputFieldStyle(style)
	return wp
}

// SetMethods sets the http methods options in the method dropdown element.
func (wp *WrkPage) SetMethods(methods []string) *WrkPage {
	wp.MiddlePanel.SetMethods(methods)
	return wp
}

// SetMethodStyles sets the style for the selected and the unselected methods in
// the *tview.DropDown list.
func (wp *WrkPage) SetMethodStyles(unselected, selected tc.Style) *WrkPage {
	wp.MiddlePanel.SetMethodStyles(unselected, selected)
	return wp
}

// SetMethodStyleFunc sets a function that is called to receive a HTTP method
// name and return a style for that method.
func (wp *WrkPage) SetMethodStyleFunc(fn func(string) tc.Style) *WrkPage {
	wp.MiddlePanel.SetMethodStyleFunc(fn)
	return wp
}

// SetTagFunc sets a function that formats text with styles into a tag.
func (wp *WrkPage) SetTagFunc(fn TagFunc) *WrkPage {
	wp.MiddlePanel.SetTagFunc(fn)
	wp.RightPanel.SetTagFunc(fn)
	return wp
}

// SetTagStyle sets the style for tags in the response header.
func (wp *WrkPage) SetTagStyle(style tc.Style) *WrkPage {
	wp.RightPanel.SetTagStyle(style)
	return wp
}

// SetResponseOptionSelectedStyle sets the style for selected response options.
func (wp *WrkPage) SetResponseOptionSelectedStyle(style tc.Style) *WrkPage {
	wp.RightPanel.SetResponseOptionSelectedStyle(style)
	return wp
}

// SetSelectedTagStyle sets the style for the selected tag in the header.
func (wp *WrkPage) SetSelectedTagStyle(style tc.Style) *WrkPage {
	wp.RightPanel.SetSelectedTagStyle(style)
	return wp
}

// SetBackgroundColor sets the background color for the page.
func (wp *WrkPage) SetBackgroundColor(bg tc.Color) *WrkPage {
	wp.bg = bg
	wp.LeftPanel.SetBackgroundColor(bg)
	wp.WrkTree.SetBackgroundColor(bg)
	wp.MiddlePanel.SetBackgroundColor(bg)
	wp.RightPanel.SetBackgroundColor(bg)
	return wp
}

// SetForegroundColor sets the foreground color for the page.
func (wp *WrkPage) SetForegroundColor(fg tc.Color) *WrkPage {
	wp.fg = fg
	wp.LeftPanel.SetBorderColor(fg)
	wp.WrkTree.SetForegroundColor(fg)
	wp.MiddlePanel.SetForegroundColor(fg)
	wp.RightPanel.SetForegroundColor(fg)
	return wp
}

// SetSelectedReqTabColor sets the color of the selected tab in the request
// editor tabs.
func (wp *WrkPage) SetSelectedReqTabColor(c tc.Color) *WrkPage {
	wp.MiddlePanel.SetSelectedReqTabColor(c)
	return wp
}

// SetUnselectedReqTabColor sets the color of the unselected tabs in the request
// editor tabs.
func (wp *WrkPage) SetUnselectedReqTabColor(c tc.Color) *WrkPage {
	wp.MiddlePanel.SetUnselectedReqTabColor(c)
	return wp
}

// SetSelectedResponseTabColor sets the color of the selected tab in the request
// editor tabs.
func (wp *WrkPage) SetSelectedResponseTabColor(c tc.Color) *WrkPage {
	wp.RightPanel.SetSelectedResponseTabColor(c)
	return wp
}

// SetUnselectedResponseTabColor sets the color of the unselected tabs in the
// request editor tabs.
func (wp *WrkPage) SetUnselectedResponseTabColor(c tc.Color) *WrkPage {
	wp.RightPanel.SetUnselectedResponseTabColor(c)
	return wp
}

// SetHeaderInputStyle sets the style for the input fields.
func (wp *WrkPage) SetHeaderInputStyle(style tc.Style) *WrkPage {
	wp.MiddlePanel.SetHeaderInputStyle(style)
	return wp
}

// SetStyleTextFunc sets a function that is called to style a text with font
// and background colors.
func (wp *WrkPage) SetStyleTextFunc(fn func(string, tc.Style) string) *WrkPage {
	wp.MiddlePanel.SetStyleTextFunc(fn)
	wp.RightPanel.SetStyleTextFunc(fn)
	return wp
}

// SetGetResponseBodyFunc sets the function that returns a response body.
func (rv *WrkPage) SetGetResponseBodyFunc(fn func(*insomnium.Response) (string, error)) *WrkPage {
	rv.RightPanel.SetGetResponseBodyFunc(fn)
	return rv
}

// SetRequestBodyInputStyle sets the style of the request body editor text area.
func (wp *WrkPage) SetRequestBodyInputStyle(style tc.Style) *WrkPage {
	wp.MiddlePanel.SetRequestBodyInputStyle(style)
	return wp
}

// SetRequestBodyTheme sets the highlighting theme of the request body editor
// text area.
func (wp *WrkPage) SetRequestBodyTheme(theme string) *WrkPage {
	wp.MiddlePanel.SetRequestBodyTheme(theme)
	return wp
}

// SetResponseBodyTheme sets the highlighting theme of the response body viewer.
func (wp *WrkPage) SetResponseBodyTheme(theme string) *WrkPage {
	wp.RightPanel.SetResponseBodyTheme(theme)
	return wp
}

// SetResponseStyleFn sets the function for customizing response styles.
func (wp *WrkPage) SetResponseStyleFn(fn func(*insomnium.Response) tc.Style) *WrkPage {
	wp.RightPanel.SetResponseStyleFn(fn)
	return wp
}

// SetResponseElapsedTimeTagStyle sets the style for the elapsed time tag.
func (wp *WrkPage) SetResponseElapsedTimeTagStyle(style tc.Style) *WrkPage {
	wp.RightPanel.SetResponseElapsedTimeTagStyle(style)
	return wp
}

// SetResponseBytesReadTagStyle sets the style for the bytes read tag.
func (wp *WrkPage) SetResponseBytesReadTagStyle(style tc.Style) *WrkPage {
	wp.RightPanel.SetResponseBytesReadTagStyle(style)
	return wp
}

// SetResponseLastCallTagStyle sets the style for the last call tag.
func (wp *WrkPage) SetResponseLastCallTagStyle(style tc.Style) *WrkPage {
	wp.RightPanel.SetResponseLastCallTagStyle(style)
	return wp
}

// SetResponseEmptyBarRune sets the rune for the empty progress bar.
func (wp *WrkPage) SetResponseEmptyBarRune(r rune) *WrkPage {
	wp.RightPanel.SetResponseEmptyBarRune(r)
	return wp
}

// SetResponseFilledBarRune sets the rune for the filled progress bar.
func (wp *WrkPage) SetResponseFilledBarRune(r rune) *WrkPage {
	wp.RightPanel.SetResponseFilledBarRune(r)
	return wp
}

// SetResponseEmptyBarStyle sets the style for the empty progress bar.
func (wp *WrkPage) SetResponseEmptyBarStyle(style tc.Style) *WrkPage {
	wp.RightPanel.SetResponseEmptyBarStyle(style)
	return wp
}

// SetResponseFilledBarStyle sets the style for the filled progress bar.
func (wp *WrkPage) SetResponseFilledBarStyle(style tc.Style) *WrkPage {
	wp.RightPanel.SetResponseFilledBarStyle(style)
	return wp
}
