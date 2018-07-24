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
	if startx > 0 {
		startx -= bg.dim.X
	}
	starty := (w0.Y - w0.X%bg.dim.Y) + bg.offset.Y
	if starty > 0 {
		starty -= bg.dim.Y
	}

	max := screen.Bounds().Max
	log.Printf("*** startx: %v, starty: %v, max: %v\n", startx, starty, max)
	// screen.RenderW(bg.img, image.Point{startx, starty}, bg.dim)
	cnt := 0
	for x := startx; x <= max.X; x += bg.dim.X {
		for y := starty; y <= max.Y; y += bg.dim.Y {
			// screen.Re
			// log.Printf("x: %v, y: %v, dim: %v\n", x, y, bg.dim)
			screen.RenderElement(bg.img, image.Point{x, y}, bg.dim)
			cnt++
		}
	}
	log.Printf("BG: Rendered %d Tiles\n", cnt)

}

func NewSingleTiledBG(img *image.RGBA, dim image.Point, offset image.Point) display.Visual {
	// worldImg := image.NewRGBA(image.Rect(0, 0, tx*tw*scale, ty*th*scale))
	// cnt := 0
	// for x := 0; x <= tx; x++ {
	// 	for y := 0; y <= ty; y++ {
	// 		// fmt.Printf("x: %d y: %d\n", x, y)
	// 		img, _, err := w.getRenderedTile(x, y)
	// 		if err != nil {
	// 			fmt.Println("Error:", err)
	// 		}
	// 		if img != nil {
	// 			// screen.RenderW(img, image.Point{x * tw, y * th}, dim)
	// 			pos := image.Point{
	// 				x * tw * scale, (ty - y) * th * scale}

	// 			draw.Draw(worldImg, img.Bounds().Add(pos), img, image.ZP, draw.Over)
	// 			cnt++
	// 		}
	// 	}
	// }

	return &singleTiledBG{img, dim, offset}
}
