package model

import "time"

type TrainLoad struct {
	ID                  int        `json:"id"`
	IsDeleted           bool       `json:"is_deleted"`
	LoadArriveTime      time.Time  `json:"load_arrive_time"`
	LoadBeginTime       time.Time  `json:"load_begin_time"`
	LoadEndTime         *time.Time `json:"load_end_time,omitempty"`
	LoadDepartTime      *time.Time `json:"load_depart_time,omitempty"`
	TrainID             int        `json:"train_id"`
	GeometryID          *int       `json:"geometry_id,omitempty"`
	UnloadID            int        `json:"unload_id"`
	ShovelID            int        `json:"shovel_id"`
	Manual              bool       `json:"manual"`
	LoadTypeIDManual    *int       `json:"load_type_id_manual,omitempty"`
	WorkTypeIDManual    *int       `json:"work_type_id_manual,omitempty"`
	LoadBeginTimeManual *time.Time `json:"load_begin_time_manual,omitempty"`
	LoadEndTimeManual   *time.Time `json:"load_end_time_manual,omitempty"`
	UnloadIDManual      *int       `json:"unload_id_manual,omitempty"`
	ShovelIDManual      *int       `json:"shovel_id_manual,omitempty"`
	VolumeManual        *float64   `json:"volume_manual,omitempty"`
	CycleIDs            []int      `json:"cycle_ids,omitempty"`
	IsCured             bool       `json:"is_cured"`
	CarriageNum         *float64   `json:"carriage_num,omitempty"`
	Source              *int16     `json:"source,omitempty"`
}
