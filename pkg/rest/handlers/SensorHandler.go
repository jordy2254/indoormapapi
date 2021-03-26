package handlers

import "github.com/jordy2254/indoormaprestapi/pkg/gorm/store"

type SensorController struct {
	sensorStore *store.SensorStore
}

func AddSensorAPI(rh *RouteHelper, sensorStore *store.SensorStore) {
	controller := SensorController{sensorStore: sensorStore}
	_ = controller
}
