package utils

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type QueryParams struct {
	Page    *int              `json:"page"`
	PerPage *int              `json:"per_page"`
	SortBy  []SortField       `json:"sort_by"`
	Filter  map[string]string `json:"filter"`
}

type MetaParams struct {
	Page    *int
	PerPage *int
	Total   int
	SortBy  []SortField
}

type Pagination struct {
	Offset int
	Limit  int
	Use    bool
}

type SortField struct {
	Key   string `json:"key"`
	Order string `json:"order"`
}

func CreatePassword(raw string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func GenerateID() string {
	return uuid.New().String()
}

func ReplacePlaceholders(query string) string {
	count := strings.Count(query, "?")
	for i := 1; i <= count; i++ {
		query = strings.Replace(query, "?", fmt.Sprintf("$%d", i), 1)
	}
	return query
}

func BuildMeta(p MetaParams) map[string]interface{} {
	meta := map[string]interface{}{
		"total": p.Total,
	}

	if p.Page != nil && p.PerPage != nil {
		meta["page"] = *p.Page
		meta["per_page"] = *p.PerPage
		meta["total_pages"] = int(math.Ceil(float64(p.Total) / float64(*p.PerPage)))
	}

	if len(p.SortBy) > 0 {
		meta["sort_by"] = p.SortBy
	}

	return meta
}

func BindJSONStrict(c *gin.Context, target interface{}) error {
	body, _ := io.ReadAll(c.Request.Body)
	log.Println("Raw JSON:", string(body))
	c.Request.Body = io.NopCloser(strings.NewReader(string(body)))

	if err := c.ShouldBindJSON(target); err != nil {
		log.Println("Bind error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return err
	}
	return nil
}

func BuildSortSQL(sortBy []SortField, whitelist map[string]string) string {
	if len(sortBy) == 0 {
		return ""
	}

	var clauses []string
	for _, s := range sortBy {
		col, ok := whitelist[s.Key]
		if !ok {
			log.Printf("Sort field '%s' is not allowed", s.Key)
			continue
		}

		order := strings.ToUpper(s.Order)
		if order != "ASC" && order != "DESC" {
			order = "ASC"
		}

		clauses = append(clauses, fmt.Sprintf("%s %s", col, order))
	}

	if len(clauses) == 0 {
		return ""
	}

	return "ORDER BY " + strings.Join(clauses, ", ")
}
func BuildFilterSQL(filter map[string]string, whitelist map[string]string) (string, []interface{}) {
	var clauses []string
	var args []interface{}

	for key, value := range filter {
		field := key
		op := "LIKE"
		arg := "%" + value + "%"

		switch {
		case strings.HasSuffix(key, ".gte"):
			field = strings.TrimSuffix(key, ".gte")
			op = ">="
			arg = value
		case strings.HasSuffix(key, ".lte"):
			field = strings.TrimSuffix(key, ".lte")
			op = "<="
			arg = value
		case strings.HasSuffix(key, ".like"):
			field = strings.TrimSuffix(key, ".like")
			op = "LIKE"
			arg = "%" + value + "%"
		}

		col, ok := whitelist[field]
		if !ok {
			log.Printf("Filter field '%s' is not allowed", key)
			continue
		}

		clauses = append(clauses, fmt.Sprintf("%s %s ?", col, op))
		args = append(args, arg)
	}

	if len(clauses) == 0 {
		return "", args
	}

	return "AND " + strings.Join(clauses, " AND "), args
}

func BuildPagination(page, perPage *int) Pagination {
	if page == nil || perPage == nil {
		return Pagination{Use: false}
	}
	p := *page
	pp := *perPage
	if p < 1 {
		p = 1
	}
	if pp < 1 {
		pp = 10
	}
	return Pagination{
		Offset: (p - 1) * pp,
		Limit:  pp,
		Use:    true,
	}
}

func BuildDynamicUpdateQuery(table string, payload map[string]interface{}, whitelist []string, audit map[string]interface{}) (string, []interface{}, error) {
	setClauses := []string{}
	args := []interface{}{}

	for key, val := range payload {
		if !contains(whitelist, key) {
			return "", nil, fmt.Errorf("field %s is not allowed to be updated", key)
		}
		setClauses = append(setClauses, fmt.Sprintf("%s = ?", key))
		args = append(args, val)
	}

	for auditKey, auditVal := range audit {
		setClauses = append(setClauses, fmt.Sprintf("%s = ?", auditKey))
		args = append(args, auditVal)
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ? AND deleted_at IS NULL", table, strings.Join(setClauses, ", "))
	return query, args, nil
}

func contains(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func IsLate(clockIn time.Time, maxRaw string, loc *time.Location) (bool, error) {
	maxParsed, err := time.ParseInLocation("15:04:05", maxRaw, loc)
	if err != nil {
		return false, err
	}
	clockInLocal := clockIn.In(loc)
	clockInOnly := time.Date(0, 1, 1, clockInLocal.Hour(), clockInLocal.Minute(), clockInLocal.Second(), 0, loc)
	return clockInOnly.After(maxParsed), nil
}

func IsEarly(clockOut time.Time, maxRaw string, loc *time.Location) (bool, error) {
	maxParsed, err := time.ParseInLocation("15:04:05", maxRaw, loc)
	if err != nil {
		return false, err
	}
	clockOutLocal := clockOut.In(loc)
	clockOutOnly := time.Date(0, 1, 1, clockOutLocal.Hour(), clockOutLocal.Minute(), clockOutLocal.Second(), 0, loc)
	return clockOutOnly.Before(maxParsed), nil
}
