package database

import (
	"log"
	"personal-finance/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func InitDB(cfg *config.Config) *gorm.DB {
	// 使用原生 SQL 创建表
	db, err := gorm.Open("sqlite3", cfg.DBPath)
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// 设置连接池
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	// 启用详细日志
	db.LogMode(cfg.GinMode == "debug")

	// 创建表
	db.Exec(`
		CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			type TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS budgets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			category_id INTEGER NOT NULL,
			amount REAL NOT NULL,
			start_date TEXT NOT NULL,
			end_date TEXT NOT NULL,
			FOREIGN KEY (category_id) REFERENCES categories(id)
		);
	`)

	return db
}

