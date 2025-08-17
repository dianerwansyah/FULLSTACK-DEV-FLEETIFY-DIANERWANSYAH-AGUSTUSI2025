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

var allowedFields = map[string]string{
	"id":              "id",
	"departementName": "departement_name",
	"maxClockInTime":  "max_Clock_In_Time",
	"maxClockOutTime": "max_Clock_Out_Time",
	"createdAt":       "createdAt",
}

// GetAllDepartements godoc
// @Summary Ambil semua departemen
// @Description Mengembalikan list semua departemen aktif. Requires valid JWT cookie named "token"
// @Tags Departement
// @Produce json
// @Success 200 {array} model.Departement
// @Failure 500 {object} map[string]string
// @Router /api/departement/GetData [POST]
func GetAllDepartements(c *gin.Context) {
	_, exists := c.Get("employee_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var params utils.QueryParams
	if err := utils.BindJSONStrict(c, &params); err != nil {
		return
	}

	sortSQL := utils.BuildSortSQL(params.SortBy, allowedFields)
	pagination := utils.BuildPagination(params.Page, params.PerPage)

	filterSQL, filterArgs := utils.BuildFilterSQL(params.Filter, allowedFields)

	// Build query
	query := fmt.Sprintf(`
		SELECT id, departement_name, max_clock_in_time, max_clock_out_time,
		       created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM departement
		WHERE deleted_at IS NULL
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
		log.Println("Departement query error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch departements"})
		return
	}
	defer rows.Close()

	var result []model.Departement

	for rows.Next() {
		var d model.Departement
		var clockInRaw, clockOutRaw string
		err := rows.Scan(
			&d.ID, &d.DepartementName,
			&clockInRaw, &clockOutRaw,
			&d.CreatedAt, &d.CreatedBy, &d.UpdatedAt, &d.UpdatedBy,
			&d.DeletedAt, &d.DeletedBy,
		)
		if err != nil {
			log.Println("Departement scan error:", err)
			continue
		}

		d.MaxClockInTime, _ = time.Parse("15:04:05", clockInRaw)
		d.MaxClockOutTime, _ = time.Parse("15:04:05", clockOutRaw)

		result = append(result, d)
	}

	// Count total with filter
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM departement
		WHERE deleted_at IS NULL
		%s
	`, filterSQL)

	var total int
	err = config.DB.QueryRow(countQuery, filterArgs...).Scan(&total)
	if err != nil {
		log.Println("Departement count error:", err)
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

// GetDepartementByID godoc
// @Summary Ambil detail departemen
// @Description Mengembalikan detail departemen berdasarkan ID. Autentikasi via JWT cookie.
// @Tags Departement
// @Produce json
// @Param id path string true "Departement ID"
// @Success 200 {object} model.Departement
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/departement/{id} [get]
func GetDepartementByID(c *gin.Context) {
	_, exists := c.Get("employee_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := c.Param("id")

	var d model.Departement
	err := config.DB.QueryRow(`
		SELECT id, departement_name, max_clock_in_time, max_clock_out_time,
		       created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM departement
		WHERE id = ? AND deleted_at IS NULL
	`, id).Scan(
		&d.ID, &d.DepartementName, &d.MaxClockInTime, &d.MaxClockOutTime,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "departement not found"})
		return
	} else if err != nil {
		log.Println("Departement detail error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, d)
}

type DepartementPayload struct {
	Name            *string `json:"departementName,omitempty"`
	MaxClockInTime  *string `json:"maxClockInTime,omitempty"`
	MaxClockOutTime *string `json:"maxClockOutTime,omitempty"`
}

// CreateDepartement godoc
// @Summary Tambah departemen baru
// @Description Menambahkan data departemen ke sistem. Hanya dapat diakses oleh user dengan role admin. Autentikasi via JWT cookie.
// @Tags Departement
// @Accept json
// @Produce json
// @Param payload body DepartementPayload true "Data departemen"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/departement [post]
func CreateDepartement(c *gin.Context) {
	employeeID, exists := c.Get("employee_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req DepartementPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	if req.Name == nil || req.MaxClockInTime == nil || req.MaxClockOutTime == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "all fields are required"})
		return
	}

	id := utils.GenerateID()
	now := time.Now()

	_, err := config.DB.Exec(`
		INSERT INTO departement (id, departement_name, max_clock_in_time, max_clock_out_time, created_at, created_by)
		VALUES (?, ?, ?, ?, ?, ?)
	`, id, req.Name, req.MaxClockInTime, req.MaxClockOutTime, now, employeeID)

	if err != nil {
		log.Println("Create departement error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create departement"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "departement created", "id": id})
}

// UpdateDepartement godoc
// @Summary Update data departemen
// @Description Mengubah data departemen berdasarkan ID. Autentikasi via JWT cookie.
// @Tags Departement
// @Accept json
// @Produce json
// @Param id path string true "ID Departemen"
// @Param payload body DepartementPayload true "Data departemen"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/departement/{id} [put]
func UpdateDepartement(c *gin.Context) {
	employeeID, exists := c.Get("employee_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := c.Param("id")
	var req DepartementPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	payload := map[string]interface{}{}
	if req.Name != nil {
		payload["departement_name"] = *req.Name
	}
	if req.MaxClockInTime != nil {
		payload["max_clock_in_time"] = *req.MaxClockInTime
	}
	if req.MaxClockOutTime != nil {
		payload["max_clock_out_time"] = *req.MaxClockOutTime
	}

	// Whitelist fields
	whitelist := []string{"departement_name", "max_clock_in_time", "max_clock_out_time"}

	// Audit fields
	audit := map[string]interface{}{
		"updated_at": time.Now(),
		"updated_by": employeeID,
	}

	query, args, err := utils.BuildDynamicUpdateQuery("departement", payload, whitelist, audit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	args = append(args, id)

	_, err = config.DB.Exec(query, args...)
	if err != nil {
		log.Println("Update departement error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update departement"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "departement updated"})
}

// DeleteDepartement godoc
// @Summary Hapus departemen (soft delete)
// @Description Menandai departemen sebagai terhapus tanpa menghapus data dari database. Autentikasi via JWT cookie.
// @Tags Departement
// @Produce json
// @Param id path string true "ID Departemen"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/departement/{id} [delete]
func DeleteDepartement(c *gin.Context) {
	userID, exists := c.Get("employee_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := c.Param("id")
	now := time.Now()

	_, err := config.DB.Exec(`
		UPDATE departement
		SET deleted_at = ?, deleted_by = ?
		WHERE id = ? AND deleted_at IS NULL
	`, now, userID, id)

	if err != nil {
		log.Println("Delete departement error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete departement"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "departement deleted"})
}
