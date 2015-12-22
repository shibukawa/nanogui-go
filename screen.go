package nanogui

import (
	"bytes"
	"fmt"
	"github.com/goxjs/gl"
	"github.com/shibukawa/glfw"
	"github.com/shibukawa/nanovgo"
	"runtime"
)

var nanoguiScreens map[*glfw.Window]*Screen = map[*glfw.Window]*Screen{}

type Screen struct {
	WidgetImplement
	window                 *glfw.Window
	context                *nanovgo.Context
	cursors                [3]int
	cursor                 Cursor
	focusPath              []Widget
	fbW, fbH               int
	pixelRatio             float32
	modifiers              glfw.ModifierKey
	mouseState             int
	mousePosX, mousePosY   int
	dragActive             bool
	dragWidget             Widget
	lastInteraction        float32
	backgroundColor        nanovgo.Color
	caption                string
	shutdownGLFWOnDestruct bool

	drawContentsCallback func()
	dropEventCallback    func([]string) bool
	resizeEventCallback  func(x, y int) bool
}

func NewScreen(width, height int, caption string, resizable, fullScreen bool) *Screen {
	screen := &Screen{
		//cursor:  glfw.CursorNormal,
		caption: caption,
	}

	if runtime.GOARCH == "js" {
		glfw.WindowHint(glfw.Hint(0x00021101), 1) // enable stencil for nanovgo
	}
	glfw.WindowHint(glfw.Samples, 4)
	//glfw.WindowHint(glfw.RedBits, 8)
	//glfw.WindowHint(glfw.GreenBits, 8)
	//glfw.WindowHint(glfw.BlueBits, 8)
	glfw.WindowHint(glfw.AlphaBits, 8)
	//glfw.WindowHint(glfw.StencilBits, 8)
	//glfw.WindowHint(glfw.DepthBits, 8)
	//glfw.WindowHint(glfw.Visible, 0)
	if resizable {
		//glfw.WindowHint(glfw.Resizable, 1)
	} else {
		//glfw.WindowHint(glfw.Resizable, 0)
	}

	var err error
	if fullScreen {
		monitor := glfw.GetPrimaryMonitor()
		mode := monitor.GetVideoMode()
		screen.window, err = glfw.CreateWindow(mode.Width, mode.Height, caption, monitor, nil)
	} else {
		screen.window, err = glfw.CreateWindow(width, height, caption, nil, nil)
	}
	if err != nil {
		panic(err)
	}
	screen.window.MakeContextCurrent()
	gl.Viewport(0, 0, screen.fbW, screen.fbH)
	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)
	glfw.SwapInterval(0)
	screen.window.SwapBuffers()

	/* Poll for events once before starting a potentially
	   lengthy loading process. This is needed to be
	   classified as "interactive" by other software such
	   as iTerm2 */

	if runtime.GOOS == "darwin" {
		glfw.PollEvents()
	}

	screen.window.SetCursorPosCallback(func(w *glfw.Window, xpos, ypos float64) {
		if screen, ok := nanoguiScreens[w]; ok {
			screen.cursorPositionCallbackEvent(xpos, ypos)
		}
	})

	screen.window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		if screen, ok := nanoguiScreens[w]; ok {
			screen.mouseButtonCallbackEvent(button, action, mods)
		}
	})

	screen.window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scanCode int, action glfw.Action, mods glfw.ModifierKey) {
		if screen, ok := nanoguiScreens[w]; ok {
			screen.keyCallbackEvent(key, scanCode, action, mods)
		}
	})

	screen.window.SetCharCallback(func(w *glfw.Window, char rune) {
		if screen, ok := nanoguiScreens[w]; ok {
			screen.charCallbackEvent(char)
		}
	})

	screen.window.SetPreeditCallback(func(w *glfw.Window, text []rune, blocks []int, focusedBlock int) {
		if screen, ok := nanoguiScreens[w]; ok {
			screen.preeditCallbackEvent(text, blocks, focusedBlock)
		}
	})

	screen.window.SetIMEStatusCallback(func(w *glfw.Window) {
		if screen, ok := nanoguiScreens[w]; ok {
			screen.imeStatusCallbackEvent()
		}
	})

	screen.window.SetDropCallback(func(w *glfw.Window, names []string) {
		if screen, ok := nanoguiScreens[w]; ok {
			screen.dropCallbackEvent(names)
		}
	})

	screen.window.SetScrollCallback(func(w *glfw.Window, xoff float64, yoff float64) {
		if screen, ok := nanoguiScreens[w]; ok {
			screen.scrollCallbackEvent(float32(xoff), float32(yoff))
		}
	})

	screen.window.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		if screen, ok := nanoguiScreens[w]; ok {
			screen.resizeCallbackEvent(width, height)
		}
	})

	screen.Initialize(screen.window, true)
	InitWidget(screen, nil)
	return screen
}

