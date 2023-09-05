package graphics

import (
	"errors"
	"fmt"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	UP    = iota
	RIGHT = iota
)

const (
	TEXT_DOWN_ALIGN   = iota
	TEXT_RIGHT_CENTER = iota
	TEXT_UP_CENTER    = iota
)

const (
	TAMANHO_ESTRUTURAS        = 9
	DIMENSAO_ESTRUTURAS       = 32
	PADX                int32 = 5
	PADY                int32 = 5
)

func (ui *uiComponents) drawSelectBox(window *_SDLWindow) error {
	var amount int32 = 3
	if ui.indexMenu < 1 {
		ui.indexMenu = 1
	}

	if ui.indexMenu > amount {
		ui.indexMenu = amount
	}

	var widthBox int32 = DIMENSAO_ESTRUTURAS * 6
	rect := sdl.Rect{
		X: WITDH/2 - widthBox/2,
		Y: HEIGHT/2 - DIMENSAO_ESTRUTURAS*amount/2,
		W: widthBox,
		H: DIMENSAO_ESTRUTURAS,
	}

	err := drawBoxList(window, rect, amount, ui.indexMenu)
	if err != nil {
		return err
	}

	err = drawText(window, []string{"Inputs", "Machine", "Run n Inputs"}, DIMENSAO_ESTRUTURAS/2, rect.X+PADX, rect.Y+rect.H/2, TEXT_DOWN_ALIGN)
	if err != nil {
		return err
	}

	return nil
}

func (ui *uiComponents) drawFita(window *_SDLWindow) error {
	// Calculo da posicao inicial da fita/texto
	var fitaWidth, thickness int32 = DIMENSAO_ESTRUTURAS, 2
	x := window.WIDTH - (fitaWidth*TAMANHO_ESTRUTURAS + DIMENSAO_ESTRUTURAS*8 + PADX*5 + fitaWidth/4)
	y := (window.HEIGHT - fitaWidth) - PADY

	// Rec representa o primeiro quadrado da fita
	fitaRec := sdl.Rect{
		X: x,
		Y: y,
		W: fitaWidth,
		H: fitaWidth, // É um quadrado
	}

	err := drawManyRects(window.renderer, thickness, int(TAMANHO_ESTRUTURAS), RIGHT, fitaRec, COLOR_DEFAULT)
	if err != nil {
		return err
	}

	// Head da fita
	headBase := fitaWidth / 2
	headHeigth := headBase / 2
	err = drawArrowDown(window.renderer, thickness, (x + fitaWidth/2), (y - PADY), headBase, headHeigth, COLOR_DEFAULT)
	if err != nil {
		return err
	}

	// Texto
	bufferFita := ui.bufferInput
	err = drawText(window, bufferFita, fitaWidth, (x + fitaWidth/2), y, TEXT_RIGHT_CENTER)
	if err != nil {
		return err
	}

	// Para manter uma referencia dos ponteiros que vou precisar liberar
	return nil
}

