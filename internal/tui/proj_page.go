package tui

import (
	"main/internal/component"

	tc "github.com/gdamore/tcell/v2"
	tv "github.com/rivo/tview"
	"github.com/vinicius-lino-figueiredo/insomnium"
)

// NewProjPage returns a *ProjPage that is used to manage the project and
// workspace selection screen.
func NewProjPage(app *tv.Application, inso *insomnium.Insomnium) *ProjPage {
	pp := &ProjPage{
		Flex:       tv.NewFlex(),
		ProjMatrix: component.NewMatrix(),
		WrkMatrix:  component.NewMatrix(),
		app:        app,
		inso:       inso,
	}

	pp.ProjMatrix.
		SetWidth(1).
		SetHeight(6).
		SetBorder(true).
		SetInputCapture(pp.ProjMatrixInputCapture)
	pp.WrkMatrix.
		SetWidth(5).
		SetHeight(5).
		SetBorder(true).
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
	openWorkspace      func(wrk *insomnium.Workspace)
	app                *tv.Application
	bg                 tc.Color
	fg                 tc.Color
	ProjMatrix         *component.Matrix
	WrkMatrix          *component.Matrix
	ResetFocus         func()
	MatrixInputCapture func(*component.Matrix) InputFn
	inso               *insomnium.Insomnium
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
			SetSelectedFunc(pp.SelectProjBtn).
			SetBorder(true).
			SetFocusFunc(pp.FocusProjBtn(proj.ID))
		pp.ProjMatrix.AddItem(projectButton)
	}
	pp.ProjMatrix.Refresh()
}

// CreateProjBtn generates a button that selects the given project.
func (pp *ProjPage) CreateProjBtn(proj insomnium.Project) *tv.Button {
	btn := tv.NewButton(proj.Name)
	btn.SetStyle(BtnStyle).
		SetLabelColor(pp.bg).
		SetActivatedStyle(ActiveBtnStyle).
		SetSelectedFunc(pp.SelectProjBtn).
		SetBorder(true).
		SetFocusFunc(pp.FocusProjBtn(proj.ID)).
		SetBackgroundColor(pp.bg)

	return btn
}

// CreateWrkBtn generates a button that selects the given workspace.
func (pp *ProjPage) CreateWrkBtn(wrk insomnium.Workspace) *tv.Button {
	btn := tv.NewButton(wrk.Name)
	btn.SetStyle(BtnStyle).
		SetLabelColor(pp.fg).
		SetActivatedStyle(ActiveBtnStyle).
		SetSelectedFunc(pp.WrkBtnSelFunc(&wrk)).
		SetBorder(true).
		SetBackgroundColor(pp.bg).
		SetBorderColor(pp.fg)
	return btn
}

// SelectProjBtn is triggered when a project button is selected. It attempts to
// set focus on the first workspace in the workspace matrix. If no workspace is
// available, it calls the ResetFocus function to shift focus away, typically
// returning it to the command input.
func (pp *ProjPage) SelectProjBtn() {
	fb, ok := pp.WrkMatrix.Get(0, 0)
	if !ok {
		pp.ResetFocus()
	} else {
		pp.app.SetFocus(fb)
	}
}

// FocusProjBtn returns a func that is called when the app sets focus to the
// given button. The function returned repopulates the workspaces grid with
// workspaces whose ParentID is equal to the ID of the focused button.
func (pp *ProjPage) FocusProjBtn(id string) func() {
	return func() {
		pp.PopulateWorkspaces(id)
	}
}

// ProjMatrixInputCapture returns a input capture func that will wait for a
// Escape or 'q' key press and will reload all the workspaces. It will ignore
// any other input and pass it to the default matrix input capture function.
func (pp *ProjPage) ProjMatrixInputCapture(event *tc.EventKey) *tc.EventKey {
	if event.Key() == tc.KeyEsc || event.Rune() == 'q' {
		pp.PopulateWorkspaces("")
		if pp.ResetFocus != nil {
			pp.ResetFocus()
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
		pp.ResetFocus()
		return nil
	}
	fn := pp.MatrixInputCapture(pp.WrkMatrix)
	return fn(event)
}

// PopulateWorkspaces uses an id to populate the workspace matrix with
// workspaces whose ParentID equals the given id. When "" is passed,
// the matrix is populated with all workspaces.
func (pp *ProjPage) PopulateWorkspaces(id string) {
	pp.WrkMatrix.Clear()
	for _, wrk := range pp.inso.Workspaces {
		if id != "" && id != wrk.ParentID {
			continue
		}
		btn := pp.CreateWrkBtn(wrk)
		pp.WrkMatrix.AddItem(btn)
	}
	pp.WrkMatrix.Refresh()
}

// WrkBtnSelFunc is called when a workspace is selected.
func (pp *ProjPage) WrkBtnSelFunc(wrk *insomnium.Workspace) func() {
	return func() {
		if pp.ResetFocus != nil {
			pp.ResetFocus()
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

// SetResetFocusFunc sets a function that will be called to neutralize
// the focus on the ProjPage.
func (pp *ProjPage) SetResetFocusFunc(fn func()) *ProjPage {
	pp.ResetFocus = fn
	return pp
}

// SetBackgroundColor sets the page's background color.
func (pp *ProjPage) SetBackgroundColor(bg tc.Color) *ProjPage {
	pp.bg = bg
	for _, p := range pp.WrkMatrix.GetAll() {
		b := p.(*tv.Button)
		b.SetBackgroundColor(bg)
	}
	for _, p := range pp.ProjMatrix.GetAll() {
		b := p.(*tv.Button)
		b.SetBackgroundColor(bg)
	}
	pp.WrkMatrix.SetBackgroundColor(bg)
	pp.ProjMatrix.SetBackgroundColor(bg)
	return pp
}

// SetForegroundColor sets the page's foreground color.
func (pp *ProjPage) SetForegroundColor(fg tc.Color) *ProjPage {
	pp.fg = fg
	for _, p := range pp.WrkMatrix.GetAll() {
		b := p.(*tv.Button)
		b.SetBorderColor(fg)
		b.SetLabelColor(fg)
	}
	for _, p := range pp.ProjMatrix.GetAll() {
		b := p.(*tv.Button)
		b.SetBorderColor(fg)
		b.SetLabelColor(fg)
	}
	pp.WrkMatrix.SetForegroundColor(fg)
	pp.ProjMatrix.SetForegroundColor(fg)
	return pp
}

// SetMatrixInputCapture sets a function that creates input capture function for
// the given matrix.
func (pp *ProjPage) SetMatrixInputCapture(fn func(*component.Matrix) InputFn) *ProjPage {
	pp.MatrixInputCapture = fn
	return pp
}
