package models

// AuthResponse 认证响应
type AuthResponse struct {
    Message      string `json:"message"`
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token,omitempty"`
    User         *User  `json:"user"`
}

// TokenResponse 令牌响应
type TokenResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresIn    int64  `json:"expires_in"`
}

// SuccessResponse 成功响应
type SuccessResponse struct {
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Error   string `json:"error,omitempty"`
}