func (s *Screen) Initialize(window *glfw.Window, shutdownGLFWOnDestruct bool) {
	s.window = window
	s.shutdownGLFWOnDestruct = shutdownGLFWOnDestruct
	s.w, s.h = window.GetSize()
	s.fbW, s.fbH = window.GetFramebufferSize()
	var err error
	s.context, err = nanovgo.NewContext(nanovgo.StencilStrokes | nanovgo.AntiAlias)
	if err != nil {
		panic(err)
	}
	s.visible = true //window.GetAttrib(glfw.Visible)
	s.theme = NewStandardTheme(s.context)
	s.mousePosX = 0
	s.mousePosY = 0
	s.mouseState = 0
	s.modifiers = 0
	s.dragActive = false
	s.lastInteraction = GetTime()
	nanoguiScreens[window] = s
	runtime.SetFinalizer(s, func(s *Screen) {
		delete(nanoguiScreens, window)
		if s.context != nil {
			s.context.Delete()
			s.context = nil
		}
		if s.window != nil && s.shutdownGLFWOnDestruct {
			s.window.Destroy()
			s.window = nil
		}
	})
}

// Caption() gets the window title bar caption
func (s *Screen) Caption() string {
	return s.caption
}

// SetCaption() sets the window title bar caption
func (s *Screen) SetCaption(caption string) {
	if s.caption != caption {
		s.window.SetTitle(caption)
		s.caption = caption
	}
}

// BackgroundColor() returns the screen's background color
func (s *Screen) BackgroundColor() nanovgo.Color {
	return s.backgroundColor
}

// SetBackgroundColor() sets the screen's background color
func (s *Screen) SetBackgroundColor(color nanovgo.Color) {
	s.backgroundColor = color
	s.backgroundColor.A = 1.0
}

// SetVisible() sets the top-level window visibility (no effect on full-screen windows)
func (s *Screen) SetVisible(flag bool) {
	if s.visible != flag {
		s.visible = flag
		if flag {
			s.window.Show()
		} else {
			s.window.Hide()
		}
	}
}

// SetSize() sets window size
func (s *Screen) SetSize(w, h int) {
	s.WidgetImplement.SetSize(w, h)
	s.window.SetSize(w, h)
}

// DrawAll() draws the Screen contents
func (s *Screen) DrawAll() {
	gl.ClearColor(s.backgroundColor.R, s.backgroundColor.G, s.backgroundColor.B, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)

	if s.drawContentsCallback != nil {
		s.drawContentsCallback()
	}
	s.drawWidgets()
	s.window.SwapBuffers()
}

// SetResizeEventCallback() sets window resize event handler
func (s *Screen) SetResizeEventCallback(callback func(x, y int) bool) {
	s.resizeEventCallback = callback
}

// SetDrawContentsCallback() sets event handler to use OpenGL draw call
func (s *Screen) SetDrawContentsCallback(callback func()) {
	s.drawContentsCallback = callback
}

// SetDropEventCallback() sets event handler to handle a file drop event
func (s *Screen) SetDropEventCallback(callback func(files []string) bool) {
	s.dropEventCallback = callback
}

// KeyboardEvent() is a default key event handler
func (s *Screen) KeyboardEvent(self Widget, key glfw.Key, scanCode int, action glfw.Action, modifiers glfw.ModifierKey) bool {
	if len(s.focusPath) > 1 {
		for i := len(s.focusPath) - 2; i >= 0; i-- {
			path := s.focusPath[i]
			if path.Focused() && path.KeyboardEvent(path, key, scanCode, action, modifiers) {
				return true
			}
		}
	}
	return false
}

