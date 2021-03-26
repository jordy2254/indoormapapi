package handlers

import (
	"github.com/jordy2254/indoormaprestapi/pkg/gorm/store"
	"github.com/op/go-logging"
)

type SensorController struct {
	sensorStore *store.SensorStore
}

func AddSensorAPI(rh *RouteHelper, sensorStore *store.SensorStore, logger *logging.Logger) {
	controller := SensorController{sensorStore: sensorStore}
	_ = controller
}
