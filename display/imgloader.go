package display

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

func LoadImage(path string) (*image.RGBA, image.Point, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, image.ZP, err
	}
	reader := bufio.NewReader(file)
	img, err := png.Decode(reader)
	if err != nil {
		return nil, image.ZP, err
	}
	if img.ColorModel() != color.NRGBAModel {
		return nil, image.ZP, fmt.Errorf("Bad Color model: %c", img.ColorModel())
	}

	nrgba := img.(*image.NRGBA)
	rgba := image.NewRGBA(nrgba.Bounds())
	for x := 0; x <= nrgba.Bounds().Max.X; x++ {
		for y := 0; y <= nrgba.Bounds().Max.Y; y++ {
			rgba.Set(x, y, nrgba.NRGBAAt(x, y))
			// fmt.Printf("nrgba: %v, rgba: %v\n", nrgba.NRGBAAt(x, y), rgba.RGBAAt(x, y))
		}
	}

	dim := image.Point{img.Bounds().Dx(), img.Bounds().Dy()}
	return rgba, dim, nil
}

func ScaleImage(img *image.RGBA, scale int) *image.RGBA {
	ret := image.NewRGBA(image.Rectangle{img.Bounds().Min, img.Bounds().Max.Mul(scale)})
	for x := 0; x <= img.Bounds().Max.X; x++ {
		for y := 0; y <= img.Bounds().Max.Y; y++ {
			col := img.RGBAAt(x, y)
			x0 := x * scale
			y0 := y * scale
			x1 := x0 + scale
			y1 := y0 + scale
			draw.Draw(ret, image.Rect(x0, y0, x1, y1), image.NewUniform(col), image.ZP, draw.Over)
		}
	}
	return ret

}

// func (d *scaledDisplay) SetPixelAt(x, y int, col color.RGBA) {
// 	x0 := x * d.scale
// 	y0 := y * d.scale
// 	x1 := x0 + d.scale
// 	y1 := y0 + d.scale
// 	draw.Draw(d.rgba, image.Rect(x0, y0, x1, y1), image.NewUniform(col), image.ZP, draw.Src)
// }
