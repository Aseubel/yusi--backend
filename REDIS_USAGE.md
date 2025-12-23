# Redis 使用指南

## 概述

Redis 已集成到 yusi-backend 项目中，可以在任何 handler 或 logic 层通过 `svcCtx.Redis` 访问。

## 配置

Redis 配置位于 `config.yaml`:

```yaml
Redis:
  Host: 127.0.0.1:6379
  Type: node
  Pass: ""
```

## 使用 RedisHelper

### 创建 RedisHelper 实例

```go
import "yusi-backend/internal/utils"

// 在 logic 或 handler 中
redisHelper := utils.NewRedisHelper(l.svcCtx.Redis)
```

### 基本操作

#### 设置和获取字符串

```go
// 设置字符串（带过期时间）
err := redisHelper.SetString(ctx, "user:token:123", "abc123xyz", 24*time.Hour)

// 获取字符串
token, err := redisHelper.GetString(ctx, "user:token:123")
```

#### 设置和获取对象（自动 JSON 序列化）

```go
type User struct {
    ID   int64  `json:"id"`
    Name string `json:"name"`
}

// 设置对象
user := &User{ID: 1, Name: "张三"}
err := redisHelper.Set(ctx, "user:1", user, time.Hour)

// 获取对象
var cachedUser User
err := redisHelper.Get(ctx, "user:1", &cachedUser)
```

### Token 黑名单示例

```go
// 将 token 加入黑名单
func (l *LogoutLogic) Logout(req *types.LogoutReq) error {
    redisHelper := utils.NewRedisHelper(l.svcCtx.Redis)

    // 从 header 获取 token
    token := l.ctx.Request.Header.Get("Authorization")

    // 添加到黑名单，过期时间与 token 相同
    expiration := time.Duration(l.svcCtx.Config.Auth.AccessExpire) * time.Second
    err := redisHelper.SetString(l.ctx, "blacklist:"+token, "1", expiration)
    if err != nil {
        return err
    }

    return nil
}

// 检查 token 是否在黑名单中
func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")

        redisHelper := utils.NewRedisHelper(m.redis)
        exists, err := redisHelper.Exists(r.Context(), "blacklist:"+token)
        if err == nil && exists {
            http.Error(w, "Token 已失效", http.StatusUnauthorized)
            return
        }

        next(w, r)
    }
}
```

### 分布式锁示例

```go
// 获取分布式锁
func (l *CreateRoomLogic) CreateRoom(req *types.CreateRoomReq) error {
    redisHelper := utils.NewRedisHelper(l.svcCtx.Redis)

    lockKey := "lock:room:create:" + strconv.FormatInt(req.UserId, 10)

    // 尝试获取锁，10 秒过期
    locked, err := redisHelper.SetNX(l.ctx, lockKey, "1", 10*time.Second)
    if err != nil {
        return err
    }

    if !locked {
        return errors.New("操作太频繁，请稍后再试")
    }

    // 确保释放锁
    defer redisHelper.Delete(l.ctx, lockKey)

    // 执行创建房间逻辑...

    return nil
}
```

### 缓存日记搜索结果

```go
// 缓存搜索结果
func (l *SearchDiaryLogic) SearchDiary(req *types.SearchDiaryReq) (*types.SearchDiaryResp, error) {
    redisHelper := utils.NewRedisHelper(l.svcCtx.Redis)

    // 生成缓存 key
    cacheKey := fmt.Sprintf("search:diary:%d:%s", req.UserId, req.Keyword)

    // 先从缓存获取
    var cachedResp types.SearchDiaryResp
    err := redisHelper.Get(l.ctx, cacheKey, &cachedResp)
    if err == nil {
        return &cachedResp, nil
    }

    // 缓存未命中，从数据库查询
    // ... 数据库查询逻辑 ...

    // 将结果缓存 5 分钟
    redisHelper.Set(l.ctx, cacheKey, resp, 5*time.Minute)

    return resp, nil
}
```

### 计数器示例

