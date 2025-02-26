package main

import (
	"fmt"
	"log"
	"personal-finance/config"
	"personal-finance/database"
	"personal-finance/handlers"
	"personal-finance/middleware"
	"personal-finance/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 设置 Gin 模式
	gin.SetMode(cfg.GinMode)

	// 连接数据库
	db := database.InitDB(cfg)
	defer db.Close()

	// 自动迁移数据库结构
	db.AutoMigrate(
		&models.Account{},
		&models.Category{},
		&models.Transaction{},
		&models.Budget{},
	)

	// 创建路由
	r := gin.New()

	// 使用日志和恢复中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	
	// 使用错误处理中间件
	r.Use(middleware.ErrorHandler())

	// 允许跨域
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 初始化处理器
	accountHandler := &handlers.AccountHandler{DB: db}
	transactionHandler := &handlers.TransactionHandler{DB: db}
	categoryHandler := &handlers.CategoryHandler{DB: db}
	budgetHandler := &handlers.BudgetHandler{DB: db}
	statisticsHandler := &handlers.StatisticsHandler{DB: db}

	// API 版本前缀
	v1 := r.Group("/api/v1")
	{
		// 账户相关路由
		accounts := v1.Group("/accounts")
		{
			accounts.POST("", accountHandler.CreateAccount)
			accounts.GET("", accountHandler.GetAccounts)
			accounts.PUT("/:id", accountHandler.UpdateAccount)
		}

		// 交易相关路由
		transactions := v1.Group("/transactions")
		{
			transactions.POST("", transactionHandler.CreateTransaction)
			transactions.GET("", transactionHandler.GetTransactions)
		}

		// 分类相关路由
		categories := v1.Group("/categories")
		{
			categories.POST("", categoryHandler.CreateCategory)
			categories.GET("", categoryHandler.GetCategories)
			categories.PUT("/:id", categoryHandler.UpdateCategory)
			categories.DELETE("/:id", categoryHandler.DeleteCategory)
		}

		// 预算相关路由
		budgets := v1.Group("/budgets")
		{
			budgets.POST("", budgetHandler.CreateBudget)
			budgets.GET("", budgetHandler.GetBudgets)
			budgets.GET("/:id/status", budgetHandler.GetBudgetStatus)
			budgets.PUT("/:id", budgetHandler.UpdateBudget)
			budgets.DELETE("/:id", budgetHandler.DeleteBudget)
		}

		// 统计相关路由
		stats := v1.Group("/statistics")
		{
			stats.GET("", statisticsHandler.GetStatistics)
			stats.GET("/budget-overview", statisticsHandler.GetBudgetOverview)
		}
	}

	// 添加健康检查端点
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 启动服务器
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on http://localhost%s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
