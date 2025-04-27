package tui

import (
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	tc "github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vinicius-lino-figueiredo/insomnium"
)

// NewResponseBody creates a new ResponseBody instance and loads its initial
// state.
func NewResponseBody() *ResponseBody {
	rb := &ResponseBody{}
	rb.Load()
	return rb
}

// ResponseBody displays a formatted HTTP response body inside a TUI text view.
type ResponseBody struct {
	*tview.TextView
	response        *insomnium.Response
	styleText       func(string, tc.Style) string
	lexer           chroma.Lexer
	getResponseBody func(*insomnium.Response) (string, error)
	theme           string
	bg              tc.Color
	resetFocus      func()
}

// Load initializes the embedded TextView and sets its input capture handler.
func (rb *ResponseBody) Load() {
	rb.TextView = tview.NewTextView()
	rb.TextView.SetDynamicColors(true).
		SetInputCapture(rb.CaptureInput)
}

// CaptureInput handles key events to allow exiting the view.
func (rb *ResponseBody) CaptureInput(event *tc.EventKey) *tc.EventKey {
	if event.Key() == tc.KeyEsc || event.Rune() == 'q' {
		rb.resetFocus()
		return nil
	}
	return event
}

// SetResponse sets the HTTP response and updates the displayed body content.
func (rb *ResponseBody) SetResponse(res *insomnium.Response) *ResponseBody {
	rb.response = res
	body := rb.getBody(res)
	rb.lexer = rb.getLexer(res, body)
	rb.SetText(rb.decorateBody(body))
	return rb
}

// getBody retrieves the body string from the response using the custom
// function.
func (rb *ResponseBody) getBody(res *insomnium.Response) string {
	if res == nil || rb.getResponseBody == nil {
		return ""
	}
	body, _ := rb.getResponseBody(res)
	return body
}

// getLexer selects a syntax lexer based on the response's Content-Type header.
func (rb *ResponseBody) getLexer(res *insomnium.Response, body string) chroma.Lexer {
	if res == nil || body == "" {
		return nil
	}
	contentype := rb.getContentType()
	return lexers.MatchMimeType(contentype)
}

// getContentType retrieves the Content-Type header from the response.
func (rb *ResponseBody) getContentType() string {
	if rb.response == nil {
		return ""
	}
	for _, header := range rb.response.Headers {
		if strings.ToLower(header.Name) == "content-type" {
			return header.Value
		}
	}
	return ""
}

// decorateBody applies syntax highlighting to the body using the selected
// lexer.
func (rb *ResponseBody) decorateBody(body string) string {
	if body == "" || rb.styleText == nil || rb.lexer == nil {
		return body
	}
	iter, err := rb.lexer.Tokenise(nil, body)
	if err != nil {
		return body
	}

	theme := styles.Get(rb.theme)
	style := tc.StyleDefault.Background(rb.bg)
	var result string
	for token := iter(); token != chroma.EOF; token = iter() {
		tokenStyle := theme.Get(token.Type)
		color := rb.GetTokenColor(tokenStyle)
		style = style.Foreground(color)
		result += rb.styleText(token.Value, style)
	}

	return result
}

// GetTokenColor converts a Chroma StyleEntry color to a tcell Color.
func (rb *ResponseBody) GetTokenColor(hex chroma.StyleEntry) tc.Color {
	return tc.NewHexColor(int32(hex.Colour))
}

// SetStyleTextFunc sets the function used to style text tokens.
func (rb *ResponseBody) SetStyleTextFunc(fn func(string, tc.Style) string) *ResponseBody {
	rb.styleText = fn
	return rb
}

// SetGetResponseBodyFunc sets the function used to extract the response body.
func (rb *ResponseBody) SetGetResponseBodyFunc(fn func(*insomnium.Response) (string, error)) *ResponseBody {
	rb.getResponseBody = fn
	return rb
}

// SetResponseBodyTheme sets the syntax highlighting theme.
func (rb *ResponseBody) SetResponseBodyTheme(theme string) *ResponseBody {
	rb.theme = theme
	rb.SetResponse(rb.response)
	return rb
}

// SetBackgroundColor sets the background color of the TextView.
func (rb *ResponseBody) SetBackgroundColor(bg tc.Color) *ResponseBody {
	rb.bg = bg
	rb.TextView.SetBackgroundColor(bg)
	return rb
}

// SetResetFocusFunc sets the function called when exiting the TextView.
func (rb *ResponseBody) SetResetFocusFunc(fn func()) *ResponseBody {
	rb.resetFocus = fn
	return rb
}
