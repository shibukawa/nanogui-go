package nanogui

import (
	"fmt"
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

	var yOffset int

	if _, ok := widget.(*Window); ok {
		if b.orientation == Vertical {
			position += widget.Theme().WindowHeaderHeight - b.margin/2
		} else {
			yOffset = widget.Theme().WindowHeaderHeight
		}
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
		pos[1] = yOffset
		pos[axis1] = position

		switch b.alignment {
		case Minimum:
			pos[axis2] += b.margin
		case Middle:
			pos[axis2] += (containerSize[axis2] - yOffset - targetSize[axis2]) / 2
		case Maximum:
			pos[axis2] += containerSize[axis2] - yOffset - targetSize[axis2] - b.margin*2
		case Fill:
			pos[axis2] += b.margin
			if fs[axis2] > 0 {
				targetSize[axis2] = fs[axis2]
			} else {
				targetSize[axis2] = containerSize[axis2] - yOffset - b.margin*2
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

	axis2Offset := 0
	if _, ok := widget.(*Window); ok {
		if b.orientation == Vertical {
			size[1] += widget.Theme().WindowHeaderHeight - b.margin/2
		} else {
			axis2Offset = widget.Theme().WindowHeaderHeight
		}
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
		size[axis2] = maxI(size[axis2], targetSize[axis2]+2*b.margin+axis2Offset)
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
	grid := g.computeLayout(widget, ctx)
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
	grid := g.computeLayout(widget, ctx)

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

func (g *GridLayout) computeLayout(widget Widget, ctx *nanovgo.Context) [][]int {
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

type Anchor struct {
	pos   [2]uint8
	size  [2]uint8
	align [2]Alignment
}

func NewAnchor(x, y int, aligns ...Alignment) Anchor {
	a := Anchor{
		pos:  [2]uint8{uint8(x), uint8(y)},
		size: [2]uint8{1, 1},
	}
	switch len(aligns) {
	case 0:
		a.align[0] = Fill
		a.align[1] = Fill
	case 1:
		a.align[0] = aligns[0]
		a.align[1] = Fill
	case 2:
		a.align[0] = aligns[0]
		a.align[1] = aligns[1]
	default:
		panic("NewAnchor can accept extra parameter upto 2 (hAlign, vAlign).")
	}
	return a
}

func NewAnchorWithSize(x, y, w, h int, aligns ...Alignment) Anchor {
	a := Anchor{
		pos:  [2]uint8{uint8(x), uint8(y)},
		size: [2]uint8{uint8(x), uint8(y)},
	}
	switch len(aligns) {
	case 0:
		a.align[0] = Fill
		a.align[1] = Fill
	case 1:
		a.align[0] = aligns[0]
		a.align[1] = Fill
	case 2:
		a.align[0] = aligns[0]
		a.align[1] = aligns[1]
	default:
		panic("NewAnchorWithSize can accept extra parameter upto 2 (hAlign, vAlign).")
	}
	return a
}

func (a *Anchor) String() string {
	return fmt.Sprintf("Format[pos=(%i, %i), size=(%i, %i), align=(%i, %i)]",
		a.pos[0], a.pos[1], a.size[0], a.size[1], int(a.align[0]), int(a.align[1]))
}

type AdvancedGridLayout struct {
	cols       []int
	rows       []int
	colStretch []float32
	rowStretch []float32
	anchors    map[Widget]Anchor
	margin     int
}

func NewAdvancedGridLayout(sizes ...[]int) *AdvancedGridLayout {
	var rows []int
	var cols []int
	switch len(sizes) {
	case 0:
	case 1:
		cols = sizes[0]
	case 2:
		cols = sizes[0]
		rows = sizes[1]
	default:
		panic("NewBoxLayout can accept extra parameter upto 2 (cols, rows).")
	}
	return &AdvancedGridLayout{
		cols:       cols,
		rows:       rows,
		colStretch: make([]float32, len(cols)),
		rowStretch: make([]float32, len(rows)),
		anchors:    make(map[Widget]Anchor),
	}
}

func (a *AdvancedGridLayout) Margin() int {
	return a.margin
}

func (a *AdvancedGridLayout) SetMargin(m int) {
	a.margin = m
}

func (a *AdvancedGridLayout) ColCount() int {
	return len(a.cols)
}

func (a *AdvancedGridLayout) RowCount() int {
	return len(a.rows)
}

func (a *AdvancedGridLayout) AppendCol(size int, defaultStretch ...float32) {
	var stretch float32
	switch len(defaultStretch) {
	case 0:
	case 1:
		stretch = defaultStretch[0]
	default:
		panic("AppendCol can accept only one extra parameter(stretch).")
	}
	a.cols = append(a.cols, size)
	a.colStretch = append(a.colStretch, stretch)
}

func (a *AdvancedGridLayout) AppendRow(size int, defaultStretch ...float32) {
	var stretch float32
	switch len(defaultStretch) {
	case 0:
	case 1:
		stretch = defaultStretch[0]
	default:
		panic("AppendRow can accept only one extra parameter(stretch).")
	}
	a.rows = append(a.rows, size)
	a.rowStretch = append(a.rowStretch, stretch)
}

func (a *AdvancedGridLayout) SetRowStretch(index int, stretch float32) {
	a.rowStretch[index] = stretch
}

func (a *AdvancedGridLayout) SetColStretch(index int, stretch float32) {
	a.colStretch[index] = stretch
}

func (a *AdvancedGridLayout) SetAnchor(widget Widget, anchor Anchor) {
	a.anchors[widget] = anchor
}

func (a *AdvancedGridLayout) Anchor(widget Widget) Anchor {
	return a.anchors[widget]
}

func (a *AdvancedGridLayout) OnPerformLayout(widget Widget, ctx *nanovgo.Context) {
	grid := a.computeLayout(widget, ctx)
	grid[0] = append([]int{a.margin}, grid[0]...)
	if _, ok := widget.(*Window); ok {
		grid[1] = append([]int{widget.Theme().WindowHeaderHeight + a.margin/2}, grid[1]...)
	} else {
		grid[1] = append([]int{a.margin}, grid[1]...)
	}
	for axis := 0; axis < 2; axis++ {
		for i := 1; i < len(grid[axis]); i++ {
			grid[axis][i] += grid[axis][i-1]
		}
		for _, w := range widget.Children() {
			if !w.Visible() {
				continue
			}
			anchor := a.Anchor(w)
			itemPos := grid[axis][anchor.pos[axis]]
			cellSize := grid[axis][anchor.pos[axis]+anchor.size[axis]] - itemPos
			pw, ph := w.PreferredSize(w, ctx)
			fw, fh := w.FixedSize()
			var targetSize int
			if axis == 0 {
				targetSize = toI(fw > 0, fw, pw)
			} else {
				targetSize = toI(fh > 0, fh, ph)
			}
			switch anchor.align[axis] {
			case Minimum:
			case Middle:
				itemPos += (cellSize - targetSize) / 2
			case Maximum:
				itemPos += cellSize - targetSize
			case Fill:
				if axis == 0 {
					targetSize = toI(fw > 0, fw, cellSize)
				} else {
					targetSize = toI(fh > 0, fh, cellSize)
				}
			}
			posX, posY := w.Position()
			sizeW, sizeH := w.Size()
			if axis == 0 {
				posX = itemPos
				sizeW = targetSize
			} else {
				posY = itemPos
				sizeH = targetSize
			}
			w.SetPosition(posX, posY)
			w.SetSize(sizeW, sizeH)
			w.OnPerformLayout(w, ctx)
		}
	}
}

func (a *AdvancedGridLayout) PreferredSize(widget Widget, ctx *nanovgo.Context) (int, int) {
	grid := a.computeLayout(widget, ctx)
	sizeW := a.margin * 2
	sizeH := a.margin * 2
	for _, size := range grid[0] {
		sizeW += size
	}
	for _, size := range grid[1] {
		sizeH += size
	}
	if _, ok := widget.(*Window); ok {
		sizeH += widget.Theme().WindowHeaderHeight - a.margin/2
	}
	return sizeW, sizeH
}

func (a *AdvancedGridLayout) computeLayout(widget Widget, ctx *nanovgo.Context) [][]int {
	var grids [][]int = [][]int{[]int{}, []int{}}
	fw, fh := widget.FixedSize()
	containerW := toI(fw > 0, fw, widget.Width())
	containerH := toI(fh > 0, fh, widget.Height())

	extraX := 2 * a.margin
	extraY := 2 * a.margin

	if _, ok := widget.(*Window); ok {
		extraY += widget.Theme().WindowHeaderHeight - a.margin/2
	}

	containerW -= extraX
	containerH -= extraY

	for axis := 0; axis < 2; axis++ {
		var sizes []int
		var stretch []float32
		if axis == 0 {
			sizes = a.cols
			stretch = a.colStretch
		} else {
			sizes = a.rows
			stretch = a.rowStretch
		}
		grid := make([]int, len(sizes))
		copy(grid, sizes)
		grids[axis] = grid

		for phase := 0; phase < 2; phase++ {
			for widget, anchor := range a.anchors {
				if !widget.Visible() {
					continue
				}
				if (anchor.size[axis]) == 1 != (phase == 0) {
					continue
				}
				pw, ph := widget.PreferredSize(widget, ctx)
				ps := toI(axis == 0, pw, ph)
				fw, fh := widget.FixedSize()
				fs := toI(axis == 0, fw, fh)
				targetSize := toI(fs > 0, fs, ps)
				if int(anchor.pos[axis])+int(anchor.size[axis]) > len(grid) {
					panic("Advanced grid layout: widget is out of bounds: " + anchor.String())
				}
				currentSize := 0
				var totalStretch float32
				for i := anchor.pos[axis]; i < anchor.pos[axis]+anchor.size[axis]; i++ {
					if sizes[i] == 0 && anchor.size[axis] == 1 {
						grid[i] = maxI(grid[i], targetSize)
					}
					currentSize += grid[i]
					totalStretch += stretch[i]
				}
				if targetSize <= currentSize {
					continue
				}
				if totalStretch == 0 {
					panic("Advanced grid layout: no space to place widget: " + anchor.String())
				}
				var amt float32 = float32(targetSize-currentSize) / totalStretch
				for i := anchor.pos[axis]; i < anchor.pos[axis]+anchor.size[axis]; i++ {
					grid[i] += int(amt * stretch[i])
				}
			}
		}
		var currentSize int
		for _, val := range grid {
			currentSize += val
		}
		var totalStretch float32
		for _, val := range stretch {
			totalStretch += val
		}
		var containerSize float32
		if axis == 0 {
			containerSize = float32(containerW)
		} else {
			containerSize = float32(containerH)
		}
		if float32(currentSize) >= containerSize || totalStretch == 0 {
			continue
		}
		amt := (containerSize - float32(currentSize)) / totalStretch
		for i := range grid {
			grid[i] += int(amt*stretch[i] + 0.5)
		}
	}
	return grids
}
