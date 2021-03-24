package main

import (
	"github.com/jordy2254/indoormaprestapi/pkg/gorm"
	"github.com/jordy2254/indoormaprestapi/pkg/rest"
	"github.com/jordy2254/indoormaprestapi/pkg/rest/wrappers"
	"net/http"
)

func main() {
	gormConnectionString := "admin:welcome@tcp(localhost:3306)/project?charset=utf8mb4&parseTime=True&loc=Local"
	dbConnection := gorm.Connect(gormConnectionString)

	restService := rest.New(dbConnection)

	http.ListenAndServe("192.168.0.28:3500", wrappers.NewCorsWrapper().Handler(restService))
}
