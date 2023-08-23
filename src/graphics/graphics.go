package graphics

import (
	"autosimulator/src/machine"
	"fmt"
	"runtime"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	WITDH, HEIGTH = 800, 600
)

// Enumeration

var (
	window       *sdl.Window
	renderer     *sdl.Renderer
	font         *ttf.Font
	states       map[string]*graphicalState
	currentState string
	LastMessage  = machine.STATE_NOT_CHANGE
)

type SDLWindow struct {
	window   *sdl.Window
	renderer *sdl.Renderer
	font     *ttf.Font
}

func NewSDLWindow() *SDLWindow {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}

	// TODO: title of the window
	window, err = sdl.CreateWindow("Autômato Simulator", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
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

	runtime.LockOSThread()

	return &SDLWindow{window: window, renderer: renderer, font: font}
}

func (w *SDLWindow) pollEvent() {

}

func Run(m machine.Machine) {
	channel := make(chan int)
	var fps uint32 = 1000 / 60
	states = machineStates(m)

	////// Mouse events
	mousePos := sdl.Point{X: 0, Y: 0}
	clickOffset := sdl.Point{X: 0, Y: 0}
	var selectedState *graphicalState
	leftMouseButtonDown := false
	//////////////////

	// Execute
	executing := false
	//////////////////

	// Mainloop
	running := true
	for running {

		if executing {
			LastMessage = <-channel
		}

		if LastMessage == machine.STATE_CHANGE {
			if currentState != "" {
				states[currentState].isCurrent = false
			}
			currentState = m.CurrentState()
			states[currentState].isCurrent = true
			LastMessage = machine.STATE_NOT_CHANGE
		}

		if LastMessage == machine.STATE_INPUT_ACCEPTED ||
			LastMessage == machine.STATE_INPUT_REJECTED {
			LastMessage = machine.STATE_CHANGE
			fmt.Printf("Finished")
			executing = false
			fps = 1000 / 60
		}

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false

			case *sdl.KeyboardEvent:
				if event.(*sdl.KeyboardEvent).Keysym.Sym == sdl.K_RETURN &&
					!executing {
					fmt.Println("executing")
					fps = 1000 / 1
					executing = true
					go machine.Execute(m, []string{"a", "a", "b"}, channel)
				}

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

						for _, state := range states {
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
		for _, state := range states {
			state.Draw(renderer, font, states)
		}
		renderer.Present()

		sdl.Delay(fps)
	}

	defer shutDown()
}

// Inits / shutdows

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
