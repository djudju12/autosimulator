package graphics

import (
	"autosimulator/src/machine"
	"fmt"
	"os"
	"runtime"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	WITDH, HEIGTH   = 800, 600
	TITLE           = "Simulador de Autômato"
	FONT_PATH       = "/home/jonathan/programacao/autosimulator/src/graphics/assets/IBMPlexMono-ExtraLight.ttf"
	FONT_ZIE        = 12
	FPS_DEFAULT     = 60
	FPS_EXECUTING   = 2
	BLACK_RING      = "BLACK_RING"
	BLACK_RING_PATH = "/home/jonathan/programacao/autosimulator/src/graphics/assets/ring.png"
	RED_RING        = "RED_RING"
	RED_RING_PATH   = "/home/jonathan/programacao/autosimulator/src/graphics/assets/red_ring.png"
	GREEN_RING      = "GREEN_RING"
	GREEN_RING_PATH = "/home/jonathan/programacao/autosimulator/src/graphics/assets/green_ring.png"
)

type (
	environment struct {
		w        *SDLWindow
		dragInfo *drag
		radio    *machineChannel
		mousePos *sdl.Point
		states   map[string]*graphicalState
	}

	SDLWindow struct {
		window       *sdl.Window
		renderer     *sdl.Renderer
		font         *ttf.Font
		cacheWords   map[string]*sdl.Surface
		cacheSprites map[string]*sdl.Surface
		terminate    bool
		fps          uint32
	}

	drag struct {
		clickOffset   *sdl.Point
		selected      *graphicalState
		leftMouseDown bool
	}

	machineChannel struct {
		activeMachine machine.Machine
		channel       chan int
		lastMsg       int
		inExecution   bool
		lastState     string
	}
)

// var _TESTE = []string{"a", "a", "c", "d", "c", "c"}

var _TESTE = []string{"a", "a", "c", "d", "d", "c"}

func Mainloop(env *environment) {
	window := env.w
	runtime.LockOSThread()

	for !window.terminate {
		talk(env)
		pollEvent(env)
		draw(env)
		sdl.Delay(window.fps)
	}

	defer env.Destroy()
}

func PopulateEnvironment(window *SDLWindow, activeMachine machine.Machine) *environment {
	dragInfo := &drag{clickOffset: &sdl.Point{X: 0, Y: 0}, selected: nil, leftMouseDown: false}
	radio := &machineChannel{channel: make(chan int), lastMsg: machine.STATE_NOT_CHANGE, inExecution: false, activeMachine: activeMachine}

	checkError := func(err error, name string) {
		if err != nil {
			fmt.Printf("erro ao carregar %s: %v\n", err, name)
			os.Exit(1)
		}
	}

	blackRing, err := img.Load(BLACK_RING_PATH)
	checkError(err, BLACK_RING)

	redRing, err := img.Load(RED_RING_PATH)
	checkError(err, RED_RING)

	greenRing, err := img.Load(GREEN_RING_PATH)
	checkError(err, GREEN_RING)

	window.cacheSprites[BLACK_RING] = blackRing
	window.cacheSprites[RED_RING] = redRing
	window.cacheSprites[GREEN_RING] = greenRing

	env := &environment{
		w:        window,
		dragInfo: dragInfo,
		radio:    radio,
		states:   machineStates(activeMachine),
	}

	return env
}

func NewSDLWindow() *SDLWindow {
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

	cacheSprites := make(map[string]*sdl.Surface)
	cacheWords := make(map[string]*sdl.Surface)

	return &SDLWindow{window: window,
		renderer:     renderer,
		font:         font,
		terminate:    false,
		fps:          1000 / FPS_DEFAULT,
		cacheWords:   cacheWords,
		cacheSprites: cacheSprites,
	}
}

func (env *environment) Destroy() {
	sdl.Quit()
	ttf.Quit()
	env.w.window.Destroy()
	env.w.renderer.Destroy()
	env.w.font.Close()
	close(env.radio.channel)

	for _, v := range env.w.cacheSprites {
		v.Free()
	}

	for _, v := range env.w.cacheWords {
		v.Free()
	}
}

func talk(env *environment) {
	radio := env.radio
	if radio.inExecution {
		radio.lastMsg = <-radio.channel
	}

	window := env.w
	msg := radio.lastMsg
	switch msg {
	case machine.STATE_CHANGE:
		if radio.lastState != "" {
			// remove a flag de "estado atual" do estado passado
			env.states[radio.lastState].spriteName = BLACK_RING
		}

		// Marca a flag "estado atual" no estado atual
		currentState := radio.activeMachine.CurrentState()
		env.states[currentState].spriteName = RED_RING

		// gravando pra saber qm foi o ultimo quando voltar aqui
		radio.lastState = currentState
		radio.lastMsg = machine.STATE_NOT_CHANGE

	case machine.STATE_INPUT_ACCEPTED:
		currentState := radio.activeMachine.CurrentState()
		env.states[currentState].spriteName = GREEN_RING
		radio.inExecution = false
		window.Fps(FPS_DEFAULT)

	case machine.STATE_INPUT_REJECTED:
		// vou fazer algo com isso depois (eu acho)
		radio.inExecution = false
		window.Fps(FPS_DEFAULT)

	default:
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
	radio := env.radio
	window := env.w

	switch event.Keysym.Sym {
	case sdl.K_RETURN:
		if !radio.inExecution {
			window.Fps(FPS_EXECUTING)
			radio.inExecution = true
			go machine.Execute(radio.activeMachine, _TESTE, radio.channel)
		}

	case sdl.K_r:
		if !radio.inExecution {
			window.Fps(FPS_DEFAULT)
			radio.inExecution = false
			env.states[radio.lastState].spriteName = BLACK_RING
		}
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
	}
}

func draw(env *environment) {
	window := env.w
	window.cleanUp()
	for _, state := range env.states {
		state.Draw(window, env.states)
	}

	window.renderer.Present()
}

func (w *SDLWindow) cleanUp() {
	w.renderer.SetDrawColor(255, 255, 255, 255)
	w.renderer.Clear()
}

func (w *SDLWindow) Fps(amout uint32) {
	w.fps = 1000 / amout
}
