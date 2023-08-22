package graphics

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const WITDH, HEIGTH = 800, 600

var (
	window   *sdl.Window
	renderer *sdl.Renderer
	font     *ttf.Font
)

func init() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}

	window, err = sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WITDH, HEIGTH, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(1)
	}

	err = ttf.Init()
	if err != nil {
		panic(err)
	}

	font, err = ttf.OpenFont("/home/jonathan/programacao/autosimulator/src/graphics/assets/IBMPlexMono-ExtraLight.ttf", 12)
	if err != nil {
		panic(1)
	}

}

func shutDown() {
	sdl.Quit()
	ttf.Quit()
	window.Destroy()
	renderer.Destroy()
	font.Close()
}

func Run() {

	////// States for testing
	state1 := newState(&sdl.Rect{X: 100, Y: 100, W: 50, H: 50}, "Q0", BLACK)
	state2 := newState(&sdl.Rect{X: WITDH / 2, Y: HEIGTH / 2, W: 50, H: 50}, "Q1", BLACK)
	states := []graphicalState{*state1, *state2}
	////////////////////

	////// Mouse events
	mousePos := sdl.Point{X: 0, Y: 0}
	clickOffset := sdl.Point{X: 0, Y: 0}
	var selectedState *graphicalState
	leftMouseButtonDown := false
	////////////////////

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false

			case *sdl.MouseMotionEvent:
				mousePos.X, mousePos.Y, _ = sdl.GetMouseState()

				if leftMouseButtonDown && selectedState != nil {
					fmt.Print("moving..")
					selectedState.X = mousePos.X - clickOffset.X
					selectedState.Y = mousePos.Y - clickOffset.Y
				}

			case *sdl.MouseButtonEvent:
				if event.(*sdl.MouseButtonEvent).Button == sdl.BUTTON_LEFT {
					if leftMouseButtonDown &&
						event.(*sdl.MouseButtonEvent).Type == sdl.MOUSEBUTTONUP {
						leftMouseButtonDown = false
						selectedState = nil
					}

					if !leftMouseButtonDown &&
						event.(*sdl.MouseButtonEvent).Type == sdl.MOUSEBUTTONDOWN {
						leftMouseButtonDown = true

						for _, state := range states {
							if mousePos.InRect(state.Rect) {
								selectedState = &state
								clickOffset.X = mousePos.X - state.X
								clickOffset.Y = mousePos.Y - state.Y
								break
							}
						}

					}

				}

			}

		}

		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.Clear()
		state1.draw(renderer, font)
		state2.draw(renderer, font)
		renderer.Present()

		sdl.Delay(1000 / 60)
	}

	defer shutDown()
}

/////////////////////////////////////////////

// func Run() {
// 	rect := &sdl.Rect{X: WITDH / 2, Y: HEIGTH / 2, W: 50, H: 50}
// 	graphicalState := NewNode(rect, "Q0", BLACK)

// 	renderer.SetDrawColor(255, 255, 255, 255)
// 	renderer.Clear()
// 	graphicalState.draw(renderer, font)
// 	renderer.Present()
// 	sdl.Delay(5000)

// }

// // func drawTransition() {
// // 	surfaceImg, _ := img.Load()
// // 	textureImg, _ := renderer.CreateTextureFromSurface(surfaceImg)

// // 	surfaceText, _ := font.RenderUTF8Solid("Q0", BLACK)
// // 	textureText, _ := renderer.CreateTextureFromSurface(surfaceText)

// // 	renderer.Copy(textureImg, nil, nil)
// // 	renderer.Copy(textureText, nil, nil)
// // }
