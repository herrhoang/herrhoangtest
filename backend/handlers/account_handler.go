package handlers

import (
	"net/http"
	"personal-finance/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type AccountHandler struct {
	DB *gorm.DB
}

// CreateAccount 创建新账户
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var account models.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Create(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, account)
}

// GetAccounts 获取所有账户
func (h *AccountHandler) GetAccounts(c *gin.Context) {
	var accounts []models.Account
	if err := h.DB.Find(&accounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 计算总余额
	var totalBalance float64
	for _, account := range accounts {
		totalBalance += account.Balance
	}

	c.JSON(http.StatusOK, gin.H{
		"accounts": accounts,
		"total_balance": totalBalance,
	})
}

// UpdateAccount 更新账户信息
func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	id := c.Param("id")
	var account models.Account
	
	if err := h.DB.First(&account, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	var input struct {
		Name    string  `json:"name"`
		Balance float64 `json:"balance"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account.Name = input.Name
	account.Balance = input.Balance

	if err := h.DB.Save(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

// DeleteAccount 删除账户
func (h *AccountHandler) DeleteAccount(c *gin.Context) {
	id := c.Param("id")
	var account models.Account
	
	if err := h.DB.First(&account, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	// 检查是否有关联的交易记录
	var transactionCount int64
	h.DB.Model(&models.Transaction{}).Where("account_id = ?", id).Count(&transactionCount)
	if transactionCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法删除有关联交易记录的账户，请先删除相关交易"})
		return
	}

	if err := h.DB.Delete(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}
