package store

import (
	"gorm.io/gorm"
	"github.com/jordy2254/indoormaprestapi/model"
)

type FloorStore struct {
	DB *gorm.DB
}

func NewFloorStore(DB *gorm.DB) FloorStore{
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