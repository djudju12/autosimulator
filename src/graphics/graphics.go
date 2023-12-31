package graphics

import (
	"autosimulator/src/collections"
	"autosimulator/src/machine"
	"autosimulator/src/machine/oneStackMachine"
	"autosimulator/src/machine/twoStackMachine"
	"autosimulator/src/reader"
	"autosimulator/src/utils"
	"fmt"
	"path/filepath"
	"runtime"
	"unicode/utf8"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	TITLE         = "Simulador de Autômato"
	FONT_PATH     = "src/graphics/assets/IBMPlexMono-ExtraLight.ttf"
	EXAMPLES_PATH = "machines"
	INPUT_PATH    = "inputs"
	FONT_SIZE     = 24
	FPS_DEFAULT   = 60
	WITDH, HEIGHT = 580, 750
)

type (
	environment struct {
		w         *_SDLWindow
		machine   machine.Machine
		input     *collections.Fita
		terminate bool
		running   bool
		typing    bool
	}

	_SDLWindow struct {
		window     *sdl.Window
		renderer   *sdl.Renderer
		font       *ttf.Font
		WIDTH      int32
		HEIGHT     int32
		cacheWords map[string]*sdl.Surface
	}

	uiComponents struct {
		states            map[string]*graphicalState
		bufferComputation machine.Computation
		waitingFile       bool
		menuMode          bool
		menuInfo          *menu
		dragInfo          *drag
		computationHist
		*stackHist
	}

	menu struct {
		currentMenu *SelectBox
		menus       map[string]*SelectBox
	}

	drag struct {
		clickOffset   *sdl.Point
		selected      *graphicalState
		leftMouseDown bool
		mousePos      *sdl.Point
	}

	computationHist struct {
		indexComputation int
		bufferInput      []string
	}

	stackHist struct {
		stackA [][]string
		stackB [][]string
	}
)

var (
	main = &SelectBox{
		Name:         "main",
		CurrentIndex: 1,
		MaxItems:     4,
		MaxLen:       13,
		Options:      []string{"Machines", "New Input", "Save Input", "Load Input"},
	}

	menus = map[string]*SelectBox{
		"main": main,
		"explorer": {
			Name:         "explorer",
			CurrentIndex: 1,
			MaxItems:     14,
			MaxLen:       25,
			Options:      nil,
		},
		"input": {
			Name: "input",
		},
		"load_input": {
			Name:         "load_input",
			CurrentIndex: 1,
			MaxItems:     14,
			MaxLen:       25,
			Options:      nil,
		},
	}

	ui *uiComponents = &uiComponents{
		menuInfo: &menu{
			currentMenu: menus["main"],
			menus:       menus,
		},
		stackHist: &stackHist{
			[][]string{},
			[][]string{},
		},
	}

	fpsTimer       uint64
	delayAnimation = 0.5

	BLACK = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	WHITE = sdl.Color{R: 255, G: 255, B: 255, A: 255}
	BLUE  = sdl.Color{R: 0, G: 0, B: 255, A: 255}
	RED   = sdl.Color{R: 255, G: 0, B: 0, A: 255}
	PINK  = sdl.Color{R: 255, G: 80, B: 80, A: 255}
	GREEN = sdl.Color{R: 0, G: 255, B: 0, A: 255}

	COLOR_DEFAULT   = sdl.Color{R: 235, G: 174, B: 52, A: 255}
	COLOR_BACKGROUD = sdl.Color{R: 18, G: 18, B: 18, A: 255}

	typedInput []string
)

func Mainloop(env *environment) {
	runtime.LockOSThread() // sdl2 precisa rodar na main thread.
	ui.init(env, true)
	for !env.terminate {
		pollEvent(env)
		draw(env)
		sdl.Delay(1000 / FPS_DEFAULT)
	}
	defer env.Destroy()
}

func PopulateEnvironment(window *_SDLWindow, activeMachine machine.Machine) *environment {
	env := &environment{
		w:         window,
		machine:   activeMachine,
		terminate: false,
		running:   false,
		input:     activeMachine.GetInput(),
	}

	return env
}

