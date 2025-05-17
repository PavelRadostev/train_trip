package model

import "time"

type TrainLoad struct {
	ID                  int        `cbor:"id"`
	IsDeleted           bool       `cbor:"is_deleted"`
	LoadArriveTime      time.Time  `cbor:"load_arrive_time"`
	LoadBeginTime       time.Time  `cbor:"load_begin_time"`
	LoadEndTime         *time.Time `cbor:"load_end_time,omitempty"`
	LoadDepartTime      *time.Time `cbor:"load_depart_time,omitempty"`
	TrainID             int        `cbor:"train_id"`
	GeometryID          *int       `cbor:"geometry_id,omitempty"`
	UnloadID            int        `cbor:"unload_id"`
	ShovelID            int        `cbor:"shovel_id"`
	Manual              bool       `cbor:"manual"`
	LoadTypeIDManual    *int       `cbor:"load_type_id_manual,omitempty"`
	WorkTypeIDManual    *int       `cbor:"work_type_id_manual,omitempty"`
	LoadBeginTimeManual *time.Time `cbor:"load_begin_time_manual,omitempty"`
	LoadEndTimeManual   *time.Time `cbor:"load_end_time_manual,omitempty"`
	UnloadIDManual      *int       `cbor:"unload_id_manual,omitempty"`
	ShovelIDManual      *int       `cbor:"shovel_id_manual,omitempty"`
	VolumeManual        *float64   `cbor:"volume_manual,omitempty"`
	CycleIDs            []int      `cbor:"cycle_ids,omitempty"`
	IsCured             bool       `cbor:"is_cured"`
	CarriageNum         *float64   `cbor:"carriage_num,omitempty"`
	Source              *int16     `cbor:"source,omitempty"`
}
