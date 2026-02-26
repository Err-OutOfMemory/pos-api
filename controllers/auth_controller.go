package controllers

import (
	"fmt"
	"net/http"
	"os"
	"pos-service/config"
	"pos-service/models"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GenerateToken(empCode string, role string) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	claims := jwt.MapClaims{
		"emp_code": empCode,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 12).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func CheckUser(c *gin.Context) {
	type CheckUserRequest struct {
		EmpCode string `json:"emp_code" binding:"required"`
	}
	var req CheckUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "กรุณากรอกรหัสพนักงาน"})
		return
	}

	var employee models.Employee
	result := config.Db.Preload("User").
		Where("emp_code = ? AND status = ?", req.EmpCode, "active").
		First(&employee)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบรหัสพนักงานนี้ในระบบ"})
		return
	}

	isFirstLogin := employee.User.PinHash == nil || *employee.User.PinHash == ""

	c.JSON(http.StatusOK, gin.H{
		"employee_id":    employee.ID,
		"emp_code":       employee.EmpCode,
		"name":           employee.Name,
		"role":           employee.Role,
		"is_first_login": isFirstLogin,
	})
}

func SetupPin(c *gin.Context) {
	type SetupPinRequest struct {
		EmployeeID int    `json:"employee_id" binding:"required"`
		Pin        string `json:"pin" binding:"required,len=6"`
	}
	var req SetupPinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("Bind Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ถูกต้อง (PIN ต้องมี 6 หลัก)"})
		return
	}

	var user models.User
	if err := config.Db.Where("employee_id = ?", req.EmployeeID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบข้อมูลผู้ใช้งาน"})
		return
	}

	if user.PinHash != nil && *user.PinHash != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "พนักงานคนนี้ตั้งรหัส PIN เรียบร้อยแล้ว"})
		return
	}

	hashedPin, err := bcrypt.GenerateFromPassword([]byte(req.Pin), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถประมวลผลรหัสผ่านได้"})
		return
	}

	var employee models.Employee
	err = config.Db.Transaction(func(tx *gorm.DB) error {
		pinStr := string(hashedPin)
		if err := tx.Model(&user).Update("pin_hash", pinStr).Error; err != nil {
			return err
		}

		if err := tx.First(&employee, req.EmployeeID).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "บันทึกข้อมูลล้มเหลว"})
		return
	}

	token, _ := GenerateToken(employee.EmpCode, employee.Role)

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"emp_code": employee.EmpCode,
			"name":     employee.Name,
			"role":     employee.Role,
		},
		"token": token,
	})
}

func Login(c *gin.Context) {
	type LoginRequest struct {
		EmpCode string `json:"emp_code" binding:"required"`
		Pin     string `json:"pin" binding:"required,len=6"`
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ถูกต้อง"})
		return
	}

	var employee models.Employee
	if err := config.Db.Preload("User").Where("emp_code = ?", req.EmpCode).First(&employee).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบรหัสพนักงานนี้ในระบบ"})
		return
	}

	user := employee.User

	if user.PinHash == nil || *user.PinHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "พนักงานยังไม่ได้ตั้งรหัส PIN กรุณาลงทะเบียนก่อน"})
		return
	}

	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
		timeLeft := time.Until(*user.LockedUntil).Minutes()
		c.JSON(http.StatusForbidden, gin.H{
			"error": fmt.Sprintf("บัญชีถูกระงับชั่วคราว กรุณาลองใหม่ในอีก %.0f นาที", timeLeft),
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(*user.PinHash), []byte(req.Pin))

	if err != nil {
		newAttempts := user.FailedAttempts + 1
		updateData := map[string]interface{}{"failed_attempts": newAttempts}

		if newAttempts >= 5 {
			lockTime := time.Now().Add(time.Minute * 5)
			updateData["locked_until"] = lockTime
			c.JSON(http.StatusForbidden, gin.H{"error": "กรอกรหัสผิดเกินกำหนด บัญชีถูกล็อก 5 นาที"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "รหัส PIN ไม่ถูกต้อง"})
		}

		config.Db.Model(&user).Updates(updateData)
		return
	}

	config.Db.Model(&user).Updates(map[string]interface{}{
		"failed_attempts": 0,
		"locked_until":    nil,
	})

	token, _ := GenerateToken(employee.EmpCode, employee.Role)

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"emp_code": employee.EmpCode,
			"name":     employee.Name,
			"role":     employee.Role,
		},
		"token": token,
	})
}

func GetProfile(c *gin.Context) {
	empCode := c.GetString("emp_code")

	if empCode == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	var employee models.Employee

	err := config.Db.
		Where("emp_code = ? AND status = ?", empCode, "active").
		First(&employee).Error

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"emp_code": employee.EmpCode,
		"name":     employee.Name,
		"role":     employee.Role,
	})
}
