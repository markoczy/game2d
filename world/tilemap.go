package world

import ()

// Tilemap map of tiles definition
type Tilemap struct {
	TilesWide, TilesHigh  int
	TileWidth, TileHeight int
	Tiles                 []*Tile
}

// Tile definition of a tile
type Tile struct {
	Index     int
	Tile      int
	X, Y, Rot int
	flipX     bool
}

// LoadTilemap loads JSON file to tilemap
func LoadTilemap(file string, layer int) (*Tilemap, error) {
	// dat, err := ioutil.ReadFile(file)
	// if err != nil {
	// 	return nil, err
	// }

	// var tilemap map[string]interface{}

	// json.Unmarshal(dat, &tilemap)
	tilemap, err := unmarshall(file)
	if err != nil {
		return nil, err
	}
	tilewidth, err := getIntValue(tilemap, "tilewidth")
	if err != nil {
		return nil, err
	}
	tileheight, err := getIntValue(tilemap, "tileheight")
	if err != nil {
		return nil, err
	}
	tileswide, err := getIntValue(tilemap, "tileswide")
	if err != nil {
		return nil, err
	}
	tileshigh, err := getIntValue(tilemap, "tileshigh")
	if err != nil {
		return nil, err
	}
	layers, err := getChildrenArray(tilemap, "layers")
	layerel := layers[layer]
	tilesel, err := getChildrenArray(layerel, "tiles")
	tiles := []*Tile{}
	for _, cur := range tilesel {
		tile, err := mapToTile(cur)
		if err != nil {
			return nil, err
		}
		tiles = append(tiles, tile)
	}
	return &Tilemap{TileWidth: tilewidth,
		TileHeight: tileheight,
		TilesWide:  tileswide,
		TilesHigh:  tileshigh,
		Tiles:      tiles}, nil
}

func mapToTile(data map[string]interface{}) (*Tile, error) {
	tile, err := getIntValue(data, "tile")
	if err != nil {
		return nil, err
	}
	if tile == -1 {
		return nil, nil
	}
	index, err := getIntValue(data, "index")
	if err != nil {
		return nil, err
	}
	rot, err := getIntValue(data, "rot")
	if err != nil {
		return nil, err
	}
	x, err := getIntValue(data, "x")
	if err != nil {
		return nil, err
	}
	y, err := getIntValue(data, "y")
	if err != nil {
		return nil, err
	}
	flipX, err := getBoolValue(data, "flipX")
	if err != nil {
		return nil, err
	}
	return &Tile{Index: index, Tile: tile,
		X: x, Y: y, Rot: rot,
		flipX: flipX}, nil
}
