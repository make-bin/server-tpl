package model

// User 用户模型
type User struct {
	BaseEntity
	Email    string `json:"email" gorm:"uniqueIndex;not null"`
	Name     string `json:"name" gorm:"not null"`
	Password string `json:"-" gorm:"not null"`
	Role     string `json:"role" gorm:"default:'user'"`
	Status   string `json:"status" gorm:"default:'active'"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// ShortTableName 获取短表名
func (User) ShortTableName() string {
	return "user"
}

// Index 获取索引信息
func (u *User) Index() map[string]interface{} {
	index := u.BaseEntity.Index()
	index["email"] = u.Email
	index["name"] = u.Name
	index["role"] = u.Role
	index["status"] = u.Status
	return index
}
