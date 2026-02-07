package controllers

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"pos-service/config"
	"pos-service/models"
)

func GetOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var orders []models.Order
	var total int64

	config.Db.Model(&models.Order{}).Count(&total)

	if err := config.Db.Limit(limit).Offset(offset).
		Order("created_at desc").
		Preload("OrderType").
		Preload("Employee").
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถดึงข้อมูลได้"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       orders,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"total_page": math.Ceil(float64(total) / float64(limit)),
	})
}

func GetOrderByID(c *gin.Context) {
	id := c.Param("id")
	var order models.Order

	if err := config.Db.
		Preload("OrderType").
		Preload("User").
		Preload("OrderDetails.Product").
		First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบข้อมูลออเดอร์"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

func CreateOrder(c *gin.Context) {
	type orderItemRequest struct {
		ProductID      int     `json:"product_id" binding:"required"`
		Quantity       int     `json:"quantity" binding:"required"`
		Price          float64 `json:"price" binding:"required"`
		DiscountAmount float64 `json:"discount_amount"`
		Description    string  `json:"description"`
	}

	type createOrderRequest struct {
		EmpID       int                `json:"emp_id" binding:"required"`
		OrderTypeID int                `json:"order_type_id" binding:"required"`
		Discount    float64            `json:"discount"`
		TotalPrice  float64            `json:"total_price" binding:"required"`
		Items       []orderItemRequest `json:"items" binding:"required"`
	}

	var req createOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ถูกต้อง"})
		return
	}

	tx := config.Db.Begin()

	var details []models.OrderDetail
	for _, item := range req.Items {
		details = append(details, models.OrderDetail{
			ProductID:      item.ProductID,
			Quantity:       item.Quantity,
			Price:          item.Price,
			DiscountAmount: item.DiscountAmount,
			Description:    item.Description,
		})
	}

	order := models.Order{
		EmpID:        req.EmpID,
		OrderTypeID:  req.OrderTypeID,
		Discount:     req.Discount,
		TotalPrice:   req.TotalPrice,
		Date:         time.Now(),
		OrderDetails: details,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถสร้างออเดอร์ได้"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "สร้างออเดอร์สำเร็จ",
		"order_id": order.ID,
	})
}

func UpdateOrder(c *gin.Context) {
	id := c.Param("id")

	type orderItemRequest struct {
		ProductID int     `json:"product_id"`
		Quantity  int     `json:"quantity"`
		Price     float64 `json:"price"`
	}

	type updateOrderRequest struct {
		EmpID       *int               `json:"emp_id"`
		OrderTypeID *int               `json:"order_type_id"`
		TotalPrice  *float64           `json:"total_price"`
		Items       []orderItemRequest `json:"items"`
	}

	var req updateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ถูกต้อง"})
		return
	}

	tx := config.Db.Begin()

	var order models.Order
	if err := tx.First(&order, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบออเดอร์"})
		return
	}

	if req.EmpID != nil {
		order.EmpID = *req.EmpID
	}
	if req.OrderTypeID != nil {
		order.OrderTypeID = *req.OrderTypeID
	}
	if req.TotalPrice != nil {
		order.TotalPrice = *req.TotalPrice
	}

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update header failed"})
		return
	}

	if len(req.Items) > 0 {
		if err := tx.Where("order_id = ?", id).Delete(&models.OrderDetail{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear old items"})
			return
		}

		for _, item := range req.Items {
			newDetail := models.OrderDetail{
				OrderID:   order.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Price,
			}
			if err := tx.Create(&newDetail).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new items"})
				return
			}
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "อัปเดตออเดอร์และรายการสินค้าสำเร็จ"})
}

func UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	type updateStatus struct {
		Status string `json:"status" binding:"required"`
	}

	var req updateStatus
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "กรุณาระบุสถานะ"})
		return
	}

	if err := config.Db.Model(&models.Order{}).Where("id = ?", id).Update("status", req.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถอัปเดตสถานะได้"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "เปลี่ยนสถานะออเดอร์เป็น " + req.Status + " เรียบร้อยแล้ว"})
}

func DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order
	if err := config.Db.First(&order, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบออเดอร์"})
		return
	}
	if err := config.Db.Delete(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถลบออเดอร์ได้"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ลบออเดอร์สำเร็จ"})
}
