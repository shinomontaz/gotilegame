package components

import (
	"fmt"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

type Menu struct {
	Items        []Item
	curr         int
	DefaultColor color.Color
	SelectColor  color.Color
	c            Controller
	Active       bool
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

func (m *Menu) AddItem(name string, atlas *text.Atlas, action func()) {
	txt := text.New(pixel.V(1, 1), atlas)
	txt.Color = m.DefaultColor
	item := Item{
		action: action,
		txt:    txt,
		title:  name,
	}
	m.Items = append(m.Items, item)
}

func (m *Menu) Up() {
	m.curr = (m.curr - 1 + len(m.Items)) % len(m.Items)
}

func (m *Menu) Down() {
	m.curr = (m.curr + 1) % len(m.Items)
}

func (m *Menu) Action() {
	m.Items[m.curr].action()
}

func (m *Menu) Draw(win *pixelgl.Window) {
	cam := pixel.IM
	win.SetMatrix(cam)

	w := win.Bounds().W()
	h := win.Bounds().H()

	itemScale := 2.0 //w / h / 3
	offsetX := w / 2
	offsetY := 0.0

	for i := range m.Items {
		m.Items[i].txt.Clear()
		m.Items[i].txt.Dot.X -= m.Items[i].txt.BoundsOf(m.Items[i].title).W() / 2
		m.Items[i].txt.Color = m.DefaultColor

		if i == m.curr {
			m.Items[i].txt.Color = m.SelectColor
		}

		fmt.Fprintln(m.Items[i].txt, m.Items[i].title)

		offsetY = h/1.6 - float64(i)*m.Items[i].txt.LineHeight*itemScale
		m.Items[i].txt.Draw(win, pixel.IM.Scaled(pixel.ZV, itemScale).Moved(pixel.V(offsetX, offsetY)))
	}
}
