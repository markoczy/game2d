package game

import (
	// "fmt"
	"github.com/markoczy/game2d/display"
	"github.com/markoczy/game2d/entity"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/lifecycle"
	"image"
	// "image"
	"image/color"
	"log"
	"math/rand"
	"time"
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
	entities      []entity.Entity
}

var (
	black = color.RGBA{0x00, 0x00, 0x00, 0xff}
	white = color.RGBA{0xff, 0xff, 0xff, 0xff}
)

func (g *game) Run() error {
	// fmt.Println("Start.")

	// k := 10 * g.scale
	// img := image.NewRGBA(image.Rect(0, 0, k, k))
	// // for i := range img.Pix {}
	// for x := 0; x < k; x++ {
	// 	for y := 0; y < k; y++ {
	// 		// if x == y {
	// 		img.SetRGBA(x, y, white)
	// 		// }
	// 	}
	// }
	img, err := display.LoadImage("./res/sample0.png")
	if err != nil {
		return err
	}
	img = display.ScaleImage(img, g.scale)
	g.addTestEntity(img, image.Point{10, 10}, image.Point{1, 1})
	g.addTestEntity(img, image.Point{50, 50}, image.Point{-1, 1})
	g.addTestEntity(img, image.Point{100, 100}, image.Point{1, -1})
	g.addTestEntity(img, image.Point{200, 100}, image.Point{-1, -1})

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
		// g.screen.SetPos(image.Point{100, 100})

		// Start Game Thread
		go func() {
			ticks := 0
			for !interrupt {
				// g.screen.SetPos(image.Point{ticks, ticks})
				tStart := time.Now().UTC().UnixNano()
				err := g.screen.InitBuffer()
				if err != nil {
					chErr <- err
					return
				}

				g.tick()
				g.render()

				// Smooth Framerate
				deltaT := (time.Now().UTC().UnixNano() - tStart) / 10e6
				g.screen.UploadBuffer()
				sleep := g.ttick - int(deltaT)
				if sleep > 0 {
					log.Printf("Sleeping %d millis", sleep)
					time.Sleep(time.Duration(sleep) * time.Millisecond)
				} else {
					log.Printf("Overdue %d millis", -sleep)
				}
				ticks++

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
	for _, ent := range g.entities {
		ent.Render(g.screen)
	}
}

func (g *game) addTestEntity(img *image.RGBA, p0 image.Point, dp image.Point) {
	ent := entity.NewUniformSprite(1, img, p0, func(e entity.Entity) {
		e.SetPos(e.Pos().Add(dp))
	})
	g.entities = append(g.entities, ent)
}
