package main

import (
	"fmt"
	"github.com/shibukawa/nanogui.go"
)

type Application struct {
	screen   *nanogui.Screen
	progress *nanogui.ProgressBar
	shader   *nanogui.GLShader
}

func (a *Application) init() {
	a.screen = nanogui.NewScreen(1024, 768, "NanoGUI.Go Test", true, false)

	window := nanogui.NewWindow(a.screen, "Button demo")
	window.SetPosition(15, 15)
	window.SetLayout(nanogui.NewGroupLayout(-1, -1, -1, -1))

	nanogui.NewLabel(window, "Push buttons").SetFont("sans-bold")

	b1 := nanogui.NewButton(window, "Plain button")
	b1.SetCallback(func() {
		fmt.Println("pushed!")
	})

	b2 := nanogui.NewButton(window, "Styled")
	b2.SetIcon(nanogui.IconFlightTakeoff)
	b2.SetCallback(func() {
		fmt.Println("pushed!")
	})

	nanogui.NewLabel(window, "Toggle button").SetFont("sans-bold")
	b3 := nanogui.NewButton(window, "Toggle me")
	b3.SetFlags(nanogui.ToggleButton)
	b3.SetChangeCallback(func(state bool) {
		fmt.Println("Toggle button state:", state)
	})

	nanogui.NewLabel(window, "Radio buttons").SetFont("sans-bold")
	b4 := nanogui.NewButton(window, "Radio button 1")
	b4.SetFlags(nanogui.RadioButton)
	b5 := nanogui.NewButton(window, "Radio button 2")
	b5.SetFlags(nanogui.RadioButton)

	nanogui.NewLabel(window, "A tool palette").SetFont("sans-bold")
	tools := nanogui.NewWidget(window)
	tools.SetLayout(nanogui.NewBoxLayout(nanogui.Horizontal, nanogui.Middle, 0, 6))

	nanogui.NewToolButton(tools, nanogui.IconCloud)
	nanogui.NewToolButton(tools, nanogui.IconFavorite)
	nanogui.NewToolButton(tools, nanogui.IconComputer)
	nanogui.NewToolButton(tools, nanogui.IconFontDownload)

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
