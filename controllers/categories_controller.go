package controllers

import (
	"net/http"
	"strconv"

	"pos-service/config"
	"pos-service/models"

	"github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	status := c.Query("status")

	var categories []models.Category
	var total int64

	query := config.Db.Model(&models.Category{})

	if search != "" {
		query = query.Where("category_name LIKE ?", "%"+search+"%")
	}
	if status != "" {
		switch status {
		case "active":
			query = query.Where("status = ?", true)
		case "inactive":
			query = query.Where("status = ?", false)
		default:
		}
	}

	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":       categories,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"total_page": int64(total)/int64(limit) + 1,
	})
}

func GetCategoryByID(c *gin.Context) {
	id := c.Param("id")
	var category models.Category

	if err := config.Db.Preload("Category").First(&category, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	c.JSON(http.StatusOK, category)
}

func CreateCategory(c *gin.Context) {
	type createCategoryRequest struct {
		CategoryName string `json:"category_name" binding:"required"`
	}

	var req createCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := models.Category{
		CategoryName: req.CategoryName,
		Status:       true,
	}

	if err := config.Db.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	type updateCategoryRequest struct {
		CategoryName *string `json:"category_name"`
		Status       *bool   `json:"status"`
	}
	var req updateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "รูปแบบข้อมูลไม่ถูกต้อง"})
		return
	}

	var category models.Category
	if err := config.Db.First(&category, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	updateData := make(map[string]any)

	if req.CategoryName != nil {
		updateData["category_name"] = *req.CategoryName
	}
	if req.Status != nil {
		updateData["status"] = *req.Status
	}

	if err := config.Db.Model(&category).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถอัปเดตข้อมูลได้"})
		return
	}
	c.JSON(
		http.StatusOK, gin.H{
			"message":  "Category updated successfully",
			"response": category,
		})
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	if err := config.Db.First(&category, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	if err := config.Db.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถลบ Category ได้"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
