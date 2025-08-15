package memory

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/make-bin/server-tpl/pkg/infrastructure/datastore"
)

// Memory Memory 数据存储实现
type Memory struct {
	mu       sync.RWMutex
	config   *datastore.Config
	entities map[string]map[string]datastore.Entity // tableName -> primaryKey -> entity
	counters map[string]uint                        // tableName -> next ID
}

// NewMemory 创建 Memory 实例
func NewMemory(config *datastore.Config) *Memory {
	return &Memory{
		config:   config,
		entities: make(map[string]map[string]datastore.Entity),
		counters: make(map[string]uint),
	}
}

// Connect 连接数据库（内存存储无需连接）
func (m *Memory) Connect(ctx context.Context) error {
	// 初始化默认数据
	m.mu.Lock()
	defer m.mu.Unlock()

	// 初始化表结构
	m.entities["users"] = make(map[string]datastore.Entity)
	m.entities["applications"] = make(map[string]datastore.Entity)
	m.entities["variables"] = make(map[string]datastore.Entity)

	// 初始化计数器
	m.counters["users"] = 1
	m.counters["applications"] = 1
	m.counters["variables"] = 1

	return nil
}

// Disconnect 断开数据库连接（内存存储无需断开）
func (m *Memory) Disconnect(ctx context.Context) error {
	return nil
}

// HealthCheck 健康检查
func (m *Memory) HealthCheck(ctx context.Context) error {
	return nil
}

// BeginTx 开始事务
func (m *Memory) BeginTx(ctx context.Context) (datastore.Transaction, error) {
	return &MemoryTransaction{datastore: m}, nil
}

