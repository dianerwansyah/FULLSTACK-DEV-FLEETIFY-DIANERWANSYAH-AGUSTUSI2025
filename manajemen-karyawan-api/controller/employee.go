package controller

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"manajemen-karyawan-api/config"
	"manajemen-karyawan-api/model"
	"manajemen-karyawan-api/utils"

	"github.com/gin-gonic/gin"
)

var allowedEmployeeFields = map[string]string{
	"employeeID":     "e.employee_id",
	"departmentName": "d.departement_name",
	"name":           "e.name",
	"address":        "e.address",
}

// GetAllEmployees godoc
// @Summary Ambil semua karyawan aktif
// @Description Mengembalikan list semua karyawan aktif. Autentikasi via JWT cookie.
// @Tags Employee
// @Produce json
// @Success 200 {array} model.Employee
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/employee [POST]
func GetAllEmployees(c *gin.Context) {
	_, exists := c.Get("employee_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var params utils.QueryParams
	if err := utils.BindJSONStrict(c, &params); err != nil {
		return
	}

	sortSQL := utils.BuildSortSQL(params.SortBy, allowedEmployeeFields)
	pagination := utils.BuildPagination(params.Page, params.PerPage)

	filterSQL, filterArgs := utils.BuildFilterSQL(params.Filter, allowedEmployeeFields)

	// Build query
	query := fmt.Sprintf(`
		SELECT 
		e.id, 
		e.employee_id,
		e.departement_id, 
		d.departement_name, 
		e.name, 
		e.address
		FROM employee e
		JOIN departement d ON e.departement_id = d.id
		WHERE e.deleted_at IS NULL
		%s
		%s
	`, filterSQL, sortSQL)

	var args []interface{}
	args = append(args, filterArgs...)

	if pagination.Use {
		query += " LIMIT ? OFFSET ?"
		args = append(args, pagination.Limit, pagination.Offset)
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		log.Println("Employee query error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch employees"})
		return
	}
	defer rows.Close()

	var result []struct {
		model.Employee
	}
	for rows.Next() {
		var row struct {
			model.Employee
		}
		err := rows.Scan(
			&row.ID, &row.EmployeeID, &row.DepartementID,
			&row.DepartementName,
			&row.Name, &row.Address,
		)
		if err != nil {
			log.Println("Employee scan error:", err)
			continue
		}
		result = append(result, row)
	}
	// Count total with filter
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM employee
		WHERE deleted_at IS NULL
		%s
	`, filterSQL)

	var total int
	err = config.DB.QueryRow(countQuery, filterArgs...).Scan(&total)
	if err != nil {
		log.Println("Employee count error:", err)
		total = 0
	}

	// Response
	meta := utils.BuildMeta(utils.MetaParams{
		Page:    params.Page,
		PerPage: params.PerPage,
		Total:   total,
		SortBy:  params.SortBy,
	})

	c.JSON(http.StatusOK, gin.H{
		"data": result,
		"meta": meta,
	})
}

// GetEmployeeByID godoc
// @Summary Ambil detail karyawan berdasarkan ID
// @Description Mengembalikan detail karyawan berdasarkan ID. Autentikasi via JWT cookie.
// @Tags Employee
// @Produce json
// @Param id path string true "Employee ID"
// @Success 200 {object} model.Employee
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/employee/{id} [get]
func GetEmployeeByID(c *gin.Context) {
	_, exists := c.Get("employee_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	id := c.Param("id")

	var e model.Employee
	query := `
		SELECT id, departement_id, name, address
		FROM employee
		WHERE id = ? AND deleted_at IS NULL
	`
	err := config.DB.QueryRow(query, id).Scan(
		&e.ID, &e.DepartementID, &e.Name, &e.Address,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
		return
	} else if err != nil {
		log.Println("Employee detail error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, e)
}

type EmployeePayload struct {
	EmployeeID    *string `json:"employeeID,omitempty"`
	Name          *string `json:"name,omitempty"`
	DepartementID *string `json:"departmentID,omitempty"`
	Address       *string `json:"address,omitempty"`
}

// CreateEmployee godoc
// @Summary Tambah karyawan baru
// @Description Menambahkan data karyawan ke sistem. Autentikasi via JWT cookie.
// @Tags Employee
// @Accept json
// @Produce json
// @Param payload body EmployeePayload true "Data karyawan"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/employee [post]
func CreateEmployee(c *gin.Context) {
	employeeID, exists := c.Get("employee_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req EmployeePayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	fmt.Println(*req.EmployeeID)

	existsEmp := false
	err := config.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM employee WHERE employee_id = ? AND deleted_at IS NULL
		)
	`, req.EmployeeID).Scan(&existsEmp)

	if err != nil {
		log.Println("Check employee_id error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	if existsEmp {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Employee ID already exists"})
		return
	}

	id := utils.GenerateID()
	now := time.Now()
	pass, err := utils.CreatePassword("password123")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	_, err = config.DB.Exec(`
		INSERT INTO employee (id, employee_id, departement_id, name, address, password, created_at, created_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, id, req.EmployeeID, req.DepartementID, req.Name, req.Address, pass, now, employeeID)

	if err != nil {
		log.Println("Create employee error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create employee"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "employee created", "id": id})
}

// UpdateEmployee godoc
// @Summary Update data karyawan
// @Description Mengubah data karyawan berdasarkan ID. Autentikasi via JWT cookie.
// @Tags Employee
// @Accept json
// @Produce json
// @Param id path string true "ID Karyawan"
// @Param payload body EmployeePayload true "Data karyawan"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/employee/{id} [put]
func UpdateEmployee(c *gin.Context) {
	employeeID, exists := c.Get("employee_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := c.Param("id")

	var req EmployeePayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	payload := map[string]interface{}{}
	if req.Name != nil {
		payload["name"] = *req.Name
	}
	if req.DepartementID != nil {
		payload["departement_id"] = *req.DepartementID
	}
	if req.Address != nil {
		payload["address"] = *req.Address
	}

	// Whitelist fields
	whitelist := []string{"name", "departement_id", "address"}

	// Audit fields
	audit := map[string]interface{}{
		"updated_at": time.Now(),
		"updated_by": employeeID,
	}

	query, args, err := utils.BuildDynamicUpdateQuery("employee", payload, whitelist, audit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	args = append(args, id)

	_, err = config.DB.Exec(query, args...)
	if err != nil {
		log.Println("Update employee error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update employee"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "employee updated"})
}

// DeleteEmployee godoc
// @Summary Hapus karyawan (soft delete)
// @Description Menandai karyawan sebagai terhapus tanpa menghapus data dari database. Autentikasi via JWT cookie.
// @Tags Employee
// @Produce json
// @Param id path string true "ID Karyawan"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/employee/{id} [delete]

func DeleteEmployee(c *gin.Context) {
	employeeID, exists := c.Get("employee_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := c.Param("id")
	now := time.Now()

	_, err := config.DB.Exec(`
		UPDATE employee
		SET deleted_at = ?, deleted_by = ?
		WHERE id = ? AND deleted_at IS NULL
	`, now, employeeID, id)

	if err != nil {
		log.Println("Delete employee error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete employee"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "employee deleted"})
}
