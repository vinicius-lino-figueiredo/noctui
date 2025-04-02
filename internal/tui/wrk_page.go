package tui

import (
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
	inso            *insomnium.Insomnium
	WrkTree         *WrkTree
	LeftPanel       *tv.Flex
	MiddlePanel     *ReqEditor
	RightPanel      *tv.Flex
	request         *insomnium.Request
	NeutralizeFocus func()
}

// GetReqTreeElement returns the tree view element.
func (wp *WrkPage) GetReqTreeElement() tv.Primitive {
	return wp.WrkTree.GetReqTreeElement()
}

// GetUrlElement returns the element that is focused to edit the request url.
func (wp *WrkPage) GetUrlElement() tv.Primitive {
	return wp.MiddlePanel.GetUrlElement()
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
	wp.RightPanel = tv.NewFlex()
	wp.RightPanel.SetBorder(true)
}

// SetRequest sets the middle pannel request.
func (wp *WrkPage) SetRequest(req *insomnium.Request) {
	wp.request = req
	wp.MiddlePanel.SetRequest(req)
}

// GetRequest returns the current loaded request.
func (wp *WrkPage) GetRequest() *insomnium.Request {
	return wp.request
}

// SetWrkUpdates the page with a new workspace.
func (wp *WrkPage) SetWrk(wrk *insomnium.Workspace) {
	wp.WrkTree.SetWrk(wrk)
}

// SetNeutralizeFocusFunc sets the func that is called to reset the app's focus.
func (wp *WrkPage) SetNeutralizeFocusFunc(fn func()) *WrkPage {
	wp.NeutralizeFocus = fn
	wp.WrkTree.SetNeutralizeFocusFunc(fn)
	wp.MiddlePanel.SetNeutralizeFocusFunc(fn)
	return wp
}

// SetFocusFunc sets a function that is called to set focus on an element.
func (wp *WrkPage) SetFocusFunc(fn func(tv.Primitive)) *WrkPage {
	wp.MiddlePanel.SetFocusFunc(fn)
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

// SetUrlInputFieldStyle sets the url input field style.
func (wp *WrkPage) SetUrlInputFieldStyle(style tc.Style) *WrkPage {
	wp.MiddlePanel.SetUrlInputFieldStyle(style)
	return wp
}

// SetMethod sets the http methods options in the method dropdown element.
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
	return wp
}
