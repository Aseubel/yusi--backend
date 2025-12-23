# Yusi Backend (Go-Zero 版本)

这是 Yusi 后端服务的 Go 语言重构版本，使用 go-zero 微服务框架实现。

## 🎯 任务目标

### 基础环境搭建
- [x] 项目结构初始化
- [x] go-zero 框架集成
- [x] 数据库配置（MySQL + GORM）
- [x] 数据模型定义（User, Diary）
- [x] 自动建表和迁移
- [x] Redis 集成
- [ ] 日志配置优化
- [ ] Docker 容器化

### 用户认证模块
- [x] JWT 工具类实现
- [x] 密码加密工具（bcrypt）
- [x] 用户注册接口
- [x] 用户登录接口
- [x] 用户登出接口
- [x] JWT 认证中间件
- [ ] Token 黑名单（Redis）
- [ ] 刷新 Token 机制
- [ ] 密码重置功能
- [ ] 邮箱验证

### 日记管理模块
- [x] 写日记接口
- [x] 编辑日记接口
- [x] 获取日记详情接口
- [x] 获取日记列表接口（分页）
- [ ] 删除日记接口
- [ ] 日记权限控制（只能编辑自己的）
- [ ] 日记搜索功能
- [ ] 日记标签系统
- [ ] 日记导出功能

### AI 功能模块
- [ ] Qwen API 集成
- [ ] 流式响应实现
- [ ] Milvus 向量数据库集成
- [ ] 日记内容向量化
- [ ] 基于 RAG 的对话系统
- [ ] 向量检索优化
- [ ] AI 聊天历史记录

### 情景房间模块
- [ ] 房间数据模型设计
- [ ] 创建房间接口
- [ ] 加入房间接口
- [ ] 开始房间接口
- [ ] 提交叙述接口
- [ ] 生成报告接口
- [ ] 房间状态管理
- [ ] WebSocket 实时通信
- [ ] 房间权限控制

### 代码质量与测试
- [ ] 单元测试覆盖
- [ ] 集成测试
- [ ] API 文档生成（Swagger）
- [ ] 参数验证增强
- [ ] 错误处理统一化
- [ ] 代码注释完善
- [ ] 性能测试
- [ ] 安全审计

### 运维与部署
- [ ] Docker 镜像构建
- [ ] docker-compose 配置
- [ ] CI/CD 流程
- [ ] 监控告警
- [ ] 日志收集
- [ ] 性能监控
- [ ] 自动化部署脚本

## 📋 项目概述

Yusi Backend 是基于 go-zero 框架的后端服务，提供日记管理、情景房间协作与基于向量检索的对话能力。本项目是从 Spring Boot + Java 版本重构而来。

## 🛠 技术栈

- **框架**: go-zero v1.6.0
- **语言**: Go 1.21+
- **数据库**: MySQL 8.x
- **缓存**: Redis
- **ORM**: GORM (推荐)
- **向量数据库**: Milvus/Zilliz Cloud (可选)

## 📁 项目结构

```
yusi-backend/
├── api/                    # API 定义文件
│   └── yusi.api           # go-zero API 定义
├── internal/              # 内部代码
│   ├── config/           # 配置结构
│   ├── handler/          # HTTP 处理器
│   ├── logic/            # 业务逻辑
│   ├── svc/              # 服务上下文
│   └── types/            # 类型定义
├── model/                 # 数据模型
│   └── models.go         # 数据库模型
├── etc/                   # 配置文件
│   └── yusi.yaml         # 主配置文件
├── go.mod                 # Go 模块定义
└── yusi.go               # 程序入口
```

## 🚀 快速开始

### 前置要求

- Go 1.21 或更高版本
- MySQL 8.x
- Redis
- goctl 工具 (可选，用于代码生成)

### 安装 goctl (推荐)

```bash
# macOS/Linux
go install github.com/zeromicro/go-zero/tools/goctl@latest

# 验证安装
goctl --version
```

### 初始化项目

1. **安装依赖**

```bash
cd yusi-backend
go mod tidy
```

2. **配置数据库**

修改 `config.yaml` 中的数据库配置：

```yaml
Mysql:
  DataSource: root:your_password@tcp(127.0.0.1:3306)/yusi?charset=utf8mb4&parseTime=true&loc=Local
```

