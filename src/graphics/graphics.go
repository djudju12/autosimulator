package graphics

import (
	"autosimulator/src/collections"
	"autosimulator/src/machine"
	"fmt"
	"os"
	"runtime"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	WITDH, HEIGTH = 800, 600
	TITLE         = "Simulador de Autômato"
	FONT_PATH     = "/home/jonathan/programacao/autosimulator/src/graphics/assets/IBMPlexMono-ExtraLight.ttf"
	FONT_ZIE      = 24
	FPS_DEFAULT   = 60
	FPS_EXECUTING = 1
)

type (
	environment struct {
		w        *_SDLWindow
		dragInfo *drag
		mousePos *sdl.Point
		machine  machine.Machine
		input    *collections.Fita
		states   map[string]*graphicalState
	}

	_SDLWindow struct {
		window     *sdl.Window
		renderer   *sdl.Renderer
		font       *ttf.Font
		cacheWords map[string]*sdl.Surface
		terminate  bool
		redraw     bool
	}

	drag struct {
		clickOffset   *sdl.Point
		selected      *graphicalState
		leftMouseDown bool
	}
)

func Mainloop(env *environment) {
	window := env.w
	runtime.LockOSThread() // sdl2 precisa rodar na main thread.

	for !window.terminate {
		pollEvent(env)
		draw(env)
		sdl.Delay(1000 / FPS_DEFAULT)
	}

	defer env.Destroy()
}

func PopulateEnvironment(window *_SDLWindow, activeMachine machine.Machine) *environment {
	dragInfo := &drag{
		clickOffset:   &sdl.Point{X: 0, Y: 0},
		selected:      nil,
		leftMouseDown: false,
	}

	env := &environment{
		w:        window,
		dragInfo: dragInfo,
		machine:  activeMachine,
		states:   machineStates(activeMachine),
	}

	return env
}

