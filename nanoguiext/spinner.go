package nanoguiext

import (
	"fmt"
	"github.com/shibukawa/glfw"
	"github.com/shibukawa/nanogui.go"
	"github.com/shibukawa/nanovgo"
	"math"
	"runtime"
)

type SpinnerState int

const (
	SpinnerStop SpinnerState = iota
	SpinnerFadeIn
	SpinnerFadeOut
)

func finalizeSpinner(spinner *Spinner) {
	if spinner.filter != nil {
		parent := spinner.filter.Parent()
		parent.RemoveChild(spinner.filter)
		spinner.filter = nil
	}
}

type Spinner struct {
	nanogui.WidgetImplement
	filter *SpinnerFilter
}

func NewSpinner(parent nanogui.Widget) *Spinner {
	spinner := &Spinner{}
	nanogui.InitWidget(spinner, parent)
	spinner.SetVisible(false)

	screen, ok := parent.(*nanogui.Screen)
	if !ok {
		screen = parent.FindWindow().Parent().(*nanogui.Screen)
	}

	filter := &SpinnerFilter{
		state:     SpinnerStop,
		c1:        18,
		c2:        24,
		num:       25,
		speed:     1,
		lineWidth: 3.0,
	}
	nanogui.InitWidget(filter, screen)
	spinner.filter = filter
	filter.SetVisible(false)
	runtime.SetFinalizer(spinner, finalizeSpinner)

	return spinner
}

func (s *Spinner) SetActive(flag bool) {
	if flag == s.Active() {
		return
	}
	s.filter.startTime = nanogui.GetTime()
	if flag {
		s.filter.state = SpinnerFadeIn
		s.filter.RequestFocus(s.filter)
		s.filter.SetVisible(true)
	} else {
		s.filter.state = SpinnerFadeOut
	}
}

func (s *Spinner) Active() bool {
	return s.filter.isActive()
}

func (s *Spinner) SetRadius(c1, c2 int) {
	s.filter.c1 = float32(c1)
	s.filter.c2 = float32(c2)
}

func (s *Spinner) Radius() (int, int) {
	return int(s.filter.c1), int(s.filter.c2)
}

func (s *Spinner) SetLineCount(num int) {
	s.filter.num = num
}

func (s *Spinner) LineCount() int {
	return s.filter.num
}

func (s *Spinner) SetSpeed(speed float32) {
	s.filter.speed = speed
}

func (s *Spinner) Speed() float32 {
	return s.filter.speed
}

func (s *Spinner) SetLineWidth(lineWidth float32) {
	s.filter.lineWidth = lineWidth
}

func (s *Spinner) LineWidth() float32 {
	return s.filter.lineWidth
}

func (s *Spinner) OnPerformLayout(self nanogui.Widget, ctx *nanovgo.Context) {
	s.filter.SetPosition(s.Parent().AbsolutePosition())
}

func (s *Spinner) PreferredSize(self nanogui.Widget, ctx *nanovgo.Context) (int, int) {
	return 0, 0
}

func (s *Spinner) Draw(self nanogui.Widget, ctx *nanovgo.Context) {
}

func (s *Spinner) IsPositionAbsolute() bool {
	return true
}

func (s *Spinner) String() string {
	return fmt.Sprintf("Spinner")
}

type SpinnerFilter struct {
	nanogui.WidgetImplement
	startTime float32
	state     SpinnerState
	c1, c2    float32
	num       int
	speed     float32
	lineWidth float32
}

func (sf *SpinnerFilter) isActive() bool {
	currentTime := nanogui.GetTime() - sf.startTime
	return sf.state == SpinnerFadeIn || (sf.state == SpinnerFadeOut && currentTime < 1.0)
}

func (sf *SpinnerFilter) IsPositionAbsolute() bool {
	return true
}

