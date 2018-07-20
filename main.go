package main

import (
	"fmt"
	"github.com/markoczy/game2d/game"
)

func main() {
	fmt.Println("Started.")
	// game.NewGame(800, 800, 10, 20).Run()
	err := game.NewGame(800, 800, 10, 1000).Run()
	if err != nil {
		fmt.Println("Error:", err.Error())
	}
}
