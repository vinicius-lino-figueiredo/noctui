package tui

import (
	tc "github.com/gdamore/tcell/v2"
	tv "github.com/rivo/tview"
	"github.com/vinicius-lino-figueiredo/insomnium"
)

// NewProjPage returns a *ProjPage that is used to manage the project and
// workspace selection screen.
func NewProjPage(app *tv.Application, inso *insomnium.Insomnium) *ProjPage {
	pp := &ProjPage{
		Flex:       tv.NewFlex(),
		ProjMatrix: NewMatrix(),
		WrkMatrix:  NewMatrix(),
		app:        app,
		inso:       inso,
	}

	pp.ProjMatrix.
		SetWidth(1).
		SetHeight(6).
		SetBorder(true).
		SetFocusFunc(pp.MatrixOnFocus(pp.ProjMatrix)).
		SetInputCapture(pp.ProjMatrixInputCapture)
	pp.WrkMatrix.
		SetWidth(5).
		SetHeight(5).
		SetBorder(true).
		SetFocusFunc(pp.MatrixOnFocus(pp.WrkMatrix)).
		SetInputCapture(pp.WrkMatrixInputCapture)

	pp.AddItem(pp.ProjMatrix, 0, 1, false).
		AddItem(pp.WrkMatrix, 0, 5, false)
	pp.PopulateProjects()
	pp.PopulateWorkspaces("")

	return pp
}

// ProjPage represents a page in the application responsible for displaying
// and managing the selection of a project and workspace.
type ProjPage struct {
	*tv.Flex
	openWorkspace   func(wrk *insomnium.Workspace)
	app             *tv.Application
	ProjMatrix      *Matrix
	WrkMatrix       *Matrix
	NeutralizeFocus func()
	inso            *insomnium.Insomnium
}

// GetProjElement returns the project matrix.
func (pp *ProjPage) GetProjElement() tv.Primitive {
	return pp.ProjMatrix
}

// GetWrkElement returns the workspace matrix.
func (pp *ProjPage) GetWrkElement() tv.Primitive {
	return pp.WrkMatrix
}

// PopulateProjects loads the project grid with all the projects found in the
// *insomnium.Insomnium global instance.
func (pp *ProjPage) PopulateProjects() {
	for _, proj := range pp.inso.Projects {
		projectButton := pp.CreateProjBtn(proj)
		projectButton.SetStyle(BtnStyle).
			SetActivatedStyle(ActiveBtnStyle).
			SetSelectedFunc(pp.SelectProjBtn(proj.ID)).
			SetBorder(true).
			SetFocusFunc(pp.FocusProjBtn(proj.ID))
		pp.ProjMatrix.itms = append(pp.ProjMatrix.itms, projectButton)
	}
	pp.ProjMatrix.Refresh()
}

// CreateProjBtn generates a button that selects the given project.
func (pp *ProjPage) CreateProjBtn(proj insomnium.Project) *tv.Button {
	btn := tv.NewButton(proj.Name)
	btn.SetStyle(BtnStyle).
		SetActivatedStyle(ActiveBtnStyle).
		SetSelectedFunc(pp.SelectProjBtn(proj.ID)).
		SetBorder(true).
		SetFocusFunc(pp.FocusProjBtn(proj.ID))
	return btn
}

// SelectProjBtn is called when a project button is selected and it sets focus
// on the fist workspace in the grid. If there is no workspace in that grid, it
// sets focus to the command input field.
func (pp *ProjPage) SelectProjBtn(id string) func() {
	return func() {
		fb, ok := pp.WrkMatrix.Get(0, 0)
		if !ok {
			pp.NeutralizeFocus()
		} else {
			pp.app.SetFocus(fb)
		}
	}
}

// FocsProjBtn returns a func that is called when the app sets focus to the
// given button. The function returned repopulates the workspaces grid with
// workspaces whose ParentID is equal to the ID of the focused button.
func (pp *ProjPage) FocusProjBtn(id string) func() {
	return func() {
		pp.PopulateWorkspaces(id)
	}
}