// KeyboardCharacterEvent() is a text input event handler: codepoint is native endian UTF-32 format
func (s *Screen) KeyboardCharacterEvent(self Widget, codePoint rune) bool {
	if len(s.focusPath) > 1 {
		for i := len(s.focusPath) - 2; i >= 0; i-- {
			path := s.focusPath[i]
			if path.Focused() && path.KeyboardCharacterEvent(path, codePoint) {
				return true
			}
		}
	}
	return false
}

// IMEPreeditEvent() handles preedit text changes of IME (default implementation: do nothing)
func (s *Screen) IMEPreeditEvent(self Widget, text []rune, blocks []int, focusedBlock int) bool {
	if len(s.focusPath) > 1 {
		for i := len(s.focusPath) - 2; i >= 0; i-- {
			path := s.focusPath[i]
			if path.Focused() && path.IMEPreeditEvent(path, text, blocks, focusedBlock) {
				return true
			}
		}
	}
	return false
}

// IMEStatusEvent() handles IME status change event (default implementation: do nothing)
func (s *Screen) IMEStatusEvent(self Widget) bool {
	if len(s.focusPath) > 1 {
		for i := len(s.focusPath) - 2; i >= 0; i-- {
			path := s.focusPath[i]
			if path.Focused() && path.IMEStatusEvent(path) {
				return true
			}
		}
	}
	return false
}

// MousePosition() returns the last observed mouse position value
func (s *Screen) MousePosition() (int, int) {
	return s.mousePosX, s.mousePosY
}

// GLFWWindow() returns a pointer to the underlying GLFW window data structure
func (s *Screen) GLFWWindow() *glfw.Window {
	return s.window
}

// NVGContext() returns a pointer to the underlying nanoVGo draw context
func (s *Screen) NVGContext() *nanovgo.Context {
	return s.context
}

func (s *Screen) SetShutdownGLFWOnDestruct(v bool) {
	s.shutdownGLFWOnDestruct = v
}

func (s *Screen) ShutdownGLFWOnDestruct() bool {
	return s.shutdownGLFWOnDestruct
}

// UpdateFocus is an internal helper function
func (s *Screen) UpdateFocus(widget Widget) {
	for _, w := range s.focusPath {
		if w.Focused() {
			w.FocusEvent(w, false)
		}
	}
	s.focusPath = s.focusPath[:0]
	var window *Window
	for widget != nil {
		s.focusPath = append(s.focusPath, widget)
		if _, ok := widget.(*Window); ok {
			window = widget.(*Window)
		}
		widget = widget.Parent()
	}
	for _, w := range s.focusPath {
		w.FocusEvent(w, true)
	}
	if window != nil {
		s.MoveWindowToFront(window)
	}
}

// DisposeWindow is an internal helper function
func (s *Screen) DisposeWindow(window *Window) {
	find := false
	for _, w := range s.focusPath {
		if w == window {
			find = true
			break
		}
	}
	if find {
		s.focusPath = s.focusPath[:0]
	}
	if s.dragWidget == window {
		s.dragWidget = nil
	}
	s.RemoveChild(window)
}

// CenterWindow is an internal helper function
func (s *Screen) CenterWindow(window *Window) {
	w, h := window.Size()
	if w == 0 && h == 0 {
		window.SetSize(window.PreferredSize(window, s.context))
		window.OnPerformLayout(window, s.context)
	}
	x, y := window.Size()
	window.SetPosition((s.x-x)/2, (s.y-y)/2)
}

// MoveWindowToFront is an internal helper function
func (s *Screen) MoveWindowToFront(window IWindow) {
	s.RemoveChild(window)
	s.children = append(s.children, window)
	window.SetParent(s)
	changed := true
	for changed {
		baseIndex := 0
		for i, child := range s.children {
			if child == window {
				baseIndex = i
			}
		}
		changed = false
		for i, child := range s.children {
			pw, ok := child.(*Popup)
			if ok && pw.ParentWindow() == window && i < baseIndex {
				s.MoveWindowToFront(pw)
				changed = true
				break
			}
		}
	}
}