// Add 添加实体
func (m *Memory) Add(ctx context.Context, entity datastore.Entity) error {
	if entity == nil {
		return datastore.ErrNilEntity
	}

	if entity.TableName() == "" {
		return datastore.ErrTableNameEmpty
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	tableName := entity.TableName()
	if m.entities[tableName] == nil {
		m.entities[tableName] = make(map[string]datastore.Entity)
		m.counters[tableName] = 1
	}

	// 设置 ID
	if entity.PrimaryKey() == "" {
		// 使用反射设置 ID
		val := reflect.ValueOf(entity)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		idField := val.FieldByName("ID")
		if idField.IsValid() && idField.CanSet() {
			idField.SetUint(uint64(m.counters[tableName]))
			m.counters[tableName]++
		}
	}

	now := time.Now()
	entity.SetCreateTime(now)
	entity.SetUpdateTime(now)

	primaryKey := entity.PrimaryKey()
	if primaryKey == "" {
		return datastore.ErrPrimaryEmpty
	}

	// 检查是否已存在
	if _, exists := m.entities[tableName][primaryKey]; exists {
		return datastore.ErrRecordExist
	}

	m.entities[tableName][primaryKey] = entity
	return nil
}

// BatchAdd 批量添加实体
func (m *Memory) BatchAdd(ctx context.Context, entities []datastore.Entity) error {
	if len(entities) == 0 {
		return nil
	}

	for _, entity := range entities {
		if err := m.Add(ctx, entity); err != nil {
			return err
		}
	}

	return nil
}

// Put 更新实体
func (m *Memory) Put(ctx context.Context, entity datastore.Entity) error {
	if entity == nil {
		return datastore.ErrNilEntity
	}

	if entity.TableName() == "" {
		return datastore.ErrTableNameEmpty
	}

	if entity.PrimaryKey() == "" {
		return datastore.ErrPrimaryEmpty
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	tableName := entity.TableName()
	primaryKey := entity.PrimaryKey()

	if m.entities[tableName] == nil {
		return datastore.ErrRecordNotExist
	}

	if _, exists := m.entities[tableName][primaryKey]; !exists {
		return datastore.ErrRecordNotExist
	}

	entity.SetUpdateTime(time.Now())
	m.entities[tableName][primaryKey] = entity

	return nil
}

// Delete 删除实体
func (m *Memory) Delete(ctx context.Context, entity datastore.Entity) error {
	if entity == nil {
		return datastore.ErrNilEntity
	}

	if entity.TableName() == "" {
		return datastore.ErrTableNameEmpty
	}

	if entity.PrimaryKey() == "" {
		return datastore.ErrPrimaryEmpty
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	tableName := entity.TableName()
	primaryKey := entity.PrimaryKey()

	if m.entities[tableName] == nil {
		return datastore.ErrRecordNotExist
	}

	if _, exists := m.entities[tableName][primaryKey]; !exists {
		return datastore.ErrRecordNotExist
	}

	delete(m.entities[tableName], primaryKey)
	return nil
}

// Get 获取实体
func (m *Memory) Get(ctx context.Context, entity datastore.Entity) error {
	if entity == nil {
		return datastore.ErrNilEntity
	}

	if entity.TableName() == "" {
		return datastore.ErrTableNameEmpty
	}

	if entity.PrimaryKey() == "" {
		return datastore.ErrPrimaryEmpty
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	tableName := entity.TableName()
	primaryKey := entity.PrimaryKey()

	if m.entities[tableName] == nil {
		return datastore.ErrRecordNotExist
	}

	storedEntity, exists := m.entities[tableName][primaryKey]
	if !exists {
		return datastore.ErrRecordNotExist
	}

	// 复制数据到目标实体
	reflect.ValueOf(entity).Elem().Set(reflect.ValueOf(storedEntity).Elem())
	return nil
}

// List 列出实体
func (m *Memory) List(ctx context.Context, query datastore.Entity, options *datastore.ListOptions) ([]datastore.Entity, error) {
	if query == nil {
		return nil, datastore.ErrNilEntity
	}

	if query.TableName() == "" {
		return nil, datastore.ErrTableNameEmpty
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	tableName := query.TableName()
	if m.entities[tableName] == nil {
		return []datastore.Entity{}, nil
	}

	var entities []datastore.Entity
	for _, entity := range m.entities[tableName] {
		// 应用过滤条件
		if options != nil && !m.matchesFilter(entity, &options.FilterOptions) {
			continue
		}
		entities = append(entities, entity)
	}

	// 应用排序
	if options != nil && len(options.SortBy) > 0 {
		m.sortEntities(entities, options.SortBy)
	}

	// 应用分页
	if options != nil && options.Page > 0 && options.PageSize > 0 {
		start := (options.Page - 1) * options.PageSize
		end := start + options.PageSize
		if start >= len(entities) {
			return []datastore.Entity{}, nil
		}
		if end > len(entities) {
			end = len(entities)
		}
		entities = entities[start:end]
	}

	return entities, nil
}

// Count 统计实体数量
func (m *Memory) Count(ctx context.Context, entity datastore.Entity, options *datastore.FilterOptions) (int64, error) {
	if entity == nil {
		return 0, datastore.ErrNilEntity
	}

	if entity.TableName() == "" {
		return 0, datastore.ErrTableNameEmpty
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	tableName := entity.TableName()
	if m.entities[tableName] == nil {
		return 0, nil
	}

	var count int64
	for _, entity := range m.entities[tableName] {
		if options == nil || m.matchesFilter(entity, options) {
			count++
		}
	}

	return count, nil
}

// IsExist 检查实体是否存在
func (m *Memory) IsExist(ctx context.Context, entity datastore.Entity) (bool, error) {
	if entity == nil {
		return false, datastore.ErrNilEntity
	}

	if entity.TableName() == "" {
		return false, datastore.ErrTableNameEmpty
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	tableName := entity.TableName()
	if m.entities[tableName] == nil {
		return false, nil
	}

	primaryKey := entity.PrimaryKey()
	if primaryKey == "" {
		// 如果没有主键，检查是否有匹配的实体
		for _, storedEntity := range m.entities[tableName] {
			if m.matchesEntity(entity, storedEntity) {
				return true, nil
			}
		}
		return false, nil
	}

	_, exists := m.entities[tableName][primaryKey]
	return exists, nil
}

// matchesFilter 检查实体是否匹配过滤条件
func (m *Memory) matchesFilter(entity datastore.Entity, options *datastore.FilterOptions) bool {
	if options == nil {
		return true
	}

	index := entity.Index()

	// 模糊查询
	for _, query := range options.Queries {
		if value, exists := index[query.Key]; exists {
			if str, ok := value.(string); ok {
				if !strings.Contains(strings.ToLower(str), strings.ToLower(query.Query)) {
					return false
				}
			}
		}
	}

	// IN 查询
	for _, inQuery := range options.In {
		if value, exists := index[inQuery.Key]; exists {
			found := false
			for _, v := range inQuery.Values {
				if fmt.Sprintf("%v", value) == v {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	}

	// 不存在查询
	for _, notExistQuery := range options.IsNotExist {
		if value, exists := index[notExistQuery.Key]; exists {
			if value != nil && value != "" {
				return false
			}
		}
	}

	return true
}

// matchesEntity 检查两个实体是否匹配
func (m *Memory) matchesEntity(query, stored datastore.Entity) bool {
	queryIndex := query.Index()
	storedIndex := stored.Index()

	for key, queryValue := range queryIndex {
		if storedValue, exists := storedIndex[key]; exists {
			if queryValue != storedValue {
				return false
			}
		} else {
			return false
		}
	}

	return true
}

// sortEntities 排序实体
func (m *Memory) sortEntities(entities []datastore.Entity, sortBy []datastore.SortOption) {
	sort.Slice(entities, func(i, j int) bool {
		for _, sort := range sortBy {
			indexI := entities[i].Index()
			indexJ := entities[j].Index()

			valueI, existsI := indexI[sort.Key]
			valueJ, existsJ := indexJ[sort.Key]

			if !existsI && !existsJ {
				continue
			}
			if !existsI {
				return sort.Order == datastore.SortOrderAscending
			}
			if !existsJ {
				return sort.Order == datastore.SortOrderDescending
			}

			// 比较值
			comparison := m.compareValues(valueI, valueJ)
			if comparison != 0 {
				if sort.Order == datastore.SortOrderAscending {
					return comparison < 0
				}
				return comparison > 0
			}
		}
		return false
	})
}

// compareValues 比较两个值
func (m *Memory) compareValues(a, b interface{}) int {
	switch aVal := a.(type) {
	case string:
		if bVal, ok := b.(string); ok {
			return strings.Compare(aVal, bVal)
		}
	case int:
		if bVal, ok := b.(int); ok {
			if aVal < bVal {
				return -1
			} else if aVal > bVal {
				return 1
			}
			return 0
		}
	case uint:
		if bVal, ok := b.(uint); ok {
			if aVal < bVal {
				return -1
			} else if aVal > bVal {
				return 1
			}
			return 0
		}
	case time.Time:
		if bVal, ok := b.(time.Time); ok {
			if aVal.Before(bVal) {
				return -1
			} else if aVal.After(bVal) {
				return 1
			}
			return 0
		}
	}

	// 默认字符串比较
	return strings.Compare(fmt.Sprintf("%v", a), fmt.Sprintf("%v", b))
}

// MemoryTransaction Memory 事务实现
type MemoryTransaction struct {
	datastore *Memory
}

// Commit 提交事务（内存存储无需提交）
func (t *MemoryTransaction) Commit() error {
	return nil
}

// Rollback 回滚事务（内存存储无需回滚）
func (t *MemoryTransaction) Rollback() error {
	return nil
}

// GetDataStore 获取数据存储
func (t *MemoryTransaction) GetDataStore() datastore.DataStore {
	return t.datastore
}
