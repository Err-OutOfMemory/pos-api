package controllers

import (
	"fmt"
	"net/http"

	"pos-service/config"
	"pos-service/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllEmployees(c *gin.Context) {
	var employees []models.Employee

	if err := config.Db.Preload("User").Find(&employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, employees)
}

func GetEmployeeByID(c *gin.Context) {
	id := c.Param("id")
	var employee models.Employee

	if err := config.Db.Preload("User").First(&employee, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}
	c.JSON(http.StatusOK, employee)
}

func CreateEmployee(c *gin.Context) {
	type createEmployeeRequest struct {
		Name     string `json:"name" binding:"required"`
		Role     string `json:"role" binding:"required"`
		PhoneNum string `json:"phonenum"`
	}

	var req createEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ครบถ้วน"})
		return
	}

	err := config.Db.Transaction(func(tx *gorm.DB) error {
		newEmployee := models.Employee{
			Name:     req.Name,
			Role:     req.Role,
			PhoneNum: req.PhoneNum,
			Status:   "active",
		}
		if err := tx.Create(&newEmployee).Error; err != nil {
			return err
		}
		generatedCode := fmt.Sprintf("E%03d", newEmployee.ID)

		if err := tx.Model(&newEmployee).Update("emp_code", generatedCode).Error; err != nil {
			return err
		}

		newUser := models.User{
			EmployeeID: newEmployee.ID,
			PinHash:    nil,
		}
		if err := tx.Create(&newUser).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถสร้างพนักงานได้: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "สร้างพนักงานสำเร็จ"})
}

func UpdateEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee models.Employee

	if err := config.Db.First(&employee, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบพนักงาน"})
		return
	}

	type updateRequest struct {
		Name     *string `json:"name"`
		Role     *string `json:"role"`
		PhoneNum *string `json:"phonenum"`
		Status   *string `json:"status"`
	}

	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ถูกต้อง"})
		return
	}

	if req.Name != nil {
		employee.Name = *req.Name
	}
	if req.Role != nil {
		employee.Role = *req.Role
	}
	if req.PhoneNum != nil {
		employee.PhoneNum = *req.PhoneNum
	}
	if req.Status != nil {
		employee.Status = *req.Status
	}

	if err := config.Db.Save(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, employee)
}

func DeleteEmployee(c *gin.Context) {
	id := c.Param("id")
	if err := config.Db.Delete(&models.Employee{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ลบพนักงานเรียบร้อยแล้ว"})
}
