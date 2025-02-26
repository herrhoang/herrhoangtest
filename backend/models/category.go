package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Category 交易分类模型
type Category struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Name      string    `json:"name" gorm:"not null"`
	Type      string    `json:"type" gorm:"not null"` // expense 或 income
	Icon      string    `json:"icon"`                 // 分类图标
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Budget 预算模型
// BudgetInput 用于创建预算的输入结构
type BudgetInput struct {
	CategoryID uint    `json:"category_id" binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
	StartDate  string  `json:"start_date" binding:"required"`
	EndDate    string  `json:"end_date" binding:"required"`
}

// Budget 预算模型
type Budget struct {
	gorm.Model
	CategoryID uint     `json:"category_id" gorm:"not null"`
	Amount     float64  `json:"amount" gorm:"not null"`
	StartDate  string   `json:"start_date" gorm:"type:date;not null"`
	EndDate    string   `json:"end_date" gorm:"type:date;not null"`
	Category   Category `json:"category" gorm:"foreignkey:CategoryID"`
}

// Statistics 统计数据结构
type Statistics struct {
	TotalIncome  float64                  `json:"total_income"`
	TotalExpense float64                  `json:"total_expense"`
	NetAmount    float64                  `json:"net_amount"`
	ByCategory   []CategoryStatistics     `json:"by_category"`
	ByMonth      []MonthlyStatistics     `json:"by_month"`
}

// CategoryStatistics 分类统计
type CategoryStatistics struct {
	CategoryID   uint    `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Amount      float64 `json:"amount"`
	Percentage  float64 `json:"percentage"`
}

// MonthlyStatistics 月度统计
type MonthlyStatistics struct {
	Year        int     `json:"year"`
	Month       int     `json:"month"`
	Income      float64 `json:"income"`
	Expense     float64 `json:"expense"`
	NetAmount   float64 `json:"net_amount"`
}