// MatrixOnFocus returns a func to be used as a callback for a Matrix focus.
func (pp *ProjPage) MatrixOnFocus(m *Matrix) func() {
	return func() {
		if len(m.itms) != 0 {
			pp.app.SetFocus(m.itms[0])
			m.currX, m.currY = 0, 0
		}
	}
}

// ProjMatrixInputCapture returns a input capture func that will wait for a
// Escape or 'q' key press and will reload all the workspaces. It will ignore
// any other input and pass it to the default matrix input capture function.
func (pp *ProjPage) ProjMatrixInputCapture(event *tc.EventKey) *tc.EventKey {
	if event.Key() == tc.KeyEsc || event.Rune() == 'q' {
		pp.PopulateWorkspaces("")
		if pp.NeutralizeFocus != nil {
			pp.NeutralizeFocus()
		}
		return nil
	}
	fn := pp.MatrixInputCapture(pp.ProjMatrix)
	return fn(event)
}

// WrkMatrixInputCapture returns a input capture func that will wait for a
// Escape or 'q' key press and will reload all the workspaces. It will ignore
// any other input and pass it to the default matrix input capture function.
func (pp *ProjPage) WrkMatrixInputCapture(event *tc.EventKey) *tc.EventKey {
	if event.Key() == tc.KeyEsc || event.Rune() == 'q' {
		pp.PopulateWorkspaces("")
		pp.NeutralizeFocus()
		return nil
	}
	fn := pp.MatrixInputCapture(pp.WrkMatrix)
	return fn(event)
}

// PopulateWorkspaces uses an id to populate the workspace matrix with
// workspaces whose ParentID equals the given id. When "" is passed,
// the matrix is populated with all workspaces.
func (pp *ProjPage) PopulateWorkspaces(id string) {
	pp.WrkMatrix.itms = []tv.Primitive{}
	for _, wrk := range pp.inso.Workspaces {
		if id != "" && id != wrk.ParentID {
			continue
		}
		btn := tv.NewButton(wrk.Name)
		btn.SetStyle(BtnStyle).
			SetActivatedStyle(ActiveBtnStyle).
			SetSelectedFunc(pp.WrkBtnSelFunc(&wrk)).
			SetBorder(true)
		pp.WrkMatrix.itms = append(pp.WrkMatrix.itms, btn)
	}
	pp.WrkMatrix.Refresh()
}

// MatrixInputCapture will return an input capture function and will accept
// vim motions to control the selected element in the grid.
func (pp *ProjPage) MatrixInputCapture(m *Matrix) InputFn {
	return func(event *tc.EventKey) *tc.EventKey {
		switch {
		case event.Key() == tc.KeyLeft || event.Rune() == 'h':
			m.Left()
		case event.Key() == tc.KeyRight || event.Rune() == 'l':
			m.Right()
		case event.Key() == tc.KeyUp || event.Rune() == 'k':
			m.Up()
		case event.Key() == tc.KeyDown || event.Rune() == 'j':
			m.Down()
		default:
			return event
		}
		p := m.GetCurrentPrimitive()
		if p != nil {
			pp.app.SetFocus(p)
		}
		return nil
	}
}

// WrkBtnSelFunc is called when a workspace is selected.
func (pp *ProjPage) WrkBtnSelFunc(wrk *insomnium.Workspace) func() {
	return func() {
		if pp.NeutralizeFocus != nil {
			pp.NeutralizeFocus()
		}
		if pp.openWorkspace != nil {
			pp.openWorkspace(wrk)
		}
	}
}

// SetWrkFunc sets a function that will be called when a workspace
// button is pressed.
func (pp *ProjPage) SetWrkFunc(fn func(wrk *insomnium.Workspace)) *ProjPage {
	pp.openWorkspace = fn
	return pp
}

// SetNeutralizeFocusFunc sets a function that will be called to neutralize
// the focus on the ProjPage.
func (pp *ProjPage) SetNeutralizeFocusFunc(fn func()) *ProjPage {
	pp.NeutralizeFocus = fn
	return pp
}
