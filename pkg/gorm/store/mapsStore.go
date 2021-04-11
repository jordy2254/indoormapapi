package store

import (
	"github.com/jordy2254/indoormaprestapi/pkg/model"
	"gorm.io/gorm"
)

type MapStore struct {
	DB *gorm.DB
}

func NewMapStore(DB *gorm.DB) MapStore {
	return MapStore{DB: DB}
}
func (mapStore *MapStore) CreateMap(ma *model.Map){
	mapStore.DB.Create(ma)
}

func (mapStore *MapStore) GetMapById(id int) model.Map {
	var ma model.Map
	mapStore.DB.Preload("Nodes").
		Preload("Edges").
		Preload("Buildings").
		Preload("Buildings.Floors").
		Preload("Buildings.Floors.Rooms").
		Preload("Buildings.Floors.Rooms.Indents").
		Preload("Buildings.Floors.Sensors").
		Preload("Buildings.Floors.Rooms.Entrances").
		Find(&ma, "maps.id=?", id)
	return ma
}

func (mapStore *MapStore) DeleteMap(id int){
	mapStore.DB.Delete(&model.Map{}, id)
}

func (mapStore *MapStore) GetMapsByUserId(id int) []model.Map {
	var auth0User model.Auth0User
	mapStore.DB.
		Preload("Maps").
		Find(&auth0User, id)

	return auth0User.Maps
}

func (mapStore *MapStore) GetOAuthUserBySub(sub string) model.Auth0User {
	var auth0User model.Auth0User

	mapStore.DB.Find(&auth0User, "auth0_users.authid=?", sub)
	if auth0User.Authid == "" {
		return createAuth0User(mapStore, sub)
	}
	return auth0User
}

func createAuth0User(mapStore *MapStore, sub string) model.Auth0User {
	newUser := model.Auth0User{
		Authid: sub,
		Maps:   []model.Map{},
	}
	mapStore.DB.Create(&newUser)
	return newUser
}

func (mapStore *MapStore) UpdateMap(ma model.Map) {
	mapStore.DB.Model(&ma).
		Omit("Buildings").
		Omit("RootNode").
		UpdateColumns(&ma)
}
