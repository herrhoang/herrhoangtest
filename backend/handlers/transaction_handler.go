package handlers

import (
	"net/http"
	"personal-finance/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type TransactionHandler struct {
	DB *gorm.DB
}

// CreateTransaction 创建新交易
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 开始事务
	tx := h.DB.Begin()

	// 查找相关账户
	var account models.Account
	if err := tx.First(&account, transaction.AccountID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	// 检查分类是否存在
	var category models.Category
	if err := tx.First(&category, transaction.CategoryID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	// 检查分类类型是否与交易类型匹配
	if category.Type != transaction.Type {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category type does not match transaction type"})
		return
	}

	// 更新账户余额
	if transaction.Type == "expense" {
		account.Balance -= transaction.Amount
	} else if transaction.Type == "income" {
		account.Balance += transaction.Amount
	}

	// 保存更新后的账户信息
	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 创建交易记录
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 提交事务
	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"transaction": transaction,
		"new_balance": account.Balance,
	})
}

// GetTransactions 获取交易记录
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	var transactions []models.Transaction
	
	query := h.DB.Preload("Account").Preload("Category").Order("created_at desc")
	
	// 支持按账户ID筛选
	if accountID := c.Query("account_id"); accountID != "" {
		query = query.Where("account_id = ?", accountID)
	}

	// 支持按类型筛选
	if transType := c.Query("type"); transType != "" {
		query = query.Where("type = ?", transType)
	}

	if err := query.Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}
