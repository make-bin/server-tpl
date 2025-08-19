package i18n

import (
	"context"
	"net/http"
	"sort"
	"strings"
	"sync"

	"golang.org/x/text/language"
)

// localeKey is the context key for storing locale values.
type localeKey struct{}

// Config controls i18n behavior.
type Config struct {
	// Default is the default locale tag string, e.g. "en" or "zh-CN".
	Default string
	// Supported lists supported locale tags. If empty, defaults to ["en", "zh-CN"].
	Supported []string
}

// Bundle stores message catalogs and performs locale matching and translation.
type Bundle struct {
	mu         sync.RWMutex
	messages   map[string]map[string]string // locale -> key -> message
	matcher    language.Matcher
	defaultTag language.Tag
	supported  []language.Tag
}

var defaultBundle *Bundle

// Default returns the default global bundle, initializing it if needed.
func Default() *Bundle {
	if defaultBundle == nil {
		defaultBundle = NewBundle(&Config{})
		// preload minimal common messages for en and zh-CN
		defaultBundle.AddMessages("en", map[string]string{
			"common.success":  "success",
			"common.created":  "created",
			"service.running": "service is running",
			"error.400":       "Bad request",
			"error.401":       "Unauthorized",
			"error.403":       "Forbidden",
			"error.404":       "Not found",
			"error.409":       "Conflict",
			"error.422":       "Validation failed",
			"error.500":       "Internal server error",
		})
		defaultBundle.AddMessages("zh-CN", map[string]string{
			"common.success":  "成功",
			"common.created":  "创建成功",
			"service.running": "服务运行中",
			"error.400":       "错误的请求",
			"error.401":       "未认证",
			"error.403":       "没有权限",
			"error.404":       "未找到",
			"error.409":       "冲突",
			"error.422":       "参数验证失败",
			"error.500":       "服务器内部错误",
		})
	}
	return defaultBundle
}

// SetDefault replaces the global default bundle.
func SetDefault(b *Bundle) {
	defaultBundle = b
}

// NewBundle creates a new Bundle with the provided configuration.
func NewBundle(cfg *Config) *Bundle {
	config := cfg
	if config == nil {
		config = &Config{}
	}
	if len(config.Supported) == 0 {
		config.Supported = []string{"en", "zh-CN"}
	}
	if config.Default == "" {
		config.Default = "en"
	}

	supported := make([]language.Tag, 0, len(config.Supported))
	for _, tag := range config.Supported {
		supported = append(supported, language.Make(tag))
	}
	matcher := language.NewMatcher(supported)

	return &Bundle{
		messages:   make(map[string]map[string]string),
		matcher:    matcher,
		defaultTag: language.Make(config.Default),
		supported:  supported,
	}
}

// AddMessages registers messages for a locale.
func (b *Bundle) AddMessages(locale string, msgs map[string]string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.messages[locale] == nil {
		b.messages[locale] = make(map[string]string)
	}
	for k, v := range msgs {
		b.messages[locale][k] = v
	}
}

// SupportedLocales returns a sorted list of supported locales present in the bundle.
func (b *Bundle) SupportedLocales() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	locales := make([]string, 0, len(b.messages))
	for loc := range b.messages {
		locales = append(locales, loc)
	}
	sort.Strings(locales)
	return locales
}

// Match determines the best locale from an Accept-Language string.
func (b *Bundle) Match(acceptLang string) string {
	if strings.TrimSpace(acceptLang) == "" {
		return b.defaultTag.String()
	}
	tag, _ := language.MatchStrings(b.matcher, acceptLang)
	return tag.String()
}

// ResolveLocaleFromRequest extracts the desired locale from HTTP request.
// Priority: query param "lang" -> header "X-Lang" -> header "Accept-Language" -> default.
func (b *Bundle) ResolveLocaleFromRequest(r *http.Request) string {
	if r == nil {
		return b.defaultTag.String()
	}
	if lang := strings.TrimSpace(r.URL.Query().Get("lang")); lang != "" {
		return lang
	}
	if lang := strings.TrimSpace(r.Header.Get("X-Lang")); lang != "" {
		return lang
	}
	return b.Match(r.Header.Get("Accept-Language"))
}

// Translate returns the localized message for key and locale.
// If not found in the requested locale, falls back to the default locale.
// If still not found, returns the key itself.
func (b *Bundle) Translate(locale, key string) string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if msgs, ok := b.messages[locale]; ok {
		if msg, ok := msgs[key]; ok {
			return msg
		}
	}
	if msgs, ok := b.messages[b.defaultTag.String()]; ok {
		if msg, ok := msgs[key]; ok {
			return msg
		}
	}
	return key
}

// Context helpers

// WithLocale returns a new context that carries the given locale.
func WithLocale(ctx context.Context, locale string) context.Context {
	return context.WithValue(ctx, localeKey{}, locale)
}

// LocaleFromContext retrieves locale from context if present, otherwise returns empty string.
func LocaleFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v := ctx.Value(localeKey{}); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
