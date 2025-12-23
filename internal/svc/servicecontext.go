package svc

import (
	"log"

	"yusi-backend/internal/config"
	"yusi-backend/internal/database"
	"yusi-backend/internal/middleware"
	"yusi-backend/internal/websocket"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	Auth   rest.Middleware
	DB     *gorm.DB
	Redis  *redis.Client
	WsHub  *websocket.Hub
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化数据库（自动创建数据库和表）
	db, err := database.InitDB(c.Mysql.DataSource)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 初始化 Redis
	rdb, err := database.InitRedis(c.Redis.Host, c.Redis.Pass)
	if err != nil {
		log.Fatalf("初始化 Redis 失败: %v", err)
	}

	// 初始化 WebSocket Hub
	hub := websocket.NewHub()
	go hub.Run()

	return &ServiceContext{
		Config: c,
		Auth:   middleware.NewAuthMiddleware(c.Auth.AccessSecret).Handle,
		DB:     db,
		Redis:  rdb,
		WsHub:  hub,
	}
}
