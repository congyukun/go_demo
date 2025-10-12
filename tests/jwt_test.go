package tests

import (
	"go_demo/internal/utils"
	"testing"
	"time"
)

func TestJWTManager(t *testing.T) {
	// 初始化JWT配置
	config := utils.JWTConfig{
		SecretKey:     "test-secret-key-for-testing",
		AccessExpire:  3600,   // 1小时
		RefreshExpire: 604800, // 7天
		Issuer:        "go_demo_test",
	}

	// 创建JWT管理器
	jwtManager := utils.NewJWTManager(config)

	// 测试数据
	userID := int64(123)
	username := "testuser"

	t.Run("生成访问token", func(t *testing.T) {
		token, err := jwtManager.GenerateAccessToken(userID, username)
		if err != nil {
			t.Fatalf("生成访问token失败: %v", err)
		}
		if token == "" {
			t.Fatal("生成的token为空")
		}
		t.Logf("生成的访问token: %s", token)
	})

	t.Run("生成刷新token", func(t *testing.T) {
		token, err := jwtManager.GenerateRefreshToken(userID)
		if err != nil {
			t.Fatalf("生成刷新token失败: %v", err)
		}
		if token == "" {
			t.Fatal("生成的token为空")
		}
		t.Logf("生成的刷新token: %s", token)
	})

	t.Run("生成token对", func(t *testing.T) {
		accessToken, refreshToken, err := jwtManager.GenerateTokenPair(userID, username)
		if err != nil {
			t.Fatalf("生成token对失败: %v", err)
		}
		if accessToken == "" || refreshToken == "" {
			t.Fatal("生成的token对包含空值")
		}
		t.Logf("访问token: %s", accessToken)
		t.Logf("刷新token: %s", refreshToken)
	})

	t.Run("验证有效token", func(t *testing.T) {
		// 生成token
		token, err := jwtManager.GenerateAccessToken(userID, username)
		if err != nil {
			t.Fatalf("生成token失败: %v", err)
		}

		// 验证token
		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			t.Fatalf("验证token失败: %v", err)
		}

		// 检查claims
		if claims.UserID != userID {
			t.Errorf("用户ID不匹配: 期望 %d, 实际 %d", userID, claims.UserID)
		}
		if claims.Username != username {
			t.Errorf("用户名不匹配: 期望 %s, 实际 %s", username, claims.Username)
		}
	})

	t.Run("验证无效token", func(t *testing.T) {
		invalidToken := "invalid.token.here"
		_, err := jwtManager.ValidateToken(invalidToken)
		if err == nil {
			t.Fatal("应该验证失败，但没有返回错误")
		}
		t.Logf("预期的错误: %v", err)
	})

	t.Run("获取token信息", func(t *testing.T) {
		// 生成token
		token, err := jwtManager.GenerateAccessToken(userID, username)
		if err != nil {
			t.Fatalf("生成token失败: %v", err)
		}

		// 获取用户ID
		extractedUserID, err := jwtManager.GetUserIDFromToken(token)
		if err != nil {
			t.Fatalf("获取用户ID失败: %v", err)
		}
		if extractedUserID != userID {
			t.Errorf("用户ID不匹配: 期望 %d, 实际 %d", userID, extractedUserID)
		}

		// 获取用户名
		extractedUsername, err := jwtManager.GetUsernameFromToken(token)
		if err != nil {
			t.Fatalf("获取用户名失败: %v", err)
		}
		if extractedUsername != username {
			t.Errorf("用户名不匹配: 期望 %s, 实际 %s", username, extractedUsername)
		}

		// 获取过期时间
		expireTime, err := jwtManager.GetTokenExpireTime(token)
		if err != nil {
			t.Fatalf("获取过期时间失败: %v", err)
		}
		if expireTime.Before(time.Now()) {
			t.Error("token已过期")
		}

		// 检查是否过期
		if jwtManager.IsTokenExpired(token) {
			t.Error("token不应该过期")
		}
	})
}

func TestGlobalJWTManager(t *testing.T) {
	// 初始化全局JWT管理器
	config := utils.JWTConfig{
		SecretKey:     "global-test-secret-key",
		AccessExpire:  3600,
		RefreshExpire: 604800,
		Issuer:        "go_demo_global_test",
	}
	utils.InitJWT(config)

	// 测试数据
	userID := int64(456)
	username := "globaluser"

	t.Run("全局函数-生成访问token", func(t *testing.T) {
		token, err := utils.GenerateAccessToken(userID, username)
		if err != nil {
			t.Fatalf("生成访问token失败: %v", err)
		}
		if token == "" {
			t.Fatal("生成的token为空")
		}
	})

	t.Run("全局函数-生成token对", func(t *testing.T) {
		accessToken, refreshToken, err := utils.GetJWTManager().GenerateTokenPair(userID, username)
		if err != nil {
			t.Fatalf("生成token对失败: %v", err)
		}
		if accessToken == "" || refreshToken == "" {
			t.Fatal("生成的token对包含空值")
		}

		// 验证访问token
		claims, err := utils.ValidateToken(accessToken)
		if err != nil {
			t.Fatalf("验证访问token失败: %v", err)
		}
		if claims.UserID != userID || claims.Username != username {
			t.Error("token claims不匹配")
		}
	})
}
