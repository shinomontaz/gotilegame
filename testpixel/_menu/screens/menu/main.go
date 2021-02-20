package menu

import (
	"fmt"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

type Menu struct {
	items        []Item
	curr         int
	defaultColor color.Color
	selectColor  color.Color
	c            Controller
}

type Item struct {
	title    string
	action   func()
	scale    float64
	selected bool
	col      color.Color
	txt      *text.Text
}

type Controller interface {
	Navigate(to string)
	Quit()
	Log(msg string)
}

func NewMain(atlas *text.Atlas, c Controller) *Menu {
	m := &Menu{
		items:        []Item{},
		defaultColor: colornames.Blue,
		selectColor:  colornames.Red,
		c:            c,
	}

	m.AddItem("New game", atlas, func() {
		m.curr = 0
		fmt.Println("New game")
	})
	m.AddItem("Load", atlas, func() {
		m.curr = 1
		fmt.Println("second")
	})
	m.AddItem("Options", atlas, func() {
		m.curr = 2
		m.c.Navigate("secondcreen")
	})
	m.AddItem("Quit", atlas, func() {
		fmt.Println("quit")
		m.curr = 3
		m.c.Quit()
	})

	return m
}

func NewSec(atlas *text.Atlas, c Controller) *Menu {
	m := &Menu{
		items:        []Item{},
		defaultColor: colornames.Blue,
		selectColor:  colornames.Red,
		c:            c,
	}

	m.AddItem("2 another", atlas, func() { m.curr = 0; fmt.Println("2 first") })
	m.AddItem("2 second", atlas, func() { m.curr = 1; fmt.Println("2 second") })
	m.AddItem("return", atlas, func() {
		fmt.Println("return")
		m.curr = 2
		m.c.Navigate("firstcreen")
	})

	return m
}

func (m *Menu) AddItem(name string, atlas *text.Atlas, action func()) {
	txt := text.New(pixel.V(1, 1), atlas)
	txt.Color = m.defaultColor
	item := Item{
		action: action,
		txt:    txt,
		title:  name,
	}
	m.items = append(m.items, item)
}

func (m *Menu) controls(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.KeyUp) {
		m.curr = (m.curr - 1 + len(m.items)) % len(m.items)
	}
	if win.JustPressed(pixelgl.KeyDown) {
		m.curr = (m.curr + 1) % len(m.items)
	}
	if win.JustPressed(pixelgl.KeyEnter) {
		m.items[m.curr].action()
	}
}

func (m *Menu) Draw(win *pixelgl.Window) {
	w := win.Bounds().W()
	h := win.Bounds().H()

	itemScale := 2.0 //w / h / 3
	offsetX := w / 2
	offsetY := 0.0

	m.controls(win)

	for i := range m.items {
		m.items[i].txt.Clear()
		m.items[i].txt.Dot.X -= m.items[i].txt.BoundsOf(m.items[i].title).W() / 2
		m.items[i].txt.Color = m.defaultColor

		if i == m.curr {
			m.items[i].txt.Color = m.selectColor
		}

		fmt.Fprintln(m.items[i].txt, m.items[i].title)

		offsetY = h/1.6 - float64(i)*m.items[i].txt.LineHeight*itemScale
		m.items[i].txt.Draw(win, pixel.IM.Scaled(pixel.ZV, itemScale).Moved(pixel.V(offsetX, offsetY)))
	}
}
