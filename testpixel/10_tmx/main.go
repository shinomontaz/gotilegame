package main

import (
	"image/png"
	"math"
	"os"
	"time"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/salviati/go-tmx/tmx"

	"golang.org/x/image/font/basicfont"

	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"
)

var clearColor = colornames.Skyblue

var platforms []pixel.Rect

func gameloop(win *pixelgl.Window, tilemap *tmx.Map) {
	batches := make([]*pixel.Batch, 0)
	batchIndices := make(map[string]int)
	batchCounter := 0

	// Load the sprites
	sprites := make(map[string]*pixel.Sprite)
	for _, tileset := range tilemap.Tilesets {
		if _, alreadyLoaded := sprites[tileset.Image.Source]; !alreadyLoaded {
			sprite, pictureData := loadSprite(tileset.Image.Source)
			sprites[tileset.Image.Source] = sprite
			batches = append(batches, pixel.NewBatch(&pixel.TrianglesData{}, pictureData))
			batchIndices[tileset.Image.Source] = batchCounter
			batchCounter++
		}
	}

	var (
		camPos       = pixel.ZV
		camSpeed     = 1000.0
		camZoom      = 0.2
		camZoomSpeed = 1.2
	)

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	//canvas
	imd := imdraw.New(nil)

	basicScreenLogger := screenLogger{
		bt:       text.New(pixel.V(10, 900), basicAtlas),
		ba:       basicAtlas,
		canvas:   imd,
		onBt:     true,
		onCanvas: true,
	}
	basicScreenLogger.initCanvas()

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		// Camera movement
		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		//		if win.JustPressed(pixelgl.KeyF1) {
		//			basicScreenLogger.onBt = !basicScreenLogger.onBt
		//		}
		//		if win.JustPressed(pixelgl.KeyF2) {
		//			basicScreenLogger.onCanvas = !basicScreenLogger.onCanvas
		//		}

		if win.Pressed(pixelgl.KeyLeft) {
			camPos.X -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyRight) {
			camPos.X += camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyDown) {
			camPos.Y -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyUp) {
			camPos.Y += camSpeed * dt
		}
		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)

		win.Clear(clearColor)

		// Draw tiles
		for _, batch := range batches {
			batch.Clear()
		}

		for _, layer := range tilemap.Layers {
			for tileIndex, tile := range layer.DecodedTiles {
				ts := layer.Tileset
				tID := int(tile.ID)

				if tID == 0 {
					// Tile ID 0 means blank, skip it.
					continue
				}

				// Calculate the framing for the tile within its tileset's source image
				numRows := ts.Tilecount / ts.Columns
				x, y := tileIDToCoord(tID, ts.Columns, numRows)
				gamePos := indexToGamePos(tileIndex, tilemap.Width, tilemap.Height)

				iX := float64(x) * float64(ts.TileWidth)
				fX := iX + float64(ts.TileWidth)
				iY := float64(y) * float64(ts.TileHeight)
				fY := iY + float64(ts.TileHeight)

				sprite := sprites[ts.Image.Source]
				sprite.Set(sprite.Picture(), pixel.R(iX, iY, fX, fY))
				pos := gamePos.ScaledXY(pixel.V(float64(ts.TileWidth), float64(ts.TileHeight)))
				sprite.Draw(batches[batchIndices[ts.Image.Source]], pixel.IM.Moved(pos))
			}
		}

		for _, batch := range batches {
			batch.Draw(win)
		}

		basicScreenLogger.draw(win, cam)

		win.Update()
	}
}

func tileIDToCoord(tID int, numColumns int, numRows int) (x int, y int) {
	x = tID % numColumns
	y = numRows - (tID / numColumns) - 1
	return
}

func indexToGamePos(idx int, width int, height int) pixel.Vec {
	gamePos := pixel.V(
		float64(idx%width)-1,
		float64(height)-float64(idx/width),
	)
	return gamePos
}

func run() {
	// Create the window with OpenGL
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Tilemaps",
		Bounds: pixel.R(0, 0, 800, 600),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	panicIfErr(err)

	// Initialize art assets (i.e. the tilemap)
	tilemap, err := tmx.ReadFile("map.tmx")
	panicIfErr(err)

	initPlatforms(tilemap)

	gameloop(win, tilemap)
}

func loadSprite(path string) (*pixel.Sprite, *pixel.PictureData) {
	f, err := os.Open(path)
	panicIfErr(err)

	img, err := png.Decode(f)
	panicIfErr(err)

	pd := pixel.PictureDataFromImage(img)
	return pixel.NewSprite(pd, pd.Bounds()), pd
}

func main() {
	pixelgl.Run(run)
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func initPlatforms(tilemap *tmx.Map) {

	platforms = make([]pixel.Rect, 0)
	levelheigth := tilemap.TileHeight * tilemap.Height
	//	levelwidth := tilemap.TileWidth * tilemap.Width

	//	fmt.Println(tilemap.ObjectGroups)

	for _, og := range tilemap.ObjectGroups {
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
			// rc := pixel.Rect{
			// 	Min: pixel.V(og.Min.X, og.Min.Y),
			// 	Max: pixel.V(og.Max.X, og.Max.Y),
			// }

			platforms = append(platforms, rc)

		}
	}

}
