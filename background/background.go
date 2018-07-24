package background

import (
	"image"
	"image/draw"
	"log"

	"github.com/markoczy/game2d/display"
)

type singleTiledBG struct {
	tile       *image.RGBA
	tileDim    image.Point
	tileOffset image.Point
	worldDim   image.Point
	worldImg   *image.RGBA
	scale      int
}

func (bg *singleTiledBG) Render(screen display.Screen) {
	// render preloaded image
	screen.RenderElement(bg.worldImg, image.ZP, bg.worldDim)
}

func (bg *singleTiledBG) initWorld() error {
	// if the offset is not 0, begin must be calculated:
	//
	// tile:       +-------+
	// screen: +------------------
	//         /---/------/
	//          ofx   tw   ->  beginx = (ofx % tw) - tw
	//
	tw := bg.tileDim.X
	th := bg.tileDim.Y
	ww := bg.worldDim.X
	wh := bg.worldDim.Y
	scale := bg.scale
	img := bg.tile

	start := image.Point{0, 0}
	modX := mod(bg.tileOffset.X, tw)
	modY := mod(bg.tileOffset.Y, th)
	if modX > 0 {
		start.X = modX - tw
	}
	if modY > 0 {
		start.Y = modY - th
	}

	cnt := 0
	worldImg := image.NewRGBA(image.Rect(0, 0, bg.worldDim.X*bg.scale, bg.worldDim.Y*bg.scale))
	for x := start.X; x <= ww; x += tw {
		for y := start.Y; y <= wh; y += th {
			pos := image.Point{
				x * scale, (wh - y) * scale}

			draw.Draw(worldImg, img.Bounds().Add(pos), img, image.ZP, draw.Over)
			cnt++
		}
	}
	log.Printf("BG: Rendered %d Tiles\n", cnt)
	bg.worldImg = worldImg
	return nil
}

func NewSingleTiledBG(img *image.RGBA, tileDim image.Point, offset image.Point, worldDim image.Point, scale int) (display.Visual, error) {
	ret := singleTiledBG{tile: img, tileDim: tileDim, tileOffset: offset, worldDim: worldDim, scale: scale}
	err := ret.initWorld()
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// strict positive modulo
func mod(a, b int) int {
	ret := a % b
	if ret < 0 {
		ret += b
	}
	return ret
}
