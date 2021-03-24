package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jordy2254/indoormaprestapi/interfaces/gorm/store"
	"github.com/jordy2254/indoormaprestapi/model"
	"io/ioutil"
	"net/http"
	"strconv"
)

type MapsController struct {
	MapStore *store.MapStore
}

type SyncMapRequest struct {
	MapKey  string `json:"mapKey"`
	MapPass string `json:"password"`
}

func AddMapAPI(router *mux.Router, mapStore *store.MapStore) {
	controller := MapsController{MapStore: mapStore}

	router.HandleFunc("/maps/sync", controller.sync).Methods("POST")
	router.HandleFunc("/maps/{id}", controller.delete).Methods("DELETE")
	router.HandleFunc("/maps/{id}", controller.get).Methods("GET")
	router.HandleFunc("/maps/{id}", controller.update).Methods("POST")
	router.HandleFunc("/maps", controller.create).Methods("POST")
	router.HandleFunc("/maps", controller.userMaps).Methods("GET")
}

func (mc *MapsController) sync(wr http.ResponseWriter, req *http.Request) {
	var syncRequest = SyncMapRequest{}

	reqBody, err := ioutil.ReadAll(req.Body)

	if err != nil {
		http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(reqBody, &syncRequest)

	if err != nil {
		http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	id, err := strconv.Atoi(syncRequest.MapKey)
	if err != nil {
		http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	m := mc.MapStore.GetMapById(id)

	if m.Password != syncRequest.MapPass {
		http.Error(wr, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	json.NewEncoder(wr).Encode(m)
}

func (mc *MapsController) delete(wr http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	mc.MapStore.DeleteMap(id)
	wr.WriteHeader(http.StatusOK)
}

func (mc *MapsController) update(wr http.ResponseWriter, req *http.Request) {
	fmt.Println("Endpoint Hit: Update Map")
	params := mux.Vars(req)
	id, err := strconv.Atoi(params["id"])
	newMap, err := unmarshalMapRequest(wr, req)
	if err != nil || newMap.Id != id {
		return
	}
	mc.MapStore.UpdateMap(newMap)
	wr.WriteHeader(http.StatusAccepted)
	json.NewEncoder(wr).Encode(mc.MapStore.GetMapById(newMap.Id))
}

func (mc *MapsController) create(wr http.ResponseWriter, req *http.Request) {
	fmt.Println("Endpoint Hit: Create Map")
	user := req.Context().Value("user")
	userId := user.(*jwt.Token).Claims.(jwt.MapClaims)["sub"]
	oAuthUser := mc.MapStore.GetOAuthUserBySub(userId.(string))

	emptyRequest, err := checkForEmptyRequestBody(req)
	if err != nil {
		return
	}
	var newMap model.Map = model.Map{Name: "Unnamed"}

	if !emptyRequest {
		unmarshalled, err := unmarshalMapRequest(wr, req)
		if err != nil {
			return
		}
		newMap = unmarshalled
	}
	if err != nil {
		return
	}
	newMap.Users = append(newMap.Users, oAuthUser)
	mc.MapStore.CreateMap(&newMap)
	wr.WriteHeader(http.StatusCreated)
	json.NewEncoder(wr).Encode(mc.MapStore.GetMapById(newMap.Id))
}

func (mc *MapsController) userMaps(wr http.ResponseWriter, req *http.Request) {
	fmt.Println("Endpoint Hit: UserMaps")
	user := req.Context().Value("user")
	userId := user.(*jwt.Token).Claims.(jwt.MapClaims)["sub"]
	oAuthUser := mc.MapStore.GetOAuthUserBySub(userId.(string))
	json.NewEncoder(wr).Encode(mc.MapStore.GetMapsByUserId(oAuthUser.Id))
}

func (mc *MapsController) get(wr http.ResponseWriter, req *http.Request) {
	fmt.Println("Endpoint Hit: Getmap")
	params := mux.Vars(req)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return
	}
	json.NewEncoder(wr).Encode(mc.MapStore.GetMapById(id))
}

func checkForEmptyRequestBody(r *http.Request) (bool, error){
	reqBody, err := ioutil.ReadAll(r.Body)
	if(err != nil){
		return false, err
	}

	return len(reqBody) == 0, nil
}

func unmarshalMapRequest(w http.ResponseWriter, r *http.Request) (model.Map, error) {
	var newMap = model.Map{}

	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return newMap, err
	}

	err = json.Unmarshal(reqBody, &newMap)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return newMap, err
	}

	return newMap, nil
}