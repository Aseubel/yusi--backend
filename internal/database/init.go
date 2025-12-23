package database

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"yusi-backend/model"
)

// InitDB 初始化数据库，自动创建数据库和表
func InitDB(dataSource string) (*gorm.DB, error) {
	// 解析 dataSource，提取数据库名
	dbName := extractDBName(dataSource)
	if dbName == "" {
		return nil, fmt.Errorf("无法从 dataSource 中提取数据库名")
	}

	// 1. 先连接到 MySQL（不指定数据库）
	baseDataSource := strings.Replace(dataSource, "/"+dbName+"?", "/?", 1)
	baseDB, err := gorm.Open(mysql.Open(baseDataSource), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("连接 MySQL 失败: %v", err)
	}

	// 2. 创建数据库（如果不存在）
	createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName)
	if err := baseDB.Exec(createDBSQL).Error; err != nil {
		return nil, fmt.Errorf("创建数据库失败: %v", err)
	}
	log.Printf("数据库 '%s' 检查完成（已存在或已创建）", dbName)

	// 关闭基础连接
	sqlDB, _ := baseDB.DB()
	sqlDB.Close()

	// 3. 连接到指定的数据库
	db, err := gorm.Open(mysql.Open(dataSource), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("连接到数据库 '%s' 失败: %v", dbName, err)
	}

	// 4. 配置连接池
	sqlDB, err = db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库实例失败: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 5. 自动迁移表结构
	if err := autoMigrate(db); err != nil {
		return nil, err
	}

	log.Println("数据库初始化完成")
	return db, nil
}

// autoMigrate 自动创建/更新表结构
func autoMigrate(db *gorm.DB) error {
	log.Println("开始自动迁移表结构...")

	// 注册所有需要迁移的模型
	models := []interface{}{
		&model.User{},
		&model.Diary{},
		&model.SituationRoom{},
		&model.RoomMember{},
		&model.RoomNarrative{},
	}

	for _, m := range models {
		if err := db.AutoMigrate(m); err != nil {
			return fmt.Errorf("迁移表失败: %v", err)
		}
	}

	log.Println("表结构迁移完成")
	return nil
}

// extractDBName 从 dataSource 中提取数据库名
// 例如: "root:password@tcp(127.0.0.1:3306)/yusi?..." -> "yusi"
func extractDBName(dataSource string) string {
	// 查找 ")/" 和 "?" 之间的内容
	start := strings.Index(dataSource, ")/")
	if start == -1 {
		return ""
	}
	start += 2 // 跳过 ")/"

	end := strings.Index(dataSource[start:], "?")
	if end == -1 {
		// 如果没有参数，就取到字符串末尾
		return dataSource[start:]
	}

	return dataSource[start : start+end]
}
