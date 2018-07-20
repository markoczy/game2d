package background

import (
	"github.com/markoczy/game2d/display"
	"image"
	"log"
)

type singleTiledBG struct {
	img    *image.RGBA
	dim    image.Point
	offset image.Point
}

func (bg *singleTiledBG) Render(screen display.Screen) {
	//
	// p := screen.Bounds().Min.Add(bg.offset)
	// p := screen.Bounds().Min.Add(bg.offset)
	// startx := p.X % bg.dim.X
	// if startx > 0 {
	// 	startx -= bg.dim.X
	// }
	// starty := p.Y % bg.dim.Y
	// if starty > 0 {
	// 	starty -= bg.dim.Y
	// }

	w0 := screen.Bounds().Min
	// startx := (bg.dim.X - w0.X%bg.dim.X) + bg.offset.X
	// starty := (bg.dim.Y - w0.X%bg.dim.Y) + bg.offset.Y
	startx := (w0.X - w0.X%bg.dim.X) + bg.offset.X
	starty := (w0.Y - w0.X%bg.dim.Y) + bg.offset.Y

	max := screen.Bounds().Max
	log.Printf("*** startx: %v, starty: %v, max: %v\n", startx, starty, max)
	// screen.RenderW(bg.img, image.Point{startx, starty}, bg.dim)
	for x := startx; x <= max.X; x += bg.dim.X {
		for y := starty; y <= max.Y; y += bg.dim.Y {
			// screen.Re
			// log.Printf("x: %v, y: %v, dim: %v\n", x, y, bg.dim)
			screen.RenderW(bg.img, image.Point{x, y}, bg.dim)
		}
	}
}

func NewSingleTiledBG(img *image.RGBA, dim image.Point, offset image.Point) display.Visual {
	return &singleTiledBG{img, dim, offset}
}