func (sf *SpinnerFilter) PreferredSize(self nanogui.Widget, ctx *nanovgo.Context) (int, int) {
	if sf.isActive() {
		fw, fh := sf.Parent().Size()
		if window, ok := sf.Parent().(*nanogui.Window); ok {
			hh := window.Theme().WindowHeaderHeight
			fh -= hh
		}
		return fw, fh
	} else {
		return 0, 0
	}
}

func (sf *SpinnerFilter) Draw(self nanogui.Widget, ctx *nanovgo.Context) {
	if sf.isActive() {
		var py int
		fw, fh := sf.Parent().Size()
		if window, ok := sf.Parent().(*nanogui.Window); ok {
			hh := window.Theme().WindowHeaderHeight
			py += hh
			fh -= hh
		}
		sf.SetPosition(0, py)
		sf.SetSize(fw, fh)

		currentTime := nanogui.GetTime() - sf.startTime

		var alpha float32
		var showSpinner bool
		if sf.state == SpinnerFadeIn {
			if currentTime > 1 {
				alpha = 0.7
				showSpinner = true
			} else {
				alpha = currentTime * 0.7
			}
		} else {
			if currentTime > 1 {
				alpha = 0.7
			} else {
				alpha = (1.0 - currentTime) * 0.7
			}
		}
		ctx.Save()
		ctx.BeginPath()
		ctx.SetFillColor(nanovgo.MONOf(0, alpha))
		ctx.Rect(0, float32(py), float32(fw), float32(fh))
		ctx.Fill()
		if showSpinner {
			cx := float32(fw / 2)
			cy := float32(py + fh/2)
			rotation := 2 * math.Pi * float64(currentTime*float32(sf.speed)*float32(sf.num)) / float64(sf.num)
			dr := float64(2 * math.Pi / float64(sf.num))
			ctx.SetStrokeWidth(sf.lineWidth)
			for i := 0; i < sf.num; i++ {
				ctx.BeginPath()
				ctx.MoveTo(cx+float32(math.Cos(rotation))*sf.c1, cy+float32(math.Sin(rotation))*sf.c1)
				ctx.LineTo(cx+float32(math.Cos(rotation))*sf.c2, cy+float32(math.Sin(rotation))*sf.c2)
				ctx.SetStrokeColor(nanovgo.MONOf(1.0, float32(i)/float32(sf.num)))
				ctx.Stroke()
				rotation += dr
			}
		}
		ctx.Restore()
	} else {
		sf.SetSize(0, 0)
		sf.SetVisible(false)
		return
	}
}

func (sf *SpinnerFilter) MouseButtonEvent(self nanogui.Widget, x, y int, button glfw.MouseButton, down bool, modifier glfw.ModifierKey) bool {
	return true
}

func (sf *SpinnerFilter) MouseMotionEvent(self nanogui.Widget, x, y, relX, relY, button int, modifier glfw.ModifierKey) bool {
	return true
}

func (sf *SpinnerFilter) MouseDragEvent(self nanogui.Widget, x, y, relX, relY, button int, modifier glfw.ModifierKey) bool {
	return true
}

func (sf *SpinnerFilter) MouseEnterEvent(self nanogui.Widget, x, y int, enter bool) bool {
	return true
}

func (sf *SpinnerFilter) ScrollEvent(self nanogui.Widget, x, y, relX, relY int) bool {
	return true
}

func (sf *SpinnerFilter) FocusEvent(self nanogui.Widget, f bool) bool {
	return true
}

func (sf *SpinnerFilter) KeyboardEvent(self nanogui.Widget, key glfw.Key, scanCode int, action glfw.Action, modifier glfw.ModifierKey) bool {
	return true
}

func (sf *SpinnerFilter) KeyboardCharacterEvent(self nanogui.Widget, codePoint rune) bool {
	return true
}

func (sf *SpinnerFilter) IMEPreeditEvent(self nanogui.Widget, text []rune, blocks []int, focusedBlock int) bool {
	return true
}

func (sf *SpinnerFilter) IMEStatusEvent(self nanogui.Widget) bool {
	return true
}

func (sf *SpinnerFilter) String() string {
	return sf.StringHelper("SpinnerFilter", "")
}
