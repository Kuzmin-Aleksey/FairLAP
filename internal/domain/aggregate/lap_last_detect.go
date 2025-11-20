package aggregate

import "time"

type LapLastDetect struct {
	LapId      int       `json:"lap_id" db:"lap_id"`
	LastDetect time.Time `json:"last_detect" db:"last_detect"`
}
