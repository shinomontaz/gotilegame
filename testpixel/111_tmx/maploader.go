package main

import (
	"image/png"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	tmx "github.com/lafriks/go-tiled"
)

type Worldmap struct {
	tm           *tmx.Map
	geom         *tmx.ObjectGroup
	scenery      *tmx.ObjectGroup
	meta         *tmx.Object
	batches      []*pixel.Batch
	batchIndices map[string]int
	sprites      map[string]*pixel.Sprite
}

func NewWorldmap(dir, source string) *Worldmap {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	defer os.Chdir(cwd)
	os.Chdir(dir)

	tm, err := tmx.LoadFile(source)
	if err != nil {
		panic(err)
	}

	w := Worldmap{
		tm:           tm,
		batches:      make([]*pixel.Batch, 0),
		batchIndices: make(map[string]int),
		sprites:      make(map[string]*pixel.Sprite),
	}

	w.init()

	return &w
}

func (w *Worldmap) init() {
	for _, og := range w.tm.ObjectGroups {
		if og.Name == "geom" {
			w.geom = og
		}
		if og.Name == "meta" {
			w.meta = og.Objects[0]
		}
		if og.Name == "scenery" {
			w.scenery = og
		}
	}

	batchCounter := 0

	// Load the sprites
	for _, tileset := range w.tm.Tilesets {
		if len(tileset.Tiles) > 0 {
			for _, tile := range tileset.Tiles {
				if _, alreadyLoaded := w.sprites[tile.Image.Source]; !alreadyLoaded {
					sprite, pictureData := loadSprite(tile.Image.Source)
					w.sprites[tile.Image.Source] = sprite
					w.batches = append(w.batches, pixel.NewBatch(&pixel.TrianglesData{}, pictureData))
					w.batchIndices[tile.Image.Source] = batchCounter
					batchCounter++
				}
			}
		} else {
			if _, alreadyLoaded := w.sprites[tileset.Image.Source]; !alreadyLoaded {
				sprite, pictureData := loadSprite(tileset.Image.Source)
				w.sprites[tileset.Image.Source] = sprite
				w.batches = append(w.batches, pixel.NewBatch(&pixel.TrianglesData{}, pictureData))
				w.batchIndices[tileset.Image.Source] = batchCounter
				batchCounter++
			}
		}
	}
}

func (w *Worldmap) Data() pixel.Rect {
	totalheight := float64(w.tm.TileHeight * w.tm.Height)
	rect := pixel.Rect{
		Min: pixel.V(
			float64(w.meta.X),
			totalheight-float64(w.meta.Y)-float64(w.meta.Height),
		),
		Max: pixel.V(
			float64(w.meta.X)+float64(w.meta.Width),
			totalheight-float64(w.meta.Y),
		),
	}

	return rect
}

func (w *Worldmap) Draw(win *pixelgl.Window) {
	for _, batch := range w.batches {
		batch.Clear()
	}

	for _, layer := range w.tm.Layers {
		for tileIndex, tile := range layer.Tiles {
			if tile.Nil {
				continue
			}
			ts := tile.Tileset
			tID := int(tile.ID)

			// if tileIndex == 285 {
			// 	fmt.Printf("\n%+v\n", ts)
			// }

			// Calculate the framing for the tile within its tileset's source image

			numRows := ts.TileCount / ts.Columns

			x, y := tileIDToCoord(tID, ts.Columns, numRows)
			gamePos := indexToGamePos(tileIndex, w.tm.Width, w.tm.Height)

			iX := float64(x) * float64(ts.TileWidth)
			fX := iX + float64(ts.TileWidth)
			iY := float64(y) * float64(ts.TileHeight)
			fY := iY + float64(ts.TileHeight)

			sprite := w.sprites[ts.Image.Source]
			sprite.Set(sprite.Picture(), pixel.R(iX, iY, fX, fY))
			pos := gamePos.ScaledXY(pixel.V(float64(ts.TileWidth), float64(ts.TileHeight)))
			sprite.Draw(w.batches[w.batchIndices[ts.Image.Source]], pixel.IM.Moved(pos))
		}
	}

	totalheight := w.tm.TileHeight * w.tm.Height
	for _, o := range w.scenery.Objects {
		dTile, err := w.tm.TileGIDToTile(o.GID)
		if err != nil {
			panic(err) // TODO!
		}
		ts := dTile.Tileset
		tID := dTile.ID

		gamePos := pixel.V(o.X+o.Width/2.0, float64(totalheight)-o.Y+o.Height/2.0)
		tile := ts.Tiles[tID]

		iX := 0.0
		fX := float64(tile.Image.Width)
		iY := 0.0
		fY := float64(tile.Image.Height)

		sprite := w.sprites[tile.Image.Source]
		sprite.Set(sprite.Picture(), pixel.R(iX, iY, fX, fY))
		sprite.Draw(w.batches[w.batchIndices[tile.Image.Source]], pixel.IM.Moved(gamePos))
	}

	for _, batch := range w.batches {
		batch.Draw(win)
	}
}

func tileIDToCoord(tID int, numColumns int, numRows int) (x int, y int) {
	x = tID % numColumns
	y = numRows - (tID / numColumns) - 1
	return
}

func indexToGamePos(idx int, width int, height int) pixel.Vec {
	gamePos := pixel.V(
		float64(idx%width),
		float64(height)-float64(idx/width)-0.5,
	)
	return gamePos
}

func (w *Worldmap) Geometry() []pixel.Rect {
	geom := make([]pixel.Rect, 0)
	levelheigth := w.tm.TileHeight * w.tm.Height
	for _, og := range w.tm.ObjectGroups {
		for _, object := range og.Objects {

			min := pixel.V(
				float64(object.X),
				float64(levelheigth)-float64(object.Y),
			)
			max := pixel.Vec{
				X: min.X + float64(object.Width),
				Y: min.Y - float64(object.Height),
			}

			rc := pixel.Rect{
				Min: min,
				Max: max,
			}

			geom = append(geom, rc)

		}
	}
	return geom
}

func loadSprite(path string) (*pixel.Sprite, *pixel.PictureData) {
	f, err := os.Open(path)
	panicIfErr(err)

	img, err := png.Decode(f)
	panicIfErr(err)

	pd := pixel.PictureDataFromImage(img)
	return pixel.NewSprite(pd, pd.Bounds()), pd
}
