package v1

import "time"

// CreateVariableRequest 创建变量请求
type CreateVariableRequest struct {
	ApplicationID uint   `json:"application_id" binding:"required"`
	Key           string `json:"key" binding:"required"`
	Value         string `json:"value" binding:"required"`
	Description   string `json:"description"`
	Type          string `json:"type"`
	IsSecret      bool   `json:"is_secret"`
}

// UpdateVariableRequest 更新变量请求
type UpdateVariableRequest struct {
	ApplicationID uint   `json:"application_id"`
	Key           string `json:"key"`
	Value         string `json:"value"`
	Description   string `json:"description"`
	Type          string `json:"type"`
	IsSecret      bool   `json:"is_secret"`
}

// VariableResponse 变量响应
type VariableResponse struct {
	ID            uint      `json:"id"`
	ApplicationID uint      `json:"application_id"`
	Key           string    `json:"key"`
	Value         string    `json:"value"`
	Description   string    `json:"description"`
	Type          string    `json:"type"`
	IsSecret      bool      `json:"is_secret"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// VariableListResponse 变量列表响应
type VariableListResponse struct {
	Variables []VariableResponse `json:"variables"`
	Total     int64              `json:"total"`
	Offset    int                `json:"offset"`
	Limit     int                `json:"limit"`
}
