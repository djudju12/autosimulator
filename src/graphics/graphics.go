package graphics

import (
	"autosimulator/src/collections"
	"autosimulator/src/machine"
	"fmt"
	"os"
	"runtime"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	WITDH, HEIGTH   = 800, 600
	TITLE           = "Simulador de Autômato"
	FONT_PATH       = "/home/jonathan/programacao/autosimulator/src/graphics/assets/IBMPlexMono-ExtraLight.ttf"
	FONT_ZIE        = 24
	FPS_DEFAULT     = 60
	FPS_EXECUTING   = 1
	BLACK_RING      = "BLACK_RING"
	BLACK_RING_PATH = "/home/jonathan/programacao/autosimulator/src/graphics/assets/ring.png"
	RED_RING        = "RED_RING"
	RED_RING_PATH   = "/home/jonathan/programacao/autosimulator/src/graphics/assets/red_ring.png"
	GREEN_RING      = "GREEN_RING"
	GREEN_RING_PATH = "/home/jonathan/programacao/autosimulator/src/graphics/assets/green_ring.png"
	BLUE_RING       = "BLUE_RING"
	BLUE_RING_PATH  = "/home/jonathan/programacao/autosimulator/src/graphics/assets/blue_ring.png"
	FITA_HEAD_PATH  = "/home/jonathan/programacao/autosimulator/src/graphics/assets/fita_head.png"
	FITA_HEAD       = "FITA_HEAD"
	FITA_PATH       = "/home/jonathan/programacao/autosimulator/src/graphics/assets/fita.png"
	FITA            = "FITA"
)

type (
	environment struct {
		w        *_SDLWindow
		dragInfo *drag
		radio    *machineChannel
		mousePos *sdl.Point
		states   map[string]*graphicalState
	}

	_SDLWindow struct {
		window       *sdl.Window
		renderer     *sdl.Renderer
		font         *ttf.Font
		cacheWords   map[string]*sdl.Surface
		cacheSprites map[string]*sdl.Surface
		ui           []*sdl.Texture
		terminate    bool
		redraw       bool
	}

	drag struct {
		clickOffset   *sdl.Point
		selected      *graphicalState
		leftMouseDown bool
	}

	machineChannel struct {
		activeMachine machine.Machine
		input         *collections.Fita
		inputToPrint  []string
		channel       chan int
		lastMsg       int
		inExecution   bool
		lastState     string
	}
)

