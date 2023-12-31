package graphics

import (
	"autosimulator/src/machine"
	"errors"
	"fmt"
	"math/rand"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

type (
	graphicalState struct {
		*sdl.Rect
		sdl.Color
		state string

		// As keys para os estados que este aponta
		statesKeys []string
	}
)

const (
	// Constantes para desenhar os estados
	WIDTH_REC  = 50
	HEIGTH_REC = 50
	X_FIRST    = 10
	Y_FIRST    = 10
)

func NewState(rect *sdl.Rect, state string, color sdl.Color, statesKeys []string) *graphicalState {
	return &graphicalState{
		Rect:       rect,
		state:      state,
		Color:      color,
		statesKeys: statesKeys,
	}
}

func (s *graphicalState) Center() sdl.Point {
	return Center(s.Rect)
}

func (s *graphicalState) Draw(w *_SDLWindow, states map[string]*graphicalState) error {
	renderer := w.renderer
	words := w.cacheWords
	var err error

	outerRadius := s.Rect.W / 2
	innerRadius := (95 * outerRadius) / 100 // % do outerRadius
	err = s.drawRing(renderer, outerRadius, innerRadius, s.Color)
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

	err = s.drawText(w, words)
	if err != nil {
		return err
	}

	return nil
}

func (s *graphicalState) drawRing(renderer *sdl.Renderer, r1, r2 int32, color sdl.Color) error {

	var i int32
	center := Center(s.Rect)
	for i = 0; i < r1-r2; i++ {
		drawCircle(renderer, center.X, center.Y, (r1 - i), color)
	}

	return nil
}

// https://github.com/k4zmu2a/SpaceCadetPinball/blob/master/SpaceCadetPinball/DebugOverlay.cpp
func drawCircle(renderer *sdl.Renderer, x, y, radius int32, color sdl.Color) error {
	var t int32 = 256
	var points []sdl.Point = make([]sdl.Point, t)
	var pointCount int32 = 0
	var offsetx int32 = 0
	var offsety int32 = radius
	var d int32 = radius - 1

	var err error
	renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	for offsety >= offsetx {
		if (pointCount + 8) > t {
			err = renderer.DrawPoints(points)
			pointCount = 0
			if err != nil {
				break
			}
		}

		points[pointCount] = sdl.Point{X: x + offsetx, Y: y + offsety}
		pointCount++
		points[pointCount] = sdl.Point{X: x + offsety, Y: y + offsetx}
		pointCount++
		points[pointCount] = sdl.Point{X: x - offsetx, Y: y + offsety}
		pointCount++
		points[pointCount] = sdl.Point{X: x - offsety, Y: y + offsetx}
		pointCount++
		points[pointCount] = sdl.Point{X: x + offsetx, Y: y - offsety}
		pointCount++
		points[pointCount] = sdl.Point{X: x + offsety, Y: y - offsetx}
		pointCount++
		points[pointCount] = sdl.Point{X: x - offsetx, Y: y - offsety}
		pointCount++
		points[pointCount] = sdl.Point{X: x - offsety, Y: y - offsetx}
		pointCount++

		if d >= 2*offsetx {
			d -= 2*offsetx + 1
			offsetx += 1
		} else if d < 2*(radius-offsety) {
			d += 2*offsety - 1
			offsety -= 1
		} else {
			d += 2 * (offsety - offsetx - 1)
			offsety -= 1
			offsetx += 1
		}
	}

	if pointCount > 0 {
		err = renderer.DrawPoints(points)
	}

	renderer.SetDrawColor(0, 0, 0, 255)
	return err
}

func (s *graphicalState) drawText(window *_SDLWindow, words map[string]*sdl.Surface) error {
	renderer := window.renderer
	// var fontRating int32 = 2

	// verifica o cache de words
	surface, err := window.textSurface(s.state, COLOR_DEFAULT)
	if err != nil {
		return err
	}

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return err
	}
	_, _, fontW, fontH, err := texture.Query()
	if err != nil {
		return err
	}

	centerS := s.Center()
	textRect := &sdl.Rect{
		X: centerS.X - fontW/2,
		Y: centerS.Y - fontH/2,
		W: fontW,
		H: fontH,
	}

	renderer.Copy(texture, nil, textRect)
	return nil
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

	record := ui.bufferComputation.History[ui.indexComputation]
	details := record.Details()
	previus := ui.states[details["LAST_STATE"]]

	color := COLOR_DEFAULT
	if previus == from && machine.INITIAL != details["RESULT"] {
		color = to.Color
	}

	if from == to {
		var start int32 = 250
		var end int32 = 100
		gfx.ArcRGBA(renderer,
			fromCenter.X+int32(radius),
			fromCenter.Y,
			int32(radius*0.7),
			start,
			end,
			color.R,
			color.G,
			color.B,
			color.A,
		)

		//Desenha o marcador de cardinalidade no final da linha
		ok := gfx.FilledCircleColor(renderer, fromCenter.X+int32(radius), fromCenter.Y+int32(radius*0.7), radiusMiniBall, color)
		if !ok {
			return errors.New("erro ao renderizar as linhas")
		}

	} else {
		// Calcula o ponto inicial e final da linha
		start, end, ok := LinePoints(fromCenter, toCenter, radius, radius+float64(radiusMiniBall))
		if !ok {
			return nil
		}

		if start.X < 0 {
			fmt.Println(fromCenter, toCenter, radius, radius+float64(radiusMiniBall))
			fmt.Println(start, end)
		}

		// Desenha a linha
		ok = gfx.ThickLineColor(renderer, start.X, start.Y, end.X, end.Y, thickness, color)
		if !ok {
			return errors.New("erro ao renderizar as linhas")
		}

		// Desenha o marcador de cardinalidade no final da linha
		ok = gfx.FilledCircleColor(renderer, end.X, end.Y, radiusMiniBall, color)
		if !ok {
			return errors.New("erro ao renderizar as linhas")
		}
	}

	return nil
}

func machineStates(env *environment) map[string]*graphicalState {
	machine := env.machine
	window := env.w
	result := make(map[string]*graphicalState)
	for _, state := range machine.GetStates() {
		rect := &sdl.Rect{
			// X: X_FIRST,
			// Y: Y_FIRST + int32(i*HEIGTH_REC),
			X: WIDTH_REC + rand.Int31n(window.WIDTH-WIDTH_REC*2),
			Y: rand.Int31n(window.HEIGHT / 2),
			W: WIDTH_REC,
			H: HEIGTH_REC,
		}

		statesKeys := make([]string, 0)
		for _, transition := range machine.GetTransitions(state) {
			statesKeys = append(statesKeys, transition.GetResultState())
		}

		result[state] = NewState(rect, state, COLOR_DEFAULT, statesKeys)
	}

	return result
}
