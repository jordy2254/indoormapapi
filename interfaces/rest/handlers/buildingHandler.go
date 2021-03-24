package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jordy2254/indoormaprestapi/interfaces/gorm/store"
	"github.com/jordy2254/indoormaprestapi/model"
	"io/ioutil"
	"net/http"
	"strconv"
)

type BuildingController struct {
	BuildingStore *store.BuildingStore
}

func AddBuildingAPI(router *mux.Router, buildingStore *store.BuildingStore) {
	controller := BuildingController{BuildingStore: buildingStore}

	router.HandleFunc("/Buildings/{mapId}/{id}", controller.getBuilding).Methods("GET", "OPTIONS")
	router.HandleFunc("/Buildings/{mapId}/{id}", controller.updateBuilding).Methods("POST")
	router.HandleFunc("/Buildings/{mapId}", controller.createBuilding).Methods("POST")
}


func (bc *BuildingController) createBuilding(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Create Building")
	params := mux.Vars(r)
	mapid, err := strconv.Atoi(params["mapId"])
	if err != nil {
		return
	}
	newBuilding, err := unmarshalBuildingRequest(w, r)
	if err != nil {
		return
	}
	if mapid != newBuilding.Id && newBuilding.Id != 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	newBuilding.MapId = mapid

	bc.BuildingStore.CreateBuilding(&newBuilding)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bc.BuildingStore.GetBuildingById(mapid, newBuilding.Id))
}

func (bc *BuildingController) updateBuilding(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println("Endpoint Hit: Update Building")
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	mapId, err := strconv.Atoi(params["mapId"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	newBuilding, err := unmarshalBuildingRequest(w, r)
	if err != nil || newBuilding.MapId != mapId || newBuilding.Id != id {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	bc.BuildingStore.UpdateBuilding(&newBuilding)
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(bc.BuildingStore.GetBuildingById(newBuilding.MapId, newBuilding.Id))
}

func (bc *BuildingController) getBuilding(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println("Endpoint Hit: getBuilding")
	mapId, err := strconv.Atoi(params["mapId"])
	if err != nil {
		return
	}
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return
	}
	json.NewEncoder(w).Encode(bc.BuildingStore.GetBuildingById(id, mapId))
}

func unmarshalBuildingRequest(w http.ResponseWriter, r *http.Request) (model.Building, error) {
	var newBuilding = model.Building{}

	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return newBuilding, err
	}

	err = json.Unmarshal(reqBody, &newBuilding)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return newBuilding, err
	}

	return newBuilding, nil
}
