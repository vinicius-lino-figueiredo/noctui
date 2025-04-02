package tui

import "github.com/rivo/tview"

// NewMaskedTextView returns a new instance of a *MaskedTextView.
func NewMaskedTextView() *MaskedTextView {
	return &MaskedTextView{
		TextView: tview.NewTextView(),
	}
}

// MaskedTextView is a tview.Primitive with all the features of a standard
// *tview.TextView, but it also includes a mask function that processes the
// content before displaying it.
type MaskedTextView struct {
	*tview.TextView
	mask func(string) string
	text string
}

// SetMask sets the function that processes the text before displaying it. The
// provided function should take the original content as input and return a
// modified version to be rendered in the widget. If the mask function is nil,
// the widget will display the original content unchanged.
func (mtv *MaskedTextView) SetMask(fn func(string) string) *MaskedTextView {
	mtv.mask = fn
	return mtv
}

// SetText sets the original text and updates the displayed content.
func (mtv *MaskedTextView) SetText(text string) *MaskedTextView {
	mtv.text = text
	return mtv.setMaskedText()
}

// setMaskedText updates the displayed text using the current mask.
func (mtv *MaskedTextView) setMaskedText() *MaskedTextView {
	text := mtv.text
	if mtv.mask != nil {
		text = mtv.mask(text)
	}
	mtv.TextView.SetText(text)
	return mtv
}
