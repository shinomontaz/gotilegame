package game

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type Controller interface {
	Navigate(to string)
	Quit()
	Log(msg string)
}

const (
	STATE_INIT = iota
	STATE_RDY
)

type Game struct {
	c           Controller
	state       int
	treesFrames []pixel.Rect
	matrices    []pixel.Matrix
	trees       []*pixel.Sprite
	spritesheet pixel.Picture
	last        time.Time
	dt          float64
	mu          sync.Mutex
}

func New(c Controller) *Game {
	return &Game{
		c:           c,
		state:       STATE_INIT,
		treesFrames: make([]pixel.Rect, 0),
		matrices:    make([]pixel.Matrix, 0),
		trees:       make([]*pixel.Sprite, 0),
		last:        time.Now(),
	}
}

func (g *Game) Prepare() {
	rand.Seed(time.Now().UTC().UnixNano())

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)

	path := dir + "\\sprites\\trees.png"
	fmt.Println(path)
	spritesheet, err := loadPicture(path)
	if err != nil {
		panic(err)
	}

	g.spritesheet = spritesheet

	for x := g.spritesheet.Bounds().Min.X; x < g.spritesheet.Bounds().Max.X; x += 32 {
		for y := spritesheet.Bounds().Min.Y; y < g.spritesheet.Bounds().Max.Y; y += 32 {
			g.treesFrames = append(g.treesFrames, pixel.R(x, y, x+32, y+32))
		}
	}

	g.mu.Lock()
	g.state = STATE_RDY
	g.mu.Unlock()

	time.Sleep(5 * time.Second)
}

func (g *Game) Draw(win *pixelgl.Window) {
	g.mu.Lock()
	state := g.state
	g.mu.Unlock()

	switch state {
	case STATE_INIT:
		g.DrawInit(win)
	case STATE_RDY:
		g.DrawRdy(win)
	default:
		log.Fatal("unknow game state")
	}
}

func (g *Game) DrawInit(win *pixelgl.Window) {
	go g.Prepare()
	win.Clear(colornames.Black)
}

func (g *Game) DrawRdy(win *pixelgl.Window) {
	g.dt = time.Since(g.last).Seconds()

	var (
		camPos       = pixel.ZV
		camSpeed     = 500.0
		camZoom      = 1.0
		camZoomSpeed = 1.2
	)
	cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
	win.SetMatrix(cam)

	if win.JustPressed(pixelgl.MouseButtonLeft) {
		tree := pixel.NewSprite(g.spritesheet, g.treesFrames[rand.Intn(len(g.treesFrames))])
		g.trees = append(g.trees, tree)
		mouse := cam.Unproject(win.MousePosition())
		g.matrices = append(g.matrices, pixel.IM.Scaled(pixel.ZV, 3).Moved(mouse))
	}

	if win.Pressed(pixelgl.KeyEscape) {
		g.c.Navigate("firstcreen")
	}
	if win.Pressed(pixelgl.KeyLeft) {
		camPos.X -= camSpeed * g.dt
	}
	if win.Pressed(pixelgl.KeyRight) {
		camPos.X += camSpeed * g.dt
	}
	if win.Pressed(pixelgl.KeyUp) {
		camPos.Y += camSpeed * g.dt
	}
	if win.Pressed(pixelgl.KeyDown) {
		camPos.Y -= camSpeed * g.dt
	}

	camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)

	for i, tree := range g.trees {
		tree.Draw(win, g.matrices[i])
	}

	g.last = time.Now()
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("sdfrsdfs")
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("sdfrsdfs 2")

		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}
