package nanogui

import (
	"github.com/shibukawa/nanovgo"
)

type ImageSizePolicy int

const (
	ImageSizePolicyFixed ImageSizePolicy = iota
	ImageSizePolicyExpand
)

type ImageView struct {
	WidgetImplement

	image  int
	policy ImageSizePolicy
}

func NewImageView(parent Widget, images ...int) *ImageView {
	var image int
	switch len(images) {
	case 0:
	case 1:
		image = images[0]
	default:
		panic("NewImageView can accept only one extra parameter (image)")
	}

	imageView := &ImageView{
		image:  image,
		policy: ImageSizePolicyFixed,
	}
	InitWidget(imageView, parent)
	return imageView
}

func (i *ImageView) Image() int {
	return i.image
}

func (i *ImageView) SetImage(image int) {
	i.image = image
}

func (i *ImageView) Policy() ImageSizePolicy {
	return i.policy
}

func (i *ImageView) SetPolicy(policy ImageSizePolicy) {
	i.policy = policy
}

func (i *ImageView) PreferredSize(self Widget, ctx *nanovgo.Context) (int, int) {
	if i.image == 0 {
		return 0, 0
	}
	w, h, _ := ctx.ImageSize(i.image)
	return w, h
}

func (i *ImageView) Draw(self Widget, ctx *nanovgo.Context) {
	if i.image == 0 {
		return
	}
	x := float32(i.x)
	y := float32(i.y)
	ow := float32(i.w)
	oh := float32(i.h)

	var w, h float32
	{
		iw, ih, _ := ctx.ImageSize(i.image)
		w = float32(iw)
		h = float32(ih)
	}

	if i.policy == ImageSizePolicyFixed {
		if ow < w {
			h = float32(int(h * ow / w))
			w = ow
		}
		if oh < h {
			w = float32(int(w * oh / h))
			h = oh
		}
	} else { // mPolicy == Expand
		// expand to width
		h = float32(int(h * ow / w))
		w = ow
		// shrink to height, if necessary
		if oh < h {
			w = float32(int(w * oh / h))
			h = oh
		}
	}

	imgPaint := nanovgo.ImagePattern(x, y, w, h, 0, i.image, 1.0)

	ctx.BeginPath()
	ctx.Rect(x, y, w, h)
	ctx.SetFillPaint(imgPaint)
	ctx.Fill()
}

func (i *ImageView) String() string {
	return i.StringHelper("ImageView", "")
}
