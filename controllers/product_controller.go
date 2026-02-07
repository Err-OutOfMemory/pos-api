package controllers

import (
	"net/http"

	"pos-service/config"
	"pos-service/models"

	"github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context) {
	var products []models.Product

	if err := config.Db.Preload("Category").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
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

	if err := config.Db.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถลบ Product ได้"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
