package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"manajemen-karyawan-api/config"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenExpired          = errors.New("token expired")
	ErrTokenMalformed        = errors.New("malformed token")
	ErrTokenInvalid          = errors.New("invalid token")
	ErrTokenSignatureInvalid = errors.New("signature invalid")
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read token from cookie
		tokenString, err := c.Cookie("access_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		// Parse and validate token
		claims, err := parseToken(tokenString, []byte(config.JWTSecret))
		if err != nil {
			log.Printf("JWT error: %v", err)

			switch err {
			case ErrTokenExpired:
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
			case ErrTokenMalformed:
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "malformed token"})
			case ErrTokenSignatureInvalid, ErrTokenInvalid:
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			default:
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			}
			return
		}

		// Extract id
		id, ok := claims["id"].(string)
		if !ok || strings.TrimSpace(id) == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "id missing in token"})
			return
		}

		// Extract employee_id
		employeeID, ok := claims["employee_id"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "employee_id missing in token"})
			return
		}

		// Inject into context
		c.Set("id", id)
		c.Set("employee_id", employeeID)
		c.Next()
	}
}

func GenerateToken(id string, employeeID string, secret []byte) (string, error) {
	if len(secret) == 0 {
		return "", ErrTokenMalformed
	}

	claims := jwt.MapClaims{
		"id":          id,
		"employee_id": employeeID,
		"exp":         time.Now().Add(24 * time.Hour).Unix(),
		"iat":         time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func parseToken(tokenString string, secret []byte) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, ErrTokenInvalid
	}

	if !token.Valid {
		return nil, ErrTokenInvalid
	}

	// Manual expiry check
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, ErrTokenMalformed
	}
	if int64(exp) < time.Now().Unix() {
		return nil, ErrTokenExpired
	}

	return claims, nil
}