func (s *Screen) PreeditCursorPos() (int, int, int) {
	return s.window.GetPreeditCursorPos()
}

func (s *Screen) SetPreeditCursorPos(x, y, h int) {
	s.window.SetPreeditCursorPos(x, y, h)
}

func (s *Screen) drawWidgets() {
	if !s.visible {
		return
	}
	s.window.MakeContextCurrent()
	s.fbW, s.fbH = s.window.GetFramebufferSize()
	s.w, s.h = s.window.GetSize()
	gl.Viewport(0, 0, s.fbW, s.fbH)

	s.pixelRatio = float32(s.fbW) / float32(s.w)
	s.context.BeginFrame(s.w, s.h, s.pixelRatio)
	s.Draw(s.context)
	elapsed := GetTime() - s.lastInteraction

	if elapsed > 0.5 {
		// Draw tooltips
		widget := s.FindWidget(s, s.mousePosX, s.mousePosY)
		if widget != nil && widget.Tooltip() != "" {
			var tooltipWidth float32 = 150
			ctx := s.context
			ctx.SetFontFace(s.theme.FontNormal)
			ctx.SetFontSize(15.0)
			ctx.SetTextAlign(nanovgo.AlignCenter | nanovgo.AlignTop)
			ctx.SetTextLineHeight(1.1)
			posX, posY := widget.AbsolutePosition()
			posX += widget.Width() / 2
			posY += widget.Height() + 10
			bounds := ctx.TextBoxBounds(float32(posX), float32(posY), tooltipWidth, widget.Tooltip())
			ctx.SetGlobalAlpha(minF(1.0, 2*(elapsed-0.5)) * 0.8)
			ctx.BeginPath()
			ctx.SetFillColor(nanovgo.MONO(0, 255))
			h := (bounds[2] - bounds[0]) / 2
			ctx.RoundedRect(bounds[0]-4-h, bounds[1]-4, bounds[2]-bounds[0]+8, bounds[3]-bounds[1]+8, 3)
			px := (bounds[2]+bounds[0])/2 - h
			ctx.MoveTo(px, bounds[1]-10)
			ctx.LineTo(px+7, bounds[1]+1)
			ctx.LineTo(px-7, bounds[1]+1)
			ctx.Fill()

			ctx.SetFillColor(nanovgo.MONO(255, 255))
			ctx.SetFontBlur(0.0)
			ctx.TextBox(float32(posX)-h, float32(posY), tooltipWidth, widget.Tooltip())

		}
	}

	s.context.EndFrame()
}

func (s *Screen) cursorPositionCallbackEvent(x, y float64) bool {
	ret := false
	s.lastInteraction = GetTime()

	px := int(x) - 1
	py := int(y) - 2
	if !s.dragActive {
		widget := s.FindWidget(s, int(x), int(y))
		if widget != nil && widget.Cursor() != s.cursor {
			//s.cursor = widget.Cursor()
			//s.window.SetCursor()
		}
	} else {
		ax, ay := s.dragWidget.Parent().AbsolutePosition()
		ret = s.dragWidget.MouseDragEvent(s.dragWidget, px-ax, py-ay, px-s.mousePosX, py-s.mousePosY, s.mouseState, s.modifiers)
	}
	if !ret {
		ret = s.MouseMotionEvent(s, px, py, px-s.mousePosX, py-s.mousePosY, s.mouseState, s.modifiers)
	}
	s.mousePosX = px
	s.mousePosY = py
	return ret
}

