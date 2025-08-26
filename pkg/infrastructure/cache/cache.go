package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/make-bin/server-tpl/pkg/infrastructure/datastore"
	"github.com/make-bin/server-tpl/pkg/utils/logger"
)

// MemoryCache implements Cache interface using in-memory storage
type MemoryCache struct {
	data   map[string]*cacheItem
	mutex  sync.RWMutex
	config *datastore.CacheConfig
}

type cacheItem struct {
	Value     interface{}
	ExpiresAt time.Time
}

// NewMemoryCache creates a new memory cache instance
func NewMemoryCache(config *datastore.CacheConfig) datastore.Cache {
	cache := &MemoryCache{
		data:   make(map[string]*cacheItem),
		config: config,
	}

	// Start cleanup goroutine
	go cache.cleanup()

	return cache
}

// Get retrieves a value from cache
func (c *MemoryCache) Get(ctx context.Context, key string) (interface{}, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return nil, datastore.ErrNotFound
	}

	// Check expiration
	if time.Now().After(item.ExpiresAt) {
		delete(c.data, key)
		return nil, datastore.ErrNotFound
	}

	return item.Value, nil
}

// Set stores a value in cache with TTL
func (c *MemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if ttl == 0 {
		ttl = c.config.TTL
	}

	c.data[key] = &cacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}

	return nil
}

// Delete removes a value from cache
func (c *MemoryCache) Delete(ctx context.Context, key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, key)
	return nil
}

// Clear removes all values from cache
func (c *MemoryCache) Clear(ctx context.Context) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = make(map[string]*cacheItem)
	return nil
}

// Exists checks if a key exists in cache
func (c *MemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return false, nil
	}

	// Check expiration
	if time.Now().After(item.ExpiresAt) {
		delete(c.data, key)
		return false, nil
	}

	return true, nil
}

// Expire sets TTL for a key
func (c *MemoryCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, exists := c.data[key]
	if !exists {
		return datastore.ErrNotFound
	}

	item.ExpiresAt = time.Now().Add(ttl)
	return nil
}

// cleanup removes expired items
func (c *MemoryCache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		now := time.Now()
		for key, item := range c.data {
			if now.After(item.ExpiresAt) {
				delete(c.data, key)
			}
		}
		c.mutex.Unlock()
	}
}

// RedisCache implements Cache interface using Redis
type RedisCache struct {
	client datastore.Cache // This would be the Redis client from infrastructure/middleware
	config *datastore.CacheConfig
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(config *datastore.CacheConfig, client datastore.Cache) datastore.Cache {
	return &RedisCache{
		client: client,
		config: config,
	}
}

// Get retrieves a value from Redis cache
func (c *RedisCache) Get(ctx context.Context, key string) (interface{}, error) {
	return c.client.Get(ctx, key)
}

// Set stores a value in Redis cache
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if ttl == 0 {
		ttl = c.config.TTL
	}
	return c.client.Set(ctx, key, value, ttl)
}

// Delete removes a value from Redis cache
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Delete(ctx, key)
}

// Clear removes all values from Redis cache
func (c *RedisCache) Clear(ctx context.Context) error {
	return c.client.Clear(ctx)
}

// Exists checks if a key exists in Redis cache
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	return c.client.Exists(ctx, key)
}

// Expire sets TTL for a key in Redis cache
func (c *RedisCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return c.client.Expire(ctx, key, ttl)
}

// CacheManager manages multiple cache layers
type CacheManager struct {
	l1Cache datastore.Cache // Memory cache (L1)
	l2Cache datastore.Cache // Redis cache (L2)
	config  *datastore.CacheConfig
}

// NewCacheManager creates a new cache manager with L1 and L2 caches
func NewCacheManager(config *datastore.CacheConfig) *CacheManager {
	manager := &CacheManager{
		config: config,
	}

	// Always create L1 (memory) cache
	manager.l1Cache = NewMemoryCache(config)

	// Create L2 (Redis) cache if configured
	if config.Type == "redis" {
		// This would need to be implemented with actual Redis client
		logger.Info("Redis cache would be initialized here")
	}

	return manager
}

