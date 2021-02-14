package types

import (
	"image"
	"image/draw"
)

type IComponent interface {
	Draw(r image.Rectangle) func(draw.Image) image.Rectangle
	Run()
}

// type IErrorHandler interface {
// 	HandleError(w http.ResponseWriter, err error)
// 	Set401(w http.ResponseWriter, str string)
// }
