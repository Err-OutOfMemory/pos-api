package controllers

import (
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"pos-service/config"
	"pos-service/models"

	"github.com/gin-gonic/gin"
)

func removeFile(fileURL string) {
	if fileURL == "" || strings.Contains(fileURL, "placehold.co") {
		return
	}
	fileName := filepath.Base(fileURL)
	filePath := filepath.Join("./uploads", fileName)
	os.Remove(filePath)
}

func GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	categoryID := c.Query("category_id")
	search := c.Query("search")

	var products []models.Product
	var total int64

	query := config.Db.Model(&models.Product{}).Preload("Category")

	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if search != "" {
		query = query.Where("product_name LIKE ?", "%"+search+"%")
	}

	query.Count(&total)

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       products,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"total_page": math.Ceil(float64(total) / float64(limit)),
	})
}

func GetProductByID(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := config.Db.Preload("Category").First(&product, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func CreateProduct(c *gin.Context) {
	type createProductRequest struct {
		ProductName string  `json:"product_name" binding:"required"`
		CategoryID  int     `json:"category_id" binding:"required"`
		Price       float64 `json:"price" binding:"required"`
		Type        string  `json:"type"`
		Description string  `json:"description"`
		ImgPath     *string `json:"img_path"`
	}

	var req createProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ครบถ้วน"})
		return
	}

	imgPath := req.ImgPath
	if imgPath == nil {
		defaultImg := "https://placehold.co/400x400/png/default-image.png"
		imgPath = &defaultImg
	}

	product := models.Product{
		ProductName: req.ProductName,
		CategoryID:  req.CategoryID,
		Price:       req.Price,
		Type:        req.Type,
		Description: req.Description,
		ImgPath:     *imgPath,
		Status:      "active",
	}

	if err := config.Db.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถสร้างสินค้าได้ (ตรวจสอบ CategoryID อีกครั้ง)"})
		return
	}

	config.Db.Preload("Category").First(&product, product.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "สร้างสินค้าสำเร็จ",
		"product": product,
	})
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	type updateProductRequest struct {
		ProductName *string  `json:"product_name"`
		CategoryID  *uint    `json:"category_id"`
		Price       *float64 `json:"price"`
		Status      *string  `json:"status"`
		ImgPath     *string  `json:"img_path"`
		Description *string  `json:"description"`
		Type        *string  `json:"type"`
	}
	var req updateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "รูปแบบข้อมูลไม่ถูกต้อง"})
		return
	}

	var product models.Product
	if err := config.Db.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบสินค้าที่ระบุ"})
		return
	}

	updateData := make(map[string]any)

	if req.ProductName != nil {
		updateData["product_name"] = *req.ProductName
	}
	if req.ImgPath != nil && *req.ImgPath != product.ImgPath {
		removeFile(product.ImgPath)
		updateData["img_path"] = *req.ImgPath
	}
	if req.CategoryID != nil {
		updateData["category_id"] = *req.CategoryID
	}
	if req.Price != nil {
		updateData["price"] = *req.Price
	}
	if req.Status != nil {
		updateData["status"] = *req.Status
	}
	if req.Description != nil {
		updateData["description"] = *req.Description
	}
	if req.Type != nil {
		updateData["type"] = *req.Type
	}

	if err := config.Db.Model(&product).Updates(updateData).Scan(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถอัปเดตข้อมูลได้"})
		return
	}

	config.Db.Preload("Category").First(&product, product.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"product": product,
	})
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := config.Db.First(&product, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	removeFile(product.ImgPath)
	if err := config.Db.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถลบ Product ได้"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
