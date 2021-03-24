package store

import (
	"gorm.io/gorm"
	"github.com/jordy2254/indoormaprestapi/model"
)

type RoomStore struct {
	DB *gorm.DB
}

func NewRoomStore(DB *gorm.DB) RoomStore{
	return RoomStore{DB: DB}
}

func (rs *RoomStore) CreateRoom(room *model.Room) {
	rs.DB.Create(room)
}

func (rs *RoomStore) GetRoomById(id int) model.Room {
	var room model.Room
	rs.DB.Preload("Indents").Find(&room, "rooms.id=?", id)
	return room
}

func (rs *RoomStore) UpdateRoom(room model.Room) {
	rs.DB.Model(&room).Omit("Indents").UpdateColumns(&room)
}