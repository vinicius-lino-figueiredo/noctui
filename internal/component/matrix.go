package component

import (
	"slices"

	tc "github.com/gdamore/tcell/v2"
	tv "github.com/rivo/tview"
)

// NewMatrix creates and returns a new Matrix instance that consists of a grid
// for displaying content and two navigation buttons (up and down) arranged
// vertically using a Flex container. The buttons are styled and positioned
// above and below the grid. This function initializes the layout, setting up
// the Flex container with buttons and a grid, but it doesn't handle the logic
// related to button actions or the dynamic content of the grid, which are
// managed separately.
func NewMatrix() *Matrix {
	grid := tv.NewGrid()
	buttonUp := tv.NewButton("△")
	buttonDown := tv.NewButton("▽")
	flex := tv.NewFlex().
		SetDirection(tv.FlexRow).
		AddItem(buttonUp, 0, 0, false).
		AddItem(grid, 0, 1, false).
		AddItem(buttonDown, 0, 0, false)
	m := &Matrix{
		h:          1,
		w:          1,
		Flex:       flex,
		Grid:       grid,
		ButtonUp:   buttonUp,
		ButtonDown: buttonDown,
		boxes:      []*tv.Box{},
	}
	m.Refresh()
	return m
}

// Matrix represents a UI component consisting of a grid (Grid) and two
// navigation buttons (up and down). The buttons serve to indicate whether there
// is more content above or below the grid, adjusting the visibility of the
// content as the user navigates. The structure keeps track of the selected
// position (currX, currY), the number of rows and columns (w, h), the list of
// items to be displayed (itms), and the skip variable, which manages the offset
// of the visible content in the grid.
type Matrix struct {
	*tv.Flex
	Grid         *tv.Grid
	ButtonUp     *tv.Button
	ButtonDown   *tv.Button
	bg           tc.Color
	fg           tc.Color
	boxes        []*tv.Box
	currX, currY int
	skip         int
	w, h         int
	itms         []tv.Primitive
}

// Focus implements the tv.Primitive interface by delegating focus to the first
// item in the matrix, if available. This allows the Matrix to be used in tview
// layouts that manage focus between components.
func (m *Matrix) Focus(delegate func(tv.Primitive)) {
	if len(m.itms) > 0 {
		delegate(m.itms[0])
		m.currX, m.currY = 0, 0
	}
}

// SetWidth sets the number of columns in the matrix grid.
func (m *Matrix) SetWidth(w int) *Matrix {
	m.w = w
	return m
}

// SetHeight sets the number of rows in the matrix grid.
func (m *Matrix) SetHeight(h int) *Matrix {
	m.h = h
	return m
}

// Refresh reloads and reorders the grid elements.
func (m *Matrix) Refresh() {
	m.Grid.Clear()
	m.boxes = m.boxes[:0]
	for n := m.skip * m.w; n < (m.h*m.w)+m.skip*m.w; n++ {
		var itm tv.Primitive
		if n < len(m.itms) {
			itm = m.itms[n]
		} else {
			b := tv.NewBox().
				SetBackgroundColor(m.bg)
			itm = b
			m.boxes = append(m.boxes, b)
		}
		x := n % m.w
		y := n / m.w
		m.Grid.AddItem(itm, y-m.skip, x, 1, 1, 1, 1, false)
	}
	if m.skip > 0 {
		m.Flex.ResizeItem(m.ButtonUp, 1, 0)
	} else {
		m.Flex.ResizeItem(m.ButtonUp, 0, 0)
	}
	a := ((len(m.itms) + m.w - 1) / m.w) - m.skip - m.h
	if a > 0 {
		m.Flex.ResizeItem(m.ButtonDown, 1, 0)
	} else {
		m.Flex.ResizeItem(m.ButtonDown, 0, 0)
	}
}

// RefreshY checks if the selected cell is within the visible range. If it is
// not, it adjusts the range so the element can be displayed and then refreshes
// the grid.
func (m *Matrix) RefreshY() {
	if m.currY > m.skip+m.h-1 {
		m.skip++
		m.Refresh()
	} else if m.currY < m.skip {
		m.skip--
		m.Refresh()
	}
}

