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

func NewBoxLayout(orientation Orientation, alignment Alignment, setting ...int) *BoxLayout {
	var margin, spacing int
	switch len(setting) {
	case 0:
	case 1:
		margin = setting[0]
	case 2:
		margin = setting[0]
		spacing = setting[1]
	default:
		panic("NewBoxLayout can accept extra parameter upto 2 (margin, spacing).")
	}
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

func NewGroupLayout(setting ...int) Layout {
	margin := -1
	spacing := -1
	groupIndent := -1
	groupSpacing := -1
	switch len(setting) {
	case 0:
	case 1:
		margin = setting[0]
	case 2:
		margin = setting[0]
		spacing = setting[1]
	case 3:
		margin = setting[0]
		spacing = setting[1]
		groupIndent = setting[2]
	case 4:
		margin = setting[0]
		spacing = setting[1]
		groupIndent = setting[2]
		groupSpacing = setting[3]
	default:
		panic("NewGroupLayout can accept extra parameter upto 4 (margin, spacing, groupIndent, groupSpacing).")
	}

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

// Grid layout
//
// Widgets are arranged in a grid that has a fixed grid resolution \c resolution
// along one of the axes. The layout orientation indicates the fixed dimension;
// widgets are also appended on this axis. The spacing between items can be
// specified per axis. The horizontal/vertical alignment can be specified per
// row and column.

type GridLayout struct {
	orientation       Orientation
	resolution        int
	defaultAlignments [2]Alignment
	alignments        [2][]Alignment
	margin            int
	spacing           []int
}

func NewGridLayout(orientation Orientation, resolution int, alignment Alignment, setting ...int) *GridLayout {
	var margin, spacing int
	switch len(setting) {
	case 0:
	case 1:
		margin = setting[0]
	case 2:
		margin = setting[0]
		spacing = setting[1]
	default:
		panic("NewGridLayout can accept extra parameter upto 2 (margin, spacing).")
	}
	return &GridLayout{
		orientation:       orientation,
		resolution:        resolution,
		defaultAlignments: [2]Alignment{alignment, alignment},
		alignments:        [2][]Alignment{{}, {}},
		margin:            margin,
		spacing:           []int{spacing, spacing},
	}
}

func (g *GridLayout) Orientation() Orientation {
	return g.orientation
}

func (g *GridLayout) SetOrientation(o Orientation) {
	g.orientation = o
}

func (g *GridLayout) Resolution() int {
	return g.resolution
}

func (g *GridLayout) SetResolution(r int) {
	g.resolution = r
}

func (g *GridLayout) ColDefaultAlignment() Alignment {
	return g.defaultAlignments[0]
}

func (g *GridLayout) RowDefaultAlignment() Alignment {
	return g.defaultAlignments[1]
}

func (g *GridLayout) SetColDefaultAlignment(a Alignment) {
	g.defaultAlignments[0] = a
}

func (g *GridLayout) SetRowDefaultAlignment(a Alignment) {
	g.defaultAlignments[1] = a
}

func (g *GridLayout) ColAlignment() []Alignment {
	return g.alignments[0]
}

func (g *GridLayout) RowAlignment() []Alignment {
	return g.alignments[1]
}

func (g *GridLayout) SetColAlignment(a ...Alignment) {
	g.alignments[0] = a
}

func (g *GridLayout) SetRowAlignment(a ...Alignment) {
	g.alignments[1] = a
}

func (g *GridLayout) Alignment(axis, item int) Alignment {
	if item < len(g.alignments[axis]) {
		return g.alignments[axis][item]
	}
	return g.defaultAlignments[axis]
}

func (g *GridLayout) Margin() int {
	return g.margin
}

func (g *GridLayout) SetMargin(m int) {
	g.margin = m
}

func (g *GridLayout) ColSpacing() int {
	return g.spacing[0]
}

func (g *GridLayout) RowSpacing() int {
	return g.spacing[1]
}

func (g *GridLayout) SetColSpacing(s int) {
	g.spacing[0] = s
}

func (g *GridLayout) SetRowSpacing(s int) {
	g.spacing[1] = s
}

func (g *GridLayout) OnPerformLayout(widget Widget, ctx *nanovgo.Context) {
	fw, fh := widget.FixedSize()
	containerSize := []int{
		toI(fw > 0, fw, widget.Width()),
		toI(fh > 0, fh, widget.Height()),
	}

	/* Compute minimum row / column sizes */
	grid := g.ComputeLayout(widget, ctx)
	dim := []int{len(grid[0]), len(grid[1])}

	extra := []int{0, 0}
	if _, ok := widget.(*Window); ok {
		extra[1] = widget.Theme().WindowHeaderHeight - g.margin/2
	}

	/* Strech to size provided by widget */
	for i := 0; i < 2; i++ {
		gridSize := g.margin*2 + extra[i]
		for _, s := range grid[i] {
			gridSize += s
			if i+1 < dim[i] {
				gridSize += g.spacing[i]
			}
		}

		if gridSize < containerSize[i] {
			gap := containerSize[i] - gridSize
			g := gap / dim[i]
			rest := gap - g*dim[i]
			for j := 0; j < dim[i]; j++ {
				grid[i][j] += g
			}
			for j := 0; rest > 0 && j < dim[i]; j++ {
				grid[i][j]++
				rest--
			}
		}
	}

	axis1 := int(g.orientation)
	axis2 := (int(g.orientation) + 1) % 2
	start := []int{g.margin, g.margin + extra[1]}
	pos := []int{start[0], start[1]}
	numChildren := widget.ChildCount()
	child := 0
	children := widget.Children()

	for i2 := 0; i2 < dim[axis2]; i2++ {
		pos[axis1] = start[axis1]
		for i1 := 0; i1 < dim[axis1]; i1++ {
			var w Widget
			for {
				if child >= numChildren {
					return
				}
				w = children[child]
				child++
				if w.Visible() {
					break
				}
			}
			pw, ph := w.PreferredSize(w, ctx)
			fw, fh := w.FixedSize()
			fs := []int{fw, fh}
			targetSize := []int{
				toI(fw > 0, fw, pw),
				toI(fh > 0, fh, ph),
			}
			itemPos := []int{pos[0], pos[1]}
			for j := 0; j < 2; j++ {
				axis := (axis1 + j) % 2
				item := toI(j == 0, i1, i2)
				align := g.Alignment(axis, item)

				switch align {
				case Minimum:
				case Middle:
					itemPos[axis] += (grid[axis][item] - targetSize[axis]) / 2
				case Maximum:
					itemPos[axis] += grid[axis][item] - targetSize[axis]
				case Fill:
					targetSize[axis] = toI(fs[axis] > 0, fs[axis], grid[axis][item])
				}
			}
			w.SetPosition(itemPos[0], itemPos[1])
			w.SetSize(targetSize[0], targetSize[1])
			w.OnPerformLayout(w, ctx)
			pos[axis1] += grid[axis1][i1] + g.spacing[axis1]
		}
		pos[axis2] += grid[axis2][i2] + g.spacing[axis2]
	}
}

func (g *GridLayout) PreferredSize(widget Widget, ctx *nanovgo.Context) (int, int) {
	grid := g.ComputeLayout(widget, ctx)

	w := g.margin*2 + maxI(len(grid[0])-1, 0)*g.spacing[0]
	for _, v := range grid[0] {
		w += v
	}
	h := g.margin*2 + maxI(len(grid[1])-1, 0)*g.spacing[1]
	for _, v := range grid[1] {
		h += v
	}
	if _, ok := widget.(*Window); ok {
		h += widget.Theme().WindowHeaderHeight - g.margin/2
	}
	return w, h
}

func (g *GridLayout) ComputeLayout(widget Widget, ctx *nanovgo.Context) [][]int {
	axis1 := int(g.orientation)
	axis2 := (int(g.orientation) + 1) % 2
	numChildren := widget.ChildCount()
	visibleChildren := 0
	for _, child := range widget.Children() {
		if child.Visible() {
			visibleChildren++
		}
	}
	dim := make([]int, 2)
	dim[axis1] = g.resolution
	dim[axis2] = int((visibleChildren + g.resolution - 1) / g.resolution)

	grid := make([][]int, 2)
	grid[axis1] = make([]int, dim[axis1])
	grid[axis2] = make([]int, dim[axis2])

	child := 0
	children := widget.Children()
	for i2 := 0; i2 < dim[axis2]; i2++ {
		for i1 := 0; i1 < dim[axis1]; i1++ {
			var w Widget
			for {
				if child >= numChildren {
					return grid
				}
				w = children[child]
				child++
				if w.Visible() {
					break
				}
			}
			pw, ph := w.PreferredSize(w, ctx)
			fw, fh := w.FixedSize()
			targetSize := []int{
				toI(fw > 0, fw, pw),
				toI(fh > 0, fh, ph),
			}
			grid[axis1][i1] = maxI(grid[axis1][i1], targetSize[axis1])
			grid[axis2][i2] = maxI(grid[axis2][i2], targetSize[axis2])
		}
	}
	return grid
}
