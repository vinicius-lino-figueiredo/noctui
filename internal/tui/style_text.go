package tui

//
// import (
// 	"fmt"
//
// 	tc "github.com/gdamore/tcell/v2"
// )
//
// // StyleText returns the styled version of the given text
// func StyleText(text string, style tc.Style) string {
// 	fg, bg, _ := style.Decompose()
// 	text = fmt.Sprintf("[%s:%s]%s", fg, bg, text)
// 	return text
// }
//
// func TextToTag(s string, start, end string, tagStyle, noStyle tc.Style) string {
// 	halfBlockStyle := GetHalfBlockStyle(tagStyle, noStyle)
// 	left := StyleText(start, halfBlockStyle)
// 	middle := StyleText(s, tagStyle)
// 	right := StyleText(end, halfBlockStyle)
// 	fixColor := StyleText("", noStyle)
// 	return left + middle + right + fixColor
// }
//
// // GetHalfBlockStyle creates the style for the half blocks that complete
// // the start and end of a tag. It uses a base style to match the widget's
// // background and applies the tag style to set the foreground, making it
// // look like a tag completion.
// func GetHalfBlockStyle(tagStyle, noStyle tc.Style) tc.Style {
// 	_, fg, _ := tagStyle.Decompose()
// 	_, bg, _ := noStyle.Decompose()
// 	return tc.StyleDefault.Background(bg).Foreground(fg)
// }
