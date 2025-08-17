package model

import "time"

type AttendanceItem struct {
	ID              string    `json:"id"`
	EmployeeID      string    `json:"employeeID"`
	EmployeeName    string    `json:"employeeName"`
	DepartementName string    `json:"departementName"`
	Clock           time.Time `json:"clock"`
	MaxClock        string    `json:"maxClock"`
	DateAttendance  time.Time `json:"dateAttendance"`
	Desc            string    `json:"description"`
	Status          string    `json:"status"`
	AttendanceType  string    `json:"attendanceType"`
}
