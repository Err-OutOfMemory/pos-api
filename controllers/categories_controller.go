package controllers

import (
	"net/http"

	"pos-service/config"
	"pos-service/models"

	"github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {
	var categories []models.Category

	if err := config.Db.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
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