func NewSDLWindow() *_SDLWindow {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow(TITLE, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WITDH, HEIGHT, sdl.WINDOW_SHOWN)
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

	font, err := ttf.OpenFont(FONT_PATH, FONT_SIZE)
	if err != nil {
		panic(err)
	}

	cacheWords := make(map[string]*sdl.Surface)

	return &_SDLWindow{window: window,
		renderer:   renderer,
		font:       font,
		cacheWords: cacheWords,
		WIDTH:      WITDH,
		HEIGHT:     HEIGHT,
	}
}

func (env *environment) Destroy() {
	env.w.font.Close()
	env.w.window.Destroy()
	env.w.renderer.Destroy()
	for _, v := range env.w.cacheWords {
		v.Free()
	}

	ttf.Quit()
	sdl.Quit()
}

func pollEvent(env *environment) {
	event := sdl.PollEvent()
	for event != nil {
		switch event := event.(type) {
		case *sdl.QuitEvent:
			fmt.Printf("Quiting....")
			env.Quit()

		case *sdl.KeyboardEvent:
			err := handleKeyboardEvents(event, env)
			if err != nil {
				env.throw(err)
			}

		case *sdl.MouseButtonEvent:
			handleMouseButtonEvents(event, env)

		case *sdl.MouseMotionEvent:
			handleMouseMotionEvent(env)

		case *sdl.TextInputEvent:
			if env.typing {
				r, _ := utf8.DecodeRune(event.Text[:])
				typedInput = append(typedInput, string(r))
			}

		default:
		}

		event = sdl.PollEvent()
	}
}

func handleKeyboardEvents(event *sdl.KeyboardEvent, env *environment) error {
	var err error
	// Por simplicidade vou lidar apenas com teclas apertadas para baixo
	if event.Type != sdl.KEYDOWN {
		return err
	}

	// Handle digitação
	if env.typing {
		switch event.Keysym.Sym {
		case sdl.K_RETURN:
			env.input = collections.FitaFromArray(typedInput)
			ui.init(env, false)

		case sdl.K_BACKSPACE:
			if len(typedInput) > 0 {
				typedInput = typedInput[:len(typedInput)-1]
			}

		case sdl.K_ESCAPE:
			previusInput := env.input.PeekAll()
			typedInput = previusInput[:len(previusInput)-1]
			ui.closeMenus(env)

		default:
		}
	}

	// Handle menus
	if ui.menuMode {
		menu := ui.menuInfo.currentMenu

		switch event.Keysym.Sym {
		case sdl.K_UP:
			menu.CurrentIndex--

		case sdl.K_DOWN:
			menu.CurrentIndex++

		case sdl.K_RETURN:
			err = ui.changeMenu(env)

		case sdl.K_m:
			ui.closeMenus(env)

		default:
		}

		return err
	}

	// Eventos fora do menu
	switch event.Keysym.Sym {
	case sdl.K_DOWN:
		ui.nextComputation()

	case sdl.K_UP:
		ui.previusComputation()

	case sdl.K_SPACE:
		// toggle running
		env.running = !env.running

	case sdl.K_r:
		ui.reset(env)

	case sdl.K_m:
		ui.menuMode = !ui.menuMode

	case sdl.K_EQUALS:
		delayAnimation += 0.1

	case sdl.K_MINUS:
		delayAnimation -= 0.1

	default:
	}

	return err
}

func (ui *uiComponents) closeMenus(env *environment) {
	ui.menuInfo.currentMenu.CurrentIndex = 1
	ui.menuInfo.currentMenu = ui.menuInfo.menus["main"]
	ui.menuInfo.currentMenu.CurrentIndex = 1
	ui.menuMode = false
	env.stopTyping()
}

