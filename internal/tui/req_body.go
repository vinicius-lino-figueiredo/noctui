package tui

import (
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	tc "github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vinicius-lino-figueiredo/insomnium"
)

// NewReqBody creates and initializes a new ReqBody component.
func NewReqBody() *ReqBody {
	rb := &ReqBody{}
	rb.Load()
	return rb
}

// ReqBody is a tview.Primitive component for editing and viewing the HTTP
// request body. It supports syntax highlighting using Chroma and toggling
// between edit/view modes.
type ReqBody struct {
	*tview.Pages
	editor     *tview.TextArea
	viewer     *MaskedTextView
	request    *insomnium.Request
	styleText  func(string, tc.Style) string
	ResetFocus func()
	bg         tc.Color
	fg         tc.Color
	theme      string
	lexer      chroma.Lexer
}

// Load initializes the ReqBody layout and components if not already set. It
// sets up the "edit" and "view" pages.
func (rb *ReqBody) Load() {
	if rb.editor == nil {
		rb.LoadEditor()
	}
	if rb.viewer == nil {
		rb.LoadViewer()
	}
	if rb.Pages == nil {
		rb.Pages = tview.NewPages()
	}
	rb.Pages.AddPage("edit", rb.editor, true, false)
	rb.Pages.AddPage("view", rb.viewer, true, true)
}

// LoadEditor sets up the text editor with change and blur callbacks.
func (rb *ReqBody) LoadEditor() {
	rb.editor = tview.NewTextArea()
	rb.editor.SetChangedFunc(rb.EditorChangeFunc).
		SetBlurFunc(rb.EditBlur)
}

// LoadViewer initializes the viewer and its input capture behavior.
func (rb *ReqBody) LoadViewer() {
	rb.viewer = NewMaskedTextView()
	rb.viewer.SetMask(rb.Mask).
		SetDynamicColors(true).
		SetInputCapture(rb.viewerCapture)
}

// viewerCapture handles key events in the viewer. Pressing Escape or 'q' calls
// the ResetFocus callback.
func (rb *ReqBody) viewerCapture(event *tc.EventKey) *tc.EventKey {
	if event.Key() == tc.KeyEsc || event.Rune() == 'q' {
		if rb.ResetFocus != nil {
			rb.ResetFocus()
			return nil
		}
	}
	return event
}

// EditorChangeFunc syncs the editor's text with the viewer and toggles the
// viewer border.
func (rb *ReqBody) EditorChangeFunc() {
	text := rb.editor.GetText()
	rb.viewer.SetText(text).
		SetBorder(text != "")

}

// Focus switches to the editor page and applies the given focus function.
func (rb *ReqBody) Focus(fn func(tview.Primitive)) {
	rb.Pages.SwitchToPage("edit")
	fn(rb.editor)
}

// EditBlur switches to the viewer page when the editor loses focus.
func (rb *ReqBody) EditBlur() {
	rb.Pages.SwitchToPage("view")
}

// Mask applies syntax highlighting to the given text based on the lexer and
// theme. If no style function or lexer is set, it returns the raw text.
func (rb *ReqBody) Mask(text string) string {
	if rb.styleText == nil {
		return text
	}
	if rb.lexer == nil {
		return text
	}
	iter, err := rb.lexer.Tokenise(nil, text)
	if err != nil {
		return text
	}

	theme := styles.Get(rb.theme)

	baseStyle := tc.StyleDefault.Background(rb.bg)
	var res string
	for token := iter(); token != chroma.EOF; token = iter() {
		tokenStyle := theme.Get(token.Type)
		color := rb.GetTokenColor(tokenStyle)
		style := baseStyle.Foreground(color)
		res += rb.styleText(token.Value, style)
	}

	return res
}

// GetTokenColor converts a Chroma color to a tcell color.
func (rb *ReqBody) GetTokenColor(hex chroma.StyleEntry) tc.Color {
	return tc.NewHexColor(int32(hex.Colour))
}

// SetText sets the text of both editor and viewer, enabling their borders if
// not empty.
func (rb *ReqBody) SetText(text string) *ReqBody {
	rb.editor.SetText(text, false).
		SetBorder(true)
	rb.viewer.SetText(text).
		SetBorder(text != "")
	return rb
}

// SetStyle applies a tcell style (background, abd foreground) to the editor and
// viewer.
func (rb *ReqBody) SetStyle(style tc.Style) *ReqBody {
	_, rb.bg, _ = style.Decompose()
	rb.editor.SetTextStyle(style)
	rb.viewer.SetTextStyle(style)
	rb.editor.GetCursor()
	return rb
}

// SetTheme applies a theme ("monokai", etc.) to the Body tab.
func (rb *ReqBody) SetTheme(theme string) *ReqBody {
	rb.theme = theme
	rb.viewer.setMaskedText()
	return rb
}

// SetStyleTextFunc sets the function used to render styled text.
func (rb *ReqBody) SetStyleTextFunc(fn func(string, tc.Style) string) *ReqBody {
	rb.styleText = fn
	return rb
}

// SetRequest binds a request to the component and sets its body and lexer based
// on Content-Type.
func (rb *ReqBody) SetRequest(req *insomnium.Request) *ReqBody {
	rb.request = req
	if req == nil {
		rb.lexer = nil
		rb.SetText("")
		return rb
	}
	var typ string
	for _, h := range req.Headers {
		if h.Name == "Content-Type" {
			typ = h.Value
			break
		}
	}
	rb.lexer = lexers.MatchMimeType(typ)
	rb.SetText(req.Body.Text)
	return rb
}

// SetResetFocus sets the callback used to release focus.
func (rb *ReqBody) SetResetFocus(fn func()) *ReqBody {
	rb.ResetFocus = fn
	return rb
}

func (rb *ReqBody) SetBackgroundColor(bg tc.Color) *ReqBody {
	rb.bg = bg
	rb.Pages.SetBackgroundColor(bg)
	rb.viewer.SetBackgroundColor(bg)
	rb.editor.SetBackgroundColor(bg)
	return rb
}

func (rb *ReqBody) SetForegroundColor(fg tc.Color) *ReqBody {
	rb.fg = fg
	rb.Pages.SetBorderColor(fg)
	rb.viewer.SetBorderColor(fg)
	rb.editor.SetBorderColor(fg)
	return rb
}
