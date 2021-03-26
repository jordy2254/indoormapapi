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

func (ps *PathStore) CreateNode(node *model.MapNode){
	ps.DB.Create(node)
}

func (ps *PathStore) LinkNodesById(mapId, node1Id, node2Id int){
	nodeLink := model.NodeEdge{
		MapId:   mapId,
		Node1Id: node1Id,
		Node2Id: node2Id,
	}
	ps.DB.Create(&nodeLink)
}

func (ps *PathStore) DeleteNode(nodeId int) error {
	return ps.DB.Delete(&model.MapNode{}, nodeId).Error
}
func (ps *PathStore) DeleteLink(mapId, id1, id2 int){
	ids := []int{id1, id2}

	ps.DB.Where("map_id=? AND (node1_id IN ?  AND node2_id IN ?)", mapId, ids, ids).
		Delete(&model.NodeEdge{})
}