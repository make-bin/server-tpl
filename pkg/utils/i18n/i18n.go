package i18n

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Supported languages
const (
	LanguageZhCN = "zh-CN"
	LanguageZhTW = "zh-TW"
	LanguageEnUS = "en-US"
	LanguageEnGB = "en-GB"
	LanguageJaJP = "ja-JP"
	LanguageKoKR = "ko-KR"
	LanguageFrFR = "fr-FR"
	LanguageDeDE = "de-DE"
	LanguageEsES = "es-ES"
	LanguagePtBR = "pt-BR"
)

// DefaultLanguage is the default language
const DefaultLanguage = LanguageZhCN

// LanguageMap maps language codes to display names
var LanguageMap = map[string]string{
	LanguageZhCN: "简体中文",
	LanguageZhTW: "繁體中文",
	LanguageEnUS: "English (US)",
	LanguageEnGB: "English (UK)",
	LanguageJaJP: "日本語",
	LanguageKoKR: "한국어",
	LanguageFrFR: "Français",
	LanguageDeDE: "Deutsch",
	LanguageEsES: "Español",
	LanguagePtBR: "Português",
}

// Translator interface for translation operations
type Translator interface {
	Translate(key string, args ...interface{}) string
	TranslateWithLang(lang, key string, args ...interface{}) string
	GetLanguage() string
	SetLanguage(lang string) error
	GetSupportedLanguages() []string
	HasTranslation(key string) bool
	Reload() error
}

// Localizer interface for localization operations
type Localizer interface {
	FormatNumber(number interface{}) string
	FormatCurrency(amount interface{}, currency string) string
	FormatDate(date time.Time, format string) string
	FormatTime(date time.Time, format string) string
	FormatDateTime(date time.Time, format string) string
	ParseDate(dateStr, format string) (time.Time, error)
	GetTimeZone() *time.Location
	SetTimeZone(tz string) error
}

// I18nManager implements both Translator and Localizer interfaces
type I18nManager struct {
	translations map[string]map[string]interface{}
	currentLang  string
	localesPath  string
	timeZone     *time.Location
	mutex        sync.RWMutex
}

// NewTranslator creates a new translator instance
func NewTranslator(localesPath string) Translator {
	manager := &I18nManager{
		translations: make(map[string]map[string]interface{}),
		currentLang:  DefaultLanguage,
		localesPath:  localesPath,
		timeZone:     time.UTC,
	}

	// Load translations
	manager.loadTranslations()

	return manager
}

// NewLocalizer creates a new localizer instance
func NewLocalizer(language, timeZone string) Localizer {
	location, err := time.LoadLocation(timeZone)
	if err != nil {
		location = time.UTC
	}

	manager := &I18nManager{
		translations: make(map[string]map[string]interface{}),
		currentLang:  language,
		timeZone:     location,
	}

	return manager
}

// Translate translates a key with the current language
func (i *I18nManager) Translate(key string, args ...interface{}) string {
	return i.TranslateWithLang(i.currentLang, key, args...)
}

// TranslateWithLang translates a key with the specified language
func (i *I18nManager) TranslateWithLang(lang, key string, args ...interface{}) string {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	langTranslations, exists := i.translations[lang]
	if !exists {
		// Fallback to default language
		if langTranslations, exists = i.translations[DefaultLanguage]; !exists {
			return key
		}
	}

	value := i.getNestedValue(langTranslations, key)
	if value == "" {
		// Fallback to default language if not found
		if lang != DefaultLanguage {
			if defaultTranslations, exists := i.translations[DefaultLanguage]; exists {
				if value = i.getNestedValue(defaultTranslations, key); value == "" {
					return key
				}
			} else {
				return key
			}
		} else {
			return key
		}
	}

	// Apply arguments if provided
	if len(args) > 0 {
		return fmt.Sprintf(value, args...)
	}

	return value
}

// GetLanguage returns the current language
func (i *I18nManager) GetLanguage() string {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	return i.currentLang
}

