package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pos-service/config"
	"pos-service/models"
)

func GetOrderTypes(c *gin.Context) {
	var orderTypes []models.OrderType
	if err := config.Db.Find(&orderTypes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orderTypes)
}

func GetOrderTypeByID(c *gin.Context) {
	id := c.Param("id")
	var orderType models.OrderType
	if err := config.Db.First(&orderType, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบประเภทคำสั่งซื้อ"})
		return
	}

	c.JSON(http.StatusOK, orderType)
}

func CreateOrderType(c *gin.Context) {
	type createOrderTypeRequest struct {
		Type string `json:"type" binding:"required"`
	}

	var req createOrderTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ครบถ้วน"})
		return
	}

	orderType := models.OrderType{
		Type: req.Type,
	}

	if err := config.Db.Create(&orderType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถสร้างประเภทคำสั่งซื้อได้"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":    "สร้างประเภทคำสั่งซื้อสำเร็จ",
		"order_type": orderType,
	})
}

func UpdateOrderType(c *gin.Context) {
	id := c.Param("id")
	type updateOrderTypeRequest struct {
		Type   *string `json:"type"`
		Status *string `json:"status"`
	}
	var req updateOrderTypeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "รูปแบบข้อมูลไม่ถูกต้อง"})
		return
	}

	var orderType models.OrderType
	if err := config.Db.First(&orderType, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบประเภทคำสั่งซื้อ"})
		return
	}

	updateData := make(map[string]any)

	if req.Type != nil {
		updateData["type"] = *req.Type
	}
	if req.Status != nil {
		updateData["status"] = *req.Status
	}

	if err := config.Db.Model(&orderType).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถอัปเดตประเภทคำสั่งซื้อได้"})
		return
	}

	if err := config.Db.Model(&orderType).Updates(updateData).Scan(&orderType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถอัปเดตข้อมูลได้"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":    "อัปเดตประเภทคำสั่งซื้อสำเร็จ",
		"order_type": orderType,
	})
}

func DeleteOrderType(c *gin.Context) {
	id := c.Param("id")
	var orderType models.OrderType

	if err := config.Db.First(&orderType, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบประเภทคำสั่งซื้อ"})
		return
	}

	if err := config.Db.Delete(&orderType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถลบประเภทคำสั่งซื้อได้"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ลบประเภทคำสั่งซื้อสำเร็จ"})
}
