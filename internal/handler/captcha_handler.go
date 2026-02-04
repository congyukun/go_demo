package handler

import (
	"go_demo/internal/models"
	"go_demo/internal/utils"
	"go_demo/pkg/captcha"

	"github.com/gin-gonic/gin"
)

// CaptchaHandler 验证码处理器
type CaptchaHandler struct {
	captchaService captcha.CaptchaService
}

// NewCaptchaHandler 创建验证码处理器实例
func NewCaptchaHandler(captchaService captcha.CaptchaService) *CaptchaHandler {
	return &CaptchaHandler{
		captchaService: captchaService,
	}
}

// GetCaptcha 获取验证码
// @Summary 获取验证码
// @Description 获取图形验证码，返回验证码ID和Base64编码的图片
// @Tags 验证码
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=models.CaptchaResponse} "获取成功"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/captcha [get]
func (h *CaptchaHandler) GetCaptcha(c *gin.Context) {
	id, b64s, err := h.captchaService.Generate()
	if err != nil {
		utils.ResponseError(c, 500, "生成验证码失败")
		return
	}

	response := models.CaptchaResponse{
		CaptchaID: id,
		Image:     b64s,
	}

	utils.ResponseSuccess(c, "获取验证码成功", response)
}

// VerifyCaptcha 验证验证码（仅用于测试）
// @Summary 验证验证码
// @Description 验证验证码是否正确（仅用于测试环境）
// @Tags 验证码
// @Accept json
// @Produce json
// @Param captcha_id query string true "验证码ID"
// @Param captcha query string true "验证码"
// @Success 200 {object} utils.Response "验证成功"
// @Failure 400 {object} utils.Response "验证失败"
// @Router /api/v1/captcha/verify [get]
func (h *CaptchaHandler) VerifyCaptcha(c *gin.Context) {
	captchaID := c.Query("captcha_id")
	captchaCode := c.Query("captcha")

	if captchaID == "" || captchaCode == "" {
		utils.ResponseError(c, 400, "验证码ID和验证码不能为空")
		return
	}

	// 使用不清除的验证方式（仅用于测试）
	if h.captchaService.VerifyWithoutClear(captchaID, captchaCode) {
		utils.ResponseSuccess(c, "验证码正确", nil)
	} else {
		utils.ResponseError(c, 400, "验证码错误或已过期")
	}
}

// GetCaptchaService 获取验证码服务（供其他处理器使用）
func (h *CaptchaHandler) GetCaptchaService() captcha.CaptchaService {
	return h.captchaService
}