// Left moves the selected x position one cell to the left. If the current row
// is the last one, it means the current x position might be beyond the last
// element in the list, because the last row is the only one that may not be
// fully populated. In this case, it adjusts the x position to refer to the last
// element in the list.
func (m *Matrix) Left() {
	if m.currX > 0 {
		m.currX--
	}
	lastRow := len(m.itms) / m.w
	lastRowLimit := (len(m.itms) % m.w) - 1
	if m.currY == lastRow {
		m.currX = min(m.currX, lastRowLimit)
	}
}

// Right moves the selected x position one cell to the right. If the new cell is
// beyond the current row limit, it does nothing.
func (m *Matrix) Right() {
	lastRow := len(m.itms) / m.w
	currLineLimit := m.w
	if m.currY == lastRow {
		currLineLimit = (len(m.itms) % m.w) - 1
	}
	if m.currX < currLineLimit {
		m.currX++
	}
}

// Up moves the selected y position one cell up, unless it would result in a
// negative y position. It also calls the function that updates the range, which
// adjusts it if necessary and refreshes the grid.
func (m *Matrix) Up() {
	if m.currY > 0 {
		m.currY--
	}
	m.RefreshY()
}

// Down moves the selected y position one cell down, unless it would exceed the
// last row. It also calls the function that updates the range, which adjusts it
// if necessary and refreshes the grid.
func (m *Matrix) Down() {
	lastRow := len(m.itms) / m.w
	if m.currY < lastRow {
		m.currY++
	}
	m.RefreshY()
}

// Get returns the item for the given location. If there is no item, it returns
// nil and the returned boolean is set to false.
func (m *Matrix) Get(x, y int) (tv.Primitive, bool) {
	if len(m.itms) == 0 {
		return nil, false
	}
	i := (y * m.w) + x
	if i >= len(m.itms) {
		if i <= m.w*((len(m.itms)+m.w-1)/m.w) {
			return m.itms[len(m.itms)-1], true
		}
		return nil, false
	}
	return m.itms[i], true
}

// GetCurrentPrimitive returns the currently selected cell. If the position
// refers to a blank spot in the grid, it adjusts it accordingly.
func (m *Matrix) GetCurrentPrimitive() tv.Primitive {
	totalRows := (len(m.itms) + m.w - 1) / m.w
	m.currX = max(m.currX, 0)
	m.currX = min(m.currX, m.w-1)
	m.currY = max(m.currY, 0)
	m.currY = min(m.currY, totalRows-1)
	curr, _ := m.Get(m.currX, m.currY)
	return curr
}

// GetAll returns all the items displayed in the matrix.
func (m *Matrix) GetAll() []tv.Primitive {
	res := make([]tv.Primitive, len(m.itms))
	copy(res, m.itms)
	return res
}

// SetBackgroundColor sets the background color for the matrix.
func (m *Matrix) SetBackgroundColor(bg tc.Color) *Matrix {
	m.bg = bg
	for _, box := range m.boxes {
		box.SetBackgroundColor(bg)
	}
	m.Flex.SetBackgroundColor(bg)
	m.ButtonUp.SetStyle(tc.StyleDefault.Foreground(m.fg).Background(bg))
	m.ButtonDown.SetStyle(tc.StyleDefault.Foreground(m.fg).Background(bg))
	return m
}

// SetForegroundColor sets the foreground color for the url matrix.
func (m *Matrix) SetForegroundColor(fg tc.Color) *Matrix {
	m.fg = fg
	for _, box := range m.boxes {
		box.SetBorderColor(fg)
	}
	m.Flex.SetBorderColor(fg)
	m.ButtonUp.SetStyle(tc.StyleDefault.Foreground(fg).Background(m.bg))
	m.ButtonDown.SetStyle(tc.StyleDefault.Foreground(fg).Background(m.bg))
	return m
}

// GetBackgroundColor returns the matrix background color.
func (m *Matrix) GetBackgroundColor() tc.Color {
	return m.bg
}

// GetForegroundColor returns the matrix foreground color.
func (m *Matrix) GetForegroundColor() tc.Color {
	return m.fg
}

// GetItems returns all the matrix items.
func (m *Matrix) GetItems() []tv.Primitive {
	return slices.Clone(m.itms)
}

// AddItem adds an item to the matrix.
func (m *Matrix) AddItem(item tv.Primitive) *Matrix {
	m.itms = append(m.itms, item)
	return m
}

// Clear removes all the matrix items.
func (m *Matrix) Clear() *Matrix {
	m.itms = m.itms[:0]
	return m
}
