package controller

import (
	"go-web/internal/dao"
	"go-web/internal/moudels"
	"go-web/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// 认证相关处理器
func RegisterHandler(c *gin.Context) {
	var req models.RegisterRequest

	// 绑定并验证请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorw("注册请求参数验证失败", "error", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	db := dao.GetDB()

	// 检查用户名是否已存在
	var existingUser models.User
	if err := db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		logger.Warnw("用户名已存在", "username", req.Username)
		c.JSON(http.StatusConflict, models.ErrorResponse{
			Code:    http.StatusConflict,
			Message: "用户名已存在",
		})
		return
	}

	// 检查邮箱是否已存在
	if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		logger.Warnw("邮箱已存在", "email", req.Email)
		c.JSON(http.StatusConflict, models.ErrorResponse{
			Code:    http.StatusConflict,
			Message: "邮箱已被注册",
		})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorw("密码加密失败", "error", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "密码加密失败",
			Error:   err.Error(),
		})
		return
	}

	// 设置用户类型，如果未提供则默认为app
	userType := models.UserTypeApp
	if req.UserType != "" {
		userType = models.UserType(req.UserType)
	}

	// 创建新用户
	newUser := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		UserType:     userType,
		Status:       models.UserStatusActive, // 默认激活状态
	}

	if err := db.Create(&newUser).Error; err != nil {
		logger.Errorw("创建用户失败", "error", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "创建用户失败",
			Error:   err.Error(),
		})
		return
	}

	// 清除密码字段后返回用户信息
	newUser.PasswordHash = ""

	logger.Infow("用户注册成功", "user_id", newUser.ID, "username", newUser.Username)

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "注册成功",
		Data:    newUser,
	})
}

func LoginHandler(c *gin.Context) {
	var req models.LoginRequest

	// 绑定并验证请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorw("登录请求参数验证失败", "error", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	// 查询用户（支持用户名或邮箱登录）
	var user models.User
	db := dao.GetDB()

	result := db.Where("username = ? OR email = ?", req.Username, req.Username).First(&user)
	if result.Error != nil {
		logger.Warnw("用户不存在", "username", req.Username)
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "用户名或密码错误",
		})
		return
	}

	// 检查用户状态
	if user.Status != models.UserStatusActive {
		logger.Warnw("用户账户未激活", "user_id", user.ID, "status", user.Status)
		c.JSON(http.StatusForbidden, models.ErrorResponse{
			Code:    http.StatusForbidden,
			Message: "账户未激活，无法登录",
		})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		logger.Warnw("密码错误", "user_id", user.ID)
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "用户名或密码错误",
		})
		return
	}

	// 清除密码字段后返回用户信息
	user.PasswordHash = ""

	logger.Infow("用户登录成功", "user_id", user.ID, "username", user.Username)

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "登录成功",
		Data:    user,
	})
}

func LogoutHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":  "登出功能待实现",
		"endpoint": "POST /api/v1/auth/logout",
	})
}

func RefreshTokenHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":  "刷新令牌功能待实现",
		"endpoint": "POST /api/v1/auth/refresh",
	})
}

func ForgotPasswordHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":  "忘记密码功能待实现",
		"endpoint": "POST /api/v1/auth/forgot-password",
	})
}

func ResetPasswordHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":  "重置密码功能待实现",
		"endpoint": "POST /api/v1/auth/reset-password",
	})
}

// 用户相关处理器
func GetProfileHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":  "获取用户信息功能待实现",
		"endpoint": "GET /api/v1/users/profile",
	})
}

func UpdateProfileHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":  "更新用户信息功能待实现",
		"endpoint": "PUT /api/v1/users/profile",
	})
}

func ChangePasswordHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":  "修改密码功能待实现",
		"endpoint": "PUT /api/v1/users/password",
	})
}

// 健康检查
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Server is running",
	})
}