// SetLanguage sets the current language
func (i *I18nManager) SetLanguage(lang string) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	if _, exists := i.translations[lang]; !exists {
		return fmt.Errorf("language %s not supported", lang)
	}

	i.currentLang = lang
	return nil
}

// GetSupportedLanguages returns all supported languages
func (i *I18nManager) GetSupportedLanguages() []string {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	var languages []string
	for lang := range i.translations {
		languages = append(languages, lang)
	}
	return languages
}

// HasTranslation checks if a translation exists for the given key
func (i *I18nManager) HasTranslation(key string) bool {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	langTranslations, exists := i.translations[i.currentLang]
	if !exists {
		return false
	}

	return i.getNestedValue(langTranslations, key) != ""
}

// Reload reloads all translations from files
func (i *I18nManager) Reload() error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	i.translations = make(map[string]map[string]interface{})
	return i.loadTranslations()
}

// loadTranslations loads translations from files
func (i *I18nManager) loadTranslations() error {
	if i.localesPath == "" {
		return nil
	}

	// Load supported languages
	for lang := range LanguageMap {
		langDir := filepath.Join(i.localesPath, lang)

		// Load all JSON files in the language directory
		files, err := filepath.Glob(filepath.Join(langDir, "*.json"))
		if err != nil {
			continue
		}

		langTranslations := make(map[string]interface{})
		for _, file := range files {
			if err := i.loadTranslationFile(file, langTranslations); err != nil {
				fmt.Printf("Failed to load translation file %s: %v\n", file, err)
			}
		}

		if len(langTranslations) > 0 {
			i.translations[lang] = langTranslations
		}
	}

	return nil
}

// loadTranslationFile loads a single translation file
func (i *I18nManager) loadTranslationFile(filename string, target map[string]interface{}) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	var translations map[string]interface{}
	if err := json.Unmarshal(data, &translations); err != nil {
		return err
	}

	// Merge translations
	for key, value := range translations {
		target[key] = value
	}

	return nil
}

// getNestedValue retrieves a nested value using dot notation
func (i *I18nManager) getNestedValue(data map[string]interface{}, key string) string {
	keys := strings.Split(key, ".")
	current := data

	for i, k := range keys {
		if value, exists := current[k]; exists {
			if i == len(keys)-1 {
				// Last key, return the string value
				if str, ok := value.(string); ok {
					return str
				}
				return ""
			} else {
				// Intermediate key, continue traversing
				if nested, ok := value.(map[string]interface{}); ok {
					current = nested
				} else {
					return ""
				}
			}
		} else {
			return ""
		}
	}

	return ""
}

// Localizer implementation

// FormatNumber formats a number according to the current locale
func (i *I18nManager) FormatNumber(number interface{}) string {
	// Simplified number formatting
	return fmt.Sprintf("%v", number)
}

// FormatCurrency formats currency according to the current locale
func (i *I18nManager) FormatCurrency(amount interface{}, currency string) string {
	// Simplified currency formatting
	switch currency {
	case "CNY":
		return fmt.Sprintf("¥%v", amount)
	case "USD":
		return fmt.Sprintf("$%v", amount)
	case "EUR":
		return fmt.Sprintf("€%v", amount)
	case "JPY":
		return fmt.Sprintf("¥%v", amount)
	default:
		return fmt.Sprintf("%s %v", currency, amount)
	}
}

// FormatDate formats a date according to the current locale
func (i *I18nManager) FormatDate(date time.Time, format string) string {
	return date.In(i.timeZone).Format(format)
}

// FormatTime formats a time according to the current locale
func (i *I18nManager) FormatTime(date time.Time, format string) string {
	return date.In(i.timeZone).Format(format)
}

// FormatDateTime formats a datetime according to the current locale
func (i *I18nManager) FormatDateTime(date time.Time, format string) string {
	return date.In(i.timeZone).Format(format)
}

// ParseDate parses a date string according to the current locale
func (i *I18nManager) ParseDate(dateStr, format string) (time.Time, error) {
	return time.ParseInLocation(format, dateStr, i.timeZone)
}

