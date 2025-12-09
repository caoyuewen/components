package util

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"
	"unicode"
)

// HasChinese 检查字符串是否包含中文
func HasChinese(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}

// RemoveSpaces 移除所有空白字符
func RemoveSpaces(input string) string {
	pattern := `[ \t\n\r]+`
	regex := regexp.MustCompile(pattern)
	return regex.ReplaceAllString(input, "")
}

// TrimAll 去除字符串首尾空白并压缩中间空白为单个空格
func TrimAll(s string) string {
	s = strings.TrimSpace(s)
	pattern := regexp.MustCompile(`\s+`)
	return pattern.ReplaceAllString(s, " ")
}

// IsEmpty 检查字符串是否为空（含空白）
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// IsNotEmpty 检查字符串是否非空
func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}

// DefaultIfEmpty 如果为空则返回默认值
func DefaultIfEmpty(s, defaultVal string) string {
	if IsEmpty(s) {
		return defaultVal
	}
	return s
}

// Truncate 截断字符串到指定长度
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

// TruncateWithSuffix 截断字符串并添加后缀（如 ...）
func TruncateWithSuffix(s string, maxLen int, suffix string) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= len(suffix) {
		return suffix[:maxLen]
	}
	return s[:maxLen-len(suffix)] + suffix
}

// Contains 检查字符串是否包含子串（忽略大小写）
func ContainsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// EqualIgnoreCase 比较两个字符串是否相等（忽略大小写）
func EqualIgnoreCase(a, b string) bool {
	return strings.EqualFold(a, b)
}

// Reverse 反转字符串
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// PadLeft 左填充字符串到指定长度
func PadLeft(s string, length int, pad string) string {
	if len(s) >= length {
		return s
	}
	padLen := length - len(s)
	return strings.Repeat(pad, padLen/len(pad)+1)[:padLen] + s
}

// PadRight 右填充字符串到指定长度
func PadRight(s string, length int, pad string) string {
	if len(s) >= length {
		return s
	}
	padLen := length - len(s)
	return s + strings.Repeat(pad, padLen/len(pad)+1)[:padLen]
}

// RandomString 生成随机字符串
func RandomString(length int) string {
	bytes := make([]byte, (length+1)/2)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}

// SplitAndTrim 分割字符串并去除每个元素的空白
func SplitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// JoinNonEmpty 连接非空字符串
func JoinNonEmpty(sep string, strs ...string) string {
	nonEmpty := make([]string, 0, len(strs))
	for _, s := range strs {
		if IsNotEmpty(s) {
			nonEmpty = append(nonEmpty, s)
		}
	}
	return strings.Join(nonEmpty, sep)
}

// FirstN 获取前 N 个字符
func FirstN(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n])
}

// LastN 获取后 N 个字符
func LastN(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[len(runes)-n:])
}

// IsNumeric 检查字符串是否只包含数字
func IsNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// IsAlpha 检查字符串是否只包含字母
func IsAlpha(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// IsAlphanumeric 检查字符串是否只包含字母和数字
func IsAlphanumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// CamelToSnake 驼峰转下划线
func CamelToSnake(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteByte('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// SnakeToCamel 下划线转驼峰
func SnakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	var result strings.Builder
	for i, part := range parts {
		if part == "" {
			continue
		}
		if i == 0 {
			result.WriteString(strings.ToLower(part))
		} else {
			result.WriteString(strings.Title(strings.ToLower(part)))
		}
	}
	return result.String()
}

// SnakeToPascal 下划线转帕斯卡（首字母大写的驼峰）
func SnakeToPascal(s string) string {
	parts := strings.Split(s, "_")
	var result strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		result.WriteString(strings.Title(strings.ToLower(part)))
	}
	return result.String()
}
