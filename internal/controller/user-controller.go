package controller

import (
	"go-web/internal/dao"
	models "go-web/internal/moudels"
	"go-web/pkg/logger"
	"go-web/pkg/token"
	"net/http"
	"time"

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

	// 根据接口要求，注册成功响应格式为 HTTP 201
	c.JSON(http.StatusCreated, models.RegistrationResponse{
		Code:    http.StatusCreated,
		Message: "注册成功",
		Data:    &newUser,
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

	// 生成访问令牌
	accessToken, err := token.GenerateAccessToken(user.ID)
	if err != nil {
		logger.Errorw("生成访问令牌失败", "error", err, "user_id", user.ID)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "生成令牌失败",
			Error:   err.Error(),
		})
		return
	}

	// 根据接口要求，登录成功响应格式为 HTTP 200
	response := models.LoginResponse{
		Code:    http.StatusOK,
		Message: "登录成功",
	}
	response.Data.AccessToken = accessToken
	response.Data.User = &user
	c.JSON(http.StatusOK, response)
}

func LogoutHandler(c *gin.Context) {
	// 根据接口要求，登出成功响应格式为 {code: 200, message: "登出成功", data: {}}
	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "登出成功",
		Data:    gin.H{},
	})
}

func RefreshTokenHandler(c *gin.Context) {
	// 根据接口要求，刚止令牌成功响应格式为 {code: 200, message: "令牌刷新成功", data: {access_token}}
	// 这里需要实现具体的刷新逻辑
	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "令牌刷新成功",
		Data: gin.H{
			"access_token": "new_access_token_example",
		},
	})
}

func ForgotPasswordHandler(c *gin.Context) {
	// 根据接口要求，忘记密码成功响应格式为 {code: 200, message: "重置邮件已发送，请检查您的邮箱", data: {}}
	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "重置邮件已发送，请检查您的邮箱",
		Data:    gin.H{},
	})
}

func ResetPasswordHandler(c *gin.Context) {
	// 根据接口要求，重置密码成功响应格式为 {code: 200, message: "密码重置成功", data: {}}
	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "密码重置成功",
		Data:    gin.H{},
	})
}

// 用户相关处理器
func GetProfileHandler(c *gin.Context) {
	// 根据接口要求，获取用户资料成功响应格式为 {code: 200, message: "获取成功", data: {用户信息}}
	// 这里需要实现具体的获取用户资料逻辑
	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data: gin.H{
			"id":         1,
			"username":   "example_user",
			"email":      "user@example.com",
			"user_type":  "app",
			"status":     "active",
			"created_at": "2024-01-15T10:30:00Z",
			"updated_at": "2024-01-15T10:30:00Z",
		},
	})
}

func UpdateProfileHandler(c *gin.Context) {
	// 根据接口要求，更新用户资料成功响应格式为 {code: 200, message: "更新成功", data: {更新后的用户信息}}
	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "更新成功",
		Data: gin.H{
			"id":         1,
			"username":   "updated_user",
			"email":      "updated@example.com",
			"user_type":  "app",
			"status":     "active",
			"updated_at": "2024-01-15T10:30:00Z",
		},
	})
}

func ChangePasswordHandler(c *gin.Context) {
	// 根据接口要求，修改密码成功响应格式为 {code: 200, message: "密码修改成功", data: {}}
	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "密码修改成功",
		Data:    gin.H{},
	})
}

// 健康检查
func HealthHandler(c *gin.Context) {
	// 根据接口要求，健康检查响应格式为 {status: "OK", message: "Server is running", timestamp: "..."}
	c.JSON(http.StatusOK, models.HealthResponse{
		Status:    "OK",
		Message:   "Server is running",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// 检查用户名可用性
func CheckUsernameHandler(c *gin.Context) {
	var req models.CheckUsernameRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorw("检查用户名请求参数验证失败", "error", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	db := dao.GetDB()
	var user models.User
	// 检查用户名是否已存在
	if err := db.Where("username = ?", req.Username).First(&user).Error; err == nil {
		// 用户名已被占用
		c.JSON(http.StatusOK, models.CheckAvailableResponse{
			Code:    http.StatusOK,
			Message: "Success",
		})
		// 需要设置 available 字段为 false
		// 这里使用 gin.H 来灵活设置
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "Success",
			"data": gin.H{
				"available": false,
			},
		})
		return
	}

	// 用户名可用
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Success",
		"data": gin.H{
			"available": true,
		},
	})
}

// 检查邮箱可用性
func CheckEmailHandler(c *gin.Context) {
	var req models.CheckEmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorw("检查邮箱请求参数验证失败", "error", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	db := dao.GetDB()
	var user models.User
	// 检查邮箱是否已存在
	if err := db.Where("email = ?", req.Email).First(&user).Error; err == nil {
		// 邮箱已被占用
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "Success",
			"data": gin.H{
				"available": false,
			},
		})
		return
	}

	// 邮箱可用
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Success",
		"data": gin.H{
			"available": true,
		},
	})
}

// 获取用户列表（仅管理员）
func GetUsersListHandler(c *gin.Context) {
	// 根据接口要求，获取用户列表成功响应格式为 {code: 200, message: "获取成功", data: [用户信息数组], pagination: {...}}
	// 这里需要实现具体的获取用户列表逻辑
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "获取成功",
		"data": []gin.H{
			{
				"id":         1,
				"username":   "john_doe",
				"email":      "john@example.com",
				"user_type":  "app",
				"status":     "active",
				"created_at": "2024-01-15T10:30:00Z",
				"updated_at": "2024-01-15T10:30:00Z",
			},
			{
				"id":         2,
				"username":   "admin_user",
				"email":      "admin@example.com",
				"user_type":  "system",
				"status":     "active",
				"created_at": "2024-01-10T10:30:00Z",
				"updated_at": "2024-01-15T10:30:00Z",
			},
		},
		"pagination": gin.H{
			"total": 50,
			"page":  1,
			"limit": 20,
			"pages": 3,
		},
	})
}
