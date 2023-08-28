package graphics

import (
	"autosimulator/src/machine"
	"errors"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

type (
	graphicalState struct {
		*sdl.Rect
		state string
		color sdl.Color

		// As keys para os estados que este aponta
		statesKeys []string
		spriteName string
	}
)

const (
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
		spriteName: BLACK_RING,
	}
}

func (s *graphicalState) Center() sdl.Point {
	return Center(s.Rect)
}

func (s *graphicalState) Draw(w *_SDLWindow, states map[string]*graphicalState) error {
	renderer := w.renderer
	words := w.cacheWords

	outerRadius := s.Rect.W / 2
	innerRadius := (96 * outerRadius) / 100 // % do outerRadius
	err := s.drawRing(renderer, s.Rect.W/2, innerRadius, BLACK)
	if err != nil {
		return err
	}

	var lineThickness int32 = 2
	if len(s.statesKeys) != 0 {
		err = s.drawLines(renderer, states, lineThickness)
		if err != nil {
			return err
		}
	}

	textTexture, err := s.drawText(w, words)
	if err != nil {
		return err
	}

	w.ui = append(w.ui, textTexture)
	// w.ui = append(w.ui, ringTexture, textTexture)
	return nil
}

// func (s *graphicalState) drawRing(renderer *sdl.Renderer, sprites map[string]*sdl.Surface) (*sdl.Texture, error) {
// 	// TODO: GLOBAL
// 	imgSurface := sprites[s.spriteName]
// 	texture, err := renderer.CreateTextureFromSurface(imgSurface)
// 	if err != nil {
// 		return nil, err
// 	}

// 	renderer.Copy(texture, nil, s.Rect)
// 	return texture, nil
// }

func (s *graphicalState) drawRing(renderer *sdl.Renderer, r1, r2 int32, color sdl.Color) error {

	var i int32
	center := Center(s.Rect)
	for i = 0; i < r1-r2; i++ {
		gfx.AACircleColor(renderer, center.X, center.Y, (r1 - i), BLACK)
	}

	return nil
}

func (s *graphicalState) drawText(window *_SDLWindow, words map[string]*sdl.Surface) (*sdl.Texture, error) {
	renderer := window.renderer
	// var fontRating int32 = 2

	// verifica o cache de words
	surface, err := window.textSurface(s.state, s.color)
	if err != nil {
		return nil, err
	}

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return nil, err
	}
	_, _, fontW, fontH, err := texture.Query()
	if err != nil {
		return nil, err
	}

	centerS := s.Center()
	textRect := &sdl.Rect{
		X: centerS.X - fontW/2,
		Y: centerS.Y - fontH/2,
		W: fontW,
		H: fontH,
	}

	renderer.Copy(texture, nil, textRect)
	return texture, nil
}

func (s *graphicalState) drawLines(renderer *sdl.Renderer, states map[string]*graphicalState, thickness int32) error {
	// Desenha os estados cujo o estado atual aponta
	var err error
	for _, next := range s.statesKeys {
		state := states[next]
		if state != nil {
			err = s.drawLine(renderer, state, thickness)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Função que desenha uma linah entre dois estados. O estato "To" recebera a bolinha (cardinalidade)!!
func (from *graphicalState) drawLine(renderer *sdl.Renderer, to *graphicalState, thickness int32) error {
	fromCenter := from.Center()
	toCenter := to.Center()
	radius := float64(from.H / 2)
	radiusMiniBall := thickness * 2

	// Calcula o ponto inicial e final da linha
	start, end := LinePoints(fromCenter, toCenter, radius, radius+float64(radiusMiniBall))

	// Desenha a linha
	ok := gfx.ThickLineColor(renderer, start.X, start.Y, end.X, end.Y, thickness, BLACK)
	if !ok {
		return errors.New("erro ao renderizar as linhas")
	}

	// Desenha o marcador de cardinalidade no final da linha
	ok = gfx.FilledCircleColor(renderer, end.X, end.Y, radiusMiniBall, BLACK)
	if !ok {
		return errors.New("erro ao renderizar as linhas")
	}

	return nil
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
