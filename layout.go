package nanogui

import (
	"github.com/shibukawa/nanovgo"
)

type Alignment uint8

const (
	Middle Alignment = iota
	Minimum
	Maximum
	Fill
)

type Orientation int

const (
	Horizontal Orientation = iota
	Vertical
)

type Layout interface {
	OnPerformLayout(widget Widget, ctx *nanovgo.Context)
	PreferredSize(widget Widget, ctx *nanovgo.Context) (int, int)
}

// Simple horizontal/vertical box layout
//
// This widget stacks up a bunch of widgets horizontally or vertically. It adds
// margins around the entire container and a custom spacing between adjacent
// widgets

type BoxLayout struct {
	orientation Orientation
	alignment   Alignment
	margin      int
	spacing     int
}

func NewBoxLayout(orientation Orientation, alignment Alignment, margin, spacing int) *BoxLayout {
	return &BoxLayout{
		orientation: orientation,
		alignment:   alignment,
		margin:      margin,
		spacing:     spacing,
	}
}

func (b *BoxLayout) Orientation() Orientation {
	return b.orientation
}

func (b *BoxLayout) SetOrientation(o Orientation) {
	b.orientation = o
}

func (b *BoxLayout) Alignment() Alignment {
	return b.alignment
}

func (b *BoxLayout) SetAlignment(a Alignment) {
	b.alignment = a
}

func (b *BoxLayout) Margin() int {
	return b.margin
}

func (b *BoxLayout) SetMargin(m int) {
	b.margin = m
}

func (b *BoxLayout) Spacing() int {
	return b.spacing
}

func (b *BoxLayout) SetSpacing(s int) {
	b.spacing = s
}

func (b *BoxLayout) OnPerformLayout(widget Widget, ctx *nanovgo.Context) {
	fX, fY := widget.FixedSize()
	var containerSize [2]int
	if fX > 0 {
		containerSize[0] = fX
	} else {
		containerSize[0] = widget.Width()
	}
	if fY > 0 {
		containerSize[1] = fY
	} else {
		containerSize[1] = widget.Height()
	}
	axis1 := int(b.orientation)
	axis2 := (int(b.orientation) + 1) % 2
	position := b.margin

	if _, ok := widget.(*Window); ok {
		position += widget.Theme().WindowHeaderHeight - b.margin/2
	}
	first := true
	for _, child := range widget.Children() {
		if !child.Visible() {
			continue
		}
		if first {
			first = false
		} else {
			position += b.spacing
		}
		var fs [2]int
		pX, pY := child.PreferredSize(child, ctx)
		fs[0], fs[1] = child.FixedSize()
		var targetSize [2]int
		if fX > 0 {
			targetSize[0] = fs[0]
		} else {
			targetSize[0] = pX
		}
		if fY > 0 {
			targetSize[1] = fs[1]
		} else {
			targetSize[1] = pY
		}
		var pos [2]int
		pos[axis1] = position

		switch b.alignment {
		case Minimum:
			pos[axis2] = b.margin
		case Middle:
			pos[axis2] = (containerSize[axis2] - targetSize[axis2]) / 2
		case Maximum:
			pos[axis2] = containerSize[axis2] - targetSize[axis2] - b.margin
		case Fill:
			pos[axis2] = b.margin
			if fs[axis2] > 0 {
				targetSize[axis2] = fs[axis2]
			} else {
				targetSize[axis2] = containerSize[axis2]
			}
		}
		child.SetPosition(pos[0], pos[1])
		child.SetSize(targetSize[0], targetSize[1])
		child.OnPerformLayout(child, ctx)
		position += targetSize[axis1]
	}
}

func (b *BoxLayout) PreferredSize(widget Widget, ctx *nanovgo.Context) (int, int) {
	size := []int{2 * b.margin, 2 * b.margin}

	if _, ok := widget.(*Window); ok {
		size[1] += widget.Theme().WindowHeaderHeight - b.margin/2
	}

	first := true
	axis1 := int(b.orientation)
	axis2 := (int(b.orientation) + 1) % 2

	for _, child := range widget.Children() {
		if !child.Visible() {
			continue
		}
		if first {
			first = false
		} else {
			size[axis1] += b.spacing
		}

		pX, pY := child.PreferredSize(child, ctx)
		fX, fY := child.FixedSize()
		var targetSize [2]int
		if fX > 0 {
			targetSize[0] = fX
		} else {
			targetSize[0] = pX
		}
		if fY > 0 {
			targetSize[1] = fY
		} else {
			targetSize[1] = pY
		}
		size[axis1] += targetSize[axis1]
		size[axis2] = maxI(size[axis2], targetSize[axis2]+2*b.margin)
	}
	return size[0], size[1]
}

