package model

import "time"

type Attendance struct {
	ID         string     `json:"id"`
	EmployeeID string     `json:"employeeID"`
	ClockIn    *time.Time `json:"clockIn"`
	ClockOut   *time.Time `json:"clockOut"`
	Audit
}
