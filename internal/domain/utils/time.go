package domain

import (
	"time"
)

// Interval - интервал времени
type Interval struct {
	TimeFrom time.Time `cbor:"time_from"`
	TimeTo   time.Time `cbor:"time_to"`
}
