package game

import (
	"fmt"
	// "fmt"
	// "github.com/markoczy/game2d/background"
	"image"
	"log"
	"math/rand"
	"time"

	"github.com/markoczy/game2d/background"
	"github.com/markoczy/game2d/display"
	"github.com/markoczy/game2d/entity"
	"github.com/markoczy/game2d/world"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
)

var (
	up      = false
	down    = false
	left    = false
	right   = false
	screenX = 0
	screenY = 0
)

func NewGame(width, height, scale, tick int) Game {
	return &game{width: width, height: height, scale: scale, ttick: tick}
}

type Game interface {
	Run() error
}

type game struct {
	width, height int
	scale, ttick  int
	screen        display.Screen
	bg            display.Visual
	world         world.World
	entities      []entity.Entity
}

func (g *game) Run() error {
	img, dim, err := display.LoadImage("./res/testbg.png")
	if err != nil {
		return err
	}
	img = display.ScaleImage(img, g.scale)
	g.bg, err = background.NewSingleTiledBG(img, dim, image.Point{1, 1}, image.Point{200, 200}, g.scale)
	if err != nil {
		return err
	}

	img, dim, err = display.LoadImage("./res/sample0.png")
	if err != nil {
		return err
	}
	img = display.ScaleImage(img, g.scale)
	g.addTestEntity(img, image.Point{0, 0}, dim, image.Point{0, 0})
	g.addTestEntity(img, image.Point{50, 50}, dim, image.Point{-1, 1})
	g.addTestEntity(img, image.Point{100, 100}, dim, image.Point{1, -1})
	g.addTestEntity(img, image.Point{200, 100}, dim, image.Point{-1, -1})

	rand.Seed(time.Now().UTC().UnixNano())

	chErr := make(chan error)
	chInterrupt := make(chan bool)
	// bounds := image.Point{width, height}
	interrupt := false

	// Start Application and Game Thread
	go driver.Main(func(s screen.Screen) {
		// Init Window
		w, err := s.NewWindow(&screen.NewWindowOptions{
			Width:  g.width,
			Height: g.height,
			Title:  "Application Window"})
		if err != nil {
			chErr <- err
			return
		}

		g.screen = display.NewScreen(w, s, g.width, g.height, g.scale)
		// g.screen.SetPos(image.Point{-100, -100})

		tilemap, err := world.LoadTilemap("./res/testtiles/testtiles.json", 0)
		if err != nil {
			chErr <- err
			return
		}
		// fmt.Println("tilemap", tilemap)
		materialmap, err := world.LoadMaterialmap("./res/testtiles/testtiles_mat.json")
		if err != nil {
			chErr <- err
			return
		}
		// fmt.Println("materialmap", materialmap)
		tiles, err := world.LoadRawTiles("./res/testtiles", "testtile")
		if err != nil {
			chErr <- err
			return
		}
		// fmt.Println("tiles", len(tiles))
		level, err := world.NewSimpleWorld(tilemap, materialmap, tiles, g.scale)
		if err != nil {
			chErr <- err
			return
		}
		g.world = level

		// Start Game Thread
		go func() {
			ticks := 0
			for !interrupt {
				if up {
					screenY++
				}
				if down {
					screenY--
				}
				if left {
					screenX--
				}
				if right {
					screenX++
				}

				g.screen.SetPos(image.Point{screenX, screenY})
				tStart := time.Now().UTC().UnixNano()

				g.tick()

				// Smooth Framerate
				deltaT := (time.Now().UTC().UnixNano() - tStart) / 10e6
				sleep := g.ttick - int(deltaT)
				if sleep > 0 {
					log.Printf("Sleeping %d millis", sleep)
					time.Sleep(time.Duration(sleep) * time.Millisecond)
				} else {
					log.Printf("Overdue %d millis", -sleep)
				}
				ticks++
				w.Send(paint.Event{})

			}
		}()

		// Loop Screen Thread until interrupted
		for !interrupt {
			switch e := w.NextEvent().(type) {
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					interrupt = true
					chInterrupt <- true
				}
			case paint.Event:
				err := g.screen.InitBuffer()
				if err != nil {
					chErr <- err
					return
				}
				g.render()
				g.screen.UploadBuffer()
			case key.Event:
				fmt.Println("Keypress:", e.Code, e.Direction)
				switch e.Code {
				case key.CodeW:
					if e.Direction == key.DirPress {
						up = true
					} else if e.Direction == key.DirRelease {
						up = false
					}
				case key.CodeS:
					if e.Direction == key.DirPress {
						down = true
					} else if e.Direction == key.DirRelease {
						down = false
					}
				case key.CodeA:
					if e.Direction == key.DirPress {
						left = true
					} else if e.Direction == key.DirRelease {
						left = false
					}
				case key.CodeD:
					if e.Direction == key.DirPress {
						right = true
					} else if e.Direction == key.DirRelease {
						right = false
					}
				}
			}
		}
	})

	// Wait for interruption or error signal
	for {
		select {
		case err := <-chErr:
			return err
		case <-chInterrupt:
			return nil
		}
	}
}

func (g *game) tick() {
	for _, ent := range g.entities {
		ent.Tick()
	}
}

func (g *game) render() {
	if g.bg != nil {
		g.bg.Render(g.screen)
	}
	if g.world != nil {
		g.world.Render(g.screen)
	}
	for _, ent := range g.entities {
		ent.Render(g.screen)
	}
}

func (g *game) addTestEntity(img *image.RGBA, p0 image.Point, dim image.Point, dp image.Point) {
	ent := entity.NewUniformSprite(1, img, p0, dim, func(e entity.Entity) {
		e.SetPos(e.Pos().Add(dp))
	})
	g.entities = append(g.entities, ent)
}
