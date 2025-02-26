package handlers

import (
	"net/http"
	"personal-finance/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type BudgetHandler struct {
	DB *gorm.DB
}

// CreateBudget 创建新预算
func (h *BudgetHandler) CreateBudget(c *gin.Context) {
	// 解析输入
	var input models.BudgetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证日期格式
	_, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
		return
	}

	_, err = time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
		return
	}

	// 创建预算对象
	budget := models.Budget{
		CategoryID: input.CategoryID,
		Amount:     input.Amount,
		StartDate:  input.StartDate,
		EndDate:    input.EndDate,
	}

	// 检查分类是否存在
	var category models.Category
	if err := h.DB.First(&category, budget.CategoryID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category not found"})
		return
	}

	if err := h.DB.Create(&budget).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 加载关联的分类信息
	h.DB.Model(&budget).Related(&budget.Category)

	c.JSON(http.StatusCreated, budget)
}

// GetBudgets 获取预算列表
func (h *BudgetHandler) GetBudgets(c *gin.Context) {
	var budgets []models.Budget
	
	query := h.DB.Preload("Category")
	
	// 支持按时间范围筛选
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	
	if startDate != "" {
		query = query.Where("start_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("end_date <= ?", endDate)
	}

	if err := query.Find(&budgets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, budgets)
}

// GetBudgetStatus 获取预算执行状况
func (h *BudgetHandler) GetBudgetStatus(c *gin.Context) {
	id := c.Param("id")
	var budget models.Budget
	
	if err := h.DB.Preload("Category").First(&budget, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found"})
		return
	}

	// 计算该预算分类下的实际支出
	var actualExpense float64
	h.DB.Model(&models.Transaction{}).
		Where("category_id = ? AND type = 'expense' AND date(created_at) BETWEEN ? AND ?",
			budget.CategoryID, budget.StartDate, budget.EndDate).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&actualExpense)

	// 计算预算使用百分比
	percentageUsed := (actualExpense / budget.Amount) * 100

	c.JSON(http.StatusOK, gin.H{
		"budget": budget,
		"status": gin.H{
			"actual_expense":   actualExpense,
			"percentage_used": percentageUsed,
			"remaining":       budget.Amount - actualExpense,
		},
	})
}

// UpdateBudget 更新预算
func (h *BudgetHandler) UpdateBudget(c *gin.Context) {
	id := c.Param("id")
	var budget models.Budget
	
	if err := h.DB.First(&budget, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found"})
		return
	}

	if err := c.ShouldBindJSON(&budget); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Save(&budget).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 加载关联的分类信息
	h.DB.Model(&budget).Related(&budget.Category)
	c.JSON(http.StatusOK, budget)
}

// DeleteBudget 删除预算
func (h *BudgetHandler) DeleteBudget(c *gin.Context) {
	id := c.Param("id")
	if err := h.DB.Delete(&models.Budget{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Budget deleted successfully"})
}
