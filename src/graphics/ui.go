package graphics

import (
	"errors"
	"fmt"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	RIGHT = iota
	LEFT  = iota
	UP    = iota
	DOWN  = iota
)

const TAMANHO_ESTRUTURAS = 9

func (ui *uiComponents) drawFita(window *_SDLWindow, padx, pady int32) error {
	// Calculo da posicao inicial da fita/texto
	var fitaWidth, thickness int32 = 32, 2
	x := padx
	y := (HEIGTH - fitaWidth) - pady

	// Rec representa o primeiro quadrado da fita
	fitaRec := sdl.Rect{
		X: x,
		Y: y,
		W: fitaWidth,
		H: fitaWidth, // É um quadrado
	}

	err := drawManyRects(window.renderer, thickness, int(TAMANHO_ESTRUTURAS), RIGHT, fitaRec, BLACK)
	if err != nil {
		return err
	}

	// Head da fita
	headBase := fitaWidth / 2
	headHeigth := headBase / 2
	err = drawArrow(window.renderer, thickness, (x + fitaWidth/2), (y - pady), headBase, headHeigth, BLACK)
	if err != nil {
		return err
	}

	// Texto
	bufferFita := ui.bufferInput
	err = drawText(window, bufferFita, fitaWidth, fitaRec.X, fitaRec.Y, RIGHT)
	if err != nil {
		return err
	}

	// Para manter uma referencia dos ponteiros que vou precisar liberar
	return nil
}

func (env *environment) drawStacks(hist *stackHist, histIndex int, padx, pady int32) error {
	var err error
	a, b := hist.get(histIndex)
	err = env.drawStack(a, 1, padx, pady)
	if err != nil {
		return err
	}

	if b != nil {
		err = env.drawStack(b, 2, padx, pady)
		if err != nil {
			return err
		}
	}

	return nil
}

func (env *environment) drawStack(stack []string, index, padx, pady int32) error {
	window := env.w

	// Calculo da posicao inicial do stack/texto
	var stackWidth, thickness int32 = 32, 2
	x := WITDH - (padx+stackWidth)*(index+1)
	y := HEIGTH - (pady + stackWidth*TAMANHO_ESTRUTURAS)

	// Esse rec represeta o primeiro quadrado do stack
	oneStackCointainer := sdl.Rect{
		X: x,
		Y: y,
		W: stackWidth,
		H: stackWidth, // É um quadrado
	}

	// Retangulos
	err := drawManyRects(window.renderer, thickness, int(TAMANHO_ESTRUTURAS), UP, oneStackCointainer, BLACK)
	if err != nil {
		return err
	}

	// Textos
	firstCharX := x
	firstCharY := y + (stackWidth * (TAMANHO_ESTRUTURAS - 1))
	err = drawText(window, stack, stackWidth, firstCharX, firstCharY, UP)
	if err != nil {
		return err
	}

	// Para manter uma referencia dos ponteiros que vou precisar liberar
	return nil
}

func drawText(window *_SDLWindow, text []string, space, x1, y1 int32, direction int) error {
	var textSurface *sdl.Surface
	var textTexture *sdl.Texture
	var x, y int32
	var err error
	for i, s := range text {
		// Essa função checa se há a palavra no cache antes de cirar a surface
		// O cache esta armazenado na SDLWindow
		textSurface, err = window.textSurface(s, BLACK)
		if err != nil {
			return err
		}

		textTexture, err = window.renderer.CreateTextureFromSurface(textSurface)
		if err != nil {
			return err
		}

		_, _, fontW, fontH, err := textTexture.Query()
		if err != nil {
			return err
		}

		switch direction {
		case UP:
			x = x1 + fontW/2
			y = y1 - (space * int32(i))
		case DOWN:
			x = x1 + fontW/2
			y = y1 + (space * int32(i))
		case RIGHT:
			x = x1 + (space * int32(i)) + fontW/2
			y = y1
		case LEFT:
			x = x1 - (space * int32(i)) - fontW/2
			y = y1
		default:
			return errors.New("direção invalida. drawManyRects()")
		}

		window.renderer.Copy(textTexture, nil, &sdl.Rect{
			X: x,
			Y: y,
			W: fontW,
			H: fontH,
		})
	}

	return nil
}

func drawManyRects(renderer *sdl.Renderer, thickness int32, amount, direction int, rect sdl.Rect, color sdl.Color) error {
	thick32 := int32(thickness)

	var newRect sdl.Rect
	var x, y, i int32
	var err error
	for i = 0; i < int32(amount); i++ {
		switch direction {
		case UP:
			x = rect.X
			y = rect.Y + (rect.H * i)
		case DOWN:
			x = rect.X
			y = rect.Y - (rect.H * i)
		case RIGHT:
			x = rect.X + (rect.W * i)
			y = rect.Y
		case LEFT:
			x = rect.X - (rect.W * i)
			y = rect.Y

		default:
			return errors.New("direção invalida. drawManyRects()")
		}

		newRect = sdl.Rect{
			X: x,
			Y: y,
			W: rect.W,
			H: rect.H,
		}

		if err = drawRect(renderer, thick32, newRect, color); err != nil {
			return err
		}
	}

	return nil
}

func drawRect(renderer *sdl.Renderer, thickness int32, rect sdl.Rect, color sdl.Color) error {
	// RectangleColor(renderer *sdl.Renderer, x1, y1, x2, y2 int32, color sdl.Color) bool {
	var i, x1, y1, x2, y2 int32
	var ok bool
	for i = 0; i < thickness; i++ {
		x1 = rect.X + i
		y1 = rect.Y + i
		x2 = x1 + rect.W
		y2 = y1 + rect.H
		if ok = gfx.RectangleColor(renderer, x1, y1, x2, y2, color); !ok {
			return fmt.Errorf("não foi possivel desenhar o retangulo: x1: %d, y1: %d, x2: %d, y2: %d", x1, y1, x2, y2)
		}
	}

	return nil
}

func drawArrow(renderer *sdl.Renderer, thickness, x, y, base, heigth int32, color sdl.Color) error {
	var ok bool
	errText := "nao foi possível desenhar a flecha"
	ok = gfx.ThickLineColor(renderer, x, y, x-base/2, y-heigth, thickness, color)
	if !ok {
		return errors.New(errText)
	}

	ok = gfx.ThickLineColor(renderer, x, y, x+base/2, y-heigth, thickness, color)
	if !ok {
		return errors.New(errText)
	}

	return nil
}
