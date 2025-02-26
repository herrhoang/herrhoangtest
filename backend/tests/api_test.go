package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"personal-finance/config"
	"personal-finance/database"
	"personal-finance/handlers"
	"personal-finance/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() (*gin.Engine, *handlers.CategoryHandler, *handlers.StatisticsHandler) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// 配置测试数据库
	cfg := &config.Config{
		DBPath:  ":memory:", // 使用内存数据库进行测试
		GinMode: "test",
	}
	db := database.InitDB(cfg)

	// 自动迁移数据库结构
	db.AutoMigrate(
		&models.Category{},
		&models.Budget{},
		&models.Transaction{},
		&models.Account{},
	)

	categoryHandler := &handlers.CategoryHandler{DB: db}
	statisticsHandler := &handlers.StatisticsHandler{DB: db}

	return r, categoryHandler, statisticsHandler
}

func TestCategoryOperations(t *testing.T) {
	r, h, _ := setupTestRouter()

	// 设置路由
	r.POST("/categories", h.CreateCategory)
	r.GET("/categories", h.GetCategories)
	r.PUT("/categories/:id", h.UpdateCategory)
	r.DELETE("/categories/:id", h.DeleteCategory)

	// 测试创建分类
	t.Run("Create Category", func(t *testing.T) {
		category := models.Category{
			Name: "食品",
			Type: "expense",
			Icon: "food",
		}
		body, _ := json.Marshal(category)
		req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		
		var response models.Category
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.Nil(t, err)
		assert.Equal(t, category.Name, response.Name)
		assert.Equal(t, category.Type, response.Type)
		assert.Equal(t, category.Icon, response.Icon)
		assert.NotZero(t, response.CreatedAt)
		assert.NotZero(t, response.UpdatedAt)
	})

	// 测试创建分类 - 错误情况
	t.Run("Create Category - Invalid Type", func(t *testing.T) {
		category := models.Category{
			Name: "食品",
			Type: "invalid_type", // 无效的类型
		}
		body, _ := json.Marshal(category)
		req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// 测试获取分类列表
	t.Run("Get Categories", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/categories", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var categories []models.Category
		err := json.Unmarshal(w.Body.Bytes(), &categories)
		assert.Nil(t, err)
		assert.NotEmpty(t, categories)
	})

	// 测试按类型筛选分类
	t.Run("Get Categories By Type", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/categories?type=expense", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var categories []models.Category
		err := json.Unmarshal(w.Body.Bytes(), &categories)
		assert.Nil(t, err)
		for _, category := range categories {
			assert.Equal(t, "expense", category.Type)
		}
	})

	// 测试更新分类
	t.Run("Update Category", func(t *testing.T) {
		// 先创建一个分类
		category := models.Category{
			Name: "原始分类",
			Type: "expense",
		}
		body, _ := json.Marshal(category)
		req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var createdCategory models.Category
		json.Unmarshal(w.Body.Bytes(), &createdCategory)

		// 更新分类
		updatedCategory := models.Category{
			Name: "更新后的分类",
			Type: "expense",
		}
		body, _ = json.Marshal(updatedCategory)
		req = httptest.NewRequest("PUT", "/categories/"+fmt.Sprint(createdCategory.ID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Category
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.Nil(t, err)
		assert.Equal(t, updatedCategory.Name, response.Name)
	})

	// 测试删除分类
	t.Run("Delete Category", func(t *testing.T) {
		// 先创建一个分类
		category := models.Category{
			Name: "待删除分类",
			Type: "expense",
		}
		body, _ := json.Marshal(category)
		req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var createdCategory models.Category
		json.Unmarshal(w.Body.Bytes(), &createdCategory)

		// 删除分类
		req = httptest.NewRequest("DELETE", "/categories/"+fmt.Sprint(createdCategory.ID), nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestBudgetOperations(t *testing.T) {
	r, categoryHandler, statisticsHandler := setupTestRouter()

	// 设置路由
	r.POST("/categories", categoryHandler.CreateCategory)
	r.POST("/budgets", statisticsHandler.CreateBudget)
	r.GET("/budgets", statisticsHandler.GetBudgets)
	r.PUT("/budgets/:id", statisticsHandler.UpdateBudget)
	r.DELETE("/budgets/:id", statisticsHandler.DeleteBudget)

	// 首先创建一个分类用于测试
	category := models.Category{
		Name: "食品",
		Type: "expense",
	}
	body, _ := json.Marshal(category)
	req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var createdCategory models.Category
	json.Unmarshal(w.Body.Bytes(), &createdCategory)

	// 测试创建预算
	t.Run("Create Budget", func(t *testing.T) {
		budget := models.BudgetInput{
			CategoryID: createdCategory.ID,
			Amount:     1000.0,
			StartDate:  "2025-02-01",
			EndDate:    "2025-02-28",
		}
		body, _ := json.Marshal(budget)
		req := httptest.NewRequest("POST", "/budgets", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.Budget
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.Nil(t, err)
		assert.Equal(t, budget.Amount, response.Amount)
		assert.Equal(t, budget.StartDate, response.StartDate)
		assert.Equal(t, budget.EndDate, response.EndDate)
		assert.Equal(t, budget.CategoryID, response.CategoryID)
		assert.NotZero(t, response.CreatedAt)
		assert.NotZero(t, response.UpdatedAt)
	})

	// 测试获取预算列表
	t.Run("Get Budgets", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/budgets", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var budgets []models.Budget
		err := json.Unmarshal(w.Body.Bytes(), &budgets)
		assert.Nil(t, err)
		assert.NotEmpty(t, budgets)
	})
}

func TestStatistics(t *testing.T) {
	r, _, h := setupTestRouter()

	// 设置路由
	r.GET("/statistics", h.GetStatistics)
	r.GET("/statistics/budget", h.GetBudgetOverview)

	// 测试获取统计数据
	t.Run("Get Statistics", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/statistics", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var stats models.Statistics
		err := json.Unmarshal(w.Body.Bytes(), &stats)
		assert.Nil(t, err)
		// 验证统计数据的结构
		assert.GreaterOrEqual(t, stats.TotalExpense, 0.0)
		assert.GreaterOrEqual(t, stats.TotalIncome, 0.0)
		assert.Equal(t, stats.NetAmount, stats.TotalIncome-stats.TotalExpense)
	})

	// 测试获取预算概览
	t.Run("Get Budget Overview", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/statistics/budget", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			CurrentMonth struct {
				StartDate string `json:"start_date"`
				EndDate   string `json:"end_date"`
			} `json:"current_month"`
			Budgets []struct {
				models.Budget
				ActualExpense  float64 `json:"actual_expense"`
				Remaining      float64 `json:"remaining"`
				PercentageUsed float64 `json:"percentage_used"`
			} `json:"budgets"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.Nil(t, err)
		// 验证预算概览的结构
		assert.NotEmpty(t, response.CurrentMonth.StartDate)
		assert.NotEmpty(t, response.CurrentMonth.EndDate)
		for _, budget := range response.Budgets {
			assert.Equal(t, budget.Remaining, budget.Amount-budget.ActualExpense)
			if budget.Amount > 0 {
				assert.Equal(t, budget.PercentageUsed, (budget.ActualExpense/budget.Amount)*100)
			}
		}
	})
}
