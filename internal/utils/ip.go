package utils

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetClientIP 获取客户端真实IP地址
func GetClientIP(c *gin.Context) string {
	// 1. 检查 X-Forwarded-For 头
	xForwardedFor := c.GetHeader("X-Forwarded-For")
	if xForwardedFor != "" {
		// X-Forwarded-For 可能包含多个IP，取第一个
		ips := net.ParseIP(xForwardedFor)
		if ips != nil {
			return xForwardedFor
		}
	}

	// 2. 检查 X-Real-IP 头
	xRealIP := c.GetHeader("X-Real-IP")
	if xRealIP != "" {
		if net.ParseIP(xRealIP) != nil {
			return xRealIP
		}
	}

	// 3. 检查 X-Forwarded 头
	xForwarded := c.GetHeader("X-Forwarded")
	if xForwarded != "" {
		if net.ParseIP(xForwarded) != nil {
			return xForwarded
		}
	}

	// 4. 检查 Forwarded 头
	forwarded := c.GetHeader("Forwarded")
	if forwarded != "" {
		if net.ParseIP(forwarded) != nil {
			return forwarded
		}
	}

	// 5. 使用 RemoteAddr
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}

	// 验证IP是否有效
	if net.ParseIP(ip) != nil {
		return ip
	}

	// 如果所有方法都失败，返回空字符串
	return ""
}

// GetClientIPFromRequest 从标准http.Request获取客户端IP
func GetClientIPFromRequest(r *http.Request) string {
	// 1. 检查 X-Forwarded-For 头
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		// X-Forwarded-For 可能包含多个IP，取第一个
		ips := net.ParseIP(xForwardedFor)
		if ips != nil {
			return xForwardedFor
		}
	}

	// 2. 检查 X-Real-IP 头
	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" {
		if net.ParseIP(xRealIP) != nil {
			return xRealIP
		}
	}

	// 3. 检查 X-Forwarded 头
	xForwarded := r.Header.Get("X-Forwarded")
	if xForwarded != "" {
		if net.ParseIP(xForwarded) != nil {
			return xForwarded
		}
	}

	// 4. 检查 Forwarded 头
	forwarded := r.Header.Get("Forwarded")
	if forwarded != "" {
		if net.ParseIP(forwarded) != nil {
			return forwarded
		}
	}

	// 5. 使用 RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	// 验证IP是否有效
	if net.ParseIP(ip) != nil {
		return ip
	}

	// 如果所有方法都失败，返回空字符串
	return ""
}
