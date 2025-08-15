package model

// Variable 变量模型
type Variable struct {
	BaseEntity
	ApplicationID uint   `json:"application_id" gorm:"not null"`
	Key           string `json:"key" gorm:"not null"`
	Value         string `json:"value" gorm:"not null"`
	Description   string `json:"description"`
	Type          string `json:"type" gorm:"default:'string'"`
	IsSecret      bool   `json:"is_secret" gorm:"default:false"`
}

// TableName 指定表名
func (Variable) TableName() string {
	return "variables"
}

// ShortTableName 获取短表名
func (Variable) ShortTableName() string {
	return "variable"
}

// Index 获取索引信息
func (v *Variable) Index() map[string]interface{} {
	index := v.BaseEntity.Index()
	index["application_id"] = v.ApplicationID
	index["key"] = v.Key
	index["type"] = v.Type
	index["is_secret"] = v.IsSecret
	return index
}
