package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/make-bin/server-tpl/pkg/utils/logger"
)

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// RedisMiddleware Redis中间件
type RedisMiddleware struct {
	client *redis.Client
	config *RedisConfig
}

// NewRedisMiddleware 创建Redis中间件
func NewRedisMiddleware(config *RedisConfig) (*RedisMiddleware, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
		PoolSize: config.PoolSize,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Redis middleware initialized successfully")

	return &RedisMiddleware{
		client: client,
		config: config,
	}, nil
}

// GetClient 获取Redis客户端
func (r *RedisMiddleware) GetClient() *redis.Client {
	return r.client
}

// Set 设置键值对
func (r *RedisMiddleware) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get 获取值
func (r *RedisMiddleware) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Del 删除键
func (r *RedisMiddleware) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func (r *RedisMiddleware) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Exists(ctx, keys...).Result()
}

// Expire 设置过期时间
func (r *RedisMiddleware) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}

// TTL 获取剩余过期时间
func (r *RedisMiddleware) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

// Incr 递增
func (r *RedisMiddleware) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

// IncrBy 按指定值递增
func (r *RedisMiddleware) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return r.client.IncrBy(ctx, key, value).Result()
}

// HSet 设置哈希字段
func (r *RedisMiddleware) HSet(ctx context.Context, key string, values ...interface{}) error {
	return r.client.HSet(ctx, key, values...).Err()
}

// HGet 获取哈希字段
func (r *RedisMiddleware) HGet(ctx context.Context, key, field string) (string, error) {
	return r.client.HGet(ctx, key, field).Result()
}

// HGetAll 获取所有哈希字段
func (r *RedisMiddleware) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.client.HGetAll(ctx, key).Result()
}

// HDel 删除哈希字段
func (r *RedisMiddleware) HDel(ctx context.Context, key string, fields ...string) error {
	return r.client.HDel(ctx, key, fields...).Err()
}

// LPush 左推入列表
func (r *RedisMiddleware) LPush(ctx context.Context, key string, values ...interface{}) error {
	return r.client.LPush(ctx, key, values...).Err()
}

// RPush 右推入列表
func (r *RedisMiddleware) RPush(ctx context.Context, key string, values ...interface{}) error {
	return r.client.RPush(ctx, key, values...).Err()
}

// LPop 左弹出列表
func (r *RedisMiddleware) LPop(ctx context.Context, key string) (string, error) {
	return r.client.LPop(ctx, key).Result()
}

// RPop 右弹出列表
func (r *RedisMiddleware) RPop(ctx context.Context, key string) (string, error) {
	return r.client.RPop(ctx, key).Result()
}

// LRange 获取列表范围
func (r *RedisMiddleware) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.client.LRange(ctx, key, start, stop).Result()
}

// SAdd 添加集合元素
func (r *RedisMiddleware) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return r.client.SAdd(ctx, key, members...).Err()
}

// SMembers 获取集合所有元素
func (r *RedisMiddleware) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.client.SMembers(ctx, key).Result()
}

// SRem 删除集合元素
func (r *RedisMiddleware) SRem(ctx context.Context, key string, members ...interface{}) error {
	return r.client.SRem(ctx, key, members...).Err()
}

// SIsMember 检查集合成员
func (r *RedisMiddleware) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return r.client.SIsMember(ctx, key, member).Result()
}

// ZAdd 添加有序集合元素
func (r *RedisMiddleware) ZAdd(ctx context.Context, key string, score float64, member string) error {
	return r.client.ZAdd(ctx, key, &redis.Z{Score: score, Member: member}).Err()
}

// ZRange 获取有序集合范围
func (r *RedisMiddleware) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.client.ZRange(ctx, key, start, stop).Result()
}

// ZRem 删除有序集合元素
func (r *RedisMiddleware) ZRem(ctx context.Context, key string, members ...interface{}) error {
	return r.client.ZRem(ctx, key, members...).Err()
}

// Close 关闭连接
func (r *RedisMiddleware) Close() error {
	return r.client.Close()
}

// HealthCheck 健康检查
func (r *RedisMiddleware) HealthCheck(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
