package store

import (
	"github.com/jordy2254/indoormaprestapi/pkg/model"
	"gorm.io/gorm"
)

type BuildingStore struct {
	DB *gorm.DB
}

func NewBuildingStore(DB *gorm.DB) BuildingStore {
	return BuildingStore{DB: DB}
}

func (bs BuildingStore) GetBuildingById(mapId, id int) model.Building {
	var building model.Building
	bs.DB.Preload("Floors").
		Preload("Floors.Sensors").
		Preload("Floors.Rooms").
		Preload("Floors.Rooms.Indents").
		Find(&building, "buildings.id=? AND buildings.map_id=?", id, mapId)
	return building
}

func (bs BuildingStore) CreateBuilding(building *model.Building) {
	bs.DB.Create(building)
}

func (bs BuildingStore) UpdateBuilding(building *model.Building) {
	bs.DB.Model(&building).Omit("Floors").UpdateColumns(&building)
}

func (bs BuildingStore) DeleteBuilding(m *model.Building) {

}
