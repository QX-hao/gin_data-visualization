package controller

import (
	"net/http"
	"gin_data-visualization/model"
	"gin_data-visualization/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthController 认证控制器
type AuthController struct {
	authService service.AuthService
}

// NewAuthController 创建认证控制器实例
func NewAuthController(db *gorm.DB) *AuthController {
	authService := service.NewAuthService(db)
	return &AuthController{
		authService: authService,
	}
}

// Login 用户登录
func (c *AuthController) Login(ctx *gin.Context) {
	var req model.LoginRequest
	
	// 绑定请求参数
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}
	
	// 调用服务层进行登录
	user, token, err := c.authService.Login(req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	
	// 返回登录成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data": model.LoginResponse{
			Token:    token,
			UserInfo: *user,
		},
	})
}

// Register 用户注册
func (c *AuthController) Register(ctx *gin.Context) {
	var req model.RegisterRequest
	
	// 绑定请求参数
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}
	
	// 调用服务层进行注册
	user, err := c.authService.Register(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	
	// 返回注册成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "注册成功",
		"data":    user,
	})
}

// Logout 用户登出
func (c *AuthController) Logout(ctx *gin.Context) {
	// 在实际项目中，这里可能需要处理token黑名单等逻辑
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登出成功",
		"data":    nil,
	})
}

// GetUserInfo 获取用户信息
func (c *AuthController) GetUserInfo(ctx *gin.Context) {
	// 从上下文中获取用户ID（需要中间件设置）
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
			"data":    nil,
		})
		return
	}
	
	// 在实际项目中，这里需要查询数据库获取用户信息
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取用户信息成功",
		"data": gin.H{
			"user_id": userID,
		},
	})
}