func (ui *uiComponents) drawStacks(window *_SDLWindow) error {
	var err error
	a, b := ui.stackHist.get(ui.indexComputation)
	err = ui.drawStack(window, a, 1)
	if err != nil {
		return err
	}

	if b != nil {
		err = ui.drawStack(window, b, 2)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ui *uiComponents) drawStack(window *_SDLWindow, stack []string, index int32) error {

	// Calculo da posicao inicial do stack/texto
	var stackWidth, thickness int32 = DIMENSAO_ESTRUTURAS, 2
	x := window.WIDTH - (PADX+stackWidth)*index
	y := window.HEIGHT - (PADY + stackWidth*TAMANHO_ESTRUTURAS)

	// Esse rec represeta o primeiro quadrado do stack
	oneStackCointainer := sdl.Rect{
		X: x,
		Y: y,
		W: stackWidth,
		H: stackWidth, // É um quadrado
	}

	// Retangulos
	err := drawManyRects(window.renderer, thickness, int(TAMANHO_ESTRUTURAS), UP, oneStackCointainer, COLOR_DEFAULT)
	if err != nil {
		return err
	}

	// Textos
	firstCharY := y + (stackWidth * (TAMANHO_ESTRUTURAS - 1))
	err = drawText(window, stack, stackWidth, x+stackWidth/2, firstCharY, TEXT_UP_CENTER)
	if err != nil {
		return err
	}

	// Para manter uma referencia dos ponteiros que vou precisar liberar
	return nil
}

func (ui *uiComponents) drawHist(window *_SDLWindow) error {
	var amount int32 = 3
	var histWidth int32 = DIMENSAO_ESTRUTURAS * 6
	x := window.WIDTH - (DIMENSAO_ESTRUTURAS*2 + PADX*3 + histWidth)
	y := window.HEIGHT - PADY - DIMENSAO_ESTRUTURAS*amount
	rect := sdl.Rect{
		X: x,
		Y: y,
		W: histWidth,
		H: DIMENSAO_ESTRUTURAS,
	}

	err := drawBoxList(window, rect, amount, 2)
	if err != nil {
		return err
	}

	yText := y + rect.H/2
	err = ui.drawHistText(window, x, yText)
	if err != nil {
		return err
	}

	return nil
}

func (ui *uiComponents) drawHistText(window *_SDLWindow, x, y int32) error {
	index := ui.indexComputation

	var upper string = "---"
	if index < len(ui.bufferComputation.History)-1 {
		upper = ui.bufferComputation.History[index+1].Stringfy()
	}

	mid := ui.bufferComputation.History[index].Stringfy()

	var bottom string = "---"
	if index > 0 {
		bottom = ui.bufferComputation.History[index-1].Stringfy()
	}

	err := drawText(window, []string{bottom, mid, upper}, DIMENSAO_ESTRUTURAS/2, x+PADX, y, TEXT_DOWN_ALIGN)
	return err
}

func drawBoxList(window *_SDLWindow, rect sdl.Rect, amount, headPos int32) error {
	var thickness int32 = 2
	if headPos > amount {
		return fmt.Errorf("a posicao da cabeça da seta não pode ser maior que a quantidade de elementos na BoxList")
	}

	yArrow := (rect.Y - rect.H/2) + (rect.H * headPos)
	headBase := rect.H / 2
	headHeigth := headBase / 2

	err := drawArrowRight(window.renderer, thickness, rect.X-PADX, yArrow, headBase, headHeigth, COLOR_DEFAULT)
	if err != nil {
		return err
	}

	err = drawManyRects(window.renderer, thickness, int(amount), UP, rect, COLOR_DEFAULT)
	if err != nil {
		return err
	}

	return nil
}

func drawArrowRight(renderer *sdl.Renderer, thickness, x, y, base, heigth int32, color sdl.Color) error {
	var ok bool
	errText := "nao foi possível desenhar a flecha"
	ok = gfx.ThickLineColor(renderer, x, y, x-heigth, y-base/2, thickness, color)
	if !ok {
		return errors.New(errText)
	}

	ok = gfx.ThickLineColor(renderer, x, y, x-heigth, y+base/2, thickness, color)
	if !ok {
		return errors.New(errText)
	}

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
		textSurface, err = window.textSurface(s, COLOR_DEFAULT)
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
		case TEXT_UP_CENTER:
			x = x1 - fontW/2
			y = y1 - (space * int32(i))
		case TEXT_RIGHT_CENTER:
			x = x1 + (space * int32(i)) - fontW/2
			y = y1
		case TEXT_DOWN_ALIGN:
			x = x1
			y = y1 + ((space + fontH/2) * int32(i)) - fontH/2
		default:
			return errors.New("direção invalida. drawText()")
		}

		textRect := &sdl.Rect{
			X: x,
			Y: y,
			W: fontW,
			H: fontH,
		}

		window.renderer.Copy(textTexture, nil, textRect)
		if err != nil {
			return err
		}
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
		case RIGHT:
			x = rect.X + (rect.W * i)
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

func drawArrowDown(renderer *sdl.Renderer, thickness, x, y, base, heigth int32, color sdl.Color) error {
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
