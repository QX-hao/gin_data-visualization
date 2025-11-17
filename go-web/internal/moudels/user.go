package models

import (
    "time"
)

// 用户类型
type UserType string

const (
    UserTypeSystem UserType = "system"
    UserTypeApp    UserType = "app"
)

// 用户状态
type UserStatus string

const (
    UserStatusActive   UserStatus = "active"
    UserStatusInactive UserStatus = "inactive"
)

// User 用户模型
type User struct {
    ID           uint       `json:"id" gorm:"primaryKey;autoIncrement"`
    Username     string     `json:"username" gorm:"size:50;uniqueIndex;not null"`
    Email        string     `json:"email" gorm:"size:100;uniqueIndex;not null"`
    PasswordHash string     `json:"-" gorm:"size:255;not null"`
    UserType     UserType   `json:"user_type" gorm:"type:ENUM('system','app');not null;default:'app'"`
    Status       UserStatus `json:"status" gorm:"type:ENUM('active','inactive');not null;default:'active'"`
    CreatedAt    time.Time  `json:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at"`
}

// UserSession 用户会话模型
type UserSession struct {
    ID                uint      `json:"id" gorm:"primaryKey;autoIncrement"`
    UserID            uint      `json:"user_id" gorm:"not null"`
    Token             string    `json:"token" gorm:"size:512;not null"`
    RefreshToken      string    `json:"refresh_token" gorm:"size:255;not null"`
    ExpiresAt         time.Time `json:"expires_at" gorm:"not null"`
    RefreshExpiresAt  time.Time `json:"refresh_expires_at" gorm:"not null"`
    UserAgent         string    `json:"user_agent,omitempty" gorm:"type:text"`
    IPAddress         string    `json:"ip_address,omitempty" gorm:"size:45"`
    CreatedAt         time.Time `json:"created_at"`
    
    User              User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
}