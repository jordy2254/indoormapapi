package main

import (
	"github.com/jordy2254/indoormaprestapi/pkg/model"
	"github.com/tfriedel6/canvas/sdlcanvas"
)

func main() {
	r := model.Room{
		Name:       "",
		Location:   model.Point2f{
			X: ptrFromVar(0),
			Y: ptrFromVar(0),
		},
		Dimensions: model.Point2f{
			X: ptrFromVar(100),
			Y: ptrFromVar(100),
		},
		Indents: []model.Indent{
			{
				WallKeyA: "TOP",
				WallKeyB: "LEFT",
				Dimensions: model.Point2f{
					X: ptrFromVar(10),
					Y: ptrFromVar(10),
				},
			},

			{
				WallKeyA: "TOP",
				Location: 30.0,
				Dimensions: model.Point2f{
					X: ptrFromVar(35),
					Y: ptrFromVar(10),
				},
			},

			{
				WallKeyA: "BOTTOM",
				WallKeyB: "LEFT",
				Dimensions: model.Point2f{
					X: ptrFromVar(35),
					Y: ptrFromVar(50),
				},
			},
		},
		Polygon:    nil,
		Walls:      nil,
		Entrances:  []model.Entrance{
			{
				WallKey:  "BOTTOM",
				Location: 65,
				Length:   35,
			},
			{
				WallKey:  "TOP",
				Location: 65,
				Length:   35,
			},
		},
	}

	r.Walls = model.CalculatePolygonEdgePairs(r, false)
	renderWalls(r.Walls)
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
		cv.SetFillStyle("#FF00FF")

		for _, wall := range walls {
			cv.BeginPath()
			cv.MoveTo(*wall.First.X + 10, *wall.First.Y + 10)
			cv.LineTo(*wall.Second.X + 10, *wall.Second.Y + 10)
			cv.Stroke()
		}
	})
}

func ptrFromVar(x float64) *float64{
	return &x
}