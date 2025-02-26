package handlers

import (
	"net/http"
	"personal-finance/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type CategoryHandler struct {
	DB *gorm.DB
}

// CreateCategory 创建新分类
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证分类类型
	if category.Type != "expense" && category.Type != "income" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category type must be either 'expense' or 'income'"})
		return
	}

	if err := h.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// GetCategories 获取所有分类
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	var categories []models.Category
	
	// 支持按类型筛选
	categoryType := c.Query("type")
	query := h.DB
	if categoryType != "" {
		query = query.Where("type = ?", categoryType)
	}

	if err := query.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// UpdateCategory 更新分类
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	
	if err := h.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	var updatedCategory models.Category
	if err := c.ShouldBindJSON(&updatedCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证分类类型
	if updatedCategory.Type != "expense" && updatedCategory.Type != "income" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category type must be either 'expense' or 'income'"})
		return
	}

	// 更新字段
	category.Name = updatedCategory.Name
	category.Type = updatedCategory.Type
	category.Icon = updatedCategory.Icon

	if err := h.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// DeleteCategory 删除分类
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	
	// 检查分类是否存在
	var category models.Category
	if err := h.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	// 检查是否有关联的交易
	var transactionCount int64
	h.DB.Model(&models.Transaction{}).Where("category_id = ?", id).Count(&transactionCount)
	if transactionCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete category with associated transactions"})
		return
	}

	// 检查是否有关联的预算
	var budgetCount int64
	h.DB.Model(&models.Budget{}).Where("category_id = ?", id).Count(&budgetCount)
	if budgetCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete category with associated budgets"})
		return
	}

	if err := h.DB.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
