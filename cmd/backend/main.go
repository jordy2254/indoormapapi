package main

import (
	"github.com/jordy2254/indoormaprestapi/pkg/gorm"
	"github.com/jordy2254/indoormaprestapi/pkg/model"
	"github.com/jordy2254/indoormaprestapi/pkg/rest"
	"github.com/jordy2254/indoormaprestapi/pkg/rest/wrappers"
	"github.com/op/go-logging"
	"github.com/tfriedel6/canvas/sdlcanvas"
	"net/http"
	"os"
)

var(
	logger  = logging.MustGetLogger("example")
	format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{color:reset} %{message}`,
	)
)

func main() {
	logging.SetBackend(logging.NewBackendFormatter(logging.NewLogBackend(os.Stdout, "", 0),format))

	gormConnectionString := "admin:welcome@tcp(localhost:3306)/project?charset=utf8mb4&parseTime=True&loc=Local"
	dbConnection := gorm.Connect(gormConnectionString)

	restService := rest.New(dbConnection, logger)

	listedAdd := "192.168.0.28:3500"

	logger.Infof("Rest Service started on %s", listedAdd)

	http.ListenAndServe(listedAdd, wrappers.NewCorsWrapper().Handler(restService))
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
