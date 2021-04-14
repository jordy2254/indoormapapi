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

type IndentController struct {
	indentStore *store.IndentStore
	logger *logging.Logger
}

func AddIndentAPI(rh *RouteHelper, indentStore *store.IndentStore, logger *logging.Logger) {
	controller := IndentController{indentStore: indentStore, logger: logger}
	rh.protectedRoute("/Indents/{id}", controller.getIndent).Methods("GET")
	rh.protectedRoute("/Indents/{id}", controller.deleteIndent).Methods("DELETE")
	rh.protectedRoute("/Indents/{id}", controller.updateIndent).Methods("POST")
	rh.protectedRoute("/Indents", controller.createIndent).Methods("POST")
}


func (ic *IndentController) createIndent(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Create Indent")

	newIndent, err := unmarshalIndentRequest(w, r)
	if err != nil {
		return
	}
	ic.indentStore.CreateIndent(&newIndent)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newIndent)
}

func (ic *IndentController) updateIndent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println("Endpoint Hit: Update Indent")
	id, err := strconv.Atoi(params["id"])
	newIndent, err := unmarshalIndentRequest(w, r)
	if err != nil || newIndent.Id != id {
		return
	}
	ic.indentStore.UpdateIndent(newIndent)
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(newIndent)
}

func (ic *IndentController) getIndent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println("Endpoint Hit: getIndent")
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return
	}
	json.NewEncoder(w).Encode(ic.indentStore.GetIndentById(id))
}

func (ic *IndentController) deleteIndent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println("Endpoint Hit: getIndent")
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return
	}
	ic.indentStore.DeleteIndent(id)
	json.NewEncoder(w).Encode(ic.indentStore.GetIndentById(id))
}

func unmarshalIndentRequest(w http.ResponseWriter, r *http.Request) (model.Indent, error) {
	var newIndent = model.Indent{}

	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return newIndent, err
	}

	err = json.Unmarshal(reqBody, &newIndent)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return newIndent, err
	}

	return newIndent, nil
}