package v1

import "time"

// BaseRequest 基础请求结构
type BaseRequest struct {
	RequestID string `json:"request_id,omitempty" example:"req_123456789"`
	Timestamp int64  `json:"timestamp,omitempty" example:"1672531200"`
}

// PageRequest 分页请求结构
// @Description 分页查询参数
type PageRequest struct {
	// @Description 页码，从1开始
	// @Example 1
	Page int `json:"page" form:"page" binding:"omitempty,min=1" example:"1"`

	// @Description 每页数量，1-100
	// @Example 10
	Size int `json:"size" form:"size" binding:"omitempty,min=1,max=100" example:"10"`

	// @Description 排序字段
	// @Example "created_at"
	SortBy string `json:"sort_by" form:"sort_by" binding:"omitempty" example:"created_at"`

	// @Description 排序方向：asc(升序) 或 desc(降序)
	// @Example "desc"
	SortDesc bool `json:"sort_desc" form:"sort_desc" binding:"omitempty" example:"true"`
}

// SearchRequest 搜索请求结构
// @Description 搜索查询参数
type SearchRequest struct {
	// @Description 搜索关键词
	// @Example "搜索关键词"
	Keyword string `json:"keyword" form:"keyword" binding:"omitempty,max=100" example:"搜索关键词"`

	// @Description 过滤条件
	Filters map[string]interface{} `json:"filters,omitempty"`
}

// Response 标准响应结构
// @Description 统一的API响应格式
type Response struct {
	// @Description 请求是否成功
	// @Example true
	Success bool `json:"success" example:"true"`

	// @Description 业务状态码
	// @Example 200
	Code int `json:"code" example:"200"`

	// @Description 响应消息
	// @Example "操作成功"
	Message string `json:"message" example:"操作成功"`

	// @Description 响应数据
	Data interface{} `json:"data,omitempty"`

	// @Description 错误详情
	Error string `json:"error,omitempty"`

	// @Description 附加信息
	Details interface{} `json:"details,omitempty"`

	// @Description 时间戳
	// @Example "2024-01-01T12:00:00Z"
	Timestamp string `json:"timestamp" example:"2024-01-01T12:00:00Z"`

	// @Description 请求ID
	// @Example "req_123456789"
	RequestID string `json:"request_id" example:"req_123456789"`
}

// PaginationResponse 分页响应结构
// @Description 分页数据响应格式
type PaginationResponse struct {
	// @Description 数据列表
	Items interface{} `json:"items"`

	// @Description 分页信息
	Pagination Pagination `json:"pagination"`
}

// Pagination 分页信息
// @Description 分页详细信息
type Pagination struct {
	// @Description 当前页码
	// @Example 1
	Page int `json:"page" example:"1"`

	// @Description 每页数量
	// @Example 10
	Size int `json:"size" example:"10"`

	// @Description 总记录数
	// @Example 100
	Total int `json:"total" example:"100"`

	// @Description 总页数
	// @Example 10
	Pages int `json:"pages" example:"10"`
}

// ErrorDetail 错误详情
// @Description 验证错误详细信息
type ErrorDetail struct {
	// @Description 错误字段
	// @Example "username"
	Field string `json:"field" example:"username"`

	// @Description 错误原因
	// @Example "用户名不能为空"
	Reason string `json:"reason" example:"用户名不能为空"`
}

// HealthCheckResponse 健康检查响应
// @Description 健康检查接口响应
type HealthCheckResponse struct {
	// @Description 服务状态
	// @Example "ok"
	Status string `json:"status" example:"ok"`

	// @Description 响应消息
	// @Example "服务运行正常"
	Message string `json:"message" example:"服务运行正常"`

	// @Description 服务版本
	// @Example "1.0.0"
	Version string `json:"version" example:"1.0.0"`

	// @Description 检查时间
	// @Example "2024-01-01T12:00:00Z"
	Timestamp time.Time `json:"timestamp" example:"2024-01-01T12:00:00Z"`

	// @Description 详细检查信息
	Details map[string]interface{} `json:"details,omitempty"`
}

