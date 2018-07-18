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
	Pos() image.Point
	SetPos(image.Point)
}

type uniformSprite struct {
	id     int
	img    *image.RGBA
	pos    image.Point
	ticker func(Entity)
}

func (s *uniformSprite) Tick() {
	if s.ticker != nil {
		s.ticker(s)
	}
}

func (s *uniformSprite) Render(screen display.Screen) {
	pos := getEntScreenPos(s, screen)
	screen.Render(s.img, pos)
}

func (s *uniformSprite) Pos() image.Point {
	return s.pos
}

func (s *uniformSprite) SetPos(pos image.Point) {
	s.pos = pos
}

func (s *uniformSprite) Id() int {
	return s.id
}

func (s *uniformSprite) Type() int {
	return UniformSpriteType
}

func NewUniformSprite(id int, img *image.RGBA, pos image.Point, ticker func(Entity)) Entity {
	return &uniformSprite{
		id:     id,
		img:    img,
		pos:    pos,
		ticker: ticker}
}

func getEntScreenPos(e Entity, s display.Screen) image.Point {
	// e0' := (e0 - s0) * scale
	return e.Pos().Sub(s.Bounds().Min).Mul(s.Scale())
}
