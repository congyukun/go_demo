package tests

import (
	"go_demo/internal/models"
	"go_demo/internal/service"
	"testing"
)

func TestUserModel(t *testing.T) {
	user := &models.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Name:     "Test User",
		Status:   1,
	}

	response := user.ToResponse()

	if response.ID != user.ID {
		t.Errorf("Expected ID %d, got %d", user.ID, response.ID)
	}

	if response.Username != user.Username {
		t.Errorf("Expected Username %s, got %s", user.Username, response.Username)
	}

	if response.Email != user.Email {
		t.Errorf("Expected Email %s, got %s", user.Email, response.Email)
	}
}

func TestAuthService(t *testing.T) {
	authService := service.NewAuthService(nil)

	// 测试token验证
	claims, err := authService.ValidateToken("test_token")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if claims == nil {
		t.Error("Expected claims, got nil")
	}
}

func TestLoginRequest(t *testing.T) {
	req := models.LoginRequest{
		Username: "testuser",
		Password: "123456",
	}

	if req.Username == "" {
		t.Error("Username should not be empty")
	}

	if req.Password == "" {
		t.Error("Password should not be empty")
	}
}

func TestRegisterRequest(t *testing.T) {
	req := models.RegisterRequest{
		Username: "testuser",
		Password: "123456",
		Email:    "test@example.com",
		Name:     "Test User",
	}

	if req.Username == "" {
		t.Error("Username should not be empty")
	}

	if req.Email == "" {
		t.Error("Email should not be empty")
	}

	if req.Name == "" {
		t.Error("Name should not be empty")
	}
}

func TestUpdateUserRequest(t *testing.T) {
	status := 1
	req := models.UpdateUserRequest{
		Email:  "newemail@example.com",
		Name:   "New Name",
		Status: &status,
	}

	if req.Email == "" {
		t.Error("Email should not be empty")
	}

	if req.Name == "" {
		t.Error("Name should not be empty")
	}

	if req.Status == nil {
		t.Error("Status should not be nil")
	}

	if *req.Status != 1 {
		t.Errorf("Expected Status 1, got %d", *req.Status)
	}
}