// GetDefaultPageRequest 获取默认分页请求
func GetDefaultPageRequest() PageRequest {
	return PageRequest{
		Page:     1,
		Size:     10,
		SortBy:   "created_at",
		SortDesc: true,
	}
}

// Validate 验证分页请求参数
func (p *PageRequest) Validate() error {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Size < 1 {
		p.Size = 10
	}
	if p.Size > 100 {
		p.Size = 100
	}
	if p.SortBy == "" {
		p.SortBy = "created_at"
	}
	return nil
}

// GetOffset 获取偏移量
func (p *PageRequest) GetOffset() int {
	return (p.Page - 1) * p.Size
}

// CalculatePages 计算总页数
func (p *PageRequest) CalculatePages(total int) int {
	if p.Size == 0 {
		return 0
	}
	return (total + p.Size - 1) / p.Size
}

// IDRequest ID请求结构
// @Description 通过ID获取资源的请求参数
type IDRequest struct {
	// @Description 资源ID
	// @Example "123"
	ID string `json:"id" uri:"id" binding:"required,min=1" example:"123"`
}

// BatchIDRequest 批量ID请求结构
// @Description 批量操作的ID列表
type BatchIDRequest struct {
	// @Description ID列表
	// @Example ["1", "2", "3"]
	IDs []string `json:"ids" binding:"required,min=1,dive,required" example:"1,2,3"`
}

// BulkOperationResponse 批量操作响应
// @Description 批量操作结果
type BulkOperationResponse struct {
	// @Description 成功处理的数量
	// @Example 5
	SuccessCount int `json:"success_count" example:"5"`

	// @Description 失败处理的数量
	// @Example 2
	FailureCount int `json:"failure_count" example:"2"`

	// @Description 总数量
	// @Example 7
	TotalCount int `json:"total_count" example:"7"`

	// @Description 失败的项目详情
	Failures []BulkFailureItem `json:"failures,omitempty"`
}

// BulkFailureItem 批量操作失败项
// @Description 批量操作中失败的单个项目
type BulkFailureItem struct {
	// @Description 项目标识
	// @Example "item_123"
	ID string `json:"id" example:"item_123"`

	// @Description 失败原因
	// @Example "项目不存在"
	Reason string `json:"reason" example:"项目不存在"`
}

// FileUploadRequest 文件上传请求
// @Description 文件上传参数
type FileUploadRequest struct {
	// @Description 文件描述
	// @Example "用户头像"
	Description string `json:"description" binding:"omitempty,max=200" example:"用户头像"`

	// @Description 文件分类
	// @Example "avatar"
	Category string `json:"category" binding:"omitempty,max=50" example:"avatar"`

	// @Description 是否公开文件
	// @Example true
	Public bool `json:"public" example:"true"`
}

// FileUploadResponse 文件上传响应
// @Description 文件上传结果
type FileUploadResponse struct {
	// @Description 文件ID
	// @Example "file_123456"
	ID string `json:"id" example:"file_123456"`

	// @Description 文件名
	// @Example "avatar.jpg"
	FileName string `json:"file_name" example:"avatar.jpg"`

	// @Description 文件大小（字节）
	// @Example 1024000
	FileSize int64 `json:"file_size" example:"1024000"`

	// @Description 文件类型
	// @Example "image/jpeg"
	ContentType string `json:"content_type" example:"image/jpeg"`

	// @Description 文件URL
	// @Example "https://cdn.example.com/files/avatar.jpg"
	URL string `json:"url" example:"https://cdn.example.com/files/avatar.jpg"`

	// @Description 上传时间
	// @Example "2024-01-01T12:00:00Z"
	UploadedAt time.Time `json:"uploaded_at" example:"2024-01-01T12:00:00Z"`
}
