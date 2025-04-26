package component

import (
	tc "github.com/gdamore/tcell/v2"
	tv "github.com/rivo/tview"
)

// NewTag creates a new Tag component with default style and delimiters.
func NewTag() *Tag {
	t := &Tag{
		tagStart: '▐',
		tagEnd:   '▌',
	}
	t.SetTagStyle(tc.StyleDefault.Background(tc.ColorDarkGray))
	return t
}

// Tag is a custom visual component that renders a text label between two
// decorative runes with configurable style.
type Tag struct {
	text            string
	x, y            int
	width, height   int
	hasFocus        bool
	inputCapture    func(*tc.EventKey) *tc.EventKey
	mouseCapture    func(tv.MouseAction, *tc.EventMouse) (tv.MouseAction, *tc.EventMouse)
	fontColor       tc.Color
	tagColor        tc.Color
	tagStyle        tc.Style
	endStyle        tc.Style
	backgroundColor tc.Color
	tagStart        rune
	tagEnd          rune
}

// Blur implements tview.Primitive.
func (t *Tag) Blur() {
	t.hasFocus = false
}

// Draw implements tview.Primitive.
func (t *Tag) Draw(screen tc.Screen) {
	if t.width <= 0 || t.height <= 0 {
		return
	}
	screen.SetContent(t.x, t.y, t.tagStart, nil, t.endStyle)
	var n int
	for _, r := range t.text {
		x := t.x + n + 1
		screen.SetContent(x, t.y, r, nil, t.tagStyle)
		n++
	}
	screen.SetContent(t.x+n+1, t.y, t.tagEnd, nil, t.endStyle)
}

// Focus implements tview.Primitive.
func (t *Tag) Focus(_ func(p tv.Primitive)) {
	t.hasFocus = true
}

// GetRect implements tview.Primitive.
func (t *Tag) GetRect() (int, int, int, int) {
	return t.x, t.y, t.width, t.height
}

// HasFocus implements tview.Primitive.
func (t *Tag) HasFocus() bool {
	return t.hasFocus
}

// inputHandler handles input events.
func (t *Tag) inputHandler(event *tc.EventKey, _ func(tv.Primitive)) {
	if t.inputCapture != nil {
		_ = t.inputCapture(event)
	}
}

// InputHandler implements tview.Primitive.
func (t *Tag) InputHandler() func(event *tc.EventKey, setFocus func(p tv.Primitive)) {
	return t.inputHandler
}

// MouseHandler implements tview.Primitive.
func (t *Tag) MouseHandler() func(tv.MouseAction, *tc.EventMouse, func(tv.Primitive)) (bool, tv.Primitive) {
	return func(action tv.MouseAction, event *tc.EventMouse, _ func(tv.Primitive)) (consumed bool, _ tv.Primitive) {
		if t.mouseCapture != nil {
			action, event = t.mouseCapture(action, event)
		}
		if event == nil && action == tv.MouseConsumed {
			consumed = true
		}
		return
	}
}

// PasteHandler implements tview.Primitive.
func (t *Tag) PasteHandler() func(text string, setFocus func(p tv.Primitive)) {
	return nil
}

// SetRect implements tview.Primitive.
func (t *Tag) SetRect(x int, y int, width int, height int) {
	t.x, t.y = x, y
	t.width, t.height = width, height
}

// SetText updates the text to be displayed in the tag.
func (t *Tag) SetText(text string) *Tag {
	t.text = text
	return t
}

// SetTagStyle sets the style of the text and tag background.
func (t *Tag) SetTagStyle(style tc.Style) *Tag {
	fg, bg, _ := style.Decompose()
	t.tagColor = bg
	t.fontColor = fg
	t.tagStyle = style
	t.endStyle = t.endStyle.Foreground(bg)
	return t
}

// SetBackgroundColor changes the background color of the delimiters.
func (t *Tag) SetBackgroundColor(bg tc.Color) *Tag {
	t.backgroundColor = bg
	t.endStyle = t.endStyle.Background(bg)
	return t
}

var _ tv.Primitive = (*Tag)(nil)