func Mainloop(env *environment) {
	window := env.w
	runtime.LockOSThread()
	env.Reset()
	for !window.terminate {
		talk(env)
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

	radio := &machineChannel{
		channel:       make(chan int),
		lastMsg:       machine.STATE_NOT_CHANGE,
		inExecution:   false,
		activeMachine: activeMachine,
		input:         collections.NewFita(),
		lastState:     activeMachine.CurrentState(),
	}

	checkError := func(err error, name string) {
		if err != nil {
			fmt.Printf("erro ao carregar %s: %v\n", err, name)
			os.Exit(1)
		}
	}

	fita, err := img.Load(FITA_PATH)
	checkError(err, FITA)

	fitaHead, err := img.Load(FITA_HEAD_PATH)
	checkError(err, FITA_HEAD_PATH)

	blackRing, err := img.Load(BLACK_RING_PATH)
	checkError(err, BLACK_RING)

	redRing, err := img.Load(RED_RING_PATH)
	checkError(err, RED_RING)

	greenRing, err := img.Load(GREEN_RING_PATH)
	checkError(err, GREEN_RING)

	blueRing, err := img.Load(BLUE_RING_PATH)
	checkError(err, GREEN_RING)

	window.cacheSprites[BLACK_RING] = blackRing
	window.cacheSprites[RED_RING] = redRing
	window.cacheSprites[BLUE_RING] = blueRing
	window.cacheSprites[GREEN_RING] = greenRing
	window.cacheSprites[FITA] = fita
	window.cacheSprites[FITA_HEAD] = fitaHead

	env := &environment{
		w:        window,
		dragInfo: dragInfo,
		radio:    radio,
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

	cacheSprites := make(map[string]*sdl.Surface)
	cacheWords := make(map[string]*sdl.Surface)

	return &_SDLWindow{window: window,
		renderer:     renderer,
		font:         font,
		terminate:    false,
		cacheWords:   cacheWords,
		cacheSprites: cacheSprites,
		redraw:       true,
	}
}

func (env *environment) Destroy() {
	sdl.Quit()
	ttf.Quit()
	env.w.window.Destroy()
	env.w.renderer.Destroy()
	env.w.font.Close()
	close(env.radio.channel)
	_ = free(env.w.ui)

	for _, v := range env.w.cacheSprites {
		v.Free()
	}

	for _, v := range env.w.cacheWords {
		v.Free()
	}
}

func talk(env *environment) {
	radio := env.radio
	if radio.lastMsg == machine.STATE_NOT_CHANGE {
		radio.inExecution = false
		return
	}

	radio.inExecution = true
	env.w.redraw = true
	radio.inputToPrint = radio.input.Peek(8)
	radio.lastMsg = <-radio.channel

	if radio.lastState != "" {
		env.states[radio.lastState].spriteName = BLACK_RING
	}

	// window := env.w
	msg := radio.lastMsg
	switch msg {
	case machine.STATE_CHANGE:
		env.changeState(RED_RING)

	case machine.STATE_INPUT_ACCEPTED:
		env.changeState(GREEN_RING)
		radio.lastMsg = machine.STATE_NOT_CHANGE

	case machine.STATE_INPUT_REJECTED:
		env.changeState(RED_RING)
		radio.lastMsg = machine.STATE_NOT_CHANGE

	default:
	}

}

func (env *environment) changeState(spriteName string) {
	radio := env.radio

	currentState := radio.activeMachine.CurrentState()
	env.states[currentState].spriteName = spriteName

	// para mudar a cor para o normal na proxima iteraçao
	radio.lastState = currentState
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
	switch event.Keysym.Sym {
	case sdl.K_RETURN:
		if !radio.inExecution && event.Type == sdl.KEYDOWN {
			env.RunMachine()
		}

	case sdl.K_r:
		if !radio.inExecution && event.Type == sdl.KEYDOWN {
			env.Reset()
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
		env.w.redraw = true
	}
}

func draw(env *environment) {
	window := env.w
	if env.w.redraw {
		if env.radio.inExecution {
			sdl.Delay(1000 / FPS_EXECUTING)
		}

		if len(window.ui) != 0 {
			free(window.ui)
		}

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

func drawUi(env *environment) error {
	var padx, pady int32 = 5, 5
	err := env.w.drawFita(env, 0, padx, pady)
	return err
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

// utilidade
func (env *environment) Reset() {
	fmt.Println("Reseting...")
	window := env.w
	radio := env.radio
	radio.activeMachine.Init()
	window.redraw = true
	radio.inExecution = false
	lastState := env.states[radio.lastState]
	initalState := env.states[radio.activeMachine.CurrentState()]
	radio.input.Reset()
	radio.inputToPrint = radio.input.Peek(8)

	if lastState != nil {
		lastState.spriteName = BLACK_RING
	}

	if initalState != nil {
		initalState.spriteName = BLUE_RING
	}
}

func (env *environment) RunMachine() {
	fmt.Println("Running...")
	radio := env.radio
	radio.lastMsg = machine.STATE_CHANGE
	go machine.Execute(radio.activeMachine, radio.input, radio.channel)
}

func (env *environment) Input(fita *collections.Fita) {
	env.radio.input = fita
	fita.Write(machine.TAIL_FITA)
}

func (w *_SDLWindow) textSurface(text string, color sdl.Color) (*sdl.Surface, error) {
	font := w.font
	words := w.cacheWords

	surface := words[text]
	var err error
	if surface == nil {
		surface, err = font.RenderUTF8Solid(text, color)
		if err != nil {
			return nil, err
		}

		words[text] = surface
	}

	return surface, err
}

func Center(rec *sdl.Rect) sdl.Point {
	return sdl.Point{
		X: rec.X + rec.W/2,
		Y: rec.Y + rec.H/2,
	}
}

func free(textures []*sdl.Texture) error {
	var err error
	for _, t := range textures {
		err = t.Destroy()
		if err != nil {
			return err
		}
	}

	return nil
}

func drawAxs(env *environment) {
	gfx.ThickLineColor(env.w.renderer, WITDH/2, 0, WITDH/2, HEIGTH, 2, BLACK)
	gfx.ThickLineColor(env.w.renderer, 0, HEIGTH/2, WITDH, HEIGTH/2, 2, BLACK)
}
