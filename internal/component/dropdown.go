package component

import (
	"regexp"
	"slices"
	"strings"

	tc "github.com/gdamore/tcell/v2"
	tv "github.com/rivo/tview"
)

const tagRegex = `\[(\w+|#[0-9abcdef]{6}):(\w+|#[0-9abcdef]{6})\]`

// NewDropDown creates a new DropDown component with default settings.
func NewDropDown(widget tv.Primitive) *DropDown {
	if widget == nil {
		panic("nil widget referenced in dropdown")
	}
	return &DropDown{
		widget:   widget,
		maxRows:  5,
		tagRegex: regexp.MustCompile(tagRegex),
		hasBar:   true,
	}
}

// DropDown represents a selectable dropdown list of options.
type DropDown struct {
	hasFocus       bool
	selected       int
	widget         tv.Primitive
	isOpen         bool
	maxRows        int
	skip           int
	width, height  int
	x, y           int
	options        []*DropDownOption
	tagRegex       *regexp.Regexp
	inputCapture   func(*tc.EventKey) *tc.EventKey
	hasBar         bool
	emptyBarRune   rune
	filledBarRune  rune
	emptyBarStyle  tc.Style
	filledBarStyle tc.Style
	selectedFn     func(int, *DropDownOption)
	mouseCapture   func(tv.MouseAction, *tc.EventMouse) (tv.MouseAction, *tc.EventMouse)
}

// Blur implements tview.Primitive.
func (d *DropDown) Blur() {
	d.hasFocus = false
	d.isOpen = false
}

// Draw implements tview.Primitive.
func (d *DropDown) Draw(screen tc.Screen) {
	d.widget.Draw(screen)
	if !d.isOpen || d.width == 0 || d.height == 0 || len(d.options) == 0 {
		return
	}
	if d.selected < 0 {
		d.skip += d.selected
		d.selected = 0
	}
	if d.selected >= d.maxRows {
		d.skip += d.selected - d.maxRows + 1
		d.selected = d.maxRows - 1
	}
	if d.skip > len(d.options)-int(d.maxRows) {
		d.skip = len(d.options) - int(d.maxRows)
	}
	if d.skip < 0 {
		d.skip = 0
	}

	maxRows := min(d.maxRows, len(d.options)-d.skip)

	options := make([]string, maxRows)
	for n, option := range d.options[d.skip : d.skip+maxRows] {
		if d.hasFocus && n == d.selected {
			options[n] = option.selected
		} else {
			options[n] = option.name
		}
	}

	maxLenName := slices.MaxFunc(options, func(a, b string) int {
		aName := d.tagRegex.ReplaceAllString(a, "")
		bName := d.tagRegex.ReplaceAllString(b, "")
		return len(aName) - len(bName)
	})

	maxLen := len([]rune(d.tagRegex.ReplaceAllString(maxLenName, "")))
	if d.hasBar {
		maxLen++
	}

	// maxRows := len(d.options)

	screenWidth, _ := screen.Size()

	x := min(d.x, screenWidth-maxLen)

	for n := range maxRows {
		text := options[n]
		var style tc.Style
		var nextStyle tc.Style
		var screenPos int
		var blockPos int
		for {
			var block string
			text, block, nextStyle = d.step(text, style)

			runes := []rune(block)
			for blockPos = range len(runes) {
				r := runes[blockPos]
				runeX := x + screenPos + blockPos
				screen.SetContent(runeX, d.y+n+1, r, nil, style)
			}
			screenPos += len(runes)
			if len(text) == 0 {
				break
			}
			style = nextStyle
		}
		for ; screenPos < maxLen; screenPos++ {
			screen.SetContent(x+screenPos, d.y+n+1, ' ', nil, style)
		}
	}
	d.drawBar(screen, maxRows, x+maxLen-1)
}

func (d *DropDown) drawBar(screen tc.Screen, maxRows int, x int) {
	sizeProportion := float64(maxRows) / float64(len(d.options))
	barSize := int(sizeProportion * float64(maxRows))
	barPosition := int(sizeProportion * float64(d.skip))

	if d.skip+maxRows == len(d.options) {
		barPosition = maxRows - barSize
	}
	var r rune
	var style tc.Style
	for i := range maxRows {
		r = d.emptyBarRune
		style = d.emptyBarStyle
		if i >= barPosition && i < barPosition+barSize {
			r = d.filledBarRune
			style = d.filledBarStyle
		}
		screen.SetContent(x, d.y+i+1, r, nil, style)
	}
}

func (d *DropDown) step(text string, prevStyle tc.Style) (string, string, tc.Style) {
	match := d.tagRegex.FindStringSubmatchIndex(text)
	if match == nil {
		return "", text, prevStyle
	}

	style := d.readTag(text[match[0]:match[1]], prevStyle)

	return text[match[1]:], text[:match[0]], style
}

