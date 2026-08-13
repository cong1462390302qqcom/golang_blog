package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"golang_blog/config"
	"golang_blog/models"
	"log"
)

func main() {
	//1. 加载 .env 文件（默认读取当前目录下的 .env）
	err := godotenv.Load() // 相当于 godotenv.Load(".env")
	if err != nil {
		log.Println("⚠️ 未找到 .env 文件，使用系统环境变量")
	}
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.Infof("Checking database tables...")
	// 初始化数据库连接
	config.InitDatabase()
	db := config.GetDB()

	// 检查表是否存在
	tables := []struct {
		model interface{}
		name  string
		desc  string
	}{
		{&models.User{}, "users", "用户表"},
		{&models.Post{}, "posts", "文章表"},
		{&models.Comment{}, "comments", "评论表"},
	}

	fmt.Println("\n=== 数据库表检查结果 ===")
	if db == nil {
		fmt.Println("❌ 数据库连接未初始化（db 是 nil）")
		return
	}
	for _, table := range tables {
		if db.Migrator().HasTable(table.model) {
			fmt.Printf("✓ %s (%s) - 存在\n", table.name, table.desc)

			// 获取表的列信息
			columnTypes, err := db.Migrator().ColumnTypes(table.model)
			if err != nil {
				logrus.WithError(err).Errorf("Failed to get column types for %s", table.name)
				continue
			}

			fmt.Printf("  列信息:\n")
			for _, col := range columnTypes {
				fmt.Printf("    - %s: %s\n", col.Name(), col.DatabaseTypeName())
			}
			fmt.Println()
		} else {
			fmt.Printf("✗ %s (%s) - 不存在\n", table.name, table.desc)
		}
	}

	// 检查外键关系
	fmt.Println("=== 外键关系检查 ===")

	// 检查posts表的user_id外键
	var postCount int64
	db.Model(&models.Post{}).Count(&postCount)
	fmt.Printf("posts表记录数: %d\n", postCount)

	// 检查comments表的外键
	var commentCount int64
	db.Model(&models.Comment{}).Count(&commentCount)
	fmt.Printf("comments表记录数: %d\n", commentCount)

	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	fmt.Printf("users表记录数: %d\n", userCount)

	logrus.Info("Database table check completed!")
}