func (s *Screen) mouseButtonCallbackEvent(button glfw.MouseButton, action glfw.Action, modifiers glfw.ModifierKey) bool {
	s.modifiers = modifiers
	s.lastInteraction = GetTime()

	if len(s.focusPath) > 1 {
		window, ok := s.focusPath[len(s.focusPath)-2].(*Window)
		if ok && window.Modal() {
			if !window.Contains(s.mousePosX, s.mousePosY) {
				return false
			}
		}
	}

	if action == glfw.Press {
		s.mouseState |= 1 << uint(button)
	} else {
		s.mouseState &= ^(1 << uint(button))
	}

	dropWidget := s.FindWidget(s, s.mousePosX, s.mousePosY)
	if s.dragActive && action == glfw.Release && dropWidget != s.dragWidget {
		ax, ay := s.dragWidget.Parent().AbsolutePosition()
		s.dragWidget.MouseButtonEvent(s.dragWidget, s.mousePosX-ax, s.mousePosY-ay, button, false, modifiers)
	}

	if dropWidget != nil && dropWidget.Cursor() != s.cursor {
		//s.cursor = widget.Cursor()
		//s.window.SetCursor()
	}

	if action == glfw.Press && button == glfw.MouseButton1 {
		s.dragWidget = s.FindWidget(s, s.mousePosX, s.mousePosY)
		if s.dragWidget == s {
			s.dragWidget = nil
		}
		s.dragActive = s.dragWidget != nil
		if !s.dragActive {
			s.UpdateFocus(nil)
		}
	} else {
		s.dragActive = false
		s.dragWidget = nil
	}
	return s.MouseButtonEvent(s, s.mousePosX, s.mousePosY, button, action == glfw.Press, modifiers)
}

func (s *Screen) keyCallbackEvent(key glfw.Key, scanCode int, action glfw.Action, modifiers glfw.ModifierKey) bool {
	s.lastInteraction = GetTime()
	return s.KeyboardEvent(s, key, scanCode, action, modifiers)
}

func (s *Screen) charCallbackEvent(codePoint rune) bool {
	s.lastInteraction = GetTime()
	return s.KeyboardCharacterEvent(s, codePoint)
}

func (s *Screen) preeditCallbackEvent(text []rune, blocks []int, focusedBlock int) {
	s.lastInteraction = GetTime()
	s.IMEPreeditEvent(s, text, blocks, focusedBlock)
}

func (s *Screen) imeStatusCallbackEvent() {
	s.lastInteraction = GetTime()
	s.IMEStatusEvent(s)
}

func (s *Screen) dropCallbackEvent(fileNames []string) bool {
	if s.dropEventCallback != nil {
		return s.dropEventCallback(fileNames)
	}
	return false
}

func (s *Screen) scrollCallbackEvent(x, y float32) bool {
	s.lastInteraction = GetTime()

	if len(s.focusPath) > 1 {
		window, ok := s.focusPath[len(s.focusPath)-2].(*Window)
		if ok && window.Modal() {
			if !window.Contains(s.mousePosX, s.mousePosY) {
				return false
			}
		}
	}
	return s.ScrollEvent(s, s.mousePosX, s.mousePosY, int(x), int(y))
}

func (s *Screen) resizeCallbackEvent(width, height int) bool {
	fbW, fbH := s.window.GetFramebufferSize()
	w, h := s.window.GetSize()

	if (fbW == 0 && fbH == 0) && (w == 0 && h == 0) {
		return false
	}
	s.fbW = fbW
	s.fbH = fbH
	s.w = w
	s.h = h
	s.lastInteraction = GetTime()
	if s.resizeEventCallback != nil {
		return s.resizeEventCallback(int(float32(fbW) / s.pixelRatio), int(float32(fbH) / s.pixelRatio))
	}
	return false
}

func (s *Screen) PerformLayout() {
	s.OnPerformLayout(s, s.context)
}

func (s *Screen) String() string {
	return fmt.Sprintf("Screen [%d,%d-%d,%d]", s.x, s.y, s.w, s.h)
}

func traverse(buffer *bytes.Buffer, w Widget, indent int) {
	for i := 0; i < indent; i++ {
		buffer.WriteString("  ")
	}
	buffer.WriteString(w.String())
	buffer.WriteByte('\n')
	for _, c := range w.Children() {
		traverse(buffer, c, indent+1)
	}
}

func (s *Screen) DebugPrint() {
	var buffer bytes.Buffer
	buffer.WriteString(s.String())
	buffer.WriteByte('\n')
	for _, c := range s.Children() {
		traverse(&buffer, c, 1)
	}
	println(buffer.String())
}
