package main

import (
	"fmt"
	"github.com/goxjs/glfw"
	"github.com/shibukawa/nanogui.go"
	"github.com/shibukawa/nanovgo"
	"math"
	"strconv"
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

	nanogui.NewLabel(popup, "Arbitrary widgets can be placed here").SetFont("sans-bold")
	nanogui.NewCheckBox(popup, "A check box")
	b7 := nanogui.NewPopupButton(popup, "Recursive popup")
	b7.SetIcon(nanogui.IconFlash)
	popup2 := b7.Popup()

	popup2.SetLayout(nanogui.NewGroupLayout())
	nanogui.NewCheckBox(popup2, "Another check box")
}

func basicWidgetsDemo(screen *nanogui.Screen) (*nanogui.PopupButton, *nanogui.ImagePanel, *nanogui.ProgressBar) {
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
	popup := imagePanelButton.Popup()
	vscroll := nanogui.NewVScrollPanel(popup)
	imgPanel := nanogui.NewImagePanel(vscroll)
	imgPanel.SetImages(nanogui.LoadImageDirectory(screen.NVGContext(), "icons"))
	popup.SetFixedSize(245, 150)

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
	nanogui.NewComboBox(window, []string{"Combo box item 1", "Combo box item 2", "Combo box item 3"})

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
	progress := nanogui.NewProgressBar(window)

	nanogui.NewLabel(window, "Slider and text box").SetFont("sans-bold")
	panel := nanogui.NewWidget(window)
	panel.SetLayout(nanogui.NewBoxLayout(nanogui.Horizontal, nanogui.Middle, 0, 20))
	slider := nanogui.NewSlider(panel)
	slider.SetValue(0.5)
	slider.SetFixedWidth(80)

	textBox := nanogui.NewTextBox(panel)
	textBox.SetFixedSize(60, 25)
	textBox.SetFontSize(20)
	textBox.SetAlignment(nanogui.TextRight)
	textBox.SetValue("50")
	textBox.SetUnits("%")

	slider.SetCallback(func(value float32) {
		textBox.SetValue(strconv.FormatInt(int64(value * 100), 10))
	})
	slider.SetFinalCallback(func(value float32) {
		fmt.Printf("Final slider value: %d\n", int(value*100))
	})

	return imagePanelButton, imgPanel, progress
}

func miscWidgetsDemo(screen *nanogui.Screen) {
	window := nanogui.NewWindow(screen, "Misc. widgets")
	window.SetPosition(445, 15)
	window.SetLayout(nanogui.NewGroupLayout())

	nanogui.NewLabel(window, "Color wheel").SetFont("sans-bold")
	nanogui.NewColorWheel(window)

	nanogui.NewLabel(window, "Color picker").SetFont("sans-bold")
	nanogui.NewColorPicker(window)

	nanogui.NewLabel(window, "Function graph").SetFont("sans-bold")
	graph := nanogui.NewGraph(window, "Some function")
	graph.SetHeader("E = 2.35e-3")
	graph.SetFooter("Iteration 89")
	fValues := make([]float32, 100)
	for i := 0; i < 100; i++ {
		x := float64(i)
		fValues[i] = 0.5 * float32(0.5*math.Sin(x/10.0)+0.5*math.Cos(x/23.0)+1.0)
	}
	graph.SetValues(fValues)
}