func NewSDLWindow() *_SDLWindow {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow(TITLE, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
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

	font, err := ttf.OpenFont(FONT_PATH, FONT_ZIE)
	if err != nil {
		panic(err)
	}

	cacheWords := make(map[string]*sdl.Surface)

	return &_SDLWindow{window: window,
		renderer:   renderer,
		font:       font,
		terminate:  false,
		cacheWords: cacheWords,
		redraw:     true,
	}
}

func (env *environment) Input(fita *collections.Fita) {
	env.input = fita
}

func (env *environment) Destroy() {
	sdl.Quit()
	ttf.Quit()
	env.w.window.Destroy()
	env.w.renderer.Destroy()
	env.w.font.Close()

	for _, v := range env.w.cacheWords {
		v.Free()
	}
}

func pollEvent(env *environment) {
	window := env.w
	event := sdl.PollEvent()
	for event != nil {
		switch event := event.(type) {
		case *sdl.QuitEvent:
			fmt.Printf("Quiting....")
			window.terminate = true

		case *sdl.KeyboardEvent:
			handleKeyboardEvents(event, env)

		case *sdl.MouseButtonEvent:
			handleMouseButtonEvents(event, env)

		case *sdl.MouseMotionEvent:
			handleMouseMotionEvent(env)

		default:
		}

		event = sdl.PollEvent()
	}
}

func handleKeyboardEvents(event *sdl.KeyboardEvent, env *environment) {
	switch event.Keysym.Sym {
	case sdl.K_RETURN:
	case sdl.K_r:
		// env.Reset()
	}
}

func handleMouseButtonEvents(event *sdl.MouseButtonEvent, env *environment) {
	dragInfo := env.dragInfo
	mousePos := env.mousePos
	states := env.states

	if event.Button != sdl.BUTTON_LEFT {
		// Nothing to do with other buttons
		return
	}

	// IF THE LEFT MOUSE BUTTON IS PRESSED \/
	switch event.Type {
	case sdl.MOUSEBUTTONDOWN:
		if !dragInfo.leftMouseDown {
			dragInfo.leftMouseDown = true

			for _, state := range states {
				if mousePos.InRect(state.Rect) {
					dragInfo.selected = state
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

func handleMouseMotionEvent(env *environment) {
	x, y, _ := sdl.GetMouseState()
	env.mousePos = &sdl.Point{X: x, Y: y}
	dragInfo := env.dragInfo
	if dragInfo.leftMouseDown && dragInfo.selected != nil {
		dragInfo.selected.X = env.mousePos.X - dragInfo.clickOffset.X
		dragInfo.selected.Y = env.mousePos.Y - dragInfo.clickOffset.Y
		env.w.redraw = true
	}
}

func draw(env *environment) {
	window := env.w
	if env.w.redraw {
		err := window.cleanUp()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = drawUi(env)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = drawNodes(env)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		env.w.redraw = false
	}

	window.renderer.Present()
}

func (w *_SDLWindow) cleanUp() error {

	err := w.renderer.SetDrawColor(255, 255, 255, 255)
	if err != nil {
		return err
	}

	err = w.renderer.Clear()
	if err != nil {
		return err
	}

	return nil
}

func drawUi(env *environment) error {
	var padx, pady int32 = 5, 5
	err := env.drawFita(padx, pady)
	if err != nil {
		return err
	}

	machineType := env.machine.Type()
	amount := 0
	switch machineType {
	case machine.ONE_STACK_MACHINE:
		amount = 1
	case machine.TWO_STACK_MACHINE:
		amount = 2
	default:
	}

	err = env.drawStacks(amount, padx, pady)
	if err != nil {
		return err
	}

	return nil
}

func drawNodes(env *environment) error {
	var err error

	for _, state := range env.states {
		err = state.Draw(env.w, env.states)
		if err != nil {
			return err
		}
	}

	return nil
}

// func (env *environment) changeState(spriteName string) {
// 	currentState := radio.activeMachine.CurrentState()
// 	env.states[currentState].spriteName = spriteName

// 	// para mudar a cor para o normal na proxima iteraçao
// 	radio.lastState = currentState
// }

// // utilidade
// func (env *environment) Reset() {
// 	fmt.Println("Reseting...")
// 	window := env.w
// 	radio := env.radio
// 	radio.activeMachine.Init()
// 	window.redraw = true
// 	radio.inExecution = false
// 	lastState := env.states[radio.lastState]
// 	initalState := env.states[radio.activeMachine.CurrentState()]
// 	radio.input.Reset()
// 	radio.inputToPrint = radio.input.Peek(TAMANHO_ESTRUTURAS)

// 	if lastState != nil {
// 		lastState.spriteName = BLACK_RING
// 	}

// 	if initalState != nil {
// 		initalState.spriteName = BLUE_RING
// 	}
// }

// func (env *environment) RunMachine() {
// 	fmt.Println("Running...")
// 	radio := env.radio
// 	radio.lastMsg = machine.STATE_CHANGE
// 	go machine.Execute(radio.activeMachine, radio.input, radio.channel)
// }

func (w *_SDLWindow) textSurface(text string, color sdl.Color) (*sdl.Surface, error) {
	font := w.font
	words := w.cacheWords

	// Já tem no cache
	surface := words[text]
	if surface != nil {
		return surface, nil
	}

	// Cria o novo testo e coloca no cache
	surface, err := font.RenderUTF8Solid(text, color)
	if err != nil {
		return nil, err
	}

	words[text] = surface
	return surface, nil
}

func Center(rec *sdl.Rect) sdl.Point {
	return sdl.Point{
		X: rec.X + rec.W/2,
		Y: rec.Y + rec.H/2,
	}
}

func drawAxs(env *environment) {
	gfx.ThickLineColor(env.w.renderer, WITDH/2, 0, WITDH/2, HEIGTH, 2, BLACK)
	gfx.ThickLineColor(env.w.renderer, 0, HEIGTH/2, WITDH, HEIGTH/2, 2, BLACK)
}
