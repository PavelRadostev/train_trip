package entity

import "time"

// TrainPassport - паспорта поезда
type TrainPassport struct {
	TrainId           int        `json:"train_id"`
	Begin             time.Time  `json:"begin"`
	End               *time.Time `json:"end,omitempty"`
	CargoId           int        `json:"cargo_id"`
	StdWeight         float64    `json:"std_weight"`
	StdVolume         float64    `json:"std_volume"`
	CarrigeNum        float64    `json:"carrige_num"`
	ReducedCarrigeNum float64    `json:"reduced_carrige_num"`
}
