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

type (
	SDLWindow struct {
		window    *sdl.Window
		renderer  *sdl.Renderer
		font      *ttf.Font
		terminate bool
		fps       int
	}

	drag struct {
		clickOffset   *sdl.Point
		selected      *objectInfo
		leftMouseDown bool
	}

	objectInfo struct {
		object any
		*sdl.Point
	}

	machineChannel struct {
		channel   chan int
		lastMsg   int
		executing bool
	}
)

var (
	states        map[string]*graphicalState
	currentState  string
	LastMessage   = machine.STATE_NOT_CHANGE
	mousePos      *sdl.Point
	dragInfo      *drag
	communication *machineChannel
)

func NewSDLWindow() *SDLWindow {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}

	// TODO: title of the window
	window, err := sdl.CreateWindow("Simulador de Autômato", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WITDH, HEIGTH, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	err = ttf.Init()
	if err != nil {
		panic(err)
	}

	font, err := ttf.OpenFont("/home/jonathan/programacao/autosimulator/src/graphics/assets/IBMPlexMono-ExtraLight.ttf", 12)
	if err != nil {
		panic(err)
	}

	runtime.LockOSThread()

	return &SDLWindow{window: window, renderer: renderer, font: font}
}

func init() {
	dragInfo = &drag{clickOffset: &sdl.Point{X: 0, Y: 0}, selected: nil, leftMouseDown: false}
	communication = &machineChannel{channel: make(chan int), lastMsg: machine.STATE_NOT_CHANGE, executing: false}
}

func (w *SDLWindow) Destroy() {
	sdl.Quit()
	ttf.Quit()
	w.window.Destroy()
	w.renderer.Destroy()
	w.font.Close()
}

func (w *SDLWindow) pollEvent() {
	event := sdl.PollEvent()
	for event != nil {
		switch event.(type) {
		case *sdl.QuitEvent:
			w.Destroy()
			w.terminate = true

		case *sdl.KeyboardEvent:
			handleKeyboardEvents(event.(*sdl.KeyboardEvent), w)

		case *sdl.MouseButtonEvent:
			handleMouseButtonEvents(event.(*sdl.MouseButtonEvent))

		case *sdl.MouseMotionEvent:
			handleMouseMotionEvent()

		default:
		}

		event = sdl.PollEvent()
	}
}

func handleKeyboardEvents(event *sdl.KeyboardEvent, w *SDLWindow) {
	if event.Keysym.Sym == sdl.K_RETURN &&
		!communication.executing {
		fmt.Println("executing")
		w.fps = 1
		communication.executing = true
		// Machine?
		go machine.Execute(m, []string{"a", "a", "b"}, communication.channel)
	}

}

func handleMouseButtonEvents(event *sdl.MouseButtonEvent) {
	if event.Button != sdl.BUTTON_LEFT {
		// Nothing to do with other buttons
		return
	}

	// IF THE LEFT MOUSE BUTTON IS PRESSED
	switch event.Type {
	case sdl.MOUSEBUTTONDOWN:
		if !dragInfo.leftMouseDown {
			dragInfo.leftMouseDown = true
			for _, state := range states {
				if mousePos.InRect(state.Rect) {
					object := &objectInfo{object: state, Point: &sdl.Point{X: state.X, Y: state.Y}}
					dragInfo.selected = object
					dragInfo.clickOffset.X = mousePos.X - state.X
					dragInfo.clickOffset.Y = mousePos.Y - state.Y
					break
				}
			}
		}

	case sdl.MOUSEBUTTONUP:
		if dragInfo.leftMouseDown {
			dragInfo.leftMouseDown = false
			dragInfo.selected = nil
		}
	}
}

func handleMouseMotionEvent() {
	mousePos.X, mousePos.Y, _ = sdl.GetMouseState()
	if dragInfo.leftMouseDown && dragInfo.selected != nil {
		dragObject()
	}
}

func dragObject() {
	dragInfo.selected.X = mousePos.X - dragInfo.clickOffset.X
	dragInfo.selected.Y = mousePos.Y - dragInfo.clickOffset.Y
}

// ///////////////////////////////////////////
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
func drawExs() {
	gfx.ThickLineColor(renderer, 0, HEIGTH/2, WITDH, HEIGTH/2, 1, BLACK)
	gfx.ThickLineColor(renderer, WITDH/2, 0, WITDH/2, HEIGTH, 1, BLACK)
}
