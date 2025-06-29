package util

import (
	"strings"
)

// MaskEmail 邮箱脱敏：前2位 + *** + @domain
func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) < 3 {
		return "***" // fallback
	}
	head := parts[0][:2]
	return head + "***@" + parts[1]
}

// MaskPhone 手机脱敏：保留前三位和后两位
func MaskPhone(phone string) string {
	if len(phone) < 5 {
		return "****"
	}
	return phone[:3] + "****" + phone[len(phone)-2:]
}
