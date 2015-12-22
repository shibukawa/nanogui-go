package nanoguiext

import (
	"github.com/shibukawa/nanogui.go"
	"github.com/shibukawa/nanovgo"
)

type FlexibleWidget interface {
	nanogui.Widget
	SetColumnWidth(columnWidth int)
}

type ExpandPolicy int

const (
	ExpandAll ExpandPolicy = iota
	ExpandLast
)

type ExpandBoxLayout struct {
	orientation nanogui.Orientation
	alignment   nanogui.Alignment
	margin      int
	spacing     int
}

func NewExpandBoxLayout(orientation nanogui.Orientation, alignment nanogui.Alignment, setting ...int) *ExpandBoxLayout {
	var margin, spacing int
	switch len(setting) {
	case 0:
	case 1:
		margin = setting[0]
	case 2:
		margin = setting[0]
		spacing = setting[1]
	default:
		panic("NewExpandBoxLayout can accept extra parameter upto 2 (margin, spacing).")
	}
	return &ExpandBoxLayout{
		orientation: orientation,
		alignment:   alignment,
		margin:      margin,
		spacing:     spacing,
	}
}

func (b *ExpandBoxLayout) Orientation() nanogui.Orientation {
	return b.orientation
}

func (b *ExpandBoxLayout) SetOrientation(o nanogui.Orientation) {
	b.orientation = o
}

func (b *ExpandBoxLayout) Alignment() nanogui.Alignment {
	return b.alignment
}

func (b *ExpandBoxLayout) SetAlignment(a nanogui.Alignment) {
	b.alignment = a
}

func (b *ExpandBoxLayout) Margin() int {
	return b.margin
}

func (b *ExpandBoxLayout) SetMargin(m int) {
	b.margin = m
}

func (b *ExpandBoxLayout) Spacing() int {
	return b.spacing
}

func (b *ExpandBoxLayout) SetSpacing(s int) {
	b.spacing = s
}

func (b *ExpandBoxLayout) OnPerformLayout(widget nanogui.Widget, ctx *nanovgo.Context) {
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

	var yOffset int

	if _, ok := widget.(*nanogui.Window); ok {
		if b.orientation == nanogui.Vertical {
			position += widget.Theme().WindowHeaderHeight - b.margin/2
		} else {
			yOffset = widget.Theme().WindowHeaderHeight
		}
	}
	childCount := 0
	fixedChildren := make([]bool, widget.ChildCount())
	fixedLength := make([][2]int, widget.ChildCount())
	remainedLength := containerSize[axis1] - position - b.margin + b.spacing
	for i, child := range widget.Children() {
		if child.Visible() {
			childCount++
			if _, isScroll := child.(*nanogui.VScrollPanel); !isScroll {
				fW, fH := child.FixedSize()
				fs := [2]int{fW, fH}
				fixedLength[i] = fs
				if fs[axis1] > 0 {
					remainedLength -= fs[axis1]
					fixedChildren[i] = true
					childCount--
				}
			}
			remainedLength -= b.spacing
		} else {
			fixedChildren[i] = true
		}
	}

	var averageSize int
	if childCount > 0 {
		averageSize = remainedLength / childCount
	}
	for i, child := range widget.Children() {
		if !child.Visible() {
			continue
		}
		var pos [2]int
		pos[1] = yOffset
		pos[axis1] = position
		var targetSize [2]int
		if fixedChildren[i] && fixedLength[i][axis1] > 0 {
			targetSize[axis1] = fixedLength[i][axis1]
		} else {
			targetSize[axis1] = averageSize
		}

		switch b.alignment {
		case nanogui.Minimum:
			pos[axis2] += b.margin
		case nanogui.Middle:
			pos[axis2] += (containerSize[axis2] - yOffset - targetSize[axis2]) / 2
		case nanogui.Maximum:
			pos[axis2] += containerSize[axis2] - yOffset - targetSize[axis2] - b.margin*2
		case nanogui.Fill:
			pos[axis2] += b.margin
			if fixedLength[i][axis2] > 0 {
				targetSize[axis2] = fixedLength[i][axis2]
			} else {
				targetSize[axis2] = containerSize[axis2] - yOffset - b.margin*2
			}
		}
		child.SetPosition(pos[0], pos[1])
		child.SetSize(targetSize[0], targetSize[1])
		child.OnPerformLayout(child, ctx)
		position += targetSize[axis1] + b.spacing
	}
}

func (b *ExpandBoxLayout) PreferredSize(widget nanogui.Widget, ctx *nanovgo.Context) (int, int) {
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
	return containerSize[0], containerSize[1]
}

type ExpandListLayout struct {
	widths       []int
	stretches    []float32
	alignments   [2]nanogui.Alignment
	expandPolicy [2]ExpandPolicy
	margin       int
	spacing      []int
}

func NewExpandListLayout(widths []int, setting ...int) *ExpandListLayout {
	var margin, spacing int
	switch len(setting) {
	case 0:
	case 1:
		margin = setting[0]
	case 2:
		margin = setting[0]
		spacing = setting[1]
	default:
		panic("NewExpandListLayout can accept extra parameter upto 2 (margin, spacing).")
	}
	return &ExpandListLayout{
		alignments: [2]nanogui.Alignment{nanogui.Minimum, nanogui.Minimum},
		widths:     widths,
		margin:     margin,
		spacing:    []int{spacing, spacing},
	}
}

func (g *ExpandListLayout) Resolution() int {
	return len(g.widths)
}

func (g *ExpandListLayout) ColAlignment() nanogui.Alignment {
	return g.alignments[0]
}

func (g *ExpandListLayout) RowAlignment() nanogui.Alignment {
	return g.alignments[1]
}

