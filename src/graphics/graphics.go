package graphics

import (
	"fmt"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const WITDH, HEIGTH = 800, 600

var (
	window   *sdl.Window
	renderer *sdl.Renderer
	texture  *sdl.Texture
	font     *ttf.Font
)

func init() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
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

	err = ttf.Init()
	if err != nil {
		panic(err)
	}

	font, err = ttf.OpenFont("/home/jonathan/programacao/autosimulator/src/graphics/assets/IBMPlexMono-ExtraLight.ttf", 12)
	if err != nil {
		panic(1)
	}
}

func Run() {
	pixels := make([]byte, WITDH*HEIGTH*4)
	for y := 0; y < HEIGTH; y++ {
		for x := 0; x < WITDH; x++ {
			setPixel(position{x: x, y: y}, WITHE, pixels)
		}
	}

	nodex1, nodey1 := WITDH/2, HEIGTH/2
	nodex2, nodey2 := WITDH/4, HEIGTH/4

	node := NewNode()
	node2 := NewNode()

	nodes := []*Node{node, node2}

	node.x, node.y = nodex1/2, nodey1/2
	node.radius = 40
	node.innerRadius = 35
	node.state = "Q0"
	defer node.textTexture.Destroy()

	node2.x, node2.y = nodex2/2, nodey2/2
	node2.radius = 40
	node2.innerRadius = 35
	node2.state = "Q1"
	defer node2.textTexture.Destroy()

	running := true
	mousePos := sdl.Point{X: 0, Y: 0}
	clickOffset := sdl.Point{X: 0, Y: 0}
	var selectedNode *Node
	leftMouseButtonDown := false
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false

			case *sdl.MouseMotionEvent:
				mousePos.X, mousePos.Y, _ = sdl.GetMouseState()

				if leftMouseButtonDown && selectedNode != nil {
					fmt.Print("moving..")
					selectedNode.x = int(mousePos.X - clickOffset.X)
					selectedNode.y = int(mousePos.Y - clickOffset.Y)
				}

			case *sdl.MouseButtonEvent:
				if event.(*sdl.MouseButtonEvent).Button == sdl.BUTTON_LEFT {
					if leftMouseButtonDown &&
						event.(*sdl.MouseButtonEvent).Type == sdl.MOUSEBUTTONUP {
						fmt.Print("up..")
						leftMouseButtonDown = false
						selectedNode = nil
					}

					if !leftMouseButtonDown &&
						event.(*sdl.MouseButtonEvent).Type == sdl.MOUSEBUTTONDOWN {
						leftMouseButtonDown = true
						fmt.Print("down..")

						for _, node := range nodes {
							if mousePos.InRect(node.textRect) {
								selectedNode = node
								clickOffset.X = mousePos.X - node.textRect.X
								clickOffset.Y = mousePos.Y - node.textRect.Y
								break
							}
						}

					}

				}

			}

		}

		node.draw(pixels, renderer, font)
		node2.draw(pixels, renderer, font)

		texture.Update(nil, unsafe.Pointer(&pixels[0]), WITDH*4)
		renderer.Copy(texture, nil, nil)

		renderer.Copy(node.textTexture, nil, node.textRect)
		renderer.Copy(node2.textTexture, nil, node2.textRect)

		renderer.Present()

		sdl.Delay(1000 / 60)
	}
	defer shutDown()
}

func shutDown() {
	sdl.Quit()
	ttf.Quit()
	window.Destroy()
	texture.Destroy()
	renderer.Destroy()
	font.Close()
}
