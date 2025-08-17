package controller

import (
	"database/sql"
	"log"
	"net/http"

	"manajemen-karyawan-api/config"
	"manajemen-karyawan-api/middleware"
	"manajemen-karyawan-api/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	EmployeeID string `json:"employeeID" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

// Login godoc
// @Summary Login dengan employee_id dan password
// @Description Login dan menyimpan JWT di cookie HttpOnly untuk autentikasi
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Employee ID dan Password"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/auth/login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee_id and password are required"})
		return
	}

	var emp model.Employee
	query := `
		SELECT id, employee_id, departement_id, name, address, password,
		       created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM employee
		WHERE employee_id = ? AND deleted_at IS NULL
	`
	err := config.DB.QueryRow(query, req.EmployeeID).Scan(
		&emp.ID, &emp.EmployeeID, &emp.DepartementID, &emp.Name, &emp.Address, &emp.Password,
		&emp.CreatedAt, &emp.CreatedBy, &emp.UpdatedAt, &emp.UpdatedBy,
		&emp.DeletedAt, &emp.DeletedBy,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid employee_id"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(emp.Password), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}

	token, err := middleware.GenerateToken(emp.ID, emp.EmployeeID, []byte(config.JWTSecret))
	if err != nil {
		log.Println("JWT generation error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.SetCookie(
		"access_token",
		token,
		3600,
		"/",
		config.CookieDomain,
		false, // true if using HTTPS
		true,  // HttpOnly
	)

	c.JSON(http.StatusOK, gin.H{"message": "login successful"})
}

// GetMe godoc
// @Summary Ambil data user dari JWT cookie
// @Description Validasi session dan mengembalikan user context
// @Tags Auth
// @Produce json
// @Success 200 {object} model.Employee
// @Failure 401 {object} map[string]string
// @Router /api/auth/me [get]
func GetMe(c *gin.Context) {
	empID, exists := c.Get("employee_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var emp model.Employee
	if err := config.DB.QueryRow(`
		SELECT id, employee_id, name
		FROM employee
		WHERE employee_id = ? AND deleted_at IS NULL
	`, empID).Scan(&emp.ID, &emp.EmployeeID, &emp.Name); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         emp.ID,
		"employeeID": emp.EmployeeID,
		"name":       emp.Name,
	})
}

// Logout godoc
// @Summary Logout user
// @Description Menghapus cookie JWT dan mengakhiri sesi user yang login
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/auth/logout [post]
func Logout(c *gin.Context) {
	c.SetCookie(
		"access_token",
		"",
		-1,
		"/",
		config.CookieDomain,
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