// Special layout for widgets grouped by labels
//
// This widget resembles a box layout in that it arranges a set of widgets
// vertically. All widgets are indented on the horizontal axis except for
// Label widgets, which are not indented.
//
// This creates a pleasing layout where a number of widgets are grouped
// under some high-level heading.
type GroupLayout struct {
	margin       int
	spacing      int
	groupIndent  int
	groupSpacing int
}

func NewGroupLayout(margin, spacing, groupSpacing, groupIndent int) Layout {
	if margin < 0 {
		margin = 15
	}
	if spacing < 0 {
		spacing = 6
	}
	if groupIndent < 0 {
		groupIndent = 20
	}
	if groupSpacing < 0 {
		groupSpacing = 14
	}
	return &GroupLayout{
		margin:       margin,
		spacing:      spacing,
		groupIndent:  groupIndent,
		groupSpacing: groupSpacing,
	}
}

func (g *GroupLayout) Margin() int {
	return g.margin
}

func (g *GroupLayout) SetMargin(m int) {
	g.margin = m
}

func (g *GroupLayout) Spacing() int {
	return g.spacing
}

func (g *GroupLayout) SetSpacing(s int) {
	g.spacing = s
}

func (g *GroupLayout) GroupIndent() int {
	return g.groupIndent
}

func (g *GroupLayout) SetGroupIndent(m int) {
	g.groupIndent = m
}

func (g *GroupLayout) GroupSpacing() int {
	return g.groupSpacing
}

func (g *GroupLayout) SetGroupSpacing(s int) {
	g.groupSpacing = s
}

func (g *GroupLayout) OnPerformLayout(widget Widget, ctx *nanovgo.Context) {
	height := g.margin
	availableWidth := -g.margin * 2
	availableWidth += toI(widget.FixedWidth() > 0, widget.FixedWidth(), widget.Width())
	window, ok := widget.(*Window)
	if ok && window.Title() != "" {
		height += widget.Theme().WindowHeaderHeight - g.margin/2
	}
	first := true
	indent := false

	for _, child := range widget.Children() {
		if !child.Visible() {
			continue
		}
		label, ok := child.(*Label)
		if !first {
			height += toI(ok, g.groupSpacing, g.spacing)
		}
		first = false
		var indentValue int

		if indent && !ok {
			indentValue = g.groupIndent
		}

		pW := availableWidth - indentValue
		_, pH := child.PreferredSize(child, ctx)
		fW, fH := child.FixedSize()
		tW := toI(fW > 0, fW, pW)
		tH := toI(fH > 0, fH, pH)
		child.SetPosition(g.margin+indentValue, height)
		child.SetSize(tW, tH)
		child.OnPerformLayout(child, ctx)
		height += tH

		if ok {
			indent = label.Caption() != ""
		}
	}
}

func (g *GroupLayout) PreferredSize(widget Widget, ctx *nanovgo.Context) (int, int) {
	height := g.margin
	width := g.margin * 2

	window, ok := widget.(*Window)
	if ok && window.Title() != "" {
		height += widget.Theme().WindowHeaderHeight - g.margin/2
	}
	first := true
	indent := false

	for _, child := range widget.Children() {
		if !child.Visible() {
			continue
		}
		label, ok := child.(*Label)
		if !first {
			height += toI(ok, g.groupSpacing, g.spacing)
		}
		first = false
		pW, pH := child.PreferredSize(child, ctx)
		fW, fH := child.FixedSize()
		tW := toI(fW > 0, fW, pW)
		tH := toI(fH > 0, fH, pH)
		var indentValue int
		if indent && !ok {
			indentValue = g.groupIndent
		}
		height += tH
		width = maxI(width, tW+2*g.margin+indentValue)

		if ok {
			indent = label.Caption() != ""
		}
	}
	height += g.margin
	return width, height
}
