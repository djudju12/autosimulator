package main

import (
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const WITDH, HEIGTH = 800, 600

type (
	color struct {
		R, G, B byte
	}

	pos struct {
		x, y int
	}

	node struct {
		pos
		radius      int
		outerRadius int
		innerRadius int
		color       color
		// xv float32
		// yv float32
	}
)

var (
	window   *sdl.Window
	renderer *sdl.Renderer
	texture  *sdl.Texture
	font     *ttf.Font
)

func setPixel(x, y int, c color, pixels []byte) {
	index := (y*WITDH + x) * 4

	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.R
		pixels[index+1] = c.G
		pixels[index+2] = c.B
	}
}

func (n *node) draw(pixels []byte) {
	for y := -n.radius; y < n.radius; y++ {
		for x := -n.radius; x < n.radius; x++ {
			distanceSquared := x*x + y*y
			if distanceSquared >= (n.innerRadius*n.innerRadius) && distanceSquared < (n.outerRadius*n.outerRadius) {
				setPixel(int(n.x)+x, int(n.y)+y, n.color, pixels)
			}
		}
	}
}

func init() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}

	err = ttf.Init()
	if err != nil {
		panic(err)
	}

	window, err = sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WITDH, HEIGTH, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(1)
	}

	texture, err = renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, WITDH, HEIGTH)
	if err != nil {
		panic(1)
	}

	font, err = ttf.OpenFont("/home/jonathan/programacao/autosimulator/src/graphics/assets/IBMPlexMono-ExtraLight.ttf", 24)
	if err != nil {
		panic(1)
	}
}

func shutDown() {
	sdl.Quit()
	ttf.Quit()
	window.Destroy()
	texture.Destroy()
	renderer.Destroy()
	font.Close()
}

func main() {
	defer shutDown()

	pixels := make([]byte, WITDH*HEIGTH*4)
	for y := 0; y < HEIGTH; y++ {
		for x := 0; x < WITDH; x++ {
			setPixel(x, y, color{255, 0, 0}, pixels)
		}
	}

	node := &node{pos{300, 300}, 100, 100, 80, color{255, 255, 255}}
	node.draw(pixels)
	texture.Update(nil, unsafe.Pointer(&pixels[0]), WITDH*4)

	textSurface, err := font.RenderUTF8Solid("Q0", sdl.Color{R: 0, G: 0, B: 0, A: 255})
	if err != nil {
		panic(1)
	}
	defer textSurface.Free()

	textTexture, err := renderer.CreateTextureFromSurface(textSurface)
	if err != nil {
		panic(1)
	}
	defer textTexture.Destroy()

	// textWidth, textHeigth, _ := font.SizeUTF8("Q0")
	// textX := 20 - textWidth/2
	// textY := 20 - textHeigth/2

	// renderer.Copy(textTexture, nil, nil)
	// renderer.Copy(texture, nil, nil)
	// renderer.Present()
	// sdl.Delay(2000)

	// running := true
	// for running {
	// 	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
	// 		switch event.(type) {
	// 		case *sdl.QuitEvent:
	// 			println("Quit")
	// 			running = false
	// 		}
	// 	}
	// }
}