```go
// 用户每日发布日记限制
func (l *CreateDiaryLogic) CreateDiary(req *types.CreateDiaryReq) error {
    redisHelper := utils.NewRedisHelper(l.svcCtx.Redis)

    today := time.Now().Format("2006-01-02")
    countKey := fmt.Sprintf("diary:count:%d:%s", req.UserId, today)

    // 获取今日发布数量
    count, err := redisHelper.Increment(l.ctx, countKey)
    if err != nil {
        return err
    }

    // 第一次计数时设置过期时间（到明天零点）
    if count == 1 {
        tomorrow := time.Now().AddDate(0, 0, 1)
        midnight := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())
        redisHelper.Expire(l.ctx, countKey, time.Until(midnight))
    }

    // 检查限制
    if count > 10 {
        return errors.New("今日发布数量已达上限")
    }

    // 执行创建逻辑...

    return nil
}
```

### 房间在线用户管理（使用 Set）

```go
// 用户加入房间
func (l *JoinRoomLogic) JoinRoom(req *types.JoinRoomReq) error {
    redisHelper := utils.NewRedisHelper(l.svcCtx.Redis)

    roomKey := fmt.Sprintf("room:users:%d", req.RoomId)
    err := redisHelper.SAdd(l.ctx, roomKey, req.UserId)
    if err != nil {
        return err
    }

    return nil
}

// 获取房间在线用户列表
func (l *GetRoomUsersLogic) GetRoomUsers(req *types.GetRoomUsersReq) ([]string, error) {
    redisHelper := utils.NewRedisHelper(l.svcCtx.Redis)

    roomKey := fmt.Sprintf("room:users:%d", req.RoomId)
    users, err := redisHelper.SMembers(l.ctx, roomKey)
    if err != nil {
        return nil, err
    }

    return users, nil
}

// 用户离开房间
func (l *LeaveRoomLogic) LeaveRoom(req *types.LeaveRoomReq) error {
    redisHelper := utils.NewRedisHelper(l.svcCtx.Redis)

    roomKey := fmt.Sprintf("room:users:%d", req.RoomId)
    err := redisHelper.SRem(l.ctx, roomKey, req.UserId)
    if err != nil {
        return err
    }

    return nil
}
```

## RedisHelper 方法列表

### 基本操作
- `Set(ctx, key, value, expiration)` - 设置键值对（自动 JSON 序列化）
- `Get(ctx, key, dest)` - 获取值（自动 JSON 反序列化）
- `SetString(ctx, key, value, expiration)` - 设置字符串
- `GetString(ctx, key)` - 获取字符串
- `Delete(ctx, keys...)` - 删除键
- `Exists(ctx, key)` - 检查键是否存在
- `Expire(ctx, key, expiration)` - 设置过期时间
- `GetTTL(ctx, key)` - 获取剩余生存时间

### 原子操作
- `Increment(ctx, key)` - 自增
- `Decrement(ctx, key)` - 自减
- `SetNX(ctx, key, value, expiration)` - 仅当键不存在时设置（分布式锁）

### Hash 操作
- `HSet(ctx, key, field, value)` - 设置哈希表字段
- `HGet(ctx, key, field)` - 获取哈希表字段值
- `HGetAll(ctx, key)` - 获取哈希表所有字段
- `HDel(ctx, key, fields...)` - 删除哈希表字段

### List 操作
- `LPush(ctx, key, values...)` - 从列表左侧插入
- `RPush(ctx, key, values...)` - 从列表右侧插入
- `LPop(ctx, key)` - 从列表左侧弹出
- `RPop(ctx, key)` - 从列表右侧弹出
- `LRange(ctx, key, start, stop)` - 获取列表指定范围的元素

### Set 操作
- `SAdd(ctx, key, members...)` - 向集合添加成员
- `SMembers(ctx, key)` - 获取集合所有成员
- `SIsMember(ctx, key, member)` - 判断元素是否是集合成员
- `SRem(ctx, key, members...)` - 移除集合成员

## 注意事项

1. 所有操作都需要传入 `context.Context`，建议使用请求的 context
2. 设置过期时间时使用 `time.Duration`，如 `time.Hour`、`24*time.Hour`
3. 对于复杂对象，使用 `Set/Get` 方法会自动处理 JSON 序列化
4. 使用 SetNX 实现分布式锁时，记得设置合理的过期时间防止死锁
5. Redis 连接池已配置，无需手动管理连接

## 启动 Redis

如果本地没有 Redis，可以使用 Docker 快速启动：

```bash
docker run -d --name redis -p 6379:6379 redis:latest
```

或使用项目中的 docker-compose：

```bash
docker-compose up -d redis
```
