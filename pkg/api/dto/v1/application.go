package v1

import "time"

// CreateApplicationRequest 创建应用请求
type CreateApplicationRequest struct {
	Name        string `json:"name" binding:"required,namevalidator"`
	Description string `json:"description"`
	Version     string `json:"version" binding:"required,versionvalidator"`
	CreatedBy   uint   `json:"created_by" binding:"required"`
}

// UpdateApplicationRequest 更新应用请求
type UpdateApplicationRequest struct {
	Name        string `json:"name" binding:"namevalidator"`
	Description string `json:"description"`
	Version     string `json:"version" binding:"versionvalidator"`
	Status      string `json:"status"`
}

// ApplicationResponse 应用响应
type ApplicationResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Version     string    `json:"version"`
	Status      string    `json:"status"`
	CreatedBy   uint      `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ApplicationListResponse 应用列表响应
type ApplicationListResponse struct {
	Applications []ApplicationResponse `json:"applications"`
	Total        int64                 `json:"total"`
	Offset       int                   `json:"offset"`
	Limit        int                   `json:"limit"`
}
