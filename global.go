package nanogui

import (
	"fmt"
	"github.com/goxjs/gl"
	"github.com/goxjs/glfw"
	"github.com/shibukawa/nanovgo"
	"io/ioutil"
	"path"
	"sync"
	"time"
)

var mainloopActive bool = false
var startTime time.Time
var debug bool

func Init() {
	err := glfw.Init(gl.ContextWatcher)
	if err != nil {
		panic(err)
	}
	startTime = time.Now()
}

func GetTime() float32 {
	return float32(time.Now().Sub(startTime)/time.Millisecond) * 0.001
}

func MainLoop() {
	mainloopActive = true

	var wg sync.WaitGroup

	/* If there are no mouse/keyboard events, try to refresh the
	view roughly every 50 ms; this is to support animations
	such as progress bars while keeping the system load
	reasonably low */
	wg.Add(1)
	go func() {
		for mainloopActive {
			time.Sleep(50 * time.Millisecond)
			glfw.PostEmptyEvent()
		}
		wg.Done()
	}()
	for mainloopActive {
		haveActiveScreen := false
		for _, screen := range nanoguiScreens {
			if !screen.Visible() {
				continue
			} else if screen.GLFWWindow().ShouldClose() {
				screen.SetVisible(false)
				continue
			}
			//screen.DebugPrint()
			screen.DrawAll()
			haveActiveScreen = true
		}
		if !haveActiveScreen {
			mainloopActive = false
			break
		}
		glfw.WaitEvents()
	}

	wg.Wait()
}

func SetDebug(d bool) {
	debug = d
}

func InitWidget(child, parent Widget) {
	child.SetVisible(true)
	child.SetEnabled(true)
	child.SetFontSize(-1)
	//w.cursor = Arrow
	if parent != nil {
		parent.AddChild(parent, child)
		child.SetTheme(parent.Theme())
	}
}

func LoadImageDirectory(ctx *nanovgo.Context, dir string) []Image {
	var images []Image
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(fmt.Sprintf("LoadImageDirectory: read error %v\n", err))
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
		images = append(images, Image{
			ImageID: img,
			Name:    fullPath[:len(fullPath)-4],
		})
	}
	return images
}
