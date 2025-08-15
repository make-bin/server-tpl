package model

import (
	"time"
)

// Entity 数据库实体接口
type Entity interface {
	SetCreateTime(time.Time)
	SetUpdateTime(time.Time)
	PrimaryKey() string
	TableName() string
	ShortTableName() string
	Index() map[string]interface{}
}

// BaseEntity 基础实体结构
type BaseEntity struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SetCreateTime 设置创建时间
func (e *BaseEntity) SetCreateTime(t time.Time) {
	e.CreatedAt = t
}

// SetUpdateTime 设置更新时间
func (e *BaseEntity) SetUpdateTime(t time.Time) {
	e.UpdatedAt = t
}

// PrimaryKey 获取主键
func (e *BaseEntity) PrimaryKey() string {
	return string(rune(e.ID))
}

// Index 获取索引信息
func (e *BaseEntity) Index() map[string]interface{} {
	return map[string]interface{}{
		"id": e.ID,
	}
}
