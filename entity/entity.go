package entity

import (
	"github.com/markoczy/game2d/display"
	"image"
)

const UniformSpriteType = 0

type Entity interface {
	// Game mechanics
	Tick()
	Render(screen display.Screen)
	// Entity id, not unique
	Id() int
	Type() int
	// World related
	Pos() image.Point
	Dim() image.Point
	SetPos(image.Point)
}

type uniformSprite struct {
	id     int
	img    *image.RGBA
	pos    image.Point
	dim    image.Point
	ticker func(Entity)
}

func (s *uniformSprite) Tick() {
	if s.ticker != nil {
		s.ticker(s)
	}
}

func (s *uniformSprite) Render(screen display.Screen) {
	// pos := getEntScreenPos(s, screen)
	// pos := screen.WorldCoordsToScreen(s.pos, s.dim.Y)
	screen.RenderElement(s.img, s.pos, s.dim)
}

func (s *uniformSprite) Pos() image.Point {
	return s.pos
}

func (s *uniformSprite) SetPos(pos image.Point) {
	s.pos = pos
}

func (s *uniformSprite) Dim() image.Point {
	return s.dim
}

func (s *uniformSprite) Id() int {
	return s.id
}

func (s *uniformSprite) Type() int {
	return UniformSpriteType
}

func NewUniformSprite(id int, img *image.RGBA, pos image.Point, dim image.Point, ticker func(Entity)) Entity {
	return &uniformSprite{
		id:     id,
		img:    img,
		pos:    pos,
		dim:    dim,
		ticker: ticker}
}

// func getEntScreenPos(e Entity, s display.Screen) image.Point {
// 	// Old impl: Top down
// 	// e0' := (e0 - s0) * scale
// 	// return e.Pos().Sub(s.Bounds().Min).Mul(s.Scale())

// 	// New impl:
// 	// e0'.x := (e0.x - s0.x) * scale
// 	// e0'.y := (s0.y - e0.y - he) * scale
// 	rect := s.Bounds()
// 	pos := e.Pos()
// 	scale := s.Scale()
// 	xNew := (pos.X - rect.Min.X) * scale
// 	yNew := ((rect.Min.Y + rect.Dy()) - (pos.Y + e.Dim().Y)) * scale
// 	log.Printf("rect: %v, pos: %v, xNew: %v, yNew: %v\n", rect, pos, xNew, yNew)

// 	return image.Point{xNew, yNew}
// }
