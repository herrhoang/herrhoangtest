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

	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Save(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}
