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

type RoomController struct {
	roomStore *store.RoomStore
	logger *logging.Logger
}

func AddRoomAPI(rh *RouteHelper, roomStore *store.RoomStore, logger *logging.Logger) {
	controller := RoomController{roomStore: roomStore, logger: logger}
	rh.protectedRoute("/Rooms/{buildingId}/{id}", controller.getRoom).Methods("GET", "OPTIONS")
	rh.protectedRoute("/Rooms/{buildingId}/{id}", controller.updateRoom).Methods("POST")
	rh.protectedRoute("/Rooms/{buildingId}", controller.createRoom).Methods("POST")
	rh.protectedRoute("/Rooms/{buildingId}/{id}/polygon", controller.generatePolygon).Methods("GET")
}


func (rc *RoomController) createRoom(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Create Room")

	params := mux.Vars(r)
	buildingId, err := strconv.Atoi(params["buildingId"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	newRoom, err := unmarshalRoomRequest(w, r)
	if err != nil || newRoom.BuildingId != buildingId{
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	rc.roomStore.CreateRoom(&newRoom)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rc.roomStore.GetRoomById(newRoom.Id))
}

func (rc *RoomController) updateRoom(w http.ResponseWriter, r *http.Request) {
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
	newRoom, err := unmarshalRoomRequest(w, r)
	if err != nil || newRoom.BuildingId != buildingId || newRoom.Id != id {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	rc.roomStore.UpdateRoom(newRoom)
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(rc.roomStore.GetRoomById(newRoom.Id))
}

func (rc *RoomController) getRoom(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println("Endpoint Hit: Update Floor")
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	floorId, err := strconv.Atoi(params["floorId"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	_ = floorId
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(rc.roomStore.GetRoomById(id))
}

func (rc *RoomController) generatePolygon(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println("Endpoint Hit: Generate Polygon")
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	floorId, err := strconv.Atoi(params["floorId"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	_ = floorId

	room := rc.roomStore.GetRoomById(id)
	points := model.CalculatePolygonPoints(room)

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(points)
}

func unmarshalRoomRequest(w http.ResponseWriter, r *http.Request) (model.Room, error) {
	var newRoom = model.Room{}

	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return newRoom, err
	}

	err = json.Unmarshal(reqBody, &newRoom)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return newRoom, err
	}

	return newRoom, nil
}
