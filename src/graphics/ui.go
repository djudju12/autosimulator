package graphics

import (
	"github.com/veandco/go-sdl2/sdl"
)

const TAMANHO_ESTRUTURAS = 9

func (env *environment) drawFita(headIndex int, padx, pady int32) error {
	fitaAparente := env.radio.inputToPrint
	window := env.w

	textures := []*sdl.Texture{}

	// body
	fitaTexture, err := window.renderer.CreateTextureFromSurface(window.cacheSprites[FITA])
	textures = append(textures, fitaTexture)
	if err != nil {
		return err
	}

	_, _, fitaWidth, fitaHeight, err := fitaTexture.Query()
	if err != nil {
		return err
	}

	/// head
	fitaHeadTexture, err := window.renderer.CreateTextureFromSurface(window.cacheSprites[FITA_HEAD])
	textures = append(textures, fitaHeadTexture)
	if err != nil {
		return err
	}
	_, _, headWidth, headHeigth, err := fitaHeadTexture.Query()
	if err != nil {
		return err
	}

	// TODO: Ajeitar o tamanho das imagens
	headWidth /= 2
	headHeigth /= 2
	fitaWidth /= 2
	fitaHeight /= 2

	fitaRec := &sdl.Rect{
		X: 0 + padx,
		Y: (HEIGTH - fitaHeight) - pady,
		W: fitaWidth,
		H: fitaHeight,
	}

	headRec := &sdl.Rect{
		X: fitaRec.X,
		// TODO: a mesma gambiarra do anterior (que na verdade é o proximo)
		Y: fitaRec.Y - (headHeigth - 2) - pady,
		W: headWidth,
		H: headHeigth,
	}

	window.renderer.Copy(fitaTexture, nil, fitaRec)
	window.renderer.Copy(fitaHeadTexture, nil, headRec)

	var textSurface *sdl.Surface
	var textTexture *sdl.Texture
	var textRect *sdl.Rect
	for i, symbol := range fitaAparente {
		textSurface, err = window.textSurface(symbol, BLACK)
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

		textRect = &sdl.Rect{
			// (headWidth-2) pois as paredes do quadrados contam 1 pixel cada
			// TODO: Isso é uma gambiarra.
			X: fitaRec.X + fontW/2 + ((headWidth - 2) * int32(i)),
			Y: fitaRec.Y,
			W: fontW,
			H: fontH,
		}

		textures = append(textures, textTexture)
		window.renderer.Copy(textTexture, nil, textRect)
	}

	window.ui = append(window.ui, textures...)
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
	var textures []*sdl.Texture

	stackTexture, err := window.renderer.CreateTextureFromSurface(window.cacheSprites[STACK])
	if err != nil {
		return err
	}

	textures = append(textures, stackTexture)

	_, _, stackWidth, stackHeigth, err := stackTexture.Query()
	if err != nil {
		return err
	}

	stackWidth /= 2
	stackHeigth /= 2

	stackRect := &sdl.Rect{
		X: WITDH - (padx+stackWidth)*(index+1),
		Y: HEIGTH - (pady + stackHeigth),
		W: stackWidth,
		H: stackHeigth,
	}

	stack := machine.Stacks()[index]
	stackAparente := stack.Peek(TAMANHO_ESTRUTURAS)
	var textSurface *sdl.Surface
	var textTexture *sdl.Texture
	var textRect *sdl.Rect
	for i, s := range stackAparente {

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

		textRect = &sdl.Rect{
			// (headWidth-2) pois as paredes do quadrados contam 1 pixel cada
			// TODO: Isso é uma gambiarra.
			X: stackRect.X + fontW/2,
			Y: stackRect.Y + stackHeigth - ((stackWidth - 2) * int32((len(stackAparente) - i))),
			W: fontW,
			H: fontH,
		}

		textures = append(textures, textTexture)
		window.renderer.Copy(textTexture, nil, textRect)
	}

	window.renderer.Copy(stackTexture, nil, stackRect)
	window.ui = append(window.ui, textures...)
	return nil
}
