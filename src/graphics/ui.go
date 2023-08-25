package graphics

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

func (w *_SDLWindow) drawFita(env *environment, headIndex int, padx, pady int32) error {
	fita := env.radio.input
	fitaAparente := fita.Peek(8)
	fmt.Printf("Fita aparente:%s\n", fitaAparente)
	textures := []*sdl.Texture{}

	// body
	fitaTexture, err := w.renderer.CreateTextureFromSurface(w.cacheSprites[FITA])
	textures = append(textures, fitaTexture)
	if err != nil {
		return err
	}

	_, _, fitaWidth, fitaHeight, err := fitaTexture.Query()
	if err != nil {
		return err
	}

	/// head
	fitaHeadTexture, err := w.renderer.CreateTextureFromSurface(w.cacheSprites[FITA_HEAD])
	textures = append(textures, fitaHeadTexture)
	if err != nil {
		return err
	}
	_, _, headWidth, headHeigth, err := fitaHeadTexture.Query()
	if err != nil {
		return err
	}

	// Ficou muito grande a imagem inicial, dai dividi por 2
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

	w.renderer.Copy(fitaTexture, nil, fitaRec)
	w.renderer.Copy(fitaHeadTexture, nil, headRec)

	var textSurface *sdl.Surface
	var textTexture *sdl.Texture
	var textRect *sdl.Rect
	for i, symbol := range fitaAparente {
		textSurface, err = w.textSurface(symbol, BLACK)
		if err != nil {
			return err
		}

		textTexture, err = w.renderer.CreateTextureFromSurface(textSurface)
		if err != nil {
			return err
		}

		_, _, fontW, fontH, err := textTexture.Query()
		if err != nil {
			return err
		}

		textRect = &sdl.Rect{
			// (headWidth-2) pois as paredes do quadrados contam 1 pixel cara
			// TODO: Isso é uma gambiarra.
			X: (fitaRec.X + ((headWidth - 2) * int32(i))) + fontW/2,
			Y: fitaRec.Y,
			W: fontW,
			H: fontH,
		}
		textures = append(textures, textTexture)
		w.renderer.Copy(textTexture, nil, textRect)
	}

	env.w.ui = append(env.w.ui, textures...)
	return nil
}
