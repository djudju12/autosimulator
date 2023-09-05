package graphics

import (
	"autosimulator/src/collections"
	"autosimulator/src/machine"
	"autosimulator/src/machine/oneStackMachine"
	"autosimulator/src/machine/twoStackMachine"
	"autosimulator/src/reader"
	"fmt"
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	TITLE           = "Simulador de Autômato"
	FONT_PATH       = "/home/jonathan/hd/programacao/autosimulator/src/graphics/assets/IBMPlexMono-ExtraLight.ttf"
	FONT_SIZE       = 24
	FPS_DEFAULT     = 60
	WITDH, HEIGHT   = 580, 750
	DELAY_ANIMATION = 0.5 * 1000
)

type (
	environment struct {
		w         *_SDLWindow
		machine   machine.Machine
		input     []string
		terminate bool
		running   bool
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
		dragInfo          *drag
		computationHist
		*stackHist
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
	ui *uiComponents = &uiComponents{
		stackHist: &stackHist{
			[][]string{},
			[][]string{},
		},
	}

	fpsTimer uint64

	BLACK = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	WHITE = sdl.Color{R: 255, G: 255, B: 255, A: 255}
	BLUE  = sdl.Color{R: 0, G: 0, B: 255, A: 255}
	RED   = sdl.Color{R: 255, G: 0, B: 0, A: 255}
	GREEN = sdl.Color{R: 0, G: 255, B: 0, A: 255}

	COLOR_DEFAULT   = sdl.Color{R: 235, G: 174, B: 52, A: 255}
	COLOR_BACKGROUD = sdl.Color{R: 18, G: 18, B: 18, A: 255}
)

func Mainloop(env *environment) {
	runtime.LockOSThread() // sdl2 precisa rodar na main thread.
	ui.init(env)
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

func (env *environment) Input(fita []string) {
	env.input = append(fita, collections.TAIL_FITA)
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
	event := sdl.PollEvent()
	for event != nil {
		switch event := event.(type) {
		case *sdl.QuitEvent:
			fmt.Printf("Quiting....")
			env.Quit()

		case *sdl.KeyboardEvent:
			handleKeyboardEvents(event, env)

		case *sdl.MouseButtonEvent:
			handleMouseButtonEvents(event, env)

		case *sdl.MouseMotionEvent:
			handleMouseMotionEvent(env)

		case *sdl.DropEvent:
			handleDropEvent(event, env)

		default:
		}

		event = sdl.PollEvent()
	}
}

func handleKeyboardEvents(event *sdl.KeyboardEvent, env *environment) {
	if event.Type == sdl.KEYDOWN {
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

		default:
		}
	}
}

func handleMouseButtonEvents(event *sdl.MouseButtonEvent, env *environment) {
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
	x, y, _ := sdl.GetMouseState()
	dragInfo := ui.dragInfo
	dragInfo.mousePos = &sdl.Point{X: x, Y: y}
	if dragInfo.leftMouseDown && dragInfo.selected != nil {
		dragInfo.selected.X = dragInfo.mousePos.X - dragInfo.clickOffset.X
		dragInfo.selected.Y = dragInfo.mousePos.Y - dragInfo.clickOffset.Y
	}
}

func handleDropEvent(event *sdl.DropEvent, env *environment) {
	path := event.File
	if path == "" {
		return
	}

	m, err := reader.ReadMachine(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Nova maquina carregada com sucesso!")
	env.machine = m
	ui.init(env)
	ui.reset(env)
}

func draw(env *environment) {
	window := env.w
	ui.update(env)
	err := window.cleanUp()
	if err != nil {
		fmt.Println(err)
		env.Quit()
	}

	err = drawUi(env)
	if err != nil {
		fmt.Println(err)
		env.Quit()
	}

	err = drawNodes(env)
	if err != nil {
		fmt.Println(err)
		env.Quit()
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
	var padx, pady int32 = 5, 5
	err := ui.drawFita(env.w, padx, pady)
	if err != nil {
		return err
	}

	machineType := env.machine.Type()
	if machineType != machine.SIMPLE_MACHINE {
		err = env.drawStacks(ui.stackHist, ui.indexComputation, padx, pady)
		if err != nil {
			return err
		}

	}

	err = ui.drawHist(env.w, padx, pady)
	if err != nil {
		return err
	}

	return nil
}

func (ui *uiComponents) update(env *environment) {
	// se estiver rodando, atualiza a cada DELAY_ANIMATION
	// Ou seja, anima a cada DELAY_ANIMATION milisegundos
	if env.running {
		now := sdl.GetTicks64()
		if now > uint64(fpsTimer+DELAY_ANIMATION) {
			ui.nextComputation()
			env.running = ui.indexComputation != (len(ui.bufferComputation.History) - 1)
			fpsTimer = now
		}
	}

	// Pinta todos os estados com a cor default
	for _, state := range ui.states {
		state.color = COLOR_DEFAULT
	}

	// Atualiza o buffer que printa a fita
	ui.bufferInput = bufferMe(env.input, ui.indexComputation)

	// Historico da computação atual
	record := ui.bufferComputation.History[ui.indexComputation]
	details := record.Details()

	// Cor do proximo estado
	nextSate := ui.states[details["NEXT_STATE"]]
	switch details["RESULT"] {
	case machine.INITIAL:
		nextSate.color = BLUE
	case machine.ACCEPTED:
		nextSate.color = GREEN
	default:
		nextSate.color = RED
	}

}

func (ui *uiComponents) init(env *environment) {
	dragInfo := &drag{
		clickOffset:   &sdl.Point{X: 0, Y: 0},
		selected:      nil,
		leftMouseDown: false,
		mousePos:      &sdl.Point{X: 0, Y: 0},
	}

	bufferInput := bufferMe(env.input, 0)
	fita := collections.FitaFromArray(env.input)
	computation := machine.Execute(env.machine, fita)

	// TODO: REFATORAR
	if env.machine.Type() == machine.TWO_STACK_MACHINE {
		machine, _ := env.machine.(*twoStackMachine.Machine)
		ui.stackA, ui.stackB = machine.StackHistory()
	} else if env.machine.Type() == machine.ONE_STACK_MACHINE {
		machine, _ := env.machine.(*oneStackMachine.Machine)
		ui.stackA = machine.StackHistory()
	}

	ui.states = machineStates(env)
	ui.dragInfo = dragInfo
	ui.indexComputation = 0
	ui.bufferComputation = *computation
	ui.bufferInput = bufferInput

	initial := ui.bufferComputation.History[0]
	initalDetails := initial.Details()
	firstState := ui.states[initalDetails["LAST_STATE"]]
	firstState.color = BLUE

	env.running = false
}

func (ui *uiComponents) reset(env *environment) {
	bufferInput := bufferMe(env.input, 0)
	ui.bufferInput = bufferInput
	ui.indexComputation = 0
	initial := ui.bufferComputation.History[0]
	initalDetails := initial.Details()
	firstState := ui.states[initalDetails["LAST_STATE"]]
	firstState.color = BLUE
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

func bufferMe(input []string, index int) []string {
	// Se há apenas 1 caracter na fita
	if len(input)-index < 1 {
		return []string{input[len(input)-1]}
	}

	// Se a fita é menor que o tamanho da estrutura
	if (len(input) - index) < TAMANHO_ESTRUTURAS {
		return input[index:]
	}

	return input[index : index+TAMANHO_ESTRUTURAS]
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
	if st.stackA != nil && i < len(st.stackA) {
		a = st.stackA[i]
		if len(a) > TAMANHO_ESTRUTURAS {
			a = a[len(a)-TAMANHO_ESTRUTURAS:]
		}
	}

	// fmt.Println(st.stackB)
	if st.stackB != nil && i < len(st.stackB) {
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
