package rest

import (
	"github.com/gorilla/mux"
	"github.com/jordy2254/indoormaprestapi/pkg/gorm/store"
	"github.com/jordy2254/indoormaprestapi/pkg/rest/handlers"
	"github.com/jordy2254/indoormaprestapi/pkg/rest/middleware"
	"gorm.io/gorm"
	"net/http"
)

func New(db *gorm.DB) http.Handler {

	auth0Middleware := middleware.CreateAuth0MiddleWare("http://192.168.0.28:3500","https://dev-e5s8h580.us.auth0.com/")

	router := mux.NewRouter()

	mapStore := store.NewMapStore(db)
	buildingStore := store.NewBuildingStore(db)
	floorStore := store.NewFloorStore(db)
	roomStore := store.NewRoomStore(db)
	indentStore := store.NewIndentStore(db)

	rh := handlers.NewRouteHelper(router, auth0Middleware)

	handlers.AddMapAPI(rh, &mapStore)
	handlers.AddBuildingAPI(rh, &buildingStore)
	handlers.AddFloorAPI(rh, &floorStore)
	handlers.AddRoomAPI(rh, &roomStore)
	handlers.AddIndentAPI(rh, &indentStore)

	return router
}

