package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Todo struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Title     string `json:"title"`
	Completed bool   `gorm:"default:false" json:"completed"`
}

func InitDB() {
	// اتصال بدون رمز به MySQL برای دیتابیس todo_db
	dsn := "root:@tcp(127.0.0.1:3306)/todo_db?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("خطا در اتصال به دیتابیس MySQL:", err)
	}

	fmt.Println("اتصال به MySQL موفقیت‌آمیز بود!")

	// ساخت خودکار جدول todos
	err = DB.AutoMigrate(&Todo{})
	if err != nil {
		log.Fatal("خطا در ساخت جدول (Migration):", err)
	}
	fmt.Println("جدول todos با موفقیت ساخته شد!")
}