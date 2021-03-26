package store

import (
	"github.com/jordy2254/indoormaprestapi/pkg/model"
	"gorm.io/gorm"
)

type SensorStore struct {
	DB *gorm.DB
}

func NewSensorStore(DB *gorm.DB) SensorStore {
	return SensorStore{DB: DB}
}

func (ss *SensorStore) createSensor(sensor model.Sensor){

}

func (ss *SensorStore) updateSensor(sensor model.Sensor){

}

func (ss *SensorStore) getSensorById(id int){

}
