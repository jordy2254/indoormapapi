package main

import (
	"github.com/jordy2254/indoormaprestapi/pkg/gorm"
	"github.com/jordy2254/indoormaprestapi/pkg/rest"
	"github.com/jordy2254/indoormaprestapi/pkg/rest/wrappers"
	"net/http"
	"github.com/op/go-logging"
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
