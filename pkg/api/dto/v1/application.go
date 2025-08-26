package v1

import "time"

// CreateApplicationRequest 创建应用请求
// @Description 创建应用的请求参数
type CreateApplicationRequest struct {
	// @Description 应用名称，1-100个字符
	// @Example "示例应用"
	Name string `json:"name" binding:"required,min=1,max=100,app_name" example:"示例应用"`

	// @Description 应用描述，最多500个字符
	// @Example "这是一个示例应用"
	Description string `json:"description" binding:"omitempty,max=500" example:"这是一个示例应用"`
}

// UpdateApplicationRequest 更新应用请求
// @Description 更新应用的请求参数
type UpdateApplicationRequest struct {
	// @Description 应用名称，1-100个字符
	// @Example "更新后的应用名称"
	Name string `json:"name" binding:"omitempty,min=1,max=100,app_name" example:"更新后的应用名称"`

	// @Description 应用描述，最多500个字符
	// @Example "更新后的应用描述"
	Description string `json:"description" binding:"omitempty,max=500" example:"更新后的应用描述"`
}

// ListApplicationsRequest 应用列表请求
// @Description 获取应用列表的请求参数
type ListApplicationsRequest struct {
	PageRequest
	SearchRequest

	// @Description 应用状态过滤
	// @Example "active"
	Status string `json:"status" form:"status" binding:"omitempty,oneof=active inactive deleted" example:"active"`
}

// ApplicationResponse 应用响应
// @Description 应用详细信息
type ApplicationResponse struct {
	// @Description 应用ID
	// @Example 1
	ID uint `json:"id" example:"1"`

	// @Description 应用名称
	// @Example "示例应用"
	Name string `json:"name" example:"示例应用"`

	// @Description 应用描述
	// @Example "这是一个示例应用"
	Description string `json:"description" example:"这是一个示例应用"`

	// @Description 应用状态
	// @Example "active"
	Status string `json:"status" example:"active"`

	// @Description 创建时间
	// @Example "2024-01-01T12:00:00Z"
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T12:00:00Z"`

	// @Description 更新时间
	// @Example "2024-01-01T12:00:00Z"
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T12:00:00Z"`
}

// ApplicationListResponse 应用列表响应（向后兼容）
// @Description 应用列表响应结构
type ApplicationListResponse struct {
	// @Description 应用列表
	Applications []ApplicationResponse `json:"applications"`

	// @Description 总数量
	// @Example 100
	Total int64 `json:"total" example:"100"`

	// @Description 页码
	// @Example 1
	Page int `json:"page" example:"1"`

	// @Description 每页数量
	// @Example 10
	PageSize int `json:"page_size" example:"10"`
}

// ApplicationStatsResponse 应用统计响应
// @Description 应用统计信息
type ApplicationStatsResponse struct {
	// @Description 总应用数
	// @Example 150
	TotalApps int64 `json:"total_apps" example:"150"`

	// @Description 活跃应用数
	// @Example 120
	ActiveApps int64 `json:"active_apps" example:"120"`

	// @Description 非活跃应用数
	// @Example 20
	InactiveApps int64 `json:"inactive_apps" example:"20"`

	// @Description 已删除应用数
	// @Example 10
	DeletedApps int64 `json:"deleted_apps" example:"10"`

	// @Description 今日新增应用数
	// @Example 5
	TodayNewApps int64 `json:"today_new_apps" example:"5"`

	// @Description 本月新增应用数
	// @Example 25
	MonthNewApps int64 `json:"month_new_apps" example:"25"`
}

// BatchDeleteApplicationsRequest 批量删除应用请求
// @Description 批量删除应用的请求参数
type BatchDeleteApplicationsRequest struct {
	// @Description 应用ID列表
	// @Example [1, 2, 3]
	IDs []uint `json:"ids" binding:"required,min=1,dive,required" example:"1,2,3"`

	// @Description 是否强制删除
	// @Example false
	Force bool `json:"force" example:"false"`
}

// ApplicationBackupRequest 应用备份请求
// @Description 应用备份的请求参数
type ApplicationBackupRequest struct {
	// @Description 备份名称
	// @Example "daily_backup_20240101"
	Name string `json:"name" binding:"required,min=1,max=100" example:"daily_backup_20240101"`

	// @Description 备份描述
	// @Example "每日自动备份"
	Description string `json:"description" binding:"omitempty,max=500" example:"每日自动备份"`

	// @Description 是否包含数据
	// @Example true
	IncludeData bool `json:"include_data" example:"true"`

	// @Description 是否压缩
	// @Example true
	Compress bool `json:"compress" example:"true"`
}

// ApplicationBackupResponse 应用备份响应
// @Description 应用备份结果
type ApplicationBackupResponse struct {
	// @Description 备份ID
	// @Example "backup_123456"
	ID string `json:"id" example:"backup_123456"`

	// @Description 备份名称
	// @Example "daily_backup_20240101"
	Name string `json:"name" example:"daily_backup_20240101"`

	// @Description 备份文件路径
	// @Example "/backups/app_1_20240101.tar.gz"
	FilePath string `json:"file_path" example:"/backups/app_1_20240101.tar.gz"`

	// @Description 备份文件大小（字节）
	// @Example 1048576
	FileSize int64 `json:"file_size" example:"1048576"`

	// @Description 备份状态
	// @Example "completed"
	Status string `json:"status" example:"completed"`

	// @Description 创建时间
	// @Example "2024-01-01T12:00:00Z"
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T12:00:00Z"`
}
