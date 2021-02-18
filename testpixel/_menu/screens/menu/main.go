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
}

type Item struct {
	title    string
	action   func()
	scale    float64
	selected bool
	col      color.Color
	txt      *text.Text
}

func NewMain(atlas *text.Atlas) *Menu {
	m := &Menu{
		items:        []Item{},
		defaultColor: colornames.Blue,
		selectColor:  colornames.Red,
	}

	m.AddItem("first", atlas, func() { m.curr = 0; fmt.Println("first") })
	m.AddItem("second", atlas, func() { m.curr = 1; fmt.Println("second") })
	m.AddItem("third", atlas, func() { m.curr = 2; fmt.Println("third") })

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
