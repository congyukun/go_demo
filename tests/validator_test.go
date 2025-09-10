package tests

import (
	"go_demo/internal/models"
	"go_demo/pkg/validator"
	"strings"
	"testing"
)

func TestValidatorInit(t *testing.T) {
	// 初始化验证器
	err := validator.Init()
	if err != nil {
		t.Fatalf("验证器初始化失败: %v", err)
	}
	t.Log("验证器初始化成功")
}

func TestMobileValidation(t *testing.T) {
	// 初始化验证器
	if err := validator.Init(); err != nil {
		t.Fatalf("验证器初始化失败: %v", err)
	}

	tests := []struct {
		name     string
		mobile   string
		expected bool
	}{
		{"有效手机号1", "13812345678", true},
		{"有效手机号2", "15987654321", true},
		{"有效手机号3", "18666666666", true},
		{"无效手机号-长度不够", "1381234567", false},
		{"无效手机号-长度过长", "138123456789", false},
		{"无效手机号-不是1开头", "23812345678", false},
		{"无效手机号-第二位不符合", "12812345678", false},
		{"空手机号", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := models.RegisterRequest{
				Username: "testuser",
				Password: "123456",
				Name:     "测试用户",
				Mobile:   tt.mobile,
			}

			err := validator.ValidateStruct(&req)
			if tt.expected && err != nil {
				t.Errorf("期望验证通过，但得到错误: %v", err)
			}
			if !tt.expected && err == nil {
				t.Errorf("期望验证失败，但验证通过了")
			}
			if err != nil {
				t.Logf("验证错误信息: %s", err.Error())
			}
		})
	}
}

func TestRegisterRequestValidation(t *testing.T) {
	// 初始化验证器
	if err := validator.Init(); err != nil {
		t.Fatalf("验证器初始化失败: %v", err)
	}

	tests := []struct {
		name     string
		req      models.RegisterRequest
		expected bool
		contains string // 期望错误信息包含的内容
	}{
		{
			name: "有效注册请求",
			req: models.RegisterRequest{
				Username: "testuser",
				Password: "123456",
				Name:     "测试用户",
				Email:    "test@example.com",
				Mobile:   "13812345678",
			},
			expected: true,
		},
		{
			name: "用户名太短",
			req: models.RegisterRequest{
				Username: "ab",
				Password: "123456",
				Name:     "测试用户",
				Mobile:   "13812345678",
			},
			expected: false,
			contains: "用户名",
		},
		{
			name: "用户名太长",
			req: models.RegisterRequest{
				Username: "abcdefghijklmnopqrstuvwxyz",
				Password: "123456",
				Name:     "测试用户",
				Mobile:   "13812345678",
			},
			expected: false,
			contains: "用户名",
		},
		{
			name: "密码为空",
			req: models.RegisterRequest{
				Username: "testuser",
				Password: "",
				Name:     "测试用户",
				Mobile:   "13812345678",
			},
			expected: false,
			contains: "密码",
		},
		{
			name: "密码太短",
			req: models.RegisterRequest{
				Username: "testuser",
				Password: "12345",
				Name:     "测试用户",
				Mobile:   "13812345678",
			},
			expected: false,
			contains: "密码",
		},
		{
			name: "姓名为空",
			req: models.RegisterRequest{
				Username: "testuser",
				Password: "123456",
				Name:     "",
				Mobile:   "13812345678",
			},
			expected: false,
			contains: "姓名",
		},
		{
			name: "无效邮箱",
			req: models.RegisterRequest{
				Username: "testuser",
				Password: "123456",
				Name:     "测试用户",
				Email:    "invalid-email",
				Mobile:   "13812345678",
			},
			expected: false,
			contains: "邮箱",
		},
		{
			name: "无效手机号",
			req: models.RegisterRequest{
				Username: "testuser",
				Password: "123456",
				Name:     "测试用户",
				Mobile:   "12345",
			},
			expected: false,
			contains: "手机号",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateStruct(&tt.req)
			if tt.expected && err != nil {
				t.Errorf("期望验证通过，但得到错误: %v", err)
			}
			if !tt.expected && err == nil {
				t.Errorf("期望验证失败，但验证通过了")
			}
			if err != nil {
				t.Logf("验证错误信息: %s", err.Error())
				if tt.contains != "" && !strings.Contains(err.Error(), tt.contains) {
					t.Errorf("期望错误信息包含 '%s'，但实际错误信息为: %s", tt.contains, err.Error())
				}
			}
		})
	}
}

func TestLoginRequestValidation(t *testing.T) {
	// 初始化验证器
	if err := validator.Init(); err != nil {
		t.Fatalf("验证器初始化失败: %v", err)
	}

	tests := []struct {
		name     string
		req      models.LoginRequest
		expected bool
		contains string
	}{
		{
			name: "有效登录请求",
			req: models.LoginRequest{
				Username: "testuser",
				Password: "123456",
			},
			expected: true,
		},
		{
			name: "用户名为空",
			req: models.LoginRequest{
				Username: "",
				Password: "123456",
			},
			expected: false,
			contains: "用户名",
		},
		{
			name: "用户名太短",
			req: models.LoginRequest{
				Username: "ab",
				Password: "123456",
			},
			expected: false,
			contains: "用户名",
		},
		{
			name: "密码为空",
			req: models.LoginRequest{
				Username: "testuser",
				Password: "",
			},
			expected: false,
			contains: "密码",
		},
		{
			name: "密码太短",
			req: models.LoginRequest{
				Username: "testuser",
				Password: "12345",
			},
			expected: false,
			contains: "密码",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateStruct(&tt.req)
			if tt.expected && err != nil {
				t.Errorf("期望验证通过，但得到错误: %v", err)
			}
			if !tt.expected && err == nil {
				t.Errorf("期望验证失败，但验证通过了")
			}
			if err != nil {
				t.Logf("验证错误信息: %s", err.Error())
				if tt.contains != "" && !strings.Contains(err.Error(), tt.contains) {
					t.Errorf("期望错误信息包含 '%s'，但实际错误信息为: %s", tt.contains, err.Error())
				}
			}
		})
	}
}

func TestUpdateUserRequestValidation(t *testing.T) {
	// 初始化验证器
	if err := validator.Init(); err != nil {
		t.Fatalf("验证器初始化失败: %v", err)
	}

	status0 := 0
	status1 := 1
	status2 := 2

	tests := []struct {
		name     string
		req      models.UpdateUserRequest
		expected bool
		contains string
	}{
		{
			name: "有效更新请求",
			req: models.UpdateUserRequest{
				Email:  "test@example.com",
				Name:   "新姓名",
				Status: &status1,
			},
			expected: true,
		},
		{
			name: "无效邮箱",
			req: models.UpdateUserRequest{
				Email: "invalid-email",
			},
			expected: false,
			contains: "邮箱",
		},
		{
			name: "无效状态",
			req: models.UpdateUserRequest{
				Status: &status2,
			},
			expected: false,
			contains: "状态",
		},
		{
			name: "有效状态0",
			req: models.UpdateUserRequest{
				Status: &status0,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateStruct(&tt.req)
			if tt.expected && err != nil {
				t.Errorf("期望验证通过，但得到错误: %v", err)
			}
			if !tt.expected && err == nil {
				t.Errorf("期望验证失败，但验证通过了")
			}
			if err != nil {
				t.Logf("验证错误信息: %s", err.Error())
				if tt.contains != "" && !strings.Contains(err.Error(), tt.contains) {
					t.Errorf("期望错误信息包含 '%s'，但实际错误信息为: %s", tt.contains, err.Error())
				}
			}
		})
	}
}