package nanogui

import (
	"fmt"
	"github.com/shibukawa/glfw"
	"github.com/shibukawa/nanovgo"
	"math"
)

var sqrt3 float32 = float32(math.Sqrt(3))

type ColorWheelRegion int

const (
	RegionNone          ColorWheelRegion = 0
	RegionInnerTriangle ColorWheelRegion = 1
	RegionOuterCircle   ColorWheelRegion = 2
	RegionBoth          ColorWheelRegion = 3
)

type ColorWheel struct {
	WidgetImplement
	dragRegion                 ColorWheelRegion
	hue, saturation, lightness float32
	callback                   func(color nanovgo.Color)
}

func NewColorWheel(parent Widget, colors ...nanovgo.Color) *ColorWheel {
	var color nanovgo.Color
	switch len(colors) {
	case 0:
		color = nanovgo.RGBAf(1.0, 0.0, 0.0, 1.0)
	case 1:
		color = colors[0]
	default:
		panic("NewColorWheel can accept only one extra parameter (color)")
	}
	colorWheel := &ColorWheel{
		dragRegion: RegionNone,
	}
	InitWidget(colorWheel, parent)
	colorWheel.SetColor(color)
	return colorWheel
}

func (c *ColorWheel) SetCallback(callback func(color nanovgo.Color)) {
	c.callback = callback
}

func (c *ColorWheel) Color() nanovgo.Color {
	return nanovgo.HSL(c.hue, c.saturation, c.lightness)
}

func (c *ColorWheel) SetColor(color nanovgo.Color) {
	c.hue, c.saturation, c.lightness, _ = color.HSLA()
	c.calculatePosition()
}

func (c *ColorWheel) MouseDragEvent(self Widget, x, y, relX, relY, button int, modifier glfw.ModifierKey) bool {
	c.adjustPosition(x, y)
	return true
}

func (c *ColorWheel) MouseButtonEvent(self Widget, x, y int, button glfw.MouseButton, down bool, modifier glfw.ModifierKey) bool {
	c.WidgetImplement.MouseButtonEvent(self, x, y, button, down, modifier)

	if !c.enabled || button != glfw.MouseButton1 {
		return false
	}
	if down {
		c.adjustRegion(x, y)
		return c.dragRegion != RegionNone
	}
	c.dragRegion = RegionNone
	return true
}

func (c *ColorWheel) PreferredSize(self Widget, ctx *nanovgo.Context) (int, int) {
	return 100, 100
}

func (c *ColorWheel) Draw(ctx *nanovgo.Context) {
	c.WidgetImplement.Draw(ctx)

	if !c.visible {
		return
	}
	x := float32(c.x)
	y := float32(c.y)
	w := float32(c.w)
	h := float32(c.h)

	ctx.Save()
	defer ctx.Restore()

	cx := x + w*0.5
	cy := y + h*0.5
	r1 := toF(w < h, w, h)*0.5 - 5.0
	r0 := r1 * 0.75

	aeps := 0.7 / r1 // half a pixel arc length in radians (2pi cancels out).
	for i := 0; i < 6; i++ {
		a0 := float32(i)/6.0*nanovgo.PI*2.0 - aeps
		a1 := float32(i+1)/6.0*nanovgo.PI*2.0 + aeps
		ctx.BeginPath()
		ctx.Arc(cx, cy, r0, a0, a1, nanovgo.Clockwise)
		ctx.Arc(cx, cy, r1, a1, a0, nanovgo.CounterClockwise)
		ctx.ClosePath()

		sin1, cos1 := sinCosF(a0)
		sin2, cos2 := sinCosF(a1)
		ax := cx + cos1*(r0+r1)*0.5
		ay := cy + sin1*(r0+r1)*0.5
		bx := cx + cos2*(r0+r1)*0.5
		by := cy + sin2*(r0+r1)*0.5
		color1 := nanovgo.HSLA(a0/(nanovgo.PI*2), 1.0, 0.55, 255)
		color2 := nanovgo.HSLA(a1/(nanovgo.PI*2), 1.0, 0.55, 255)
		paint := nanovgo.LinearGradient(ax, ay, bx, by, color1, color2)
		ctx.SetFillPaint(paint)
		ctx.Fill()
	}

	ctx.BeginPath()
	ctx.Circle(cx, cy, r0-0.5)
	ctx.Circle(cx, cy, r1+0.5)
	ctx.SetStrokeColor(nanovgo.MONO(0, 64))
	ctx.Stroke()

	// Selector
	ctx.Save()
	defer ctx.Restore()
	ctx.Translate(cx, cy)
	ctx.Rotate(c.hue * nanovgo.PI * 2)

	// Marker on
	u := clampF(r1/50, 1.5, 4.0)
	ctx.SetStrokeWidth(u)
	ctx.BeginPath()
	ctx.Rect(r0-1, -2*u, r1-r0+2, 4*u)
	ctx.SetStrokeColor(nanovgo.MONO(255, 192))
	ctx.Stroke()

	paint := nanovgo.BoxGradient(r0-3, -5, r1-r0+6, 10, 2, 4, nanovgo.MONO(0, 128), nanovgo.MONO(0, 0))
	ctx.BeginPath()
	ctx.Rect(r0-2-10, -4-10, r1-r0+4+20, 8+20)
	ctx.Rect(r0-2, -4, r1-r0+4, 8)
	ctx.PathWinding(nanovgo.Hole)
	ctx.SetFillPaint(paint)
	ctx.Fill()

	// Center triangle
	r := r0 - 6
	sin1, cos1 := sinCosF(120.0 / 180.0 * nanovgo.PI)
	sin2, cos2 := sinCosF(-120.0 / 180.0 * nanovgo.PI)
	ax := cos1 * r
	ay := sin1 * r
	bx := cos2 * r
	by := sin2 * r
	ctx.BeginPath()
	ctx.MoveTo(r, 0)
	ctx.LineTo(ax, ay)
	ctx.LineTo(bx, by)
	ctx.ClosePath()
	triPaint1 := nanovgo.LinearGradient(r, 0, ax, ay, nanovgo.HSL(c.hue, 1.0, 0.5), nanovgo.MONO(255, 255))
	ctx.SetFillPaint(triPaint1)
	ctx.Fill()
	triPaint2 := nanovgo.LinearGradient((r+ax)*0.5, ay*0.5, bx, by, nanovgo.MONO(0, 0), nanovgo.MONO(0, 255))
	ctx.SetFillPaint(triPaint2)
	ctx.Fill()

	// selector circle on triangle
	px, py := c.calculatePosition()
	ctx.SetStrokeWidth(u)
	ctx.BeginPath()
	ctx.Circle(px, py, 2*u)
	ctx.SetStrokeColor(nanovgo.MONO(255, 192))
	ctx.Stroke()
}