func (ui *uiComponents) changeMenu(env *environment) error {
	var err error
	switch ui.menuInfo.currentMenu.Name {
	case "main":
		menuSelected := ui.menuInfo.currentMenu.CurrentIndex
		switch menuSelected {
		case 1: // MACHINES
			ui.menuInfo.currentMenu = ui.menuInfo.menus["explorer"]

		case 2: // NEW INPUT
			ui.menuInfo.currentMenu = ui.menuInfo.menus["input"]
			env.startTyping()

		case 3: // SAVE INPUT
			err = env.saveInput()
			if err != nil {
				return err
			}

			ui.closeMenus(env)

		case 4: // LOAD INPUT
			ui.menuInfo.currentMenu = ui.menuInfo.menus["load_input"]

		default:
		}

	case "explorer":
		selectedPath := ui.menuInfo.currentMenu.CurrentIndex
		m, err := reader.ReadMachine(filepath.Join(EXAMPLES_PATH, ui.menuInfo.currentMenu.Options[selectedPath-1]))
		if err != nil {
			return err
		}

		ui.menuInfo.currentMenu.Options = nil
		env.loadMachine(m)
		ui.init(env, true)

	case "load_input":
		selectedPath := ui.menuInfo.currentMenu.CurrentIndex
		i, err := reader.ReadInput(filepath.Join(INPUT_PATH, ui.menuInfo.currentMenu.Options[selectedPath-1]))
		if err != nil {
			return err
		}

		ui.menuInfo.currentMenu.Options = nil
		env.input = i
		ui.init(env, false)

	default:
	}

	return err
}

func handleMouseButtonEvents(event *sdl.MouseButtonEvent, env *environment) {
	if ui.menuMode || ui.waitingFile {
		return
	}

	dragInfo := ui.dragInfo
	mousePos := dragInfo.mousePos
	states := ui.states

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
	if ui.menuMode || ui.waitingFile {
		return
	}

	x, y, _ := sdl.GetMouseState()
	dragInfo := ui.dragInfo
	dragInfo.mousePos = &sdl.Point{X: x, Y: y}
	if dragInfo.leftMouseDown && dragInfo.selected != nil {
		dragInfo.selected.X = dragInfo.mousePos.X - dragInfo.clickOffset.X
		dragInfo.selected.Y = dragInfo.mousePos.Y - dragInfo.clickOffset.Y
	}
}

func draw(env *environment) {
	window := env.w
	ui.update(env)
	err := window.cleanUp()
	if err != nil {
		env.throw(err)
	}

	err = drawNodes(env)
	if err != nil {
		env.throw(err)
	}

	err = drawUi(env)
	if err != nil {
		env.throw(err)
	}

	window.renderer.Present()
}

