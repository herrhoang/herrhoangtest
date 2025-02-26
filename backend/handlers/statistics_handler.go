package handlers

import (
	"net/http"
	"personal-finance/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type StatisticsHandler struct {
	DB *gorm.DB
}

// GetStatistics 获取统计数据
func (h *StatisticsHandler) GetStatistics(c *gin.Context) {
	// 获取查询参数
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	
	if startDate == "" {
		// 默认查询最近一年的数据
		now := time.Now()
		startDate = now.AddDate(-1, 0, 0).Format("2006-01-02")
	}
	
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	// 查询总收入和支出
	var stats models.Statistics
	h.DB.Model(&models.Transaction{}).
		Where("date(created_at) BETWEEN ? AND ? AND type = ?", startDate, endDate, "income").
		Select("COALESCE(SUM(amount), 0)").Row().
		Scan(&stats.TotalIncome)

	h.DB.Model(&models.Transaction{}).
		Where("date(created_at) BETWEEN ? AND ? AND type = ?", startDate, endDate, "expense").
		Select("COALESCE(SUM(amount), 0)").Row().
		Scan(&stats.TotalExpense)

	stats.NetAmount = stats.TotalIncome - stats.TotalExpense

	// 按分类统计
	rows, err := h.DB.Table("transactions").
		Select("categories.id, categories.name, SUM(transactions.amount) as amount").
		Joins("JOIN categories ON transactions.category_id = categories.id").
		Where("date(transactions.created_at) BETWEEN ? AND ?", startDate, endDate).
		Group("categories.id, categories.name").
		Rows()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var stat models.CategoryStatistics
		rows.Scan(&stat.CategoryID, &stat.CategoryName, &stat.Amount)
		
		// 计算百分比
		if stat.Amount > 0 {
			stat.Percentage = (stat.Amount / stats.TotalExpense) * 100
		}
		
		stats.ByCategory = append(stats.ByCategory, stat)
	}

	// 按月份统计
	rows, err = h.DB.Table("transactions").
		Select("strftime('%Y', created_at) as year, strftime('%m', created_at) as month, " +
			"SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END) as income, " +
			"SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END) as expense").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Group("year, month").
		Order("year DESC, month DESC").
		Rows()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var stat models.MonthlyStatistics
		rows.Scan(&stat.Year, &stat.Month, &stat.Income, &stat.Expense)
		stat.NetAmount = stat.Income - stat.Expense
		stats.ByMonth = append(stats.ByMonth, stat)
	}

	c.JSON(http.StatusOK, stats)
}

// GetBudgetOverview 获取预算概览
// CreateBudget 创建预算
func (h *StatisticsHandler) CreateBudget(c *gin.Context) {
	var input models.BudgetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证分类是否存在
	var category models.Category
	if err := h.DB.First(&category, input.CategoryID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category not found"})
		return
	}

	// 验证日期格式
	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use YYYY-MM-DD"})
		return
	}

	// 验证日期范围
	if endDate.Before(startDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End date must be after start date"})
		return
	}

	// 验证金额
	if input.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be greater than 0"})
		return
	}

	// 创建预算
	budget := models.Budget{
		CategoryID: input.CategoryID,
		Amount:     input.Amount,
		StartDate:  input.StartDate,
		EndDate:    input.EndDate,
	}

	if err := h.DB.Create(&budget).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 加载分类信息
	h.DB.Model(&budget).Related(&budget.Category)

	c.JSON(http.StatusCreated, budget)
}

// GetBudgets 获取预算列表
func (h *StatisticsHandler) GetBudgets(c *gin.Context) {
	var budgets []models.Budget

	query := h.DB.Preload("Category")

	// 支持按分类ID筛选
	if categoryID := c.Query("category_id"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	// 支持按日期范围筛选
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("start_date >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("end_date <= ?", endDate)
	}

	if err := query.Find(&budgets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, budgets)
}

// UpdateBudget 更新预算
func (h *StatisticsHandler) UpdateBudget(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid budget ID"})
		return
	}

	var budget models.Budget
	if err := h.DB.First(&budget, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found"})
		return
	}

	// 绑定更新数据
	var input models.BudgetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证分类是否存在
	if input.CategoryID != budget.CategoryID {
		var category models.Category
		if err := h.DB.First(&category, input.CategoryID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category not found"})
			return
		}
	}

	// 验证日期格式
	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use YYYY-MM-DD"})
		return
	}

	// 验证日期范围
	if endDate.Before(startDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End date must be after start date"})
		return
	}

	// 验证金额
	if input.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be greater than 0"})
		return
	}

	// 更新预算字段
	budget.CategoryID = input.CategoryID
	budget.Amount = input.Amount
	budget.StartDate = input.StartDate
	budget.EndDate = input.EndDate

	if err := h.DB.Save(&budget).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 加载分类信息
	h.DB.Model(&budget).Related(&budget.Category)

	c.JSON(http.StatusOK, budget)
}

// DeleteBudget 删除预算
func (h *StatisticsHandler) DeleteBudget(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid budget ID"})
		return
	}

	var budget models.Budget
	if err := h.DB.First(&budget, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found"})
		return
	}

	// 删除预算
	if err := h.DB.Delete(&budget).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Budget deleted successfully"})
}

// GetBudgetOverview 获取预算概览
func (h *StatisticsHandler) GetBudgetOverview(c *gin.Context) {
	currentDate := time.Now()
	startOfMonth := time.Date(currentDate.Year(), currentDate.Month(), 1, 0, 0, 0, 0, currentDate.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	type BudgetOverview struct {
		models.Budget
		ActualExpense   float64 `json:"actual_expense"`
		Remaining       float64 `json:"remaining"`
		PercentageUsed  float64 `json:"percentage_used"`
	}

	var budgets []BudgetOverview

	rows, err := h.DB.Table("budgets").
		Select("budgets.*, categories.name as category_name, " +
			"COALESCE((SELECT SUM(amount) FROM transactions " +
			"WHERE category_id = budgets.category_id " +
			"AND type = 'expense' " +
			"AND created_at BETWEEN ? AND ?), 0) as actual_expense", startOfMonth, endOfMonth).
		Joins("JOIN categories ON budgets.category_id = categories.id").
		Where("budgets.start_date <= ? AND budgets.end_date >= ?", endOfMonth, startOfMonth).
		Rows()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var budget BudgetOverview
		if err := h.DB.ScanRows(rows, &budget); err != nil {
			continue
		}
		
		budget.Remaining = budget.Amount - budget.ActualExpense
		if budget.Amount > 0 {
			budget.PercentageUsed = (budget.ActualExpense / budget.Amount) * 100
		}
		budgets = append(budgets, budget)
	}

	c.JSON(http.StatusOK, gin.H{
		"current_month": gin.H{
			"start_date": startOfMonth.Format("2006-01-02"),
			"end_date":   endOfMonth.Format("2006-01-02"),
		},
		"budgets": budgets,
	})
}
