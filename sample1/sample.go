// +build !js

package main

import (
	"fmt"
	"github.com/goxjs/glfw"
	"github.com/shibukawa/nanogui.go"
	"github.com/shibukawa/nanogui.go/sample1/demo"
	"github.com/shibukawa/nanovgo"
	"io/ioutil"
	"math"
	"path"
)

type Application struct {
	screen   *nanogui.Screen
	progress *nanogui.ProgressBar
	shader   *nanogui.GLShader
}

func (a *Application) init() {
	glfw.WindowHint(glfw.Samples, 4)
	a.screen = nanogui.NewScreen(1024, 768, "NanoGUI.Go Test", true, false)

	a.screen.NVGContext().CreateFont("japanese", "font/GenShinGothic-P-Regular.ttf")

	demo.ButtonDemo(a.screen)
	images := loadImageDirectory(a.screen.NVGContext(), "icons")
	imageButton, imagePanel, progressBar := demo.BasicWidgetsDemo(a.screen, images)
	a.progress = progressBar
	demo.SelectedImageDemo(a.screen, imageButton, imagePanel)
	demo.MiscWidgetsDemo(a.screen)
	demo.GridDemo(a.screen)

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

func loadImageDirectory(ctx *nanovgo.Context, dir string) []nanogui.Image {
	var images []nanogui.Image
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(fmt.Sprintf("loadImageDirectory: read error %v\n", err))
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		ext := path.Ext(file.Name())
		if ext != ".png" {
			continue
		}
		fullPath := path.Join(dir, file.Name())
		img := ctx.CreateImage(fullPath, 0)
		if img == 0 {
			panic("Could not open image data!")
		}
		images = append(images, nanogui.Image{
			ImageID: img,
			Name:    fullPath[:len(fullPath)-4],
		})
	}
	return images
}
