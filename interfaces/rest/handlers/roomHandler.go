package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jordy2254/indoormaprestapi/interfaces/application"
	"github.com/jordy2254/indoormaprestapi/interfaces/gorm/store"
	"github.com/jordy2254/indoormaprestapi/model"
	"io/ioutil"
	"net/http"
	"strconv"
)

type RoomController struct {
	roomStore *store.RoomStore
}

func AddRoomAPI(router *mux.Router, roomStore *store.RoomStore) {
	controller := RoomController{roomStore: roomStore}
	router.HandleFunc("/Rooms/{floorId}/{id}", controller.getRoom).Methods("GET", "OPTIONS")
	router.HandleFunc("/Rooms/{floorId}/{id}", controller.updateRoom).Methods("POST")
	router.HandleFunc("/Rooms/{floorId}", controller.createRoom).Methods("POST")
	router.HandleFunc("/Rooms/{floorId}/{id}/polygon", controller.generatePolygon).Methods("GET")
}


func (rc *RoomController) createRoom(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Create Room")

	params := mux.Vars(r)
	floorId, err := strconv.Atoi(params["floorId"])
	if err != nil {
		return
	}
	newRoom, err := unmarshalRoomRequest(w, r)
	if err != nil {
		return
	}
	if floorId != newRoom.FloorId && newRoom.FloorId != 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	newRoom.FloorId = floorId

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
	floorId, err := strconv.Atoi(params["floorId"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	newRoom, err := unmarshalRoomRequest(w, r)
	if err != nil || newRoom.FloorId != floorId || newRoom.Id != id {
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
	points := application.CalculatePolygonPoints(room)

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
