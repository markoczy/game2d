package world

import (
	"fmt"
	"github.com/markoczy/game2d/display"
	"image"
	"log"
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
}

func (w *simpleWorld) Tick() {}

func (w *simpleWorld) Render(screen display.Screen) {
	// 1) Get lowest left tile
	w0 := screen.Bounds().Min
	w1 := screen.Bounds().Max
	th := w.tileMap.TileHeight
	tw := w.tileMap.TileWidth
	tile0 := image.Point{w0.X / tw, w0.Y / th}
	tile1 := image.Point{w1.X / tw, w1.Y / th}
	// fmt.Printf("th: %d tw: %d\n", th, tw)
	// fmt.Printf("w0: %v tile0: %v\n", w0, tile0)
	// fmt.Printf("w1: %v tile1: %v\n", w1, tile1)
	cnt := 0
	for x := tile0.X; x <= tile1.X; x++ {
		for y := tile0.Y; y <= tile1.Y; y++ {
			// fmt.Printf("x: %d y: %d\n", x, y)
			img, dim, err := w.getRenderedTile(x, y)
			if err != nil {
				fmt.Println("Error:", err)
			}
			if img != nil {
				screen.RenderW(img, image.Point{x * tw, y * th}, dim)
				cnt++
			}
		}
	}
	log.Printf("World: Rendered %d tiles\n", cnt)
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
	for _, img := range rawTiles {
		cur := newRenderedTile(img, scale)
		tiles = append(tiles, cur)
	}
	return &simpleWorld{
		tileMap:   tileMap,
		materials: materials,
		tiles:     tiles}, nil
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
