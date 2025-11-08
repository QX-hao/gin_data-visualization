package api

import (
	"gin_data-visualization/controller"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthAPI 认证API接口
type AuthAPI struct {
	authController *controller.AuthController
}

// NewAuthAPI 创建认证API实例
func NewAuthAPI(db *gorm.DB) *AuthAPI {
	authController := controller.NewAuthController(db)
	return &AuthAPI{
		authController: authController,
	}
}

// RegisterRoutes 注册认证路由
func (a *AuthAPI) RegisterRoutes(router *gin.RouterGroup) {
	authGroup := router.Group("/auth")
	{
		// 公开路由（无需认证）
		authGroup.POST("/login", a.Login)
		authGroup.POST("/register", a.Register)
		
		// 需要认证的路由
		authGroup.POST("/logout", a.Logout)
		authGroup.GET("/userinfo", a.GetUserInfo)
	}
}

// Login 登录API
// @Summary 用户登录
// @Description 用户登录接口
// @Tags 认证
// @Accept json
// @Produce json
// @Param loginRequest body model.LoginRequest true "登录请求参数"
// @Success 200 {object} model.LoginResponse "登录成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "用户名或密码错误"
// @Router /auth/login [post]
func (a *AuthAPI) Login(c *gin.Context) {
	a.authController.Login(c)
}

// Register 注册API
// @Summary 用户注册
// @Description 用户注册接口
// @Tags 认证
// @Accept json
// @Produce json
// @Param registerRequest body model.RegisterRequest true "注册请求参数"
// @Success 200 {object} model.User "注册成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误或用户已存在"
// @Router /auth/register [post]
func (a *AuthAPI) Register(c *gin.Context) {
	a.authController.Register(c)
}

// Logout 登出API
// @Summary 用户登出
// @Description 用户登出接口
// @Tags 认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "登出成功"
// @Router /auth/logout [post]
func (a *AuthAPI) Logout(c *gin.Context) {
	a.authController.Logout(c)
}

// GetUserInfo 获取用户信息API
// @Summary 获取用户信息
// @Description 获取当前登录用户信息
// @Tags 认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "获取用户信息成功"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Router /auth/userinfo [get]
func (a *AuthAPI) GetUserInfo(c *gin.Context) {
	a.authController.GetUserInfo(c)
}

// HealthCheck 健康检查API
// @Summary 健康检查
// @Description 服务健康检查接口
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "服务正常"
// @Router /health [get]
func (a *AuthAPI) HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"code":    200,
		"message": "服务运行正常",
		"data":    nil,
	})
}