func (d *DropDown) readTag(tag string, prevStyle tc.Style) tc.Style {
	tag = strings.Trim(tag, "[]")
	parts := strings.Split(tag, ":")
	fg, bg := parts[0], parts[1]
	if fg != "" {
		prevStyle = prevStyle.Foreground(tc.GetColor(fg))
	}
	if bg != "" {
		prevStyle = prevStyle.Background(tc.GetColor(bg))
	}
	return prevStyle
}

// SetOpen sets the dropdown's open state and returns the dropdown.
func (d *DropDown) SetOpen(open bool) *DropDown {
	d.isOpen = open
	return d
}

// IsOpen returns true if the dropdown is currently open.
func (d *DropDown) IsOpen() bool {
	return d.isOpen
}

// Clear removes all options and resets the scroll.
func (d *DropDown) Clear() *DropDown {
	d.options = d.options[:0]
	d.skip = 0
	return d
}

// AddOption adds a new option to the dropdown.
func (d *DropDown) AddOption(option *DropDownOption) *DropDown {
	d.options = append(d.options, option)
	return d
}

// SetEmptyBarRune sets the rune used for the empty scrollbar area.
func (d *DropDown) SetEmptyBarRune(r rune) *DropDown {
	d.emptyBarRune = r
	return d
}

// SetFilledBarRune sets the rune used for the filled scrollbar area.
func (d *DropDown) SetFilledBarRune(r rune) *DropDown {
	d.filledBarRune = r
	return d
}

// SetEmptyBarStyle sets the style of the empty scrollbar area.
func (d *DropDown) SetEmptyBarStyle(style tc.Style) *DropDown {
	d.emptyBarStyle = style
	return d
}

// SetFilledBarStyle sets the style of the filled scrollbar area.
func (d *DropDown) SetFilledBarStyle(style tc.Style) *DropDown {
	d.filledBarStyle = style
	return d
}

// Focus implements tview.Primitive.
func (d *DropDown) Focus(_ func(tv.Primitive)) {
	d.selected = 0
	d.skip = 0
	d.hasFocus = true
}

// GetRect implements tview.Primitive.
func (d *DropDown) GetRect() (int, int, int, int) {
	return d.x, d.y, d.width, d.height
}

// HasFocus implements tview.Primitive.
func (d *DropDown) HasFocus() bool {
	return d.hasFocus
}

func (d *DropDown) inputHandler(event *tc.EventKey, _ func(tv.Primitive)) {
	if d.inputCapture != nil {
		event = d.inputCapture(event)
	}
	switch {
	case event.Key() == tc.KeyUp || event.Rune() == 'k':
		d.selected--
	case event.Key() == tc.KeyDown || event.Rune() == 'j':
		d.selected++
	case event.Key() == tc.KeyEnter:
		if d.selectedFn != nil {
			d.selectedFn(d.GetCurrent())
			d.SetOpen(false)
		}
	}
}

// InputHandler implements tview.Primitive.
func (d *DropDown) InputHandler() func(*tc.EventKey, func(tv.Primitive)) {
	return d.inputHandler
}

// MouseHandler implements tview.Primitive.
func (d *DropDown) MouseHandler() func(action tv.MouseAction, event *tc.EventMouse, setFocus func(p tv.Primitive)) (consumed bool, capture tv.Primitive) {
	return func(action tv.MouseAction, event *tc.EventMouse, _ func(tv.Primitive)) (consumed bool, _ tv.Primitive) {
		if d.mouseCapture != nil {
			action, event = d.mouseCapture(action, event)
		}
		if event == nil && action == tv.MouseConsumed {
			consumed = true
		}
		return
	}
}

// PasteHandler implements tview.Primitive.
func (d *DropDown) PasteHandler() func(text string, setFocus func(p tv.Primitive)) {
	return nil
}

// SetRect implements tview.Primitive.
func (d *DropDown) SetRect(x int, y int, width int, height int) {
	d.x, d.y, d.width, d.height = x, y, width, height
	d.widget.SetRect(x, y, width, height)
}

// SetMaxRows sets the maximum number of visible rows.
func (d *DropDown) SetMaxRows(n int) *DropDown {
	if n <= 0 {
		panic("max rows must be positive")
	}
	d.maxRows = n
	return d
}

// SetSelectedFunc sets the callback when an option is selected.
func (d *DropDown) SetSelectedFunc(fn func(int, *DropDownOption)) *DropDown {
	d.selectedFn = fn
	return d
}

// GetCurrent returns the index and option currently selected.
func (d *DropDown) GetCurrent() (int, *DropDownOption) {
	index := d.skip + d.selected
	if len(d.options) == 0 || index >= len(d.options) {
		return -1, nil
	}
	o := d.options[index]
	return index, o
}

// NewDropDownOption creates a new DropDownOption with text and reference.
func NewDropDownOption(name string, sel string, ref any) *DropDownOption {
	return &DropDownOption{
		name:     name,
		selected: sel,
		ref:      ref,
	}
}

// DropDownOption represents a selectable item in a dropdown list.
type DropDownOption struct {
	name     string
	selected string
	ref      any
}

// GetRef returns the reference value associated with this option.
func (do *DropDownOption) GetRef() any {
	return do.ref
}

var _ tv.Primitive = (*DropDown)(nil)
