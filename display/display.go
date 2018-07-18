package display

import (
	"golang.org/x/exp/shiny/screen"
	"image"
	// "image/color"
	"image/draw"
)

type Screen interface {
	InitBuffer() error
	UploadBuffer()
	// World related
	Bounds() image.Rectangle
	SetPos(image.Point)
	Scale() int // or somewhere else...
	Render(img *image.RGBA, pos image.Point)
}

func NewScreen(window screen.Window, screen screen.Screen, width, height, scale int) Screen {
	return &defaultScreen{
		window: window,
		screen: screen,
		width:  width,
		height: height,
		scale:  scale}
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
}

func (d *defaultScreen) InitBuffer() error {
	buf, err := d.screen.NewBuffer(image.Point{d.width, d.height})
	if err != nil {
		return err
	}
	d.buffer = buf
	d.rgba = buf.RGBA()
	return nil
}

func (d *defaultScreen) UploadBuffer() {
	d.window.Upload(image.ZP, d.buffer, d.buffer.Bounds())
	d.buffer.Release()
}

func (d *defaultScreen) SetPos(pos image.Point) {
	d.bounds = image.Rectangle{
		Min: pos,
		Max: image.Point{
			X: pos.X + d.bounds.Dx(),
			Y: pos.Y + d.bounds.Dy()}}
}

func (d *defaultScreen) Bounds() image.Rectangle {
	return d.bounds
}

func (d *defaultScreen) Scale() int {
	return d.scale
}

// AddImage pos is world-related
func (d *defaultScreen) Render(img *image.RGBA, pos image.Point) {
	draw.Draw(d.rgba, img.Bounds().Add(pos), img, image.ZP, draw.Src)
}
