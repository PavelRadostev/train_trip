package model

import "time"

type Passport struct {
	ID                 int        `cbor:"id"`
	TrainID            int        `cbor:"train_id"`
	BeginTime          time.Time  `cbor:"begin_time"`
	EndTime            *time.Time `cbor:"end_time,omitempty"`
	LoadTypeID         int        `cbor:"load_type_id"`
	WorkPlaceID        *int       `cbor:"work_type_id,omitempty"`
	Volume             float64    `cbor:"volume"`
	Weight             float64    `cbor:"weight"`
	CarriageNum        *int       `cbor:"carriage_num,omitempty"`
	ReducedCarriageNum *float64   `cbor:"reduced_carriage_num,omitempty"`
}
