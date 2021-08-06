package main

import (
	"encoding/json"
	"github.com/jordy2254/indoormaprestapi/pkg/gorm"
	"github.com/jordy2254/indoormaprestapi/pkg/gorm/store"
	"github.com/jordy2254/indoormaprestapi/pkg/model"
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/sdlcanvas"
	"os"
)
var(
	 TOP = "TOP"
	 BOTTOM = "BOTTOM"
	 LEFT = "LEFT"
	 RIGHT = "RIGHT"
)
var (

	testData model.Room = model.Room{
		Location:   model.NewPoint2f(0,0),
		Dimensions: model.NewPoint2f(100,100),

		Indents: []model.Indent{
			{
				WallKeyA:   &TOP,
				WallKeyB:   &LEFT,
				Dimensions: model.NewPoint2f(10,10),
			},
			{
				WallKeyA:   &TOP,
				WallKeyB:   &RIGHT,
				Dimensions: model.NewPoint2f(10,10),
			},
			{
				WallKeyA:   &BOTTOM,
				WallKeyB:   &LEFT,
				Dimensions: model.NewPoint2f(10,10),
			},
			{
				WallKeyA:   &BOTTOM,
				WallKeyB:   &RIGHT,
				Dimensions: model.NewPoint2f(10,10),
			},
			{
				WallKeyA:   &BOTTOM,
				Location: 25,
				Dimensions: model.NewPoint2f(10,10),
			},
			{
				WallKeyA:   &TOP,
				Location: 25,
				Dimensions: model.NewPoint2f(10,10),
			},
			{
				WallKeyA:   &LEFT,
				Location: 25,
				Dimensions: model.NewPoint2f(10,10),
			},
			{
				WallKeyA:   &RIGHT,
				Location: 25,
				Dimensions: model.NewPoint2f(10,10),
			},
		},
		Polygon:    nil,
		Walls:      nil,

	}
)

func main() {
	gormConnectionString := "admin:welcome@tcp(localhost:3306)/project?charset=utf8mb4&parseTime=True&loc=Local"
	dbConnection := gorm.Connect(gormConnectionString)

	rs := store.NewRoomStore(dbConnection)
	r := rs.GetRoomById(3)
	_ = r

	ms := store.NewMapStore(dbConnection)
	m := ms.GetMapById(1)
	_ = ms
	_ = m

	testData.Walls = model.CalculatePolygonEdgePairs(testData, false)
	testData.Polygon = model.CalculatePolygonPoints(testData)
	json.NewEncoder(os.Stdout).Encode(testData.Polygon)
	m2 := model.Map{
		Buildings: []model.Building{{
			Floors: []model.Floor{{
				Rooms: []model.Room{testData},
			}},
		}},
	}

	_ = m2

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
	offs := model.Point2f{
		X: ptrFromVar(0),
		Y: ptrFromVar(0),
	}

	wnd.MouseWheel = func(x, y int) {
		if y > 0 {
			*scale.X = 1.1 * *scale.X
			*scale.Y = 1.1 * *scale.Y
		}

		if y < 0 {
			*scale.X = 0.9 * *scale.X
			*scale.Y = 0.9 * *scale.Y
		}

	}

	initialClick := model.Point2f{
		X: ptrFromVar(0),
		Y: ptrFromVar(0),
	}
	initialOff := model.Point2f{
		X: ptrFromVar(0),
		Y: ptrFromVar(0),
	}
	down := false

	wnd.MouseDown = func(button, x, y int) {
		*initialClick.X = float64(x)
		*initialClick.Y = float64(y)
		*initialOff.X = *offs.X
		*initialOff.Y = *offs.Y
		down = true
	}

	wnd.MouseUp = func(button, x, y int) {
		down = false
	}

	wnd.MouseMove = func(x, y int) {
		if !down{
			return
		}
		nx := float64(x) - *initialClick.X + *initialOff.X
		ny := float64(y) - *initialClick.Y + *initialOff.Y

		*offs.X = nx
		*offs.Y = ny
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

				offset := model.Point2f{
					X: ptrFromVar(*offs.X + *room.Location.X),
					Y: ptrFromVar(*offs.Y + *room.Location.Y),
				}
				cv.BeginPath()
				moveToPoint(cv, room.Polygon[0], offset, scale)
				for _, f := range room.Polygon {
					lineToPoint(cv, f, offset, scale)
				}
				lineToPoint(cv, room.Polygon[0], offset, scale)
				cv.Fill()

				for _, wall := range room.Walls {
					cv.BeginPath()
					moveToPoint(cv, wall.First, offset, scale)
					lineToPoint(cv, wall.Second, offset, scale)
					cv.Stroke()
				}
			}
		}
	})
}

func moveToPoint(cv *canvas.Canvas, point, offset, scale model.Point2f) {
	cv.MoveTo(((*point.X) * (*scale.X)) + (*offset.X * (*scale.X)), ((*point.Y) * (*scale.Y)) + (*offset.Y * (*scale.X)))
}

func lineToPoint(cv *canvas.Canvas, point, offset, scale model.Point2f){
	cv.LineTo(((*point.X) * (*scale.X)) + (*offset.X * (*scale.X)), ((*point.Y) * (*scale.Y)) + (*offset.Y * (*scale.X)))
}

func ptrFromVar(x float64) *float64{
	return &x
}