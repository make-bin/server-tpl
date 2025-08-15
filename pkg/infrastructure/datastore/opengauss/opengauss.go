package opengauss

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/make-bin/server-tpl/pkg/domain/model"
	"github.com/make-bin/server-tpl/pkg/infrastructure/datastore"
)

// OpenGauss OpenGauss 数据存储实现
type OpenGauss struct {
	db     *gorm.DB
	config *datastore.Config
}

// NewOpenGauss 创建 OpenGauss 实例
func NewOpenGauss(config *datastore.Config) *OpenGauss {
	return &OpenGauss{
		config: config,
	}
}

// Connect 连接数据库
func (o *OpenGauss) Connect(ctx context.Context) error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Shanghai",
		o.config.Host, o.config.Port, o.config.User, o.config.Password, o.config.Database, o.config.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to OpenGauss: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(o.config.MaxIdle)
	sqlDB.SetMaxOpenConns(o.config.MaxOpen)
	sqlDB.SetConnMaxLifetime(time.Duration(o.config.Timeout) * time.Second)

	o.db = db

	// 自动迁移表结构
	if err := o.autoMigrate(); err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	return nil
}

// Disconnect 断开数据库连接
func (o *OpenGauss) Disconnect(ctx context.Context) error {
	if o.db != nil {
		sqlDB, err := o.db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// HealthCheck 健康检查
func (o *OpenGauss) HealthCheck(ctx context.Context) error {
	if o.db == nil {
		return fmt.Errorf("database connection is nil")
	}
	sqlDB, err := o.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// BeginTx 开始事务
func (o *OpenGauss) BeginTx(ctx context.Context) (datastore.Transaction, error) {
	tx := o.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &OpenGaussTransaction{tx: tx, datastore: o}, nil
}

// Add 添加实体
func (o *OpenGauss) Add(ctx context.Context, entity datastore.Entity) error {
	if entity == nil {
		return datastore.ErrNilEntity
	}

	if entity.TableName() == "" {
		return datastore.ErrTableNameEmpty
	}

	now := time.Now()
	entity.SetCreateTime(now)
	entity.SetUpdateTime(now)

	result := o.db.WithContext(ctx).Create(entity)
	if result.Error != nil {
		return datastore.NewDBError(result.Error)
	}

	return nil
}

// BatchAdd 批量添加实体
func (o *OpenGauss) BatchAdd(ctx context.Context, entities []datastore.Entity) error {
	if len(entities) == 0 {
		return nil
	}

	now := time.Now()
	for _, entity := range entities {
		if entity == nil {
			return datastore.ErrNilEntity
		}
		if entity.TableName() == "" {
			return datastore.ErrTableNameEmpty
		}
		entity.SetCreateTime(now)
		entity.SetUpdateTime(now)
	}

	result := o.db.WithContext(ctx).CreateInBatches(entities, 100)
	if result.Error != nil {
		return datastore.NewDBError(result.Error)
	}

	return nil
}

// Put 更新实体
func (o *OpenGauss) Put(ctx context.Context, entity datastore.Entity) error {
	if entity == nil {
		return datastore.ErrNilEntity
	}

	if entity.TableName() == "" {
		return datastore.ErrTableNameEmpty
	}

	if entity.PrimaryKey() == "" {
		return datastore.ErrPrimaryEmpty
	}

	entity.SetUpdateTime(time.Now())

	result := o.db.WithContext(ctx).Save(entity)
	if result.Error != nil {
		return datastore.NewDBError(result.Error)
	}

	if result.RowsAffected == 0 {
		return datastore.ErrRecordNotExist
	}

	return nil
}

// Delete 删除实体
func (o *OpenGauss) Delete(ctx context.Context, entity datastore.Entity) error {
	if entity == nil {
		return datastore.ErrNilEntity
	}

	if entity.TableName() == "" {
		return datastore.ErrTableNameEmpty
	}

	if entity.PrimaryKey() == "" {
		return datastore.ErrPrimaryEmpty
	}

	result := o.db.WithContext(ctx).Delete(entity)
	if result.Error != nil {
		return datastore.NewDBError(result.Error)
	}

	if result.RowsAffected == 0 {
		return datastore.ErrRecordNotExist
	}

	return nil
}

// Get 获取实体
func (o *OpenGauss) Get(ctx context.Context, entity datastore.Entity) error {
	if entity == nil {
		return datastore.ErrNilEntity
	}

	if entity.TableName() == "" {
		return datastore.ErrTableNameEmpty
	}

	if entity.PrimaryKey() == "" {
		return datastore.ErrPrimaryEmpty
	}

	result := o.db.WithContext(ctx).First(entity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return datastore.ErrRecordNotExist
		}
		return datastore.NewDBError(result.Error)
	}

	return nil
}

// List 列出实体
func (o *OpenGauss) List(ctx context.Context, query datastore.Entity, options *datastore.ListOptions) ([]datastore.Entity, error) {
	if query == nil {
		return nil, datastore.ErrNilEntity
	}

	if query.TableName() == "" {
		return nil, datastore.ErrTableNameEmpty
	}

	db := o.db.WithContext(ctx).Model(query)

	// 应用过滤条件
	if options != nil {
		db = o.applyFilterOptions(db, &options.FilterOptions)
	}

	// 应用排序
	if options != nil && len(options.SortBy) > 0 {
		for _, sort := range options.SortBy {
			order := "ASC"
			if sort.Order == datastore.SortOrderDescending {
				order = "DESC"
			}
			db = db.Order(fmt.Sprintf("%s %s", sort.Key, order))
		}
	}

	// 应用分页
	if options != nil && options.Page > 0 && options.PageSize > 0 {
		offset := (options.Page - 1) * options.PageSize
		db = db.Offset(offset).Limit(options.PageSize)
	}

	var entities []datastore.Entity
	result := db.Find(&entities)
	if result.Error != nil {
		return nil, datastore.NewDBError(result.Error)
	}

	return entities, nil
}

// Count 统计实体数量
func (o *OpenGauss) Count(ctx context.Context, entity datastore.Entity, options *datastore.FilterOptions) (int64, error) {
	if entity == nil {
		return 0, datastore.ErrNilEntity
	}

	if entity.TableName() == "" {
		return 0, datastore.ErrTableNameEmpty
	}

	db := o.db.WithContext(ctx).Model(entity)

	// 应用过滤条件
	if options != nil {
		db = o.applyFilterOptions(db, options)
	}

	var count int64
	result := db.Count(&count)
	if result.Error != nil {
		return 0, datastore.NewDBError(result.Error)
	}

	return count, nil
}

// IsExist 检查实体是否存在
func (o *OpenGauss) IsExist(ctx context.Context, entity datastore.Entity) (bool, error) {
	if entity == nil {
		return false, datastore.ErrNilEntity
	}

	if entity.TableName() == "" {
		return false, datastore.ErrTableNameEmpty
	}

	var count int64
	result := o.db.WithContext(ctx).Model(entity).Count(&count)
	if result.Error != nil {
		return false, datastore.NewDBError(result.Error)
	}

	return count > 0, nil
}

// applyFilterOptions 应用过滤选项
func (o *OpenGauss) applyFilterOptions(db *gorm.DB, options *datastore.FilterOptions) *gorm.DB {
	if options == nil {
		return db
	}

	// 模糊查询
	for _, query := range options.Queries {
		db = db.Where(fmt.Sprintf("%s ILIKE ?", query.Key), "%"+query.Query+"%")
	}

	// IN 查询
	for _, inQuery := range options.In {
		db = db.Where(fmt.Sprintf("%s IN ?", inQuery.Key), inQuery.Values)
	}

	// 不存在查询
	for _, notExistQuery := range options.IsNotExist {
		db = db.Where(fmt.Sprintf("%s IS NULL OR %s = ''", notExistQuery.Key, notExistQuery.Key))
	}

	return db
}

// autoMigrate 自动迁移表结构
func (o *OpenGauss) autoMigrate() error {
	return o.db.AutoMigrate(
		&model.User{},
		&model.Application{},
		&model.Variable{},
	)
}

// OpenGaussTransaction OpenGauss 事务实现
type OpenGaussTransaction struct {
	tx        *gorm.DB
	datastore *OpenGauss
}

// Commit 提交事务
func (t *OpenGaussTransaction) Commit() error {
	return t.tx.Commit().Error
}

// Rollback 回滚事务
func (t *OpenGaussTransaction) Rollback() error {
	return t.tx.Rollback().Error
}

// GetDataStore 获取数据存储
func (t *OpenGaussTransaction) GetDataStore() datastore.DataStore {
	return t.datastore
}
