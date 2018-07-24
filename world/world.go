package world

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/markoczy/game2d/display"
	// "log"
)

type tileDisplay struct {
	tile0   *rotatedTile
	tile90  *rotatedTile
	tile180 *rotatedTile
	tile270 *rotatedTile
	dim     image.Point
}

type rotatedTile struct {
	noFlip *image.RGBA
	flipX  *image.RGBA
	flipY  *image.RGBA
}

type World interface {
	Tick()
	Render(screen display.Screen)
}

type simpleWorld struct {
	tileMap   *Tilemap
	materials *Materialmap
	tiles     []tileDisplay
	worldImg  *image.RGBA
}

func (w *simpleWorld) Tick() {}

func (w *simpleWorld) Render(screen display.Screen) {
	// TODO any more inefficient? ;-)
	th := w.tileMap.TileHeight
	tw := w.tileMap.TileWidth
	tx := w.tileMap.TilesWide
	ty := w.tileMap.TilesHigh
	dim := image.Point{tx * tw, ty * th}
	// fmt.Printf("bounds: %v, screen: %v, dim: %v\n", w.worldImg.Bounds(), screen.Bounds(), dim)
	//

	// fmt.Println("w.worldImg", w.worldImg)
	// screen.RenderFull(w.worldImg)
	screen.RenderElement(w.worldImg, image.ZP, dim)
}

func (w *simpleWorld) getRenderedTile(x, y int) (*image.RGBA, image.Point, error) {
	if y > w.tileMap.TilesHigh-1 || x > w.tileMap.TilesWide-1 {
		return nil, image.ZP, nil
	}
	yoffset := w.tileMap.TilesWide * (w.tileMap.TilesHigh - y - 1)
	// fmt.Println("Resulting id:", yoffset+x)
	if yoffset+x > len(w.tileMap.Tiles)-1 {
		return nil, image.ZP, fmt.Errorf("Rendered tile at (%d; %d) not found, offset was: %d", x, y, yoffset+x)
	}
	tile := w.tileMap.Tiles[yoffset+x]
	if tile == nil {
		return nil, image.ZP, nil
	}
	// tile := w.tiles[yoffset+x-1]

	// verify x and y
	if tile.X != x || tile.Y != (w.tileMap.TilesHigh-y-1) {
		fmt.Printf("Warning: Wrong x,y for tile x: %d y: %d id: %d, tile: %v\n", x, y, yoffset+x, tile)
		for _, t := range w.tileMap.Tiles {
			if t.X == x && t.Y == y {
				tile = t
			}
		}
	}

	tileID := tile.Tile
	// fmt.Println("Tile id:", tileID)
	if tileID < 0 {
		return nil, image.ZP, nil
	}
	if tileID > len(w.tiles)-1 {
		return nil, image.ZP, fmt.Errorf("Tile with id %d not found", tileID)
	}
	// fmt.Printf("x: %d, y: %d, tileId: %d\n", x, y, tileID)
	return w.tiles[tileID].tile0.noFlip, w.tiles[tileID].dim, nil

}

func NewSimpleWorld(tileMap *Tilemap, materials *Materialmap, rawTiles []*image.RGBA, scale int) (World, error) {
	tiles := []tileDisplay{}
	for _, rawtile := range rawTiles {
		cur := newRenderedTile(rawtile, scale)
		tiles = append(tiles, cur)
	}
	th := tileMap.TileHeight
	tw := tileMap.TileWidth
	tx := tileMap.TilesWide
	ty := tileMap.TilesHigh

	w := simpleWorld{
		tileMap:   tileMap,
		materials: materials,
		tiles:     tiles}

	worldImg := image.NewRGBA(image.Rect(0, 0, tx*tw*scale, ty*th*scale))
	cnt := 0
	for x := 0; x <= tx; x++ {
		for y := 0; y <= ty; y++ {
			// fmt.Printf("x: %d y: %d\n", x, y)
			img, _, err := w.getRenderedTile(x, y)
			if err != nil {
				fmt.Println("Error:", err)
			}
			if img != nil {
				// screen.RenderW(img, image.Point{x * tw, y * th}, dim)
				pos := image.Point{
					x * tw * scale, (ty - y) * th * scale}

				draw.Draw(worldImg, img.Bounds().Add(pos), img, image.ZP, draw.Over)
				cnt++
			}
		}
	}
	fmt.Printf("World: Rendered %d Tiles\n", cnt)

	w.worldImg = worldImg
	for x := 0; x < w.worldImg.Bounds().Max.X; x++ {
		for y := 0; y < w.worldImg.Bounds().Max.Y; y++ {
			if x == y {
				w.worldImg.SetRGBA(x, y, color.RGBA{255, 0, 0, 255})
			}
		}
	}
	// fmt.Println("img", worldImg)

	return &w, nil
}

func newRenderedTile(img *image.RGBA, scale int) tileDisplay {
	scaled := img
	dim := img.Bounds().Max.Sub(img.Bounds().Min)
	if scale > 1 {
		scaled = display.ScaleImage(img, scale)
	}
	return tileDisplay{
		dim: dim,
		tile0: &rotatedTile{
			noFlip: scaled}}
}
