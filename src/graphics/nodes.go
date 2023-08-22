package graphics

import (
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
		nextStates []*graphicalState
	}
)

const RING_PATH = "/home/jonathan/programacao/autosimulator/src/graphics/assets/ring.png"

var (
	BLACK = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	WHITE = sdl.Color{R: 255, G: 255, B: 255, A: 255}
)

func newState(rect *sdl.Rect, state string, colour sdl.Color) *graphicalState {
	return &graphicalState{
		Rect:       rect,
		state:      state,
		color:      colour,
		nextStates: []*graphicalState{},
	}
}

func drawLines(renderer *sdl.Renderer, font *ttf.Font, states []*graphicalState, thickness int32) {
	for _, state := range states {
		for _, nextState := range state.nextStates {
			state.drawLine(renderer, nextState, 5)
		}
	}
}

func (n *graphicalState) draw(renderer *sdl.Renderer, font *ttf.Font) {
	n.drawRing(renderer)
	n.drawText(renderer, font)
}

func (n *graphicalState) drawRing(renderer *sdl.Renderer) {
	// TODO: GLOBAL
	imgSurface, err := img.Load(RING_PATH)
	if err != nil {
		panic(err)
	}

	texture, err := renderer.CreateTextureFromSurface(imgSurface)
	if err != nil {
		panic(err)
	}

	renderer.Copy(texture, nil, n.Rect)
}

// TODO: Caching fotns
func (n *graphicalState) drawText(renderer *sdl.Renderer, font *ttf.Font) {
	var fontRating int32 = 2
	surface, err := font.RenderUTF8Solid(n.state, n.color)
	if err != nil {
		panic(err)
	}
	defer surface.Free()

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err)
	}

	textRect := &sdl.Rect{
		X: n.X + n.W/(fontRating*2),
		Y: n.Y + n.H/(fontRating*2),
		W: n.W / fontRating,
		H: n.H / fontRating,
	}

	renderer.Copy(texture, nil, textRect)
}

func (n *graphicalState) Center() sdl.Point {
	return sdl.Point{
		X: n.X + n.W/2,
		Y: n.Y + n.H/2,
	}
}

func (from *graphicalState) drawLine(renderer *sdl.Renderer, to *graphicalState, thickness int32) {
	fromCenter := from.Center()
	toCenter := to.Center()
	start, end := PointsFromRadius(fromCenter, toCenter, float64(from.H/2))
	gfx.ThickLineColor(renderer, start.X, start.Y, end.X, end.Y, thickness, BLACK)
}

func (s *graphicalState) addNextState(nextState *graphicalState) {
	s.nextStates = append(s.nextStates, nextState)
}
