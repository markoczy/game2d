package game

import (
	"fmt"
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
	return &game{width, height, scale, tick}
}

type Game interface {
	Run() error
}

type game struct {
	width, height, scale, tick int
}

var (
	black = color.RGBA{0x00, 0x00, 0x00, 0xff}
	white = color.RGBA{0xff, 0xff, 0xff, 0xff}
)

func (g *game) Run() error {
	fmt.Println("Start.")

	k := 10 * g.scale
	img := image.NewRGBA(image.Rect(0, 0, k, k))
	// for i := range img.Pix {}
	for x := 0; x < k; x++ {
		for y := 0; y < k; y++ {
			// if x == y {
			img.SetRGBA(x, y, white)
			// }
		}
	}
	ent := entity.NewUniformSprite(1, img, image.Point{10, 10})

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

		disp := display.NewScreen(w, s, g.width, g.height, g.scale)

		// Start Game Thread
		go func() {
			ticks := 0
			for !interrupt {
				disp.SetPos(image.Point{ticks, ticks})
				tStart := time.Now().UTC().UnixNano()
				err := disp.InitBuffer()
				if err != nil {
					chErr <- err
					return
				}
				ent.Render(disp)

				// Smooth Framerate
				deltaT := (time.Now().UTC().UnixNano() - tStart) / 10e6
				disp.UploadBuffer()
				sleep := g.tick - int(deltaT)
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