func (g *ExpandListLayout) SetColAlignment(a nanogui.Alignment) {
	g.alignments[0] = a
}

func (g *ExpandListLayout) SetRowAlignment(a nanogui.Alignment) {
	g.alignments[1] = a
}

func (g *ExpandListLayout) Margin() int {
	return g.margin
}

func (g *ExpandListLayout) SetMargin(m int) {
	g.margin = m
}

func (g *ExpandListLayout) ColSpacing() int {
	return g.spacing[0]
}

func (g *ExpandListLayout) RowSpacing() int {
	return g.spacing[1]
}

func (g *ExpandListLayout) SetColSpacing(s int) {
	g.spacing[0] = s
}

func (g *ExpandListLayout) SetRowSpacing(s int) {
	g.spacing[1] = s
}

func (g *ExpandListLayout) ColumnWidths() []int {
	return g.widths
}

func (g *ExpandListLayout) SeColumnWidths(lengths []int) {
	g.widths = lengths
}

func (g *ExpandListLayout) SetStretches(stretches []float32) {
	g.stretches = stretches
}

func (g *ExpandListLayout) Stretches() []float32 {
	return g.stretches
}

func (g *ExpandListLayout) ExpandPolicy(axis int) ExpandPolicy {
	return g.expandPolicy[axis]
}

func (g *ExpandListLayout) SetExpandPolicy(axis int, policy ExpandPolicy) {
	g.expandPolicy[axis] = policy
}

func (g *ExpandListLayout) OnPerformLayout(widget nanogui.Widget, ctx *nanovgo.Context) {
	widths, heights, _, _ := g.computeSize(widget, ctx)

	nCols := len(g.widths)

	xOffset := g.margin
	yOffset := g.margin
	window, ok := widget.(*nanogui.Window)
	if ok && window.Title() != "" {
		yOffset += widget.Theme().WindowHeaderHeight - g.margin/2
	}

	row := 0
	for i, child := range widget.Children() {
		column := i % nCols
		width := widths[column]
		height := heights[row]
		pw, ph := child.PreferredSize(child, ctx)
		childXOffset, childWidth := alignment(g.alignments[0], width, pw)
		childYOffset, childHeight := alignment(g.alignments[1], height, ph)
		child.SetPosition(xOffset+childXOffset, yOffset+childYOffset)
		child.SetSize(childWidth, childHeight)
		if column+1 == nCols {
			yOffset += g.spacing[1] + height
			xOffset = g.margin
			row++
		} else {
			xOffset += g.spacing[0] + width
		}
		child.OnPerformLayout(child, ctx)
	}
}

func (g *ExpandListLayout) PreferredSize(widget nanogui.Widget, ctx *nanovgo.Context) (int, int) {
	_, _, totalWidth, totalHeight := g.computeSize(widget, ctx)
	return totalWidth, totalHeight
}

func (g *ExpandListLayout) computeSize(widget nanogui.Widget, ctx *nanovgo.Context) (widths, heights []int, totalWidth, totalHeight int) {
	if widget.ChildCount() == 0 {
		return nil, nil, 0, 0
	}
	nCols := len(g.widths)
	var stretches []float32
	if len(g.stretches) < nCols {
		stretches = append(append([]float32{}, g.stretches...), make([]float32, nCols-len(g.stretches))...)
	}
	isEmpty := true
	for _, stretch := range stretches {
		if stretch > 0 {
			isEmpty = false
			break
		}
	}
	if isEmpty {
		if g.expandPolicy[0] == ExpandLast {
			stretches[nCols-1] = 100.0
		} else {
			for i := 0; i < nCols; i++ {
				stretches[i] = 1.0
			}
		}
	}

	widths = make([]int, nCols)
	totalWidth = 2*g.margin + (nCols-1)*g.spacing[0]
	var totalStretch float32
	for i, columnWidth := range g.widths {
		totalWidth += columnWidth
		totalStretch += stretches[i]
	}
	tableWidth := widget.Width()
	remainedWidth := tableWidth - totalWidth
	if remainedWidth < 0 {
		remainedWidth = 0
	} else {
		totalWidth = tableWidth
	}
	nRows := (widget.ChildCount() + nCols - 1) / nCols
	for i, columnWidth := range g.widths {
		widths[i] = columnWidth + int(float32(remainedWidth)*stretches[i]/totalStretch)
	}
	totalHeight = 2*g.margin + (nRows-1)*g.spacing[1]
	window, ok := widget.(*nanogui.Window)
	if ok && window.Title() != "" {
		totalHeight += widget.Theme().WindowHeaderHeight - g.margin/2
	}
	maxRowHeight := 0
	heights = make([]int, nRows)
	row := 0
	for i, child := range widget.Children() {
		column := i % nCols
		if fWidget, ok := child.(FlexibleWidget); ok {
			fWidget.SetColumnWidth(widths[column])
		}
		_, h := child.PreferredSize(child, ctx)
		if h > maxRowHeight {
			maxRowHeight = h
		}
		if column+1 == nCols {
			heights[row] = maxRowHeight
			totalHeight += maxRowHeight
			maxRowHeight = 0
			row++
		}
	}
	/*if totalHeight < widget.Height() {
		heights[nRows-1] += widget.Height() - totalHeight
		totalHeight = widget.Height()
	}*/
	return
}

func alignment(align nanogui.Alignment, space, preferredSize int) (offset, size int) {
	if space < preferredSize {
		align = nanogui.Fill
	}
	switch align {
	case nanogui.Middle:
		offset = (space - preferredSize) / 2
		size = preferredSize
	case nanogui.Minimum:
		offset = 0
		size = preferredSize
	case nanogui.Maximum:
		offset = space - preferredSize
		size = preferredSize
	case nanogui.Fill:
		offset = 0
		size = space
	}
	return
}
