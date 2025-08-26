package model

import (
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// BaseModel contains common fields for all domain models
type BaseModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Entity interface defines common methods for all entities
type Entity interface {
	GetID() uint
	SetID(id uint)
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	SetCreateTime(time.Time)
	SetUpdateTime(time.Time)
	PrimaryKey() string
	TableName() string
	ShortTableName() string
	Index() map[string]interface{}
}

// GetID returns the ID of the entity
func (b *BaseModel) GetID() uint {
	return b.ID
}

// SetID sets the ID of the entity
func (b *BaseModel) SetID(id uint) {
	b.ID = id
}

// GetCreatedAt returns the creation time of the entity
func (b *BaseModel) GetCreatedAt() time.Time {
	return b.CreatedAt
}

// GetUpdatedAt returns the last update time of the entity
func (b *BaseModel) GetUpdatedAt() time.Time {
	return b.UpdatedAt
}

// SetCreateTime sets the creation time
func (b *BaseModel) SetCreateTime(t time.Time) {
	b.CreatedAt = t
}

// SetUpdateTime sets the update time
func (b *BaseModel) SetUpdateTime(t time.Time) {
	b.UpdatedAt = t
}

// PrimaryKey returns the primary key as string
func (b *BaseModel) PrimaryKey() string {
	return strconv.FormatUint(uint64(b.ID), 10)
}

// TableName returns the table name (to be overridden by specific models)
func (b *BaseModel) TableName() string {
	return "base_model"
}

// ShortTableName returns abbreviated table name
func (b *BaseModel) ShortTableName() string {
	tableName := b.TableName()
	parts := strings.Split(tableName, "_")
	if len(parts) <= 1 {
		return tableName
	}

	var short strings.Builder
	for _, part := range parts {
		if len(part) > 0 {
			short.WriteByte(part[0])
		}
	}
	return short.String()
}

// Index returns indexable fields (to be overridden by specific models)
func (b *BaseModel) Index() map[string]interface{} {
	return map[string]interface{}{
		"id":         b.ID,
		"created_at": b.CreatedAt,
		"updated_at": b.UpdatedAt,
	}
}

// BeforeCreate GORM hook
func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now
	return nil
}

// BeforeUpdate GORM hook
func (b *BaseModel) BeforeUpdate(tx *gorm.DB) error {
	b.UpdatedAt = time.Now()
	return nil
}