func (w *_SDLWindow) cleanUp() error {

	err := w.renderer.SetDrawColor(COLOR_BACKGROUD.R, COLOR_BACKGROUD.G, COLOR_BACKGROUD.B, COLOR_BACKGROUD.A)
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
	err := ui.drawFita(env.w)
	if err != nil {
		return err
	}

	machineType := env.machine.Type()
	if machineType != machine.SIMPLE_MACHINE {
		err = ui.drawStacks(env.w)
		if err != nil {
			return err
		}

	}

	err = ui.drawHist(env.w)
	if err != nil {
		return err
	}

	if ui.menuMode {
		if ui.menuInfo.currentMenu == nil {
			ui.menuInfo.currentMenu = ui.menuInfo.menus["main"]
		}

		ui.drawMenu(env.w)
		if err != nil {
			return err
		}
	}

	if ui.waitingFile {
		err = ui.waitForFile(env.w)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ui *uiComponents) update(env *environment) {
	// se estiver rodando, atualiza a cada delayAnimation
	// Ou seja, anima a cada delayAnimation milisegundos
	if env.running {
		now := sdl.GetTicks64()
		if now > fpsTimer+uint64(delayAnimation*1000) {
			ui.nextComputation()
			env.running = ui.indexComputation != (len(ui.bufferComputation.History) - 1)
			fpsTimer = now
		}
	}

	// Pinta todos os estados com a cor default
	for _, state := range ui.states {
		state.Color = COLOR_DEFAULT
	}

	// Atualiza o buffer que printa a fita
	ui.bufferInput = ajustBufferInput(env.machine.GetInput(), ui.indexComputation)

	// Historico da computação atual
	record := ui.bufferComputation.History[ui.indexComputation]
	details := record.Details()

	// Cor do proximo estado
	nextSate := ui.states[details["NEXT_STATE"]]
	switch details["RESULT"] {
	case machine.INITIAL:
		nextSate.Color = BLUE
	case machine.ACCEPTED:
		nextSate.Color = GREEN
	default:
		nextSate.Color = RED
	}
}

func (ui *uiComponents) init(env *environment, redraw bool) {
	dragInfo := &drag{
		clickOffset:   &sdl.Point{X: 0, Y: 0},
		selected:      nil,
		leftMouseDown: false,
		mousePos:      &sdl.Point{X: 0, Y: 0},
	}

	ui.closeMenus(env)
	bufferInput := ajustBufferInput(env.input, 0)
	computation := machine.Execute(env.machine, env.input)

	if env.machine.Type() == machine.TWO_STACK_MACHINE {
		machine, _ := env.machine.(*twoStackMachine.Machine)
		ui.stackA, ui.stackB = machine.StackHistory()
	} else if env.machine.Type() == machine.ONE_STACK_MACHINE {
		machine, _ := env.machine.(*oneStackMachine.Machine)
		ui.stackA = machine.StackHistory()
		if ui.stackB != nil {
			ui.stackB = nil
		}
	}

	if redraw {
		ui.states = machineStates(env)
		ui.dragInfo = dragInfo
	}

	ui.indexComputation = 0
	ui.bufferComputation = *computation
	ui.bufferInput = bufferInput
	env.machine.GetInput().Reset()

	initial := ui.bufferComputation.History[0]
	initalDetails := initial.Details()
	firstState := ui.states[initalDetails["LAST_STATE"]]
	firstState.Color = BLUE

	env.running = false
}

func (ui *uiComponents) reset(env *environment) {
	bufferInput := ajustBufferInput(env.machine.GetInput(), 0)
	ui.bufferInput = bufferInput
	ui.indexComputation = 0
	initial := ui.bufferComputation.History[0]
	initalDetails := initial.Details()
	firstState := ui.states[initalDetails["LAST_STATE"]]
	firstState.Color = BLUE
	env.running = false
}

func drawNodes(env *environment) error {
	var err error

	for _, state := range ui.states {
		err = state.Draw(env.w, ui.states)
		if err != nil {
			return err
		}
	}

	return nil
}

func ajustBufferInput(input *collections.Fita, index int) []string {
	arrayInput := input.ToArray()
	return utils.AjustMaxLen(arrayInput, index, TAMANHO_ESTRUTURAS)
}

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

func (st *stackHist) get(i int) ([]string, []string) {
	var a []string
	var b []string
	if st.stackA != nil {
		a = st.stackA[i]
		if len(a) > TAMANHO_ESTRUTURAS {
			a = a[len(a)-TAMANHO_ESTRUTURAS:]
		}
	}

	if st.stackB != nil {
		b = st.stackB[i]
		if len(b) > TAMANHO_ESTRUTURAS {
			b = b[len(b)-TAMANHO_ESTRUTURAS:]
		}
	}

	return a, b
}

func (env *environment) Quit() {
	env.terminate = true
}

func (ui *uiComponents) previusComputation() {
	if ui.indexComputation > 0 {
		ui.indexComputation--
	}
}

func (ui *uiComponents) nextComputation() {
	if ui.indexComputation < len(ui.bufferComputation.History)-1 {
		ui.indexComputation++
	}
}

func (env *environment) stopTyping() {
	env.typing = false
}

func (env *environment) startTyping() {
	if typedInput == nil {
		currentInput := env.input.ToArray()
		typedInput = currentInput[:len(currentInput)-1]
	}

	env.typing = true
}

func (env *environment) loadMachine(machine machine.Machine) {
	env.machine = machine
	env.input = machine.GetInput()
}

func (env *environment) saveInput() error {
	return reader.WriteInput(env.input, INPUT_PATH)
}

func (env *environment) throw(err error) {
	fmt.Println(err)
	env.Quit()
}
