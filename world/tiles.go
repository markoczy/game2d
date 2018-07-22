package world

import (
	"fmt"
	"github.com/markoczy/game2d/display"
	"image"
	"path"
)

func LoadRawTiles(folder, prefix string) ([]*image.RGBA, error) {
	// TODO error handling
	ret := []*image.RGBA{}
	cur := 0
	var err error
	for err == nil {
		var img *image.RGBA
		img, _, err = display.LoadImage(path.Join(folder, fmt.Sprintf("%s_%d.png", prefix, cur)))
		cur++
		if err == nil {
			ret = append(ret, img)
		}
	}
	return ret, nil
}
