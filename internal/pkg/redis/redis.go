package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client Redis 客户端接口
type Client interface {
	// 基础操作
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, keys ...string) (int64, error)
	Exists(ctx context.Context, keys ...string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) (bool, error)
	TTL(ctx context.Context, key string) (time.Duration, error)

	// 列表操作
	LPush(ctx context.Context, key string, values ...interface{}) (int64, error)
	RPush(ctx context.Context, key string, values ...interface{}) (int64, error)
	LPop(ctx context.Context, key string) (string, error)
	RPop(ctx context.Context, key string) (string, error)
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	LLen(ctx context.Context, key string) (int64, error)

	// 哈希操作
	HSet(ctx context.Context, key string, values ...interface{}) (int64, error)
	HGet(ctx context.Context, key, field string) (string, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HDel(ctx context.Context, key string, fields ...string) (int64, error)
	HExists(ctx context.Context, key, field string) (bool, error)

	// 集合操作
	SAdd(ctx context.Context, key string, members ...interface{}) (int64, error)
	SRem(ctx context.Context, key string, members ...interface{}) (int64, error)
	SMembers(ctx context.Context, key string) ([]string, error)
	SIsMember(ctx context.Context, key string, member interface{}) (bool, error)
	SCard(ctx context.Context, key string) (int64, error)

	// 有序集合操作
	ZAdd(ctx context.Context, key string, members ...redis.Z) (int64, error)
	ZRem(ctx context.Context, key string, members ...interface{}) (int64, error)
	ZRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error)
	ZCard(ctx context.Context, key string) (int64, error)

	// 连接管理
	Ping(ctx context.Context) (string, error)
	Close() error

	// 原生客户端访问
	Raw() *redis.Client
}

// redisClient Redis 客户端实现
type redisClient struct {
	client *redis.Client
}

// NewClient 创建 Redis 客户端
func NewClient(client *redis.Client) Client {
	return &redisClient{
		client: client,
	}
}

// Get 获取值
func (c *redisClient) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Set 设置值
func (c *redisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

// Del 删除键
func (c *redisClient) Del(ctx context.Context, keys ...string) (int64, error) {
	return c.client.Del(ctx, keys...).Result()
}

// Exists 检查键是否存在
func (c *redisClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	return c.client.Exists(ctx, keys...).Result()
}

// Expire 设置过期时间
func (c *redisClient) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return c.client.Expire(ctx, key, expiration).Result()
}

// TTL 获取剩余过期时间
func (c *redisClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.client.TTL(ctx, key).Result()
}

// LPush 从左侧推入列表
func (c *redisClient) LPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return c.client.LPush(ctx, key, values...).Result()
}

// RPush 从右侧推入列表
func (c *redisClient) RPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return c.client.RPush(ctx, key, values...).Result()
}

// LPop 从左侧弹出列表
func (c *redisClient) LPop(ctx context.Context, key string) (string, error) {
	return c.client.LPop(ctx, key).Result()
}

// RPop 从右侧弹出列表
func (c *redisClient) RPop(ctx context.Context, key string) (string, error) {
	return c.client.RPop(ctx, key).Result()
}

// LRange 获取列表范围
func (c *redisClient) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.client.LRange(ctx, key, start, stop).Result()
}

// LLen 获取列表长度
func (c *redisClient) LLen(ctx context.Context, key string) (int64, error) {
	return c.client.LLen(ctx, key).Result()
}

// HSet 设置哈希字段
func (c *redisClient) HSet(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return c.client.HSet(ctx, key, values...).Result()
}

// HGet 获取哈希字段
func (c *redisClient) HGet(ctx context.Context, key, field string) (string, error) {
	return c.client.HGet(ctx, key, field).Result()
}

// HGetAll 获取所有哈希字段
func (c *redisClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.client.HGetAll(ctx, key).Result()
}

// HDel 删除哈希字段
func (c *redisClient) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	return c.client.HDel(ctx, key, fields...).Result()
}

// HExists 检查哈希字段是否存在
func (c *redisClient) HExists(ctx context.Context, key, field string) (bool, error) {
	return c.client.HExists(ctx, key, field).Result()
}

// SAdd 添加集合成员
func (c *redisClient) SAdd(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return c.client.SAdd(ctx, key, members...).Result()
}

// SRem 移除集合成员
func (c *redisClient) SRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return c.client.SRem(ctx, key, members...).Result()
}

// SMembers 获取集合所有成员
func (c *redisClient) SMembers(ctx context.Context, key string) ([]string, error) {
	return c.client.SMembers(ctx, key).Result()
}

// SIsMember 检查成员是否在集合中
func (c *redisClient) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return c.client.SIsMember(ctx, key, member).Result()
}

// SCard 获取集合大小
func (c *redisClient) SCard(ctx context.Context, key string) (int64, error) {
	return c.client.SCard(ctx, key).Result()
}

// ZAdd 添加有序集合成员
func (c *redisClient) ZAdd(ctx context.Context, key string, members ...redis.Z) (int64, error) {
	return c.client.ZAdd(ctx, key, members...).Result()
}

// ZRem 移除有序集合成员
func (c *redisClient) ZRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return c.client.ZRem(ctx, key, members...).Result()
}

// ZRange 获取有序集合范围
func (c *redisClient) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.client.ZRange(ctx, key, start, stop).Result()
}

// ZRangeByScore 按分数范围获取有序集合
func (c *redisClient) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error) {
	return c.client.ZRangeByScore(ctx, key, opt).Result()
}

// ZCard 获取有序集合大小
func (c *redisClient) ZCard(ctx context.Context, key string) (int64, error) {
	return c.client.ZCard(ctx, key).Result()
}

// Ping 测试连接
func (c *redisClient) Ping(ctx context.Context) (string, error) {
	return c.client.Ping(ctx).Result()
}

// Close 关闭连接
func (c *redisClient) Close() error {
	return c.client.Close()
}

// Raw 获取原生 Redis 客户端
func (c *redisClient) Raw() *redis.Client {
	return c.client
}
