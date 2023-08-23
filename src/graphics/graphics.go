package graphics

import (
	"autosimulator/src/reader"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	WITDH, HEIGTH = 800, 600
)

var (
	window   *sdl.Window
	renderer *sdl.Renderer
	font     *ttf.Font
	STATES   = make(map[string]*graphicalState)
)

func Run() {

	////// States for testing
	// state1 := NewState(&sdl.Rect{X: 100, Y: 100, W: 50, H: 50}, "Q0", BLACK)
	// state2 := NewState(&sdl.Rect{X: WITDH / 2, Y: 100, W: 50, H: 50}, "Q1", BLACK)
	// state1.AddNextState(state2)
	machine := reader.ReadMachine("/home/jonathan/programacao/autosimulator/src/machine/afdMachine/machine_example.json")
	STATES = machineToStates(machine)
	////////////////////

	////// Mouse events
	mousePos := sdl.Point{X: 0, Y: 0}
	clickOffset := sdl.Point{X: 0, Y: 0}
	var selectedState *graphicalState
	leftMouseButtonDown := false
	//////////////////

	// Mainloop
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

						for _, state := range STATES {
							if mousePos.InRect(state.Rect) {
								selectedState = state
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
		//
		// drawExs()
		//

		for _, state := range STATES {
			state.Draw(renderer, font)
		}
		renderer.Present()

		sdl.Delay(1000 / 60)
	}

	defer shutDown()
}

// Inits / shutdows
func init() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}

	// TODO: title of the window
	window, err = sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WITDH, HEIGTH, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	err = ttf.Init()
	if err != nil {
		panic(err)
	}

	font, err = ttf.OpenFont("/home/jonathan/programacao/autosimulator/src/graphics/assets/IBMPlexMono-ExtraLight.ttf", 12)
	if err != nil {
		panic(err)
	}

}

func shutDown() {
	sdl.Quit()
	ttf.Quit()
	window.Destroy()
	renderer.Destroy()
	font.Close()
}

/////////////////////////////////////////////

func drawExs() {
	gfx.ThickLineColor(renderer, 0, HEIGTH/2, WITDH, HEIGTH/2, 1, BLACK)
	gfx.ThickLineColor(renderer, WITDH/2, 0, WITDH/2, HEIGTH, 1, BLACK)
}
