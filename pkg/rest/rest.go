package rest

import (
	"github.com/gorilla/mux"
	"github.com/jordy2254/indoormaprestapi/pkg/gorm/store"
	"github.com/jordy2254/indoormaprestapi/pkg/rest/handlers"
	"github.com/jordy2254/indoormaprestapi/pkg/rest/middleware"
	"github.com/op/go-logging"
	"gorm.io/gorm"
	"net/http"
)

func New(db *gorm.DB, logger *logging.Logger) http.Handler {

	auth0Middleware := middleware.CreateAuth0MiddleWare("http://192.168.0.28:3500","https://dev-e5s8h580.us.auth0.com/")
	loggerMiddleware := middleware.NewRouteLogger(logger)

	router := mux.NewRouter()
	router.Use(loggerMiddleware.Handler)

	mapStore := store.NewMapStore(db)
	buildingStore := store.NewBuildingStore(db)
	floorStore := store.NewFloorStore(db)
	roomStore := store.NewRoomStore(db)
	indentStore := store.NewIndentStore(db)
	pathStore := store.NewPathStore(db)
	sensorStore := store.NewSensorStore(db)



	rh := handlers.NewRouteHelper(router, auth0Middleware)

	handlers.AddMapAPI(rh, &mapStore, logger)
	handlers.AddBuildingAPI(rh, &buildingStore, logger)
	handlers.AddFloorAPI(rh, &floorStore, logger)
	handlers.AddRoomAPI(rh, &roomStore, logger)
	handlers.AddIndentAPI(rh, &indentStore, logger)
	handlers.AddPathsAPI(rh, &pathStore, logger)
	handlers.AddSensorAPI(rh, &sensorStore, logger)

	return router
}

