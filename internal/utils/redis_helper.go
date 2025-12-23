package utils

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisHelper Redis 辅助工具类
type RedisHelper struct {
	client *redis.Client
}

// NewRedisHelper 创建 Redis 辅助工具
func NewRedisHelper(client *redis.Client) *RedisHelper {
	return &RedisHelper{client: client}
}

// Set 设置键值对
func (r *RedisHelper) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, expiration).Err()
}

// Get 获取值
func (r *RedisHelper) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// GetString 获取字符串值
func (r *RedisHelper) GetString(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// SetString 设置字符串值
func (r *RedisHelper) SetString(ctx context.Context, key string, value string, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Delete 删除键
func (r *RedisHelper) Delete(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func (r *RedisHelper) Exists(ctx context.Context, key string) (bool, error) {
	n, err := r.client.Exists(ctx, key).Result()
	return n > 0, err
}

// Expire 设置过期时间
func (r *RedisHelper) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}

// Increment 自增
func (r *RedisHelper) Increment(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

// Decrement 自减
func (r *RedisHelper) Decrement(ctx context.Context, key string) (int64, error) {
	return r.client.Decr(ctx, key).Result()
}

// SetNX 仅当键不存在时设置 (分布式锁)
func (r *RedisHelper) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, value, expiration).Result()
}

// GetTTL 获取键的剩余生存时间
func (r *RedisHelper) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

// HSet 设置哈希表字段
func (r *RedisHelper) HSet(ctx context.Context, key string, field string, value interface{}) error {
	return r.client.HSet(ctx, key, field, value).Err()
}

// HGet 获取哈希表字段值
func (r *RedisHelper) HGet(ctx context.Context, key string, field string) (string, error) {
	return r.client.HGet(ctx, key, field).Result()
}

// HGetAll 获取哈希表所有字段
func (r *RedisHelper) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.client.HGetAll(ctx, key).Result()
}

// HDel 删除哈希表字段
func (r *RedisHelper) HDel(ctx context.Context, key string, fields ...string) error {
	return r.client.HDel(ctx, key, fields...).Err()
}

// LPush 从列表左侧插入
func (r *RedisHelper) LPush(ctx context.Context, key string, values ...interface{}) error {
	return r.client.LPush(ctx, key, values...).Err()
}

// RPush 从列表右侧插入
func (r *RedisHelper) RPush(ctx context.Context, key string, values ...interface{}) error {
	return r.client.RPush(ctx, key, values...).Err()
}

// LPop 从列表左侧弹出
func (r *RedisHelper) LPop(ctx context.Context, key string) (string, error) {
	return r.client.LPop(ctx, key).Result()
}

// RPop 从列表右侧弹出
func (r *RedisHelper) RPop(ctx context.Context, key string) (string, error) {
	return r.client.RPop(ctx, key).Result()
}

// LRange 获取列表指定范围的元素
func (r *RedisHelper) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.client.LRange(ctx, key, start, stop).Result()
}

// SAdd 向集合添加成员
func (r *RedisHelper) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return r.client.SAdd(ctx, key, members...).Err()
}

// SMembers 获取集合所有成员
func (r *RedisHelper) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.client.SMembers(ctx, key).Result()
}

// SIsMember 判断元素是否是集合成员
func (r *RedisHelper) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return r.client.SIsMember(ctx, key, member).Result()
}

// SRem 移除集合成员
func (r *RedisHelper) SRem(ctx context.Context, key string, members ...interface{}) error {
	return r.client.SRem(ctx, key, members...).Err()
}
