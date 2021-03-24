package rest

import (
	"github.com/gorilla/mux"
	"github.com/jordy2254/indoormaprestapi/interfaces/gorm/store"
	"github.com/jordy2254/indoormaprestapi/interfaces/rest/handlers"
	"github.com/jordy2254/indoormaprestapi/interfaces/rest/middleware"
	"gorm.io/gorm"
	"net/http"
)

func New(db *gorm.DB) http.Handler {

	auth0Middleware := middleware.CreateAuth0MiddleWare("http://192.168.0.28:3500","https://dev-e5s8h580.us.auth0.com/")

	router := mux.NewRouter()
	router.Use(auth0Middleware.Handler)

	mapStore := store.NewMapStore(db)
	buildingStore := store.NewBuildingStore(db)
	floorStore := store.NewFloorStore(db)
	roomStore := store.NewRoomStore(db)
	indentStore := store.NewIndentStore(db)


	handlers.AddMapAPI(router, &mapStore)
	handlers.AddBuildingAPI(router, &buildingStore)
	handlers.AddFloorAPI(router, &floorStore)
	handlers.AddRoomAPI(router, &roomStore)
	handlers.AddIndentAPI(router, &indentStore)

	return router
}

