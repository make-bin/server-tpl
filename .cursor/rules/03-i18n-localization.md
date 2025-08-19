# 3. 国际化与本地化（i18n/L10n）

## 3.1 目标
- 为 API 返回提供多语言消息（成功/错误），支持通过请求选择语言。

## 3.2 目录结构
```
pkg/utils/i18n/
└── i18n.go        # i18n 基础实现：Bundle、配置、上下文工具
```

## 3.3 语言解析策略
- 优先级：query 参数 `lang` -> 头 `X-Lang` -> 头 `Accept-Language` -> 默认语言
- 默认语言：`en`
- 默认支持：`en`, `zh-CN`

## 3.4 使用方式
```go
// 路由中启用中间件
engine.Use(middleware.I18n())

// 在响应辅助中自动本地化
middleware.SuccessResponse(c, data)        // message: success/成功
middleware.CreatedResponse(c, data)        // message: created/创建成功
middleware.ErrorResponse(c, err)           // message: 根据错误码映射到 error.xxx
```

## 3.5 API 返回规范（带 i18n）
- 统一成功响应：
```json
{
  "code": 200,
  "message": "success",  // 或 "成功"
  "data": { ... },
  "locale": "en"          // 实际解析的语言
}
```

- 统一错误响应：
```json
{
  "code": 404,
  "message": "Not found", // 或 "未找到"
  "details": "...",       // 可选
  "data": null,            // 可选
  "locale": "zh-CN"
}
```

## 3.6 文本与键
- 公共键：`common.success`, `common.created`, `service.running`
- 错误键：`error.400`, `error.401`, `error.403`, `error.404`, `error.409`, `error.422`, `error.500`
- 自定义业务错误可映射到上述通用键或新增键。

## 3.7 扩展语言
```go
bundle := i18n.Default()
bundle.AddMessages("ja", map[string]string{
    "common.success": "成功",
})
```
