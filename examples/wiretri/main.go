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
	pixelscale = 1

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
	// Draw a triangle, concurrently
	multirender.WireTriangle(cores, pixels, rw(), rh(), rw(), rh(), rw(), rh(), color.RGBA{rb(), rb(), rb(), opaque}, pitch)
}

func run() int {

	sdl.Init(sdl.INIT_VIDEO)

	var (
		window   *sdl.Window
		renderer *sdl.Renderer
		err      error
	)

	window, err = sdl.CreateWindow("Triangles!", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(width*pixelscale), int32(height*pixelscale), sdl.WINDOW_SHOWN)
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

	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, width, height)
	if err != nil {
		panic(err)
	}

	//texture.SetBlendMode(sdl.BLENDMODE_BLEND) // sdl.BLENDMODE_ADD is also possible

	rand.Seed(time.Now().UnixNano())

	var (
		pixels = make([]uint32, width*height)
		cores  = runtime.NumCPU()
		event  sdl.Event
		quit   bool
		pause  bool
	)

	// Fill the pixel buffer with opaque color #808080
	for i := range pixels {
		pixels[i] = 0xff808080
	}

	// Innerloop
	for !quit {

		if !pause {

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
					case sdl.K_SPACE, sdl.K_p:
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
						// save the image
						multirender.SavePixelsToPNG(pixels, pitch, "screenshot.png", true)
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