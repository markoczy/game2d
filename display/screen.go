package display

import (
	"fmt"
	"golang.org/x/exp/shiny/screen"
	"image"
	// "image/color"
	"image/draw"
	// "log"
	"sync"
)

type Screen interface {
	InitBuffer() error
	UploadBuffer()
	// World related
	Bounds() image.Rectangle
	SetPos(image.Point)
	Scale() int // or somewhere else...
	RenderDirect(img *image.RGBA, pos image.Point)
	RenderElement(img *image.RGBA, posw image.Point, dimw image.Point)
	// RenderFull(img *image.RGBA)
}

func NewScreen(window screen.Window, screen screen.Screen, width, height, scale int) Screen {
	return &defaultScreen{
		window: window,
		screen: screen,
		width:  width,
		height: height,
		scale:  scale,
		bounds: image.Rect(0, 0, width/scale, height/scale),
		mut:    &sync.Mutex{}}
}

type defaultScreen struct {
	// Screen related dims
	width, height, scale int
	// World related dims
	bounds image.Rectangle
	// Buffering
	window screen.Window
	screen screen.Screen
	buffer screen.Buffer
	rgba   *image.RGBA
	mut    *sync.Mutex
}

func (d *defaultScreen) InitBuffer() error {
	d.mut.Lock()
	fmt.Println("InitBuffer")
	buf, err := d.screen.NewBuffer(image.Point{d.width, d.height})
	if err != nil {
		return err
	}
	d.buffer = buf
	d.rgba = buf.RGBA()
	d.mut.Unlock()
	return nil
}

func (d *defaultScreen) UploadBuffer() {
	// go func() {
	d.mut.Lock()
	fmt.Println("UploadBuffer")
	// if d.buffer != nil {
	d.window.Upload(image.ZP, d.buffer, d.buffer.Bounds())
	// d.buffer.Release()
	// d.buffer = nil
	// d.rgba = nil
	// }
	d.mut.Unlock()
	// }()
}

func (d *defaultScreen) SetPos(pos image.Point) {
	d.mut.Lock()
	d.bounds = image.Rectangle{
		Min: pos,
		Max: image.Point{
			X: pos.X + d.bounds.Dx(),
			Y: pos.Y + d.bounds.Dy()}}
	d.mut.Unlock()
}

func (d *defaultScreen) Bounds() image.Rectangle {
	d.mut.Lock()
	ret := d.bounds
	d.mut.Unlock()
	return ret
}

func (d *defaultScreen) Scale() int {
	return d.scale
}

// Render pos is screen-related
func (d *defaultScreen) RenderDirect(img *image.RGBA, pos image.Point) {
	d.mut.Lock()
	draw.Draw(d.rgba, img.Bounds().Add(pos), img, image.ZP, draw.Over)
	d.mut.Unlock()
}

// RenderW pos is world-related
func (d *defaultScreen) RenderElement(img *image.RGBA, pos image.Point, dim image.Point) {
	d.mut.Lock()
	if d.rgba != nil {
		posS := d.worldCoordsToScreen(pos, dim.Y)
		// fmt.Printf("img.Bounds().Add(posS): %v\n", img.Bounds().Add(posS))
		draw.Draw(d.rgba, img.Bounds().Add(posS), img, image.ZP, draw.Over)
	}
	d.mut.Unlock()
}

// // RenderW pos is world-related
// func (d *defaultScreen) RenderFull(img *image.RGBA) {
// 	rect := d.bounds //image.Rectangle{image.ZP, d.bounds.Max.Div(d.scale)}
// 	xNew := (-rect.Min.X) * d.scale
// 	yNew := ((rect.Min.Y + rect.Dy()) - 320) * d.scale //(rect.Min.Y + rect.Dy()) * d.scale
// 	pos := image.Point{xNew, yNew}
// 	fmt.Printf("bounds: %v, rpos: %v\n", d.bounds, pos)
// 	draw.Draw(d.rgba, img.Bounds().Add(pos), img, image.ZP, draw.Over)
// }

func (d *defaultScreen) worldCoordsToScreen(pos image.Point, height int) image.Point {
	rect := d.bounds
	xNew := (pos.X - rect.Min.X) * d.scale
	// fmt.Printf("rec: %v, rect.Dy(): %v, pos.Y: %d, height: %v\n", rect, rect.Dy(), pos.Y, height)
	yNew := ((rect.Min.Y + rect.Dy()) - (pos.Y + height)) * d.scale
	// log.Printf("rect: %v, pos: %v, xNew: %v, yNew: %v\n", rect, pos, xNew, yNew)
	return image.Point{xNew, yNew}
}