// Get retrieves value from cache (L1 first, then L2)
func (m *CacheManager) Get(ctx context.Context, key string) (interface{}, error) {
	// Try L1 cache first
	if value, err := m.l1Cache.Get(ctx, key); err == nil {
		return value, nil
	}

	// Try L2 cache if available
	if m.l2Cache != nil {
		if value, err := m.l2Cache.Get(ctx, key); err == nil {
			// Store in L1 cache for faster access
			m.l1Cache.Set(ctx, key, value, time.Minute*5)
			return value, nil
		}
	}

	return nil, datastore.ErrNotFound
}

// Set stores value in both L1 and L2 caches
func (m *CacheManager) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	// Store in L1 cache
	if err := m.l1Cache.Set(ctx, key, value, ttl); err != nil {
		logger.Error("Failed to set L1 cache: %v", err)
	}

	// Store in L2 cache if available
	if m.l2Cache != nil {
		if err := m.l2Cache.Set(ctx, key, value, ttl); err != nil {
			logger.Error("Failed to set L2 cache: %v", err)
		}
	}

	return nil
}

// Delete removes value from both caches
func (m *CacheManager) Delete(ctx context.Context, key string) error {
	// Delete from L1 cache
	m.l1Cache.Delete(ctx, key)

	// Delete from L2 cache if available
	if m.l2Cache != nil {
		m.l2Cache.Delete(ctx, key)
	}

	return nil
}

// Clear removes all values from both caches
func (m *CacheManager) Clear(ctx context.Context) error {
	// Clear L1 cache
	m.l1Cache.Clear(ctx)

	// Clear L2 cache if available
	if m.l2Cache != nil {
		m.l2Cache.Clear(ctx)
	}

	return nil
}

// Exists checks if key exists in any cache
func (m *CacheManager) Exists(ctx context.Context, key string) (bool, error) {
	// Check L1 cache first
	if exists, err := m.l1Cache.Exists(ctx, key); err == nil && exists {
		return true, nil
	}

	// Check L2 cache if available
	if m.l2Cache != nil {
		return m.l2Cache.Exists(ctx, key)
	}

	return false, nil
}

// Expire sets TTL for key in both caches
func (m *CacheManager) Expire(ctx context.Context, key string, ttl time.Duration) error {
	// Set TTL in L1 cache
	m.l1Cache.Expire(ctx, key, ttl)

	// Set TTL in L2 cache if available
	if m.l2Cache != nil {
		m.l2Cache.Expire(ctx, key, ttl)
	}

	return nil
}

// CachedService provides caching wrapper for services
type CachedService struct {
	cache datastore.Cache
}

// NewCachedService creates a new cached service
func NewCachedService(cache datastore.Cache) *CachedService {
	return &CachedService{
		cache: cache,
	}
}

// GetOrSet retrieves from cache or executes function and caches result
func (s *CachedService) GetOrSet(ctx context.Context, key string, ttl time.Duration, fn func() (interface{}, error)) (interface{}, error) {
	// Try to get from cache first
	if value, err := s.cache.Get(ctx, key); err == nil {
		return value, nil
	}

	// Execute function to get value
	value, err := fn()
	if err != nil {
		return nil, err
	}

	// Store in cache
	if err := s.cache.Set(ctx, key, value, ttl); err != nil {
		logger.Error("Failed to cache value for key %s: %v", key, err)
	}

	return value, nil
}

// CacheKey generates a cache key with prefix
func CacheKey(prefix, identifier string) string {
	return fmt.Sprintf("%s:%s", prefix, identifier)
}

// SerializeValue serializes a value for caching
func SerializeValue(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

// DeserializeValue deserializes a cached value
func DeserializeValue(data []byte, target interface{}) error {
	return json.Unmarshal(data, target)
}
