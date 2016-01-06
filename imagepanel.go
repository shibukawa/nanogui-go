package nanogui

import (
	"github.com/shibukawa/glfw"
	"github.com/shibukawa/nanovgo"
)

type Image struct {
	ImageID int
	Name    string
}

type ImagePanel struct {
	WidgetImplement

	images     []Image
	callback   func(int)
	thumbSize  int
	spacing    int
	margin     int
	mouseIndex int
}

func NewImagePanel(parent Widget) *ImagePanel {
	panel := &ImagePanel{
		thumbSize:  64,
		spacing:    10,
		margin:     10,
		mouseIndex: -1,
	}

	InitWidget(panel, parent)
	return panel
}

func (i *ImagePanel) Images() []Image {
	return i.images
}

func (i *ImagePanel) SetImages(images []Image) {
	i.images = images
}

func (i *ImagePanel) SetCallback(callback func(int)) {
	i.callback = callback
}

func (i *ImagePanel) MouseButtonEvent(self Widget, x, y int, button glfw.MouseButton, down bool, modifier glfw.ModifierKey) bool {
	index := i.indexForPosition(x, y)
	if index >= 0 && i.callback != nil && down {
		i.callback(index)
	}
	return true
}

func (i *ImagePanel) MouseMotionEvent(self Widget, x, y, relX, relY, button int, modifier glfw.ModifierKey) bool {
	i.mouseIndex = i.indexForPosition(x, y)
	return true
}

func (i *ImagePanel) PreferredSize(self Widget, ctx *nanovgo.Context) (int, int) {
	cols, rows := i.gridSize()
	w := cols*i.thumbSize + (cols-1)*i.spacing + 2*i.margin
	h := rows*i.thumbSize + (rows-1)*i.spacing + 2*i.margin
	return w, h
}

func (i *ImagePanel) Draw(self Widget, ctx *nanovgo.Context) {
	cols, _ := i.gridSize()

	x := float32(i.x)
	y := float32(i.y)
	thumbSize := float32(i.thumbSize)

	for j, image := range i.images {
		pX := x + float32(i.margin+(j%cols)*(i.thumbSize+i.spacing))
		pY := y + float32(i.margin+(j/cols)*(i.thumbSize+i.spacing))

		imgW, imgH, _ := ctx.ImageSize(image.ImageID)
		var iw, ih, ix, iy float32
		if imgW < imgH {
			iw = thumbSize
			ih = iw * float32(imgH) / float32(imgW)
			ix = 0
			iy = -(ih - thumbSize) * 0.5
		} else {
			ih = thumbSize
			iw = ih * float32(imgH) / float32(imgW)
			iy = 0
			ix = -(iw - thumbSize) * 0.5
		}
		imgPaint := nanovgo.ImagePattern(pX+ix, pY+iy, iw, ih, 0, image.ImageID, toF(i.mouseIndex == j, 1.0, 0.7))
		ctx.BeginPath()
		ctx.RoundedRect(pX, pY, thumbSize, thumbSize, 5)
		ctx.SetFillPaint(imgPaint)
		ctx.Fill()

		shadowPaint := nanovgo.BoxGradient(pX-1, pY, thumbSize+2, thumbSize+2, 5, 3, nanovgo.MONO(0, 128), nanovgo.MONO(0, 0))

		ctx.BeginPath()
		ctx.Rect(pX-5, pY-5, thumbSize+10, thumbSize+10)
		ctx.RoundedRect(pX, pY, thumbSize, thumbSize, 6)
		ctx.PathWinding(nanovgo.Hole)
		ctx.SetFillPaint(shadowPaint)
		ctx.Fill()

		ctx.BeginPath()
		ctx.RoundedRect(pX+0.5, pY+0.5, thumbSize-1, thumbSize-1, 4-0.5)
		ctx.SetStrokeWidth(1.0)
		ctx.SetStrokeColor(nanovgo.MONO(255, 80))
		ctx.Stroke()
	}
}

func (i *ImagePanel) String() string {
	return i.StringHelper("ImagePanel", "")
}

func (i *ImagePanel) gridSize() (int, int) {
	nCols := 1 + maxI(0, int(float32(i.w-2*i.margin-i.thumbSize)/float32(i.thumbSize+i.spacing)))
	nRows := (len(i.images) + nCols - 1) / nCols
	return nCols, nRows

}

func (i *ImagePanel) indexForPosition(x, y int) int {
	pX := float32(x-i.margin) / float32(i.thumbSize+i.margin)
	pY := float32(y-i.margin) / float32(i.thumbSize+i.margin)
	iconRegion := float32(i.thumbSize) / float32(i.thumbSize+i.spacing)
	overImage := pX-floorF(pX) < iconRegion && pY-floorF(pY) < iconRegion
	gridPosX := int(pX)
	gridPosY := int(pY)
	gridCols, gridRows := i.gridSize()
	overImage = overImage && gridPosX >= 0 && gridPosY >= 0 && gridPosX < gridCols && gridPosY < gridRows
	if overImage {
		return gridPosX + gridPosY*gridCols
	}
	return -1
}
