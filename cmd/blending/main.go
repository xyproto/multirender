package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/xyproto/multirender"
	"github.com/xyproto/sdl2utils"
)

const (
	// Size of "worldspace pixels", measured in "screenspace pixels"
	pixelscale = 4

	// The resolution (worldspace)
	width  = 320
	height = 200

	// The width of the pixel buffer, used when calculating where to place pixels (y*pitch+x)
	pitch = width

	// Target framerate
	frameRate = 60

	// Alpha value for opaque colors
	opaque = 255
)

// rb returns a random byte
func rb() uint8 {
	return uint8(rand.Intn(255))
}

// rw returns a random int32 in the range [0,width)
func rw() int32 {
	return rand.Int31n(width)
}

// rh returns a random int32 in the range [0,height)
func rh() int32 {
	return rand.Int31n(height)
}

// DrawAll fills the pixel buffer with pixels.
// "cores" is how many CPU cores should be targeted when drawing triangles,
// by launching the same number of goroutines.
func DrawAll(pixels []uint32, cores int) {

	var opaqueness uint8 = 1 // Almost completely transparent

	// Draw a triangle, concurrently
	multirender.Triangle(cores, pixels, rw(), rh(), rw(), rh(), rw(), rh(), color.RGBA{rb(), rb(), rb(), opaqueness}, pitch)

	// Draw a line and a red pixel, without caring about which order they appear in, or if they will complete before the next frame is drawn
	go multirender.Line(pixels, rw(), rh(), rw(), rh(), color.RGBA{rb(), rb(), rb(), opaqueness}, pitch)
	go multirender.Pixel(pixels, rw(), rh(), color.RGBA{0xff, 0, 0, opaqueness}, pitch)
}

func run() int {

	sdl.Init(sdl.INIT_VIDEO)
	defer sdl.Quit()

	var (
		window   *sdl.Window
		renderer *sdl.Renderer
		err      error
	)

	window, err = sdl.CreateWindow("Pixels!", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(width*pixelscale), int32(height*pixelscale), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return 1
	}
	defer window.Destroy()

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return 1
	}
	defer renderer.Destroy()

	// Fill the render buffer with black
	renderer.SetDrawColor(0, 0, 0, opaque)
	renderer.Clear()

	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, width, height)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create texture: %s\n", err)
		return 1
	}

	texture.SetBlendMode(sdl.BLENDMODE_ADD)

	rand.Seed(time.Now().UnixNano())

	var (
		pixels = make([]uint32, width*height)
		cores  = runtime.NumCPU()
		event  sdl.Event
		quit   bool
		pause  bool
	)

	// Innerloop
	for !quit {

		if !pause {
			// Draw to pixel buffer
			DrawAll(pixels, cores)

			// Draw pixel buffer to screen
			texture.UpdateRGBA(nil, pixels, width)

			renderer.Copy(texture, nil, nil)
			renderer.Present()
		}

		// Check for events
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				quit = true
			case *sdl.KeyboardEvent:
				ke := event.(*sdl.KeyboardEvent)
				if ke.Type == sdl.KEYDOWN {
					ks := ke.Keysym
					switch ks.Sym {
					case sdl.K_ESCAPE:
						quit = true
					case sdl.K_q:
						quit = true
					case sdl.K_RETURN:
						altHeldDown := ks.Mod == sdl.KMOD_LALT || ks.Mod == sdl.KMOD_RALT
						if !altHeldDown {
							// alt+enter is not pressed
							break
						}
						// alt+enter is pressed
						fallthrough
					case sdl.K_f, sdl.K_F11:
						sdl2utils.ToggleFullscreen(window)
						// Clear pixels at fullscreen toggle, because of blend mode and no double buffering
						renderer.Clear()
						for i := range pixels {
							pixels[i] = 0
						}
					case sdl.K_p, sdl.K_SPACE:
						// pause toggle
						pause = !pause
					case sdl.K_s:
						ctrlHeldDown := ks.Mod == sdl.KMOD_LCTRL || ks.Mod == sdl.KMOD_RCTRL
						if !ctrlHeldDown {
							// ctrl+s is not pressed
							break
						}
						// ctrl+s is pressed
						fallthrough
					case sdl.K_F12:
						// screenshot
						sdl2utils.Screenshot(renderer, "screenshot.png", true)
					}
				}
			}
		}
		sdl.Delay(1000 / frameRate)
	}
	return 0
}

func main() {
	// This is to allow the deferred functions in run() to kick in at exit
	os.Exit(run())
}
