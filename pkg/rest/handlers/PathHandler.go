package handlers

import "github.com/jordy2254/indoormaprestapi/pkg/gorm/store"

type PathsController struct {
	pathStore *store.PathStore
}

func AddPathsAPI(rh *RouteHelper, pathStore *store.PathStore) {
	controller := PathsController{pathStore: pathStore}
	_ = controller
}
