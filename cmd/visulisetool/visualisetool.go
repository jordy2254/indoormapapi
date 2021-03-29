package main

import (
	"github.com/jordy2254/indoormaprestapi/pkg/gorm"
	"github.com/jordy2254/indoormaprestapi/pkg/gorm/store"
	"github.com/jordy2254/indoormaprestapi/pkg/model"
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/sdlcanvas"
)

func main() {
	gormConnectionString := "admin:welcome@tcp(localhost:3306)/project?charset=utf8mb4&parseTime=True&loc=Local"
	dbConnection := gorm.Connect(gormConnectionString)
	ms := store.NewMapStore(dbConnection)
	m := ms.GetMapById(1)

	mapRenderer(m)
}


func mapRenderer(m model.Map) {
	wnd, cv, err := sdlcanvas.CreateWindow(1280, 720, "Hello")
	if err != nil {
		panic(err)
	}
	defer wnd.Destroy()
	scale := model.Point2f{
		X: ptrFromVar(1),
		Y: ptrFromVar(1),
	}

	wnd.MainLoop(func() {
		w, h := float64(cv.Width()), float64(cv.Height())
		cv.SetFillStyle("#888888")
		cv.FillRect(0, 0, w, h)
		cv.SetStrokeStyle("#000000")
		cv.SetFillStyle("#f0ae5d")
		cv.SetLineWidth(2.0)

		for _, building := range m.Buildings {
			for _, room := range building.Floors[0].Rooms {
				cv.BeginPath()
				moveToPoint(cv, room.Polygon[0], room.Location, scale)
				for _, f := range room.Polygon {
					lineToPoint(cv, f, room.Location, scale)
				}
				lineToPoint(cv, room.Polygon[0], room.Location, scale)
				cv.Fill()

				for _, wall := range room.Walls {
					cv.BeginPath()
					moveToPoint(cv, wall.First, room.Location, scale)
					lineToPoint(cv, wall.Second, room.Location, scale)
					cv.Stroke()
				}
			}
		}
	})
}

func moveToPoint(cv *canvas.Canvas, point, offset, scale model.Point2f) {
	cv.MoveTo(((*point.X) * (*scale.X)) + (*offset.X), ((*point.Y) * (*scale.Y)) + (*offset.Y))
}

func lineToPoint(cv *canvas.Canvas, point, offset, scale model.Point2f){
	cv.LineTo(((*point.X) * (*scale.X)) + (*offset.X), ((*point.Y) * (*scale.Y)) + (*offset.Y))
}

func ptrFromVar(x float64) *float64{
	return &x
}