package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"path/filepath"
)

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ไม่พบไฟล์ที่อัปโหลด"})
		return
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "รองรับเฉพาะไฟล์รูปภาพ (jpg, png, webp) เท่านั้น"})
		return
	}

	newFileName := uuid.New().String() + ext
	dst := "./uploads/" + newFileName

	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "เกิดข้อผิดพลาดในการบันทึกไฟล์"})
		return
	}

	fileURL := fmt.Sprintf("/uploads/%s", newFileName)

	c.JSON(http.StatusOK, gin.H{
		"url": fileURL,
	})
}
