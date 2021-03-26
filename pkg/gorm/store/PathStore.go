package store

import (
	"github.com/jordy2254/indoormaprestapi/pkg/model"
	"gorm.io/gorm"
)

type PathStore struct {
	DB *gorm.DB
}

func NewPathStore(DB *gorm.DB) PathStore {
	return PathStore{DB: DB}
}

func (ps *PathStore) getNodeById(id int){

}

func (ps *PathStore) createNode(node model.MapNode){

}

func (ps *PathStore) linkNodesById(id1, id2 int){

}