func (c *ColorWheel) String() string {
	return fmt.Sprintf("ColorWheel [%d,%d-%d,%d] - h:%f s:%f l:%f", c.x, c.y, c.w, c.h, c.hue, c.saturation, c.lightness)
}

var sinOneThird float32 = float32(math.Sin(math.Pi * 2.0 / 3.0))
var cosOneThird float32 = float32(math.Cos(math.Pi * 2.0 / 3.0))
var sinTwoThird float32 = float32(math.Sin(-math.Pi * 2.0 / 3.0))
var cosTwoThird float32 = float32(math.Cos(-math.Pi * 2.0 / 3.0))

func (c *ColorWheel) calculatePosition() (float32, float32) {
	w := float32(c.w)
	h := float32(c.h)
	hw := w * 0.5
	hh := h * 0.5
	r1 := toF(w < h, hw, hh) - 5.0
	radius := r1*0.75 - 6

	// Colored point
	hx := radius
	// Black point
	sx := cosTwoThird * radius
	sy := -sinTwoThird * radius
	// White point
	vx := cosOneThird * radius
	vy := -sinOneThird * radius
	// Current point
	mx := (sx + vx) / 2.0
	my := (sy + vy) / 2.0
	a := (1.0 - 2.0*absF(c.lightness-0.5)) * c.saturation
	var px, py float32
	px = sx + (hx-mx)*a
	py = -(sy + (vy-sy)*c.lightness + (-my)*a)
	return px, py
}

func (c *ColorWheel) adjustRegion(px, py int) {
	x := float32(px - c.x)
	y := float32(py - c.y)
	w := float32(c.w)
	h := float32(c.h)
	hw := w * 0.5
	hh := h * 0.5
	r1 := toF(w < h, hw, hh) - 5.0
	radius := r1*0.75 - 6
	x -= hw
	y -= hh

	mr := sqrtF(x*x + y*y)

	if mr >= radius && mr <= r1 {
		c.dragRegion = RegionOuterCircle
		c.adjustPosition(px-c.x, py-c.y)
	} else if mr < radius {
		c.dragRegion = RegionInnerTriangle
		c.adjustPosition(px-c.x, py-c.y)
	} else {
		c.dragRegion = RegionNone
	}
}

func (c *ColorWheel) adjustPosition(px, py int) {
	if c.dragRegion == RegionNone {
		return
	}
	x := float32(px)
	y := float32(py)
	w := float32(c.w)
	h := float32(c.h)
	hw := w * 0.5
	hh := h * 0.5
	r1 := toF(w < h, hw, hh) - 5.0
	r0 := r1 * 0.75
	x -= hw
	y -= hh
	radius := r0 - 6

	rad := math.Atan2(float64(y), float64(x))
	if rad < 0 {
		rad += 2 * math.Pi
	}

	if c.dragRegion == RegionOuterCircle {
		c.hue = float32(rad / (2 * math.Pi))
		if c.callback != nil {
			c.callback(c.Color())
		}
	} else if c.dragRegion == RegionInnerTriangle {
		rad0 := math.Mod(rad+2*math.Pi*float64(1-c.hue), 2*math.Pi)
		rad1 := math.Mod(rad0, 2/3*math.Pi) - math.Pi/3
		a := 0.5 * radius
		b := float32(math.Tan(rad1)) * a
		r := sqrtF(x*x + y*y)
		maxR := sqrtF(a*a + b*b)

		if r > maxR {
			dx := float32(math.Tan(rad1)) * r
			rad2 := math.Atan(float64(dx / maxR))
			if rad2 > math.Pi/3 {
				rad2 = math.Pi / 3
			} else if rad2 < -math.Pi/3 {
				rad2 = -math.Pi / 3
			}
			rad += rad2 - rad1

			rad0 = math.Mod(rad+2*math.Pi-float64(c.hue)*2*math.Pi, 2*math.Pi)
			rad1 = math.Mod(rad0, (2/3)*math.Pi) - (math.Pi / 3)
			b = float32(math.Tan(rad1)) * a
			maxR = sqrtF(a*a + b*b) // Pythagoras
			r = maxR
		}
		sin, cos := math.Sincos(rad0)

		c.lightness = clampF(((float32(sin)*r)/radius/sqrt3)+0.5, 0.0, 1.0)

		widthShare := 1 - (absF(c.lightness-0.5) * 2)
		s := (((float32(cos) * r) + (radius / 2)) / (1.5 * radius)) / widthShare
		c.saturation = clampF(s, 0.0, 1.0)

		if c.callback != nil {
			c.callback(c.Color())
		}
	}
}

// https://github.com/timjb/colortriangle/blob/master/colortriangle.js
