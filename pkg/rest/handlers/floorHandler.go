package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jordy2254/indoormaprestapi/pkg/gorm/store"
	"github.com/jordy2254/indoormaprestapi/pkg/model"
	"github.com/op/go-logging"
	"io/ioutil"
	"net/http"
	"strconv"
)

type FloorController struct {
	floorStore *store.FloorStore
}

func AddFloorAPI(rh *RouteHelper, buildingStore *store.FloorStore, logger *logging.Logger) {
	controller := FloorController{floorStore: buildingStore}

	rh.protectedRoute("/Floors/{buildingId}/{id}", controller.getFloor).Methods("GET")
	rh.protectedRoute("/Floors/{buildingId}/{id}", controller.updateFloor).Methods("POST")
	rh.protectedRoute("/Floors/{buildingId}", controller.createFloor).Methods("POST")
}

func (fc *FloorController) createFloor(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Create Floor")

	params := mux.Vars(r)
	buildingId, err := strconv.Atoi(params["buildingId"])
	if err != nil {
		return
	}
	newFloor, err := unmarshalFloorRequest(w, r)
	if err != nil {
		return
	}
	if buildingId != newFloor.BuildingId && newFloor.BuildingId != 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	newFloor.BuildingId = buildingId

	fc.floorStore.CreateFloor(&newFloor)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(fc.floorStore.GetFloorById(newFloor.Id))
}

func (fc *FloorController) updateFloor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println("Endpoint Hit: Update Floor")
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	buildingId, err := strconv.Atoi(params["buildingId"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	newFloor, err := unmarshalFloorRequest(w, r)
	if err != nil || newFloor.BuildingId != buildingId || newFloor.Id != id {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	fc.floorStore.UpdateFloor(newFloor)
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(fc.floorStore.GetFloorById(newFloor.Id))
}

func (fc *FloorController) getFloor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println("Endpoint Hit: Get Floor")
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	buildingId, err := strconv.Atoi(params["buildingId"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	_ = buildingId
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(fc.floorStore.GetFloorById(id))
}

func unmarshalFloorRequest(w http.ResponseWriter, r *http.Request) (model.Floor, error) {
	var newFloor = model.Floor{}

	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return newFloor, err
	}

	err = json.Unmarshal(reqBody, &newFloor)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return newFloor, err
	}

	return newFloor, nil
}
