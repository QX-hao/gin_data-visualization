package models

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	PasswordHash string `json:"password" binding:"required"`                    // SHA-256 哈希值
	// UserType string `json:"user_type" binding:"omitempty,oneof=system app"` // 可选，默认为'app'
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名或邮箱
	PasswordHash string `json:"password" binding:"required"` // SHA-256 哈希值
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ForgotPasswordRequest 忘记密码请求
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"` // SHA-256 哈希值
}

// UpdateProfileRequest 更新用户资料请求
type UpdateProfileRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=20"`
	Email    string `json:"email" binding:"omitempty,email"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword    string `json:"old_password" binding:"required"`     // SHA-256 哈希值
	NewPassword    string `json:"new_password" binding:"required"`     // SHA-256 哈希值
	TargetUsername string `json:"target_username" binding:"omitempty"` // 可选，仅管理员
}

// CheckUsernameRequest 检查用户名可用性请求
type CheckUsernameRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
}

// CheckEmailRequest 检查邮箱可用性请求
type CheckEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}
