package store

import (
	"github.com/jordy2254/indoormaprestapi/pkg/model"
	"gorm.io/gorm"
)

type FloorStore struct {
	DB *gorm.DB
}

func NewFloorStore(DB *gorm.DB) FloorStore {
	return FloorStore{DB: DB}
}

func (fs *FloorStore) CreateFloor(floor *model.Floor) {
	fs.DB.Create(floor)
}

func (fs *FloorStore) GetFloorById(id int) model.Floor {
	var floor model.Floor
	fs.DB.Preload("Sensors").
		Preload("Rooms").
		Preload("Rooms.Indents").
		Find(&floor, "floors.id=?", id)
	return floor
}

func (fs *FloorStore) UpdateFloor(floor model.Floor) {
	fs.DB.Model(&floor).Omit("Rooms").UpdateColumns(&floor)
}