// GetTimeZone returns the current timezone
func (i *I18nManager) GetTimeZone() *time.Location {
	return i.timeZone
}

// SetTimeZone sets the current timezone
func (i *I18nManager) SetTimeZone(tz string) error {
	location, err := time.LoadLocation(tz)
	if err != nil {
		return err
	}
	i.timeZone = location
	return nil
}

// Middleware functions

// LanguageMiddleware detects and sets the language for requests
func LanguageMiddleware(translator Translator) gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := detectLanguage(c)

		// Set language in translator
		translator.SetLanguage(lang)

		// Store translator in context
		c.Set("translator", translator)
		c.Set("language", lang)

		c.Next()
	}
}

// detectLanguage detects the language from request
func detectLanguage(c *gin.Context) string {
	// 1. Check query parameter
	if lang := c.Query("lang"); lang != "" && isValidLanguage(lang) {
		return lang
	}

	// 2. Check Accept-Language header
	if acceptLang := c.GetHeader("Accept-Language"); acceptLang != "" {
		if lang := parseAcceptLanguage(acceptLang); lang != "" && isValidLanguage(lang) {
			return lang
		}
	}

	// 3. Check cookie
	if lang, err := c.Cookie("lang"); err == nil && isValidLanguage(lang) {
		return lang
	}

	// 4. Return default language
	return DefaultLanguage
}

// parseAcceptLanguage parses Accept-Language header
func parseAcceptLanguage(acceptLang string) string {
	// Simplified parsing - take the first language
	parts := strings.Split(acceptLang, ",")
	if len(parts) > 0 {
		lang := strings.TrimSpace(strings.Split(parts[0], ";")[0])
		return normalizeLanguage(lang)
	}
	return ""
}

// normalizeLanguage normalizes language code
func normalizeLanguage(lang string) string {
	lang = strings.ToLower(lang)
	switch {
	case strings.HasPrefix(lang, "zh-cn") || strings.HasPrefix(lang, "zh_cn"):
		return LanguageZhCN
	case strings.HasPrefix(lang, "zh-tw") || strings.HasPrefix(lang, "zh_tw"):
		return LanguageZhTW
	case strings.HasPrefix(lang, "en-us") || strings.HasPrefix(lang, "en_us"):
		return LanguageEnUS
	case strings.HasPrefix(lang, "en-gb") || strings.HasPrefix(lang, "en_gb"):
		return LanguageEnGB
	case strings.HasPrefix(lang, "ja"):
		return LanguageJaJP
	case strings.HasPrefix(lang, "ko"):
		return LanguageKoKR
	case strings.HasPrefix(lang, "fr"):
		return LanguageFrFR
	case strings.HasPrefix(lang, "de"):
		return LanguageDeDE
	case strings.HasPrefix(lang, "es"):
		return LanguageEsES
	case strings.HasPrefix(lang, "pt"):
		return LanguagePtBR
	default:
		return ""
	}
}

// isValidLanguage checks if a language is supported
func isValidLanguage(lang string) bool {
	_, exists := LanguageMap[lang]
	return exists
}

// Helper functions for use in handlers

// T translates a key using the translator from context
func T(c *gin.Context, key string, args ...interface{}) string {
	if translator, exists := c.Get("translator"); exists {
		if t, ok := translator.(Translator); ok {
			return t.Translate(key, args...)
		}
	}
	return key
}

// TWithLang translates a key using a specific language
func TWithLang(c *gin.Context, lang, key string, args ...interface{}) string {
	if translator, exists := c.Get("translator"); exists {
		if t, ok := translator.(Translator); ok {
			return t.TranslateWithLang(lang, key, args...)
		}
	}
	return key
}

// GetLanguage returns the current language from context
func GetLanguage(c *gin.Context) string {
	if lang, exists := c.Get("language"); exists {
		if langStr, ok := lang.(string); ok {
			return langStr
		}
	}
	return DefaultLanguage
}
