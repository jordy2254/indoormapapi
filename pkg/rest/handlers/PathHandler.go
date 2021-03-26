package handlers

import (
	"github.com/jordy2254/indoormaprestapi/pkg/gorm/store"
	"github.com/op/go-logging"
	"net/http"
)

type PathsController struct {
	pathStore *store.PathStore
}

func AddPathsAPI(rh *RouteHelper, pathStore *store.PathStore, logger *logging.Logger) {
	controller := PathsController{pathStore: pathStore}

	rh.protectedRoute("paths/nodes", controller.createNode).Methods("POST")
	rh.protectedRoute("paths/nodes/{id}", controller.updateNode).Methods("POST")
	rh.protectedRoute("paths/nodes/{id}", controller.deleteNode).Methods("DELETE")


	rh.protectedRoute("paths/nodes/link", controller.createLink).Methods("POST")
	rh.protectedRoute("paths/nodes/link/{id}", controller.updateLink).Methods("POST")
	rh.protectedRoute("paths/nodes/link/{id}", controller.deleteLink).Methods("DELETE")

	_ = controller
}

func (pc *PathsController) updateNode(w http.ResponseWriter, r *http.Request) {

}

func (pc *PathsController) createNode(w http.ResponseWriter, r *http.Request) {

}

func (pc *PathsController) deleteNode(w http.ResponseWriter, r *http.Request) {

}

func (pc *PathsController) updateLink(w http.ResponseWriter, r *http.Request) {

}

func (pc *PathsController) createLink(w http.ResponseWriter, r *http.Request) {

}

func (pc *PathsController) deleteLink(w http.ResponseWriter, r *http.Request) {

}
