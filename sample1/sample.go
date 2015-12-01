package main

import (
	"fmt"
	"github.com/shibukawa/nanogui.go"
	"github.com/shibukawa/nanovgo"
)

type Application struct {
	screen   *nanogui.Screen
	progress *nanogui.ProgressBar
	shader   *nanogui.GLShader
}

func buttonDemo(screen *nanogui.Screen) {
	window := nanogui.NewWindow(screen, "Button demo")
	window.SetPosition(15, 15)
	window.SetLayout(nanogui.NewGroupLayout())

	nanogui.NewLabel(window, "Push buttons").SetFont("sans-bold")

	b1 := nanogui.NewButton(window, "Plain button")
	b1.SetCallback(func() {
		fmt.Println("pushed!")
	})

	b2 := nanogui.NewButton(window, "Styled")
	b2.SetBackgroundColor(nanovgo.RGBA(0, 0, 255, 25))
	b2.SetIcon(nanogui.IconRocket)
	b2.SetCallback(func() {
		fmt.Println("pushed!")
	})

	nanogui.NewLabel(window, "Toggle button").SetFont("sans-bold")
	b3 := nanogui.NewButton(window, "Toggle me")
	b3.SetFlags(nanogui.ToggleButtonType)
	b3.SetChangeCallback(func(state bool) {
		fmt.Println("Toggle button state:", state)
	})

	nanogui.NewLabel(window, "Radio buttons").SetFont("sans-bold")
	b4 := nanogui.NewButton(window, "Radio button 1")
	b4.SetFlags(nanogui.RadioButtonType)
	b5 := nanogui.NewButton(window, "Radio button 2")
	b5.SetFlags(nanogui.RadioButtonType)

	nanogui.NewLabel(window, "A tool palette").SetFont("sans-bold")
	tools := nanogui.NewWidget(window)
	tools.SetLayout(nanogui.NewBoxLayout(nanogui.Horizontal, nanogui.Middle, 0, 6))

	nanogui.NewToolButton(tools, nanogui.IconCloud)
	nanogui.NewToolButton(tools, nanogui.IconFastForward)
	nanogui.NewToolButton(tools, nanogui.IconCompass)
	nanogui.NewToolButton(tools, nanogui.IconInstall)

	nanogui.NewLabel(window, "Popup buttons").SetFont("sans-bold")
	b6 := nanogui.NewPopupButton(window, "Popup")
	b6.SetIcon(nanogui.IconExport)
	popup := b6.Popup()
	popup.SetLayout(nanogui.NewGroupLayout())
	nanogui.NewCheckBox(popup, "Another check box")
}

func basicWidgetsDemo(screen *nanogui.Screen) {
	window := nanogui.NewWindow(screen, "Basic widgets")
	window.SetPosition(230, 15)
	window.SetLayout(nanogui.NewGroupLayout())

	nanogui.NewLabel(window, "Message dialog").SetFont("sans-bold")

	tools := nanogui.NewWidget(window)
	tools.SetLayout(nanogui.NewBoxLayout(nanogui.Horizontal, nanogui.Middle, 0, 6))

	b1 := nanogui.NewButton(tools, "Info")
	b1.SetCallback(func() {

	})
	b2 := nanogui.NewButton(tools, "Warn")
	b2.SetCallback(func() {

	})
	b3 := nanogui.NewButton(tools, "Ask")
	b3.SetCallback(func() {

	})

	nanogui.NewLabel(window, "Image panel & scroll panel").SetFont("sans-bold")
	imagePanelButton := nanogui.NewPopupButton(window, "Image Panel")
	imagePanelButton.SetIcon(nanogui.IconFolder)

	nanogui.NewLabel(window, "File dialog").SetFont("sans-bold")

	tools2 := nanogui.NewWidget(window)
	tools2.SetLayout(nanogui.NewBoxLayout(nanogui.Horizontal, nanogui.Middle, 0, 6))

	b4 := nanogui.NewButton(tools2, "Open")
	b4.SetCallback(func() {

	})
	b5 := nanogui.NewButton(tools2, "Save")
	b5.SetCallback(func() {

	})

	nanogui.NewLabel(window, "Combo box").SetFont("sans-bold")

	nanogui.NewLabel(window, "Check box").SetFont("sans-bold")
	cb1 := nanogui.NewCheckBox(window, "Flag 1")
	cb1.SetCallback(func(checked bool) {
		fmt.Println("Check box 1 state:", checked)
	})
	cb1.SetChecked(true)

	cb2 := nanogui.NewCheckBox(window, "Flag 2")
	cb2.SetCallback(func(checked bool) {
		fmt.Println("Check box 2 state:", checked)
	})
	nanogui.NewLabel(window, "Progress bar").SetFont("sans-bold")

	nanogui.NewLabel(window, "Slider and text box").SetFont("sans-bold")
}

func miscWidgetsDemo(screen *nanogui.Screen) {
	window := nanogui.NewWindow(screen, "Misc. widgets")
	window.SetPosition(455, 15)
	window.SetLayout(nanogui.NewGroupLayout())
}

func gridDemo(screen *nanogui.Screen) {
	window := nanogui.NewWindow(screen, "Grid of small widgets")
	window.SetPosition(455, 288)
	window.SetLayout(nanogui.NewGroupLayout())

	nanogui.NewLabel(window, "Floating point :").SetFont("sans-bold")
	nanogui.NewLabel(window, "Positive integer :").SetFont("sans-bold")
	nanogui.NewLabel(window, "Checkbox :").SetFont("sans-bold")
	nanogui.NewLabel(window, "Bombobox :").SetFont("sans-bold")
	nanogui.NewLabel(window, "Color button :").SetFont("sans-bold")

	popupButton := nanogui.NewPopupButton(window, "")
	popupButton.SetBackgroundColor(nanovgo.RGBA(255, 120, 0, 255))
	popupButton.SetFontSize(16)
	popupButton.SetFixedSize(100, 20)
	popup := popupButton.Popup()
	popup.SetLayout(nanogui.NewGroupLayout())

	colorButton := nanogui.NewButton(popup, "Pick")
	colorButton.SetFixedSize(100, 25)

	colorButton.SetChangeCallback(func(pushed bool) {
		if pushed {
			popupButton.SetPushed(false)
		}
	})
}

func selectedImageDemo(screen *nanogui.Screen) {
	window := nanogui.NewWindow(screen, "Basic widgets")
	window.SetPosition(705, 15)
	window.SetLayout(nanogui.NewGroupLayout())
	cb := nanogui.NewCheckBox(window, "Expand")
	cb.SetCallback(func(checked bool) {

	})
}

func (a *Application) init() {
	a.screen = nanogui.NewScreen(1024, 768, "NanoGUI.Go Test", true, false)

	buttonDemo(a.screen)
	basicWidgetsDemo(a.screen)
	selectedImageDemo(a.screen)
	miscWidgetsDemo(a.screen)
	gridDemo(a.screen)

	a.screen.DebugPrint()

	a.screen.PerformLayout()

	/* All NanoGUI widgets are initialized at this point. Now
	create an OpenGL shader to draw the main window contents.

	NanoGUI comes with a simple Eigen-based wrapper around OpenGL 3,
	which eliminates most of the tedious and error-prone shader and
	buffer object management.
	*/
}

func main() {
	nanogui.Init()
	//nanogui.SetDebug(true)
	app := Application{}
	app.init()
	app.screen.DrawAll()
	app.screen.SetVisible(true)
	nanogui.MainLoop()
}
