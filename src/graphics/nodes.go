package graphics

import (
	"autosimulator/src/machine"
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type (
	graphicalState struct {
		*sdl.Rect
		state      string
		color      sdl.Color
		statesKeys []string
		isCurrent  bool
	}
)

const (
	BLACK_RING_PATH = "/home/jonathan/programacao/autosimulator/src/graphics/assets/ring.png"
	RED_RING_PATH   = "/home/jonathan/programacao/autosimulator/src/graphics/assets/read_ring.png"
	ARROW_HEAD_PATH = "/home/jonathan/programacao/autosimulator/src/graphics/assets/arrow_head.png"

	// Constantes para desenhar os estados
	WIDTH_REC  = 50
	HEIGTH_REC = 50
	X_FIRST    = 10
	Y_FIRST    = 10
)

var (
	BLACK = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	WHITE = sdl.Color{R: 255, G: 255, B: 255, A: 255}
)

func NewState(rect *sdl.Rect, state string, colour sdl.Color, statesKeys []string) *graphicalState {
	return &graphicalState{
		Rect:       rect,
		state:      state,
		color:      colour,
		statesKeys: statesKeys,
		isCurrent:  false,
	}
}

// func (s *graphicalState) AddNextState(nextState *graphicalState) {
// 	s.nextStates = append(s.nextStates, *nextState)
// }

func (s *graphicalState) Center() sdl.Point {
	return sdl.Point{
		X: s.X + s.W/2,
		Y: s.Y + s.H/2,
	}
}

func (s *graphicalState) Draw(renderer *sdl.Renderer, font *ttf.Font, states map[string]*graphicalState) {
	s.drawRing(renderer)
	//TODO: what i gonna do with this fontRating thing that i created?
	s.drawText(renderer, font, 2)

	if len(s.statesKeys) != 0 {
		s.drawLines(renderer, font, states, 2)
	}
}

func (s *graphicalState) drawRing(renderer *sdl.Renderer) {
	// TODO: GLOBAL
	var path string
	if s.isCurrent {
		path = RED_RING_PATH
	} else {
		path = BLACK_RING_PATH
	}

	imgSurface, err := img.Load(path)
	if err != nil {
		fmt.Printf("Erro ao carregar a imagem: %v\n", err)
		os.Exit(1)
	}

	texture, err := renderer.CreateTextureFromSurface(imgSurface)
	if err != nil {
		fmt.Printf("Erro ao criar a textura da imagem: %v\n", err)
		os.Exit(1)
	}

	renderer.Copy(texture, nil, s.Rect)
}

func (s *graphicalState) drawText(renderer *sdl.Renderer, font *ttf.Font, fontRating int32) {
	if fontRating == 0 {
		fmt.Printf("Erro: fontRating não pode ser 0\n")
		os.Exit(1)
	}

	// TODO: Caching fonts
	surface, err := font.RenderUTF8Solid(s.state, s.color)
	if err != nil {
		fmt.Printf("Erro ao renderizar a fonte: %v\n", err)
		os.Exit(1)
	}
	defer surface.Free()

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		fmt.Printf("Erro ao criar a textura da fonte: %v\n", err)
		os.Exit(1)
	}

	textRect := &sdl.Rect{
		X: s.X + s.W/(fontRating*2),
		Y: s.Y + s.H/(fontRating*2),
		W: s.W / fontRating,
		H: s.H / fontRating,
	}

	renderer.Copy(texture, nil, textRect)
}

func (s *graphicalState) drawLines(renderer *sdl.Renderer, font *ttf.Font, states map[string]*graphicalState, thickness int32) {
	// Desenha os estados cujo o estado atual aponta
	for _, next := range s.statesKeys {
		state := states[next]
		if state != nil {
			s.drawLine(renderer, state, thickness)
		}

	}
}

// Função que desenha uma linah entre dois estados. O estato "To" recebera a bolinha (cardinalidade)!!
func (from *graphicalState) drawLine(renderer *sdl.Renderer, to *graphicalState, thickness int32) {
	fromCenter := from.Center()
	toCenter := to.Center()
	radius := float64(from.H / 2)
	radiusMiniBall := thickness * 2

	// Calcula o ponto inicial e final da linha
	start, end := LinePoints(fromCenter, toCenter, radius, radius+float64(radiusMiniBall))

	// Desenha a linha
	ok := gfx.ThickLineColor(renderer, start.X, start.Y, end.X, end.Y, thickness, BLACK)
	if !ok {
		fmt.Printf("Erro ao renderizar as linhas")
		os.Exit(1)
	}

	// Desenha o marcador de cardinalidade no final da linha
	ok = gfx.FilledCircleColor(renderer, end.X, end.Y, radiusMiniBall, BLACK)
	if !ok {
		fmt.Printf("Erro ao renderizar as mini bolas")
		os.Exit(1)
	}
}

func machineStates(machine machine.Machine) map[string]*graphicalState {
	result := make(map[string]*graphicalState)
	for i, state := range machine.GetStates() {
		rect := &sdl.Rect{
			X: X_FIRST,
			Y: Y_FIRST + int32(i*HEIGTH_REC),
			W: WIDTH_REC,
			H: HEIGTH_REC,
		}

		statesKeys := make([]string, 0)
		for _, transition := range machine.GetTransitions(state) {
			statesKeys = append(statesKeys, transition.GetResultState())
		}

		result[state] = NewState(rect, state, BLACK, statesKeys)
	}

	return result
}
