package main

import (
	"github.com/jordy2254/indoormaprestapi/pkg/model"
	"github.com/tfriedel6/canvas/sdlcanvas"
)

func main() {
	
}


func renderWalls(walls []*model.PairPoint2f) {
	wnd, cv, err := sdlcanvas.CreateWindow(1280, 720, "Hello")
	if err != nil {
		panic(err)
	}
	defer wnd.Destroy()

	wnd.MainLoop(func() {
		w, h := float64(cv.Width()), float64(cv.Height())
		cv.SetFillStyle("#000")
		cv.FillRect(0, 0, w, h)
		cv.SetStrokeStyle("#FFF")
		for _, wall := range walls {
			cv.BeginPath()
			cv.MoveTo(*wall.First.X, *wall.First.Y)
			cv.LineTo(*wall.Second.X, *wall.Second.Y)
			cv.Stroke()
		}
	})
}

func ptrFromVar(x float64) *float64{
	return &x
}