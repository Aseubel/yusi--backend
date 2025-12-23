package database

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// InitRedis 初始化 Redis 客户端
func InitRedis(host, password string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         host,
		Password:     password,
		DB:           0, // 使用默认数据库
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("Redis 连接失败: %w", err)
	}

	fmt.Println("✅ Redis 连接成功")
	return client, nil
}
