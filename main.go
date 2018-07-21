package main

import (
	"fmt"
	"github.com/markoczy/game2d/game"
	"github.com/markoczy/game2d/world"
)

func main() {
	fmt.Println("Started.")
	// game.NewGame(800, 800, 10, 20).Run()
	err := game.NewGame(800, 800, 10, 1000).Run()
	if err != nil {
		fmt.Println("Error:", err.Error())
	}

	tilemap, err := world.LoadTilemap("./res/tiles.json", 0)
	fmt.Println("tilemap", tilemap)

	materialmap, err := world.LoadMaterialmap("./res/materials.json")
	fmt.Println("materialmap", materialmap)
}
