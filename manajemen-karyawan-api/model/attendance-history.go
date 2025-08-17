package model

import "time"

type AttendanceHistory struct {
	ID             string    `json:"id"`
	EmployeeID     string    `json:"employeeID"`
	AttendanceID   string    `json:"attendanceID"`
	DateAttendance time.Time `json:"dateAttendance"`
	AttendanceType int       `json:"attendanceType"`
	Description    string    `json:"description"`
	Audit
}