3. **设置环境变量**

```bash
# Linux/macOS
export QWEN_API_KEY="your-qwen-api-key"
export YUSI_ENCRYPTION_KEY="your-encryption-key-min-16-chars"

# Windows PowerShell
$env:QWEN_API_KEY = "your-qwen-api-key"
$env:YUSI_ENCRYPTION_KEY = "your-encryption-key-min-16-chars"
```

4. **生成代码** (如果已安装 goctl)

```bash
# 从 API 定义生成代码
goctl api go -api api/yusi.api -dir .
```

5. **运行服务**

```bash
go run yusi.go -f yusi.yaml
```

服务将在 `http://localhost:8088` 启动。

### 健康检查

```bash
curl http://localhost:8088/health
```

## 📝 API 接口

### 用户模块 (`/api/user`)

- `POST /api/user/register` - 用户注册
- `POST /api/user/login` - 用户登录
- `POST /api/user/logout` - 用户登出 (需要认证)

### 日记模块 (`/api/diary`)

- `GET /api/diary/list` - 获取日记列表 (需要认证)
- `POST /api/diary` - 写日记 (需要认证)
- `PUT /api/diary` - 编辑日记 (需要认证)
- `GET /api/diary/:diaryId` - 获取日记详情 (需要认证)

### 情景房间模块 (`/api/room`)

- `POST /api/room/create` - 创建房间 (需要认证)
- `POST /api/room/join` - 加入房间 (需要认证)
- `POST /api/room/start` - 开始房间 (需要认证)
- `POST /api/room/submit` - 提交叙述 (需要认证)
- `GET /api/room/report/:code` - 获取报告 (需要认证)

### AI 模块 (`/api/ai`)

- `POST /api/ai/chat/stream` - AI 流式聊天 (需要认证)

## 🔧 开发指南

### 使用 goctl 生成代码

```bash
# 生成 API 代码
goctl api go -api api/yusi.api -dir .

# 生成 Model 代码
goctl model mysql datasource -url "root:password@tcp(127.0.0.1:3306)/yusi" -table "user,diary" -dir ./model
```

### 添加新接口

1. 在 `api/yusi.api` 中定义新的接口
2. 运行 `goctl api go -api api/yusi.api -dir .` 生成代码
3. 在生成的 `logic` 文件中实现业务逻辑

### 数据库迁移

使用 GORM 的 AutoMigrate 或者手动执行 SQL 脚本：

```go
// 在初始化时自动迁移
db.AutoMigrate(&model.User{}, &model.Diary{})
```

## 🔐 认证与授权

项目使用 JWT 进行认证，需要在 `config.yaml` 中配置密钥：

```yaml
Auth:
  AccessSecret: your-secret-key-here
  AccessExpire: 86400  # 24小时
```

## 🌟 特性对比

| 特性 | Java 版本 | Go 版本 |
|------|----------|---------|
| 启动速度 | ~10s | ~1s |
| 内存占用 | ~300MB | ~50MB |
| 并发性能 | 中 | 高 |
| 部署大小 | ~50MB JAR | ~15MB 二进制 |

## 📦 构建与部署

### 编译

```bash
# 编译为二进制文件
go build -o yusi-server yusi.go

# 交叉编译 (Linux)
GOOS=linux GOARCH=amd64 go build -o yusi-server-linux yusi.go
```

### 运行

```bash
./yusi-server -f config.yaml
```

### Docker 部署 (TODO)

```bash
# 构建镜像
docker build -t yusi-backend:latest .

# 运行容器
docker run -p 20611:20611 yusi-backend:latest
```

## 🔄 从 Java 版本迁移

详见 [MIGRATION.md](./MIGRATION.md)

## 📚 相关资源

- [go-zero 官方文档](https://go-zero.dev/)
- [go-zero GitHub](https://github.com/zeromicro/go-zero)
- [goctl 工具使用](https://go-zero.dev/docs/tutorials/cli/overview)

## ⚠️ 注意事项

1. 确保设置 `YUSI_ENCRYPTION_KEY` 环境变量以启用字段加密
2. 生产环境请修改 `Auth.AccessSecret` 为强密钥
3. 不要在配置文件中提交真实的 API Key 和密码
4. Redis 配置根据实际环境调整

## 📄 许可证

与主项目保持一致

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！
