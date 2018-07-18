package main

import (
	"fmt"
	"github.com/markoczy/game2d/game"
)

func main() {
	fmt.Println("Started.")
	// game.NewGame(800, 800, 10, 20).Run()
	game.NewGame(800, 800, 10, 1000).Run()
}
