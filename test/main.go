package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type YourModel struct {
	ID   int
	Name string
	// 其他欄位...
}

func main() {
	// 連線到 MySQL 資料庫
	dsn := "db_user:db_password@tcp(127.0.0.1:3306)/db_name?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		// 正確的方式是使用 sql.DB 的 Close 方法
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatal(err)
		}
		sqlDB.Close()
	}()

	// 自動遷移（將模型結構映射到資料庫表）
	err = db.AutoMigrate(&YourModel{})
	if err != nil {
		log.Fatal(err)
	}

	// 查詢
	var result YourModel
	err = db.First(&result, 1).Error
	if err != nil {
		log.Fatal(err)
	}

	// 打印查詢結果
	fmt.Printf("ID: %d, Name: %s\n", result.ID, result.Name)
}
