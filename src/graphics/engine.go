package graphics

type (
	Engine interface {
		Init()
		Destroy()
		Renderer()
		PollEvents()
		GetMouseState()
	}

	Renderer interface {
		Clear(color Color)
		DrawCircle(x, y, radius int32, color Color) error
		FilledCircleColor(x, y, radius int32, color Color) error
		ThickLineColor(x1, y1, x2, y2, width int32, color Color) error
		RectangleColor(x1, y1, x2, y2 int32, color Color) error
	}
)

type (
	Color struct {
		R, G, B, A int32
	}

	Point struct {
		X, Y int32
	}

	Rect struct {
		X, Y, W, H int32
	}
)
