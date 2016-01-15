package nanogui

import (
	"github.com/shibukawa/nanovgo"
)

// Text label widget
// The font and color can be customized. When SetFixedWidth()
// is used, the text is wrapped when it surpasses the specified width
//
type Label struct {
	WidgetImplement
	caption     string
	fontFace    string
	color       nanovgo.Color
	columnWidth int
	wrap        bool
}

func NewLabel(parent Widget, caption string) *Label {
	label := &Label{
		caption: caption,
		color:   parent.Theme().TextColor,
		wrap:    true,
	}
	InitWidget(label, parent)
	return label
}

// Caption() gets the label's text caption
func (l *Label) Caption() string {
	return l.caption
}

// SetCaption() sets the label's text caption
func (l *Label) SetCaption(caption string) {
	l.caption = caption
}

// Font() gets the currently active font
func (l *Label) Font() string {
	if l.fontFace == "" {
		return l.theme.FontBold
	}
	return l.fontFace
}

// SetFont() sets the currently active font (2 are available by default: 'sans' and 'sans-bold')
func (l *Label) SetFont(fontFace string) {
	l.fontFace = fontFace
}

// Color() gets the label color
func (l *Label) Color() nanovgo.Color {
	return l.color
}

// SetColor() sets the label color
func (l *Label) SetColor(color nanovgo.Color) {
	l.color = color
}

func (l *Label) ColumnWidth() int {
	return l.columnWidth
}

func (l *Label) SetColumnWidth(width int) {
	l.columnWidth = width
}

func (l *Label) Wrap() bool {
	return l.wrap
}

func (l *Label) SetWrap(wrap bool) {
	l.wrap = wrap
}

func (l *Label) PreferredSize(self Widget, ctx *nanovgo.Context) (int, int) {
	if l.caption == "" {
		return 0, 0
	}
	ctx.SetFontSize(float32(l.FontSize()))
	ctx.SetFontFace(l.Font())

	width := 0
	if l.FixedWidth() > 0 {
		width = l.FixedWidth()
	} else if l.columnWidth > 0 && l.wrap {
		width = l.columnWidth
	}

	if width > 0 {
		ctx.SetTextAlign(nanovgo.AlignLeft | nanovgo.AlignTop)
		bounds := ctx.TextBoxBounds(0, 0, float32(width), l.caption)
		return width, int(bounds[3] - bounds[1])
	} else {
		ctx.SetTextAlign(nanovgo.AlignLeft | nanovgo.AlignTop)
		w, _ := ctx.TextBounds(0, 0, l.caption)
		return int(w), l.Theme().StandardFontSize
	}
}

func (l *Label) Draw(self Widget, ctx *nanovgo.Context) {
	l.WidgetImplement.Draw(self, ctx)
	ctx.SetFontSize(float32(l.FontSize()))
	ctx.SetFontFace(l.Font())
	ctx.SetFillColor(l.color)

	width := 0
	if l.FixedWidth() > 0 {
		width = l.FixedWidth()
	} else if l.columnWidth > 0 && l.wrap {
		width = l.columnWidth
	}

	if width > 0 {
		ctx.SetTextAlign(nanovgo.AlignLeft | nanovgo.AlignTop)
		ctx.TextBox(float32(l.x), float32(l.y), float32(width), l.caption)
	} else {
		ctx.SetTextAlign(nanovgo.AlignLeft | nanovgo.AlignMiddle)
		ctx.Text(float32(l.x), float32(l.y)+float32(l.h)*0.5, l.caption)
	}
}

func (l *Label) String() string {
	return l.StringHelper("Label", l.caption)
}