func gridDemo(screen *nanogui.Screen) {
	window := nanogui.NewWindow(screen, "Grid of small widgets")
	window.SetPosition(445, 358)
	layout := nanogui.NewGridLayout(nanogui.Horizontal, 2, nanogui.Middle, 15, 5)
	layout.SetColAlignment(nanogui.Maximum, nanogui.Fill)
	layout.SetColSpacing(10)
	window.SetLayout(layout)

	{
		nanogui.NewLabel(window, "Floating point :").SetFont("sans-bold")
		textBox := nanogui.NewTextBox(window, "50.0")
		textBox.SetEditable(true)
		textBox.SetFixedSize(100, 20)
		textBox.SetUnits("GiB")
		textBox.SetDefaultValue("0.0")
		textBox.SetFontSize(16)
		textBox.SetFormat(`^[-]?[0-9]*\.?[0-9]+$`)
	}

	{
		nanogui.NewLabel(window, "Positive integer :").SetFont("sans-bold")
		textBox := nanogui.NewTextBox(window, "50")
		textBox.SetEditable(true)
		textBox.SetFixedSize(100, 20)
		textBox.SetUnits("MHz")
		textBox.SetDefaultValue("0.0")
		textBox.SetFontSize(16)
		textBox.SetFormat(`^[1-9][0-9]*$`)
	}
	{
		nanogui.NewLabel(window, "Float box :").SetFont("sans-bold")
		floatBox := nanogui.NewFloatBox(window, 10.0)
		floatBox.SetEditable(true)
		floatBox.SetFixedSize(100, 20)
		floatBox.SetUnits("GiB")
		floatBox.SetDefaultValue(0.0)
		floatBox.SetFontSize(16)
	}

	{
		nanogui.NewLabel(window, "Int box :").SetFont("sans-bold")
		intBox := nanogui.NewIntBox(window, true, 50)
		intBox.SetEditable(true)
		intBox.SetFixedSize(100, 20)
		intBox.SetUnits("MHz")
		intBox.SetDefaultValue(0)
		intBox.SetFontSize(16)
	}
	{
		nanogui.NewLabel(window, "Checkbox :").SetFont("sans-bold")
		checkbox := nanogui.NewCheckBox(window, "Check me")
		checkbox.SetFontSize(16)
		checkbox.SetChecked(true)
	}
	{
		nanogui.NewLabel(window, "Combobox :").SetFont("sans-bold")
		combobox := nanogui.NewComboBox(window, []string{"Item 1", "Item 2", "Item 3"})
		combobox.SetFontSize(16)
		combobox.SetFixedSize(100, 20)
	}
	{
		nanogui.NewLabel(window, "Color button :").SetFont("sans-bold")

		popupButton := nanogui.NewPopupButton(window, "")
		popupButton.SetBackgroundColor(nanovgo.RGBA(255, 120, 0, 255))
		popupButton.SetFontSize(16)
		popupButton.SetFixedSize(100, 20)
		popup := popupButton.Popup()
		popup.SetLayout(nanogui.NewGroupLayout())

		colorWheel := nanogui.NewColorWheel(popup)
		colorWheel.SetColor(popupButton.BackgroundColor())

		colorButton := nanogui.NewButton(popup, "Pick")
		colorButton.SetFixedSize(100, 25)
		colorButton.SetBackgroundColor(colorWheel.Color())

		colorWheel.SetCallback(func(color nanovgo.Color) {
			colorButton.SetBackgroundColor(color)
		})

		colorButton.SetChangeCallback(func(pushed bool) {
			if pushed {
				popupButton.SetBackgroundColor(colorButton.BackgroundColor())
				popupButton.SetPushed(false)
			}
		})
	}
}

func selectedImageDemo(screen *nanogui.Screen, imageButton *nanogui.PopupButton, imagePanel *nanogui.ImagePanel) {
	window := nanogui.NewWindow(screen, "Selected image")
	window.SetPosition(685, 15)
	window.SetLayout(nanogui.NewGroupLayout())

	img := nanogui.NewImageView(window)
	img.SetPolicy(nanogui.ImageSizePolicyExpand)
	img.SetFixedSize(300, 300)
	img.SetImage(imagePanel.Images()[0].ImageID)

	imagePanel.SetCallback(func(index int) {
		img.SetImage(imagePanel.Images()[index].ImageID)
	})

	cb := nanogui.NewCheckBox(window, "Expand")
	cb.SetCallback(func(checked bool) {
		if checked {
			img.SetPolicy(nanogui.ImageSizePolicyExpand)
		} else {
			img.SetPolicy(nanogui.ImageSizePolicyFixed)
		}
	})
	cb.SetChecked(true)
}

func (a *Application) init() {
	glfw.WindowHint(glfw.Samples, 4)
	a.screen = nanogui.NewScreen(1024, 768, "NanoGUI.Go Test", true, false)

	buttonDemo(a.screen)
	imageButton, imagePanel, progressBar := basicWidgetsDemo(a.screen)
	a.progress = progressBar
	selectedImageDemo(a.screen, imageButton, imagePanel)
	miscWidgetsDemo(a.screen)
	gridDemo(a.screen)

	a.screen.SetDrawContentsCallback(func() {
		a.progress.SetValue(float32(math.Mod(float64(nanogui.GetTime())/10, 1.0)))
	})

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