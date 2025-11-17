package models

// RegisterRequest 注册请求
type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    UserType string `json:"user_type" binding:"oneof=system app"` // 可选
}

// LoginRequest 登录请求
type LoginRequest struct {
    Username string `json:"username" binding:"required"` // 可以是用户名或邮箱
    Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
    RefreshToken string `json:"refresh_token" binding:"required"`
}