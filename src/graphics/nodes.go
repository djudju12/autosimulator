package graphics

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type (
	position struct {
		x, y int
	}

	Node struct {
		position
		ringRect    *sdl.Rect
		radius      int
		innerRadius int
		state       string
		color       sdl.Color
		textTexture *sdl.Texture
		textRect    *sdl.Rect
	}
)

var BLACK = sdl.Color{R: 0, G: 0, B: 0, A: 255}
var WITHE = sdl.Color{R: 255, G: 255, B: 255, A: 255}

func NewNode() *Node {
	return &Node{
		position:    position{0, 0},
		radius:      0,
		innerRadius: 0,
		state:       "",
		color:       BLACK,
		textTexture: nil,
		textRect:    nil,
	}
}

func (n *Node) draw(pixels []byte, renderer *sdl.Renderer, font *ttf.Font) {
	// n.drawRing(pixels)
	// n.ringRect = &sdl.Rect{X: int32(n.x), Y: int32(n.y), W: int32(n.radius) * 2, H: int32(n.radius) * 2}
	// n.ringTexture =

	err := n.drawText(renderer, font)
	if err != nil {
		fmt.Printf("Erro ao desenhar o texto do anel: %v", err)
		os.Exit(1)
	}
}

// func (n *node) drawRing(pixels []byte) {
// 	// Draw the ring
// 	for y := -n.radius; y < n.radius; y++ {
// 		for x := -n.radius; x < n.radius; x++ {
// 			distanceSquared := x*x + y*y
// 			if distanceSquared >= (n.innerRadius*n.innerRadius) && distanceSquared < (n.radius*n.radius) {
// 				setPixel(position{int(n.x) + x, int(n.y) + y}, n.color, pixels)
// 			}
// 		}
// 	}
// }

func (n *Node) drawText(renderer *sdl.Renderer, font *ttf.Font) error {
	w, h := int32(n.radius), int32(n.radius)
	textRec := &sdl.Rect{X: int32(n.x) - (w / 2), Y: int32(n.y) - (h / 2), W: w, H: h}

	textSurface, err := font.RenderUTF8Solid(n.state, BLACK)
	if err != nil {
		return err
	}
	defer textSurface.Free()

	textTexture, err := renderer.CreateTextureFromSurface(textSurface)
	if err != nil {
		return err
	}

	n.textTexture = textTexture
	n.textRect = textRec
	return nil
}

func setPixel(pos position, c sdl.Color, pixels []byte) {
	index := (pos.y*WITDH + pos.x) * 4
	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.R
		pixels[index+1] = c.G
		pixels[index+2] = c.B
		pixels[index+3] = c.A
	}
}
