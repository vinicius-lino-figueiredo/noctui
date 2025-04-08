package tui

import (
	"slices"

	tc "github.com/gdamore/tcell/v2"
	tv "github.com/rivo/tview"
	"github.com/vinicius-lino-figueiredo/insomnium"
)

// NewWrkTree creates a widget that contains a *tview.TreeView, representing a
// set of RequestGroups and Requests within a workspace.
func NewWrkTree(inso *insomnium.Insomnium) *WrkTree {
	workspaceTree := &WrkTree{
		TreeView: tv.NewTreeView(),
		root:     tv.NewTreeNode(""),
		inso:     inso,
	}
	workspaceTree.
		SetRoot(workspaceTree.root).
		SetTopLevel(1).
		SetChangedFunc(workspaceTree.treeChangedFunc).
		SetFocusFunc(workspaceTree.focusFunc).
		SetInputCapture(workspaceTree.inputCapture)
	return workspaceTree
}

// WrkTree represents a tree structure that displays the hierarchy of
// RequestGroups and Requests within a workspace.
type WrkTree struct {
	*tv.TreeView
	root              *tv.TreeNode
	inso              *insomnium.Insomnium
	bg                tc.Color
	fg                tc.Color
	textStyle         tc.Style
	selectedTextStyle tc.Style
	NeutralizeFocus   func()
	OpenRequest       func(*insomnium.Request)
}

// GetReqTreeElement returns the tree view element.
func (wt *WrkTree) GetReqTreeElement() tv.Primitive {
	return wt.TreeView
}

// focusFunc is called when the tree view receives focus.
func (wt *WrkTree) focusFunc() {
	wt.TreeView.SetCurrentNode(wt.root)
}

// inputCapture resets focus when escape key or q is pressed.
func (wt *WrkTree) inputCapture(event *tc.EventKey) *tc.EventKey {
	if event.Key() == tc.KeyEsc || event.Rune() == 'q' {
		wt.TreeView.SetCurrentNode(nil)
		if wt.OpenRequest != nil {
			wt.OpenRequest(nil)
		}
		if wt.NeutralizeFocus != nil {
			wt.NeutralizeFocus()
		}
		return nil
	}
	return event
}

// SetNeutralizeFocusFunc sets a func that is called to reset the app's focus.
func (wt *WrkTree) SetNeutralizeFocusFunc(fn func()) *WrkTree {
	wt.NeutralizeFocus = fn
	return wt
}

// treeChangedFunc is called when the selected tree node changes. When the new
// selected node is a request node, the the OpenRequest function is called.
func (wt *WrkTree) treeChangedFunc(node *tv.TreeNode) {
	ref := node.GetReference()
	if req, ok := ref.(*insomnium.Request); ok {
		if wt.OpenRequest != nil {
			wt.OpenRequest(req)
		}
	}
}

// SetWrk reloads the TreeView with the provided workspace data. If the
// workspace is nil, it clears the tree view, effectively resetting its content.
func (wt *WrkTree) SetWrk(wrk *insomnium.Workspace) {
	wt.root.ClearChildren()
	if wrk == nil {
		return
	}
	nodes := map[string]*tv.TreeNode{wrk.ID: wt.root}
	groups := slices.Clone(wt.inso.RequestGroups)
	unused := make([]insomnium.RequestGroup, 0, len(wt.inso.RequestGroups))
	var modified int
	for {
		modified = 0
		for _, g := range groups {
			parent, ok := nodes[g.ParentID]
			if !ok {
				unused = append(unused, g)
				continue
			}
			node := wt.NewGroupTreeNode(&g)
			nodes[g.ID] = node
			parent.AddChild(node)
			modified++
		}
		if modified == 0 {
			break
		}
		groups = groups[:len(unused)]
		copy(groups, unused)
		unused = unused[:0]
	}
	for _, req := range wt.inso.Requests {
		parent, ok := nodes[req.ParentID]
		if !ok {
			continue
		}
		node := wt.NewRequestTreeNode(&req)
		parent.AddChild(node)
	}
}

// SetBackgroundColor sets the background color for the tree view.
func (wt *WrkTree) SetBackgroundColor(bg tc.Color) *WrkTree {
	wt.bg = bg
	wt.textStyle = tc.StyleDefault.Background(bg).Foreground(wt.fg)
	wt.selectedTextStyle = tc.StyleDefault.Background(wt.fg).Foreground(bg)
	wt.TreeView.SetBackgroundColor(bg)
	return wt
}

// SetForegroundColor sets the foreground color for the tree view.
func (wt *WrkTree) SetForegroundColor(fg tc.Color) *WrkTree {
	wt.fg = fg
	wt.SetGraphicsColor(fg)
	wt.textStyle = tc.StyleDefault.Background(wt.bg).Foreground(fg)
	wt.selectedTextStyle = tc.StyleDefault.Background(fg).Foreground(wt.bg)
	wt.UpdateNodesStyle()
	return wt
}

// SetSelectedTextStyle sets the text style for selected nodes.
func (wt *WrkTree) SetSelectedTextStyle(style tc.Style) *WrkTree {
	wt.selectedTextStyle = style
	wt.UpdateNodesStyle()
	return wt
}

// UpdateNodesStyle updates the style of all nodes in the tree.
func (wt *WrkTree) UpdateNodesStyle() *WrkTree {
	lookup := []*tv.TreeNode{wt.GetRoot()}
	newLookup := []*tv.TreeNode{}
	for len(lookup) > 0 {
		newLookup = newLookup[:0]
		for _, node := range lookup {
			newLookup = append(newLookup, node.GetChildren()...)
			node.SetTextStyle(wt.textStyle).
				SetSelectedTextStyle(wt.selectedTextStyle)
		}
		lookup = make([]*tv.TreeNode, len(newLookup))
		copy(lookup, newLookup)
	}
	return wt
}

// NewGroupTreeNode creates a new tree node representing a RequestGroup. The
// node displays the group's name and is associated with the given RequestGroup
// object.
func (wt *WrkTree) NewGroupTreeNode(wrk *insomnium.RequestGroup) *tv.TreeNode {
	node := tv.NewTreeNode(wrk.Name).
		SetReference(wrk).
		SetSelectedFunc(wt.ReqGroupNodeSelectedFunc).
		SetTextStyle(wt.textStyle).
		SetSelectedTextStyle(wt.selectedTextStyle)
	return node
}

// NewRequestTreeNode creates a new tree node representing a Request. The node
// displays the request's name and is associated with the given Request object.
func (wt *WrkTree) NewRequestTreeNode(req *insomnium.Request) *tv.TreeNode {
	requestNode := tv.NewTreeNode(req.Name)
	requestNode.
		SetReference(req).
		SetSelectedFunc(wt.ReqNodeSelectedFunc).
		SetTextStyle(wt.textStyle).
		SetSelectedTextStyle(wt.selectedTextStyle)
	return requestNode
}

// ReqGroupNodeSelectedFunc is called when a request group tree node is
// selected. It changes the open/close state of the node.
func (wt *WrkTree) ReqGroupNodeSelectedFunc() {
	currNode := wt.TreeView.GetCurrentNode()
	currNode.SetExpanded(!currNode.IsExpanded())
}

// ReqNodeSelectedFunc is called when a request tree node is selected. It then
// sets focus back to the original position.
func (wt *WrkTree) ReqNodeSelectedFunc() {
	if wt.NeutralizeFocus != nil {
		wt.NeutralizeFocus()
	}
}

// SetOpenRequestFunc sets the function thaw is called when the current focused
// node changes.
func (wt *WrkTree) SetOpenRequestFunc(fn func(*insomnium.Request)) *WrkTree {
	wt.OpenRequest = fn
	return wt
}
