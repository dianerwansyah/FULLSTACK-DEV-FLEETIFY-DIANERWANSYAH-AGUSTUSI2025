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

type ClockRequest struct {
	Type        string `json:"type" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// ClockIn godoc
// @Summary Clock-in karyawan
// @Description Menyimpan waktu clock-in untuk karyawan yang login (via JWT cookie)
// @Tags Attendance
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/attendance [POST]
func ClockHandler(c *gin.Context) {
	employeeID := c.GetString("employee_id")

	var req ClockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type and description required"})
		return
	}

	now := time.Now()

	var attendanceID string
	var err error

	if req.Type == "clock_in" {
		// Check if already clocked in today
		err = config.DB.QueryRow(`
			SELECT id FROM attendance
			WHERE employee_id = ? AND DATE(clock_in) = CURDATE() AND deleted_at IS NULL
		`, employeeID).Scan(&attendanceID)

		if err != sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "already clocked in today"})
			return
		}

		attendanceID = utils.GenerateID()

		tx, err := config.DB.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction error"})
			return
		}

		_, err = tx.Exec(`
			INSERT INTO attendance (id, employee_id, clock_in, created_at, created_by)
			VALUES (?, ?, ?, ?, ?)
		`, attendanceID, employeeID, now, now, employeeID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to clock in"})
			return
		}

		_, err = tx.Exec(`
			INSERT INTO attendance_history (id, employee_id, attendance_id, date_attendance, attendance_type, description, created_at, created_by)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, utils.GenerateID(), employeeID, attendanceID, now, 1, req.Description, now, employeeID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save history"})
			return
		}

		tx.Commit()
		c.JSON(http.StatusOK, gin.H{"message": "clock-in successful"})
		return
	}

	if req.Type == "clock_out" {
		// Get today's attendance record
		err = config.DB.QueryRow(`
			SELECT id FROM attendance
			WHERE employee_id = ? AND DATE(clock_in) = CURDATE() AND deleted_at IS NULL
		`, employeeID).Scan(&attendanceID)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no clock-in record found"})
			return
		}

		tx, err := config.DB.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction error"})
			return
		}

		_, err = tx.Exec(`
			UPDATE attendance SET clock_out = ?, updated_at = ?, updated_by = ?
			WHERE id = ?
		`, now, now, employeeID, attendanceID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to clock out"})
			return
		}

		_, err = tx.Exec(`
			INSERT INTO attendance_history (id, employee_id, attendance_id, date_attendance, attendance_type, description, created_at, created_by)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, utils.GenerateID(), employeeID, attendanceID, now, 2, req.Description, now, employeeID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save history"})
			return
		}

		tx.Commit()
		c.JSON(http.StatusOK, gin.H{"message": "clock-out successful"})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid type"})
}

var allowedAttendanceFields = map[string]string{
	"dateAttendance":      "h.date_attendance",
	"employeeName":        "e.name",
	"departementID":       "e.departement_id",
	"attendanceType ":     "h.attendance_type ",
	"description":         "h.description",
	"date_attendance.gte": "date_attendance.gte",
}

// GetAttendanceLogs godoc
// @Summary List log absensi karyawan yang login
// @Description Menampilkan log absensi milik karyawan yang sedang login, berdasarkan tanggal dan departemen. Autentikasi via JWT cookie.
// @Tags Attendance
// @Produce json
// @Param date query string false "Tanggal (YYYY-MM-DD)"
// @Param departement_id query string false "ID Departemen"
// @Success 200 {object} model.AttendanceItem
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/attendance/logs [POST]
func GetAttendanceLogs(c *gin.Context) {
	employeeID, exists := c.Get("employee_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var params utils.QueryParams
	if err := utils.BindJSONStrict(c, &params); err != nil {
		return
	}

	sortSQL := utils.BuildSortSQL(params.SortBy, allowedAttendanceFields)
	pagination := utils.BuildPagination(params.Page, params.PerPage)
	filterSQL, filterArgs := utils.BuildFilterSQL(params.Filter, allowedAttendanceFields)

	query := fmt.Sprintf(`
		SELECT 
			a.id,
			a.employee_id,
			e.name AS employee_name,
			d.departement_name,
			a.clock_in,
			a.clock_out,
			d.max_clock_in_time,
			d.max_clock_out_time,
			h.date_attendance,
			h.attendance_type,
			h.description
		FROM attendance a
		JOIN employee e ON a.employee_id = e.employee_id
		JOIN departement d ON e.departement_id = d.id
		JOIN attendance_history h ON h.attendance_id = a.id
		WHERE a.deleted_at IS NULL AND a.employee_id = ?
		%s
		%s
	`, filterSQL, sortSQL)

	args := append([]interface{}{employeeID}, filterArgs...)
	if pagination.Use {
		query += " LIMIT ? OFFSET ?"
		args = append(args, pagination.Limit, pagination.Offset)
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		log.Println("Attendance query error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch attendance logs"})
		return
	}
	defer rows.Close()

	loc, _ := time.LoadLocation("Asia/Singapore")
	var logs []model.AttendanceItem

	for rows.Next() {
		var (
			id              string
			employeeID      string
			employeeName    string
			departementName string
			clockIn         time.Time
			clockOut        time.Time
			maxInRaw        string
			maxOutRaw       string
			dateAttendance  time.Time
			attendanceType  int
			description     string
		)

		err := rows.Scan(
			&id,
			&employeeID,
			&employeeName,
			&departementName,
			&clockIn,
			&clockOut,
			&maxInRaw,
			&maxOutRaw,
			&dateAttendance,
			&attendanceType,
			&description,
		)
		if err != nil {
			log.Println("Attendance scan error:", err)
			continue
		}

		var item model.AttendanceItem
		item.ID = id
		item.EmployeeID = employeeID
		item.EmployeeName = employeeName
		item.DepartementName = departementName
		item.DateAttendance = dateAttendance
		item.Desc = description

		if attendanceType == 1 {
			item.Clock = clockIn
			item.MaxClock = maxInRaw
			item.AttendanceType = "in"

			isLate, err := utils.IsLate(clockIn, maxInRaw, loc)
			if err != nil {
				item.Status = "Unknown"
			} else if isLate {
				item.Status = "Terlambat"
			} else {
				item.Status = "Tepat"
			}
		} else if attendanceType == 2 {
			item.Clock = clockOut
			item.MaxClock = maxOutRaw
			item.AttendanceType = "out"

			isEarly, err := utils.IsEarly(clockOut, maxOutRaw, loc)
			if err != nil {
				item.Status = "Unknown"
			} else if isEarly {
				item.Status = "Pulang Cepat"
			} else {
				item.Status = "Tepat"
			}
		}

		logs = append(logs, item)
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM attendance a
		JOIN employee e ON a.employee_id = e.employee_id
		JOIN departement d ON e.departement_id = d.id
		JOIN attendance_history h ON h.attendance_id = a.id
		WHERE a.deleted_at IS NULL AND a.employee_id = ?
		%s
	`, filterSQL)

	var total int
	countArgs := append([]interface{}{employeeID}, filterArgs...)
	err = config.DB.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		log.Println("Attendance count error:", err)
		total = 0
	}

	meta := utils.BuildMeta(utils.MetaParams{
		Page:    params.Page,
		PerPage: params.PerPage,
		Total:   total,
		SortBy:  params.SortBy,
	})

	c.JSON(http.StatusOK, gin.H{
		"data": logs,
		"meta": meta,
	})
}

