package model

import "time"

type TrainUnload struct {
	ID                    int        `json:"id"`
	IsDeleted             bool       `json:"is_deleted"`
	UnloadArriveTime      *time.Time `json:"unload_arrive_time,omitempty"`
	UnloadBeginTime       time.Time  `json:"unload_begin_time"`
	UnloadEndTime         *time.Time `json:"unload_end_time,omitempty"`
	UnloadDepartTime      *time.Time `json:"unload_depart_time,omitempty"`
	TrainID               int        `json:"train_id"`
	GeometryID            *int       `json:"geometry_id,omitempty"`
	UnloadID              *int       `json:"unload_id,omitempty"`
	UnloadBeginTimeButton *time.Time `json:"unload_begin_time_button,omitempty"`
	Manual                bool       `json:"manual"`
	UnloadBeginTimeManual *time.Time `json:"unload_begin_time_manual,omitempty"`
	UnloadEndTimeManual   *time.Time `json:"unload_end_time_manual,omitempty"`
	UnloadIDManual        *int       `json:"unload_id_manual,omitempty"`
	CarriageNumManual     *float64   `json:"carriage_num_manual,omitempty"`
	ManualVolume          *float64   `json:"manual_volume,omitempty"`
}
