package graphics

import (
	"autosimulator/src/utils"
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

func (env *environment) drawFita(headIndex int, padx, pady int32) error {
	window := env.w

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
	arrowBase := fitaWidth / 2
	arrowHeigth := arrowBase / 2
	err = drawArrow(window.renderer, thickness, (x + fitaWidth/2), (y - pady), arrowBase, arrowHeigth, BLACK)
	if err != nil {
		return err
	}

	// Texto
	bufferFita := env.radio.inputToPrint
	textTextures, err := drawText(window, bufferFita, fitaWidth, fitaRec.X, fitaRec.Y, RIGHT)
	if err != nil {
		return err
	}

	window.ui = append(window.ui, textTextures...)
	return nil
}

func (env *environment) drawStacks(amount int, padx, pady int32) error {
	var err error
	for i := 0; i < amount; i++ {
		err = env.drawStack(int32(i), padx, pady)
		if err != nil {
			return err
		}
	}

	return nil
}

func (env *environment) drawStack(index, padx, pady int32) error {
	window := env.w
	machine := env.radio.activeMachine

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
	stack := machine.Stacks()[index]
	stackAparente := utils.Reserve(stack.Peek(TAMANHO_ESTRUTURAS))
	firstCharX := x
	firstCharY := y + (stackWidth * (TAMANHO_ESTRUTURAS - 1))
	textures, err := drawText(window, stackAparente, stackWidth, firstCharX, firstCharY, UP)
	if err != nil {
		return err
	}

	// Para manter uma referencia dos ponteiros que vou precisar liberar
	window.ui = append(window.ui, textures...)
	return nil
}

func drawText(window *_SDLWindow, text []string, space, x1, y1 int32, direction int) ([]*sdl.Texture, error) {
	var textSurface *sdl.Surface
	var textTexture *sdl.Texture
	var textTextures []*sdl.Texture
	var x, y int32
	var err error
	for i, s := range text {
		// Essa função checa se há a palavra no cache antes de cirar a surface
		// O cache esta armazenado na SDLWindow
		textSurface, err = window.textSurface(s, BLACK)
		if err != nil {
			return textTextures, err
		}

		textTexture, err = window.renderer.CreateTextureFromSurface(textSurface)
		if err != nil {
			return textTextures, err
		}

		_, _, fontW, fontH, err := textTexture.Query()
		if err != nil {
			return textTextures, err
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
			return nil, errors.New("direção invalida. drawManyRects()")
		}

		textTextures = append(textTextures, textTexture)
		window.renderer.Copy(textTexture, nil, &sdl.Rect{
			X: x,
			Y: y,
			W: fontW,
			H: fontH,
		})
	}

	return textTextures, nil
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
