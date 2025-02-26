package models

import (
	"time"
)

type Account struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Name      string    `json:"name" gorm:"not null"`
	Balance   float64   `json:"balance" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Transaction struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	AccountID   uint      `json:"account_id" gorm:"not null"`
	Amount      float64   `json:"amount" gorm:"not null"`
	Type        string    `json:"type" gorm:"not null"` // "income" or "expense"
	CategoryID  uint      `json:"category_id" gorm:"not null"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Account     Account   `json:"account" gorm:"foreignkey:AccountID"`
	Category    Category  `json:"category" gorm:"foreignkey:CategoryID"`
}