// GetAllAttendanceLogs godoc
// @Summary List semua log absensi karyawan
// @Description Menampilkan seluruh data absensi karyawan, bisa difilter berdasarkan tanggal dan departemen. Hanya bisa diakses oleh user dengan role tertentu.
// @Tags Attendance
// @Produce json
// @Param date query string false "Tanggal (YYYY-MM-DD)"
// @Param departement_id query string false "ID Departemen"
// @Success 200 {array} model.AttendanceItem
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/attendance/GetData [POST]
func GetAllAttendanceLogs(c *gin.Context) {
	var params utils.QueryParams
	if err := utils.BindJSONStrict(c, &params); err != nil {
		return
	}

	sortSQL := utils.BuildSortSQL(params.SortBy, allowedAttendanceFields)
	pagination := utils.BuildPagination(params.Page, params.PerPage)
	filterSQL, filterArgs := utils.BuildFilterSQL(params.Filter, allowedAttendanceFields)

	query := fmt.Sprintf(`
		SELECT 
			a.id,
			a.employee_id,
			e.name AS employee_name,
			d.departement_name,
			a.clock_in,
			a.clock_out,
			d.max_clock_in_time,
			d.max_clock_out_time,
			h.date_attendance,
			h.attendance_type,
			h.description
		FROM attendance a
		JOIN employee e ON a.employee_id = e.employee_id
		JOIN departement d ON e.departement_id = d.id
		JOIN attendance_history h ON h.attendance_id = a.id
		WHERE a.deleted_at IS NULL
		%s
		%s
	`, filterSQL, sortSQL)

	args := filterArgs
	if pagination.Use {
		query += " LIMIT ? OFFSET ?"
		args = append(args, pagination.Limit, pagination.Offset)
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		log.Println("Attendance query error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch attendance logs"})
		return
	}
	defer rows.Close()

	loc, _ := time.LoadLocation("Asia/Singapore")
	var logs []model.AttendanceItem

	for rows.Next() {
		var (
			id              string
			employeeID      string
			employeeName    string
			departementName string
			clockIn         time.Time
			clockOut        time.Time
			maxInRaw        string
			maxOutRaw       string
			dateAttendance  time.Time
			attendanceType  int
			description     string
		)

		err := rows.Scan(
			&id,
			&employeeID,
			&employeeName,
			&departementName,
			&clockIn,
			&clockOut,
			&maxInRaw,
			&maxOutRaw,
			&dateAttendance,
			&attendanceType,
			&description,
		)
		if err != nil {
			log.Println("Attendance scan error:", err)
			continue
		}

		var item model.AttendanceItem
		item.ID = id
		item.EmployeeID = employeeID
		item.EmployeeName = employeeName
		item.DepartementName = departementName
		item.DateAttendance = dateAttendance
		item.Desc = description

		if attendanceType == 1 {
			item.Clock = clockIn
			item.MaxClock = maxInRaw
			item.AttendanceType = "in"

			isLate, err := utils.IsLate(clockIn, maxInRaw, loc)
			if err != nil {
				item.Status = "Unknown"
			} else if isLate {
				item.Status = "Terlambat"
			} else {
				item.Status = "Tepat"
			}
		} else if attendanceType == 2 {
			item.Clock = clockOut
			item.MaxClock = maxOutRaw
			item.AttendanceType = "out"

			isEarly, err := utils.IsEarly(clockOut, maxOutRaw, loc)
			if err != nil {
				item.Status = "Unknown"
			} else if isEarly {
				item.Status = "Pulang Cepat"
			} else {
				item.Status = "Tepat"
			}
		}

		logs = append(logs, item)
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM attendance a
		JOIN employee e ON a.employee_id = e.employee_id
		JOIN departement d ON e.departement_id = d.id
		JOIN attendance_history h ON h.attendance_id = a.id
		WHERE a.deleted_at IS NULL
		%s
	`, filterSQL)

	var total int
	err = config.DB.QueryRow(countQuery, filterArgs...).Scan(&total)
	if err != nil {
		log.Println("Attendance count error:", err)
		total = 0
	}

	meta := utils.BuildMeta(utils.MetaParams{
		Page:    params.Page,
		PerPage: params.PerPage,
		Total:   total,
		SortBy:  params.SortBy,
	})

	c.JSON(http.StatusOK, gin.H{
		"data": logs,
		"meta": meta,
	})
}

var result struct {
	ClockIn  *time.Time `json:"clock_in"`
	ClockOut *time.Time `json:"clock_out"`
	Location *string    `json:"location"`
}

// GetTodayAttendance godoc
// @Summary Ambil data absensi hari ini milik user yang login
// @Description Menampilkan data absensi karyawan yang sedang login untuk tanggal hari ini. Autentikasi via JWT cookie.
// @Tags Attendance
// @Produce json
// @Success 200 {object} model.AttendanceItem
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/attendance/today [get]
func GetTodayAttendance(c *gin.Context) {
	employeeID, exists := c.Get("employee_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get today's date range
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := `
		SELECT 
			clock_in, clock_out
		FROM attendance
		WHERE employee_id = ?
		AND created_at >= ? AND created_at < ?
		LIMIT 1
	`

	err := config.DB.QueryRow(query, employeeID, startOfDay, endOfDay).Scan(
		&result.ClockIn,
		&result.ClockOut,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, gin.H{})
			return
		}
		log.Println("Attendance query error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch attendance"})
		return
	}

	c.JSON(http.StatusOK, result)
}
