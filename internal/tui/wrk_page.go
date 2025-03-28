package tui

import (
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
	MiddlePanel     *tv.Flex
	RightPanel      *tv.Flex
	NeutralizeFocus func()
}

// GetReqTreeElement returns the tree view element.
func (wp *WrkPage) GetReqTreeElement() tv.Primitive {
	return wp.WrkTree.GetReqTreeElement()
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
	wp.WrkTree = NewWrkTree(wp.inso)
	wp.LeftPanel = tv.NewFlex()
	wp.LeftPanel.
		AddItem(wp.WrkTree, 0, 1, false).
		SetBorder(true)
}

// CreateMiddlePanel creates the *tview.Flex that holds the request editor.
func (wp *WrkPage) CreateMiddlePanel() {
	wp.MiddlePanel = tv.NewFlex()
	wp.MiddlePanel.SetBorder(true)
}

// CreateRightPanel creates the *tview.Flex that holds the response viewer.
func (wp *WrkPage) CreateRightPanel() {
	wp.RightPanel = tv.NewFlex()
	wp.RightPanel.SetBorder(true)
}

// SetWrkUpdates the page with a new workspace.
func (wp *WrkPage) SetWrk(wrk *insomnium.Workspace) {
	wp.WrkTree.SetWrk(wrk)
}

// SetNeutralizeFocusFunc sets the func that is called to reset the app focus.
func (wp *WrkPage) SetNeutralizeFocusFunc(fn func()) *WrkPage {
	wp.NeutralizeFocus = fn
	wp.WrkTree.SetNeutralizeFocusFunc(fn)
	return wp
}
