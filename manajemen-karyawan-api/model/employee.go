package model

type Employee struct {
	ID              string `json:"id"`
	EmployeeID      string `json:"employeeID"`
	DepartementID   string `json:"departementID"`
	DepartementName string `json:"departementName"`
	Name            string `json:"name"`
	Password        string `json:"-"`
	Address         string `json:"address"`
	Audit
}
