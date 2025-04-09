package tui

import (
	"slices"

	tc "github.com/gdamore/tcell/v2"
	tv "github.com/rivo/tview"
)

// NewTabs creates a new instance of *Tabs.
func NewTabs() *Tabs {
	t := &Tabs{}
	t.Load()
	return t
}

// Tabs is a tview.Primitive that holds multiple pages and has a *tview.Table
// that is used to display the tab names.
type Tabs struct {
	*tv.Flex
	header             *tv.Table
	content            *tv.Pages
	bg                 tc.Color
	fg                 tc.Color
	selectedTabColor   tc.Color
	unselectedTabColor tc.Color
	tabs               []string
	focus              func(tv.Primitive)
}

// Load loads the widgets.
func (t *Tabs) Load() {
	if t.header == nil {
		t.LoadHeader()
	}
	if t.content == nil {
		t.LoadContent()
	}
	if t.Flex == nil {
		t.Flex = tv.NewFlex().SetDirection(tv.FlexRow)
	}
	t.Flex.Clear()
	t.AddItem(t.header, 3, 0, false).
		AddItem(t.content, 0, 1, false)
}

// LoadHeader loads the *tview.Table.
func (t *Tabs) LoadHeader() {
	t.header = tv.NewTable()
	t.header.SetBorders(true)
}

// LoadContent loads the *tview.Pages.
func (t *Tabs) LoadContent() {
	t.content = tv.NewPages()
}

// AddTab adds a tab by adding a *tview.TableCell to the *tview.Table and a
// tview.Primitive to the *tview.Pages.
func (t *Tabs) AddTab(name string, widget tv.Primitive) *Tabs {
	if !slices.Contains(t.tabs, name) {
		t.tabs = append(t.tabs, name)
		cell := tv.NewTableCell(name).
			SetExpansion(1).
			SetAlign(tv.AlignCenter)
		t.header.SetCell(0, len(t.tabs)-1, cell)
	}

	t.content.AddPage(name, widget, true, false)
	return t
}

// OpenTab switches to the given tab.
func (t *Tabs) OpenTab(name string) *Tabs {
	if !slices.Contains(t.tabs, name) {
		return t
	}
	curr, _ := t.content.GetFrontPage()
	i := slices.Index(t.tabs, curr)
	cell := t.header.GetCell(0, i)
	cell.SetTextColor(t.unselectedTabColor)

	i = slices.Index(t.tabs, name)
	cell = t.header.GetCell(0, i)
	cell.SetTextColor(t.selectedTabColor)

	t.content.SwitchToPage(name)
	return t
}

// SetForegroundColor sets the widget foreground color.
func (t *Tabs) SetForegroundColor(fg tc.Color) *Tabs {
	t.fg = fg
	t.header.SetBordersColor(fg)
	return t
}

// SetBackgroundColor sets the widget background color.
func (t *Tabs) SetBackgroundColor(bg tc.Color) *Tabs {
	t.bg = bg
	t.header.SetBackgroundColor(bg)
	return t
}

// SetSelectedTabColor sets the color of the selected tab.
func (t *Tabs) SetSelectedTabColor(c tc.Color) *Tabs {
	t.selectedTabColor = c
	curr, _ := t.content.GetFrontPage()
	i := slices.Index(t.tabs, curr)
	cell := t.header.GetCell(0, i)
	cell.SetTextColor(c)
	return t
}

// SetUnselectedTabColor sets the color of the unselected tabs.
func (t *Tabs) SetUnselectedTabColor(c tc.Color) *Tabs {
	t.unselectedTabColor = c
	curr, _ := t.content.GetFrontPage()
	i := slices.Index(t.tabs, curr)
	for n := range t.header.GetColumnCount() + 1 {
		if n == i {
			continue
		}
		cell := t.header.GetCell(0, n)
		cell.SetTextColor(c)
	}
	return t
}
