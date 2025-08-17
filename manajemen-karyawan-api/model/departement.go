package model

import "time"

type Departement struct {
	ID              string    `json:"id"`
	DepartementName string    `json:"departementName"`
	MaxClockInTime  time.Time `json:"maxClockInTime"`
	MaxClockOutTime time.Time `json:"maxClockOutTime"`
	Audit
}
