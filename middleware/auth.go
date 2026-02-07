package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "กรุณาเข้าสู่ระบบ"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token ไม่ถูกต้องหรือหมดอายุ"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("emp_code", claims["emp_code"])
			c.Set("role", claims["role"])

			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "ไม่สามารถอ่านข้อมูลจาก Token ได้"})
			c.Abort()
		}
	}
}

func AuthorizeRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "คุณไม่มีสิทธิ์เข้าถึงส่วนนี้"})
			c.Abort()
			return
		}
		c.Next()
	}
}
