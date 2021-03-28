package store

import (
	"github.com/jordy2254/indoormaprestapi/pkg/model"
	"gorm.io/gorm"
)

type RoomStore struct {
	DB *gorm.DB
}

func NewRoomStore(DB *gorm.DB) RoomStore {
	return RoomStore{DB: DB}
}

func (rs *RoomStore) CreateRoom(room *model.Room) {
	rs.DB.Create(room)
}

func (rs *RoomStore) GetRoomById(id int) model.Room {
	var room model.Room
	rs.DB.Preload("Indents").Preload("Entrances").Find(&room, "rooms.id=?", id)
	return room
}

func (rs *RoomStore) UpdateRoom(room model.Room) {
	rs.DB.Model(&room).Omit("Indents").UpdateColumns(&room)
}