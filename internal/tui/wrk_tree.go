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
		SetFocusFunc(workspaceTree.focusFunc).
		SetInputCapture(workspaceTree.inputCapture)
	return workspaceTree
}

// WrkTree represents a tree structure that displays the hierarchy of
// RequestGroups and Requests within a workspace.
type WrkTree struct {
	*tv.TreeView
	root            *tv.TreeNode
	inso            *insomnium.Insomnium
	NeutralizeFocus func()
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
		wt.SetCurrentNode(nil)
		if wt.NeutralizeFocus != nil {
			wt.NeutralizeFocus()
		}
		return nil
	}
	return event
}

// SetNeutralizeFocusFunc sets a func that is called to reset the app focus.
func (wt *WrkTree) SetNeutralizeFocusFunc(fn func()) *WrkTree {
	wt.NeutralizeFocus = fn
	return wt
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

// NewGroupTreeNode creates a new tree node representing a RequestGroup. The
// node displays the group's name and is associated with the given RequestGroup
// object.
func (wt *WrkTree) NewGroupTreeNode(wrk *insomnium.RequestGroup) *tv.TreeNode {
	node := tv.NewTreeNode(wrk.Name).
		SetReference(wrk)
	return node
}

// NewRequestTreeNode creates a new tree node representing a Request. The node
// displays the request's name and is associated with the given Request object.
func (wt *WrkTree) NewRequestTreeNode(req *insomnium.Request) *tv.TreeNode {
	node := tv.NewTreeNode(req.Name).
		SetReference(req)
	return node
}
