package model

// Application 应用模型
type Application struct {
	BaseEntity
	Name        string `json:"name" gorm:"uniqueIndex;not null"`
	Description string `json:"description"`
	Version     string `json:"version" gorm:"not null"`
	Status      string `json:"status" gorm:"default:'active'"`
	CreatedBy   uint   `json:"created_by" gorm:"not null"`
}

// TableName 指定表名
func (Application) TableName() string {
	return "applications"
}

// ShortTableName 获取短表名
func (Application) ShortTableName() string {
	return "application"
}

// Index 获取索引信息
func (a *Application) Index() map[string]interface{} {
	index := a.BaseEntity.Index()
	index["name"] = a.Name
	index["version"] = a.Version
	index["status"] = a.Status
	index["created_by"] = a.CreatedBy
	